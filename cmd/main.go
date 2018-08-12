package main

import (
	"log"
	"net"
	"net/http"
	"time"

	"github.com/gorilla/rpc/v2"
	"github.com/gorilla/rpc/v2/json"
	"github.com/matthewdale/matthewrdale.com/dowithgo/ethereum"
	"github.com/matthewdale/matthewrdale.com/dowithgo/flipflop"
	"github.com/matthewdale/matthewrdale.com/dowithgo/messenger"
	"github.com/matthewdale/matthewrdale.com/dowithgo/metrics"
	"github.com/matthewdale/matthewrdale.com/dowithgo/pprof"
	"github.com/matthewdale/matthewrdale.com/dowithgo/proxy"
	"github.com/matthewdale/matthewrdale.com/dowithgo/ratelimit"
	"github.com/matthewdale/matthewrdale.com/dowithgo/shuffle"
	"github.com/matthewdale/matthewrdale.com/dowithgo/tictactoe"
	"github.com/matthewdale/matthewrdale.com/random"
	gometrics "github.com/rcrowley/go-metrics"
	"github.com/urfave/negroni"
	"golang.org/x/net/xsrftoken"
)

const xsrfCookieName = "xsrf"

var xsrfKey = random.String(60)

func setXSRF(writer http.ResponseWriter, request *http.Request, next http.HandlerFunc) {
	next(writer, request)

	host, _, err := net.SplitHostPort(request.RemoteAddr)
	if err != nil {
		log.Printf("Error parsing host from remote addr for setting XSRF cookie: %s", err)
		return
	}
	user := host + request.UserAgent()
	http.SetCookie(writer, &http.Cookie{
		Name:    xsrfCookieName,
		Value:   xsrftoken.Generate(xsrfKey, user, ""),
		Expires: time.Now().Add(xsrftoken.Timeout),
	})
}

func validateXSRF(next http.Handler) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		host, _, err := net.SplitHostPort(request.RemoteAddr)
		if err != nil {
			log.Printf("Error parsing host from addr for XSRF validation: %s", err)
			http.Error(writer, "error parsing host from addr", http.StatusInternalServerError)
			return
		}
		user := host + request.UserAgent()
		cookie, err := request.Cookie(xsrfCookieName)
		if err != nil {
			http.Error(writer, "failed to get XSRF token", http.StatusUnauthorized)
			return
		}
		if !xsrftoken.Valid(cookie.Value, xsrfKey, user, "") {
			http.Error(writer, "invalid XSRF token", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(writer, request)
	}
}

var registry = gometrics.NewRegistry()

func metricsMiddleware(writer http.ResponseWriter, request *http.Request, next http.HandlerFunc) {
	start := time.Now()
	next(writer, request)
	gometrics.GetOrRegisterTimer("http.request", registry).UpdateSince(start)
}

// TODO: This should be unnecessary. Figure out why rpc doesn't add it.
func jsonContentType(next http.Handler) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(writer, request)
	}
}

func main() {
	mux := http.NewServeMux()

	// RPC service for "Things You Can Do With Go"
	doWithGoRPC := rpc.NewServer()
	doWithGoRPC.RegisterCodec(json.NewCodec(), "application/json")
	doWithGoRPC.RegisterService(ethereum.New(), "ethereum")
	doWithGoRPC.RegisterService(flipflop.New(), "flipflop")
	doWithGoRPC.RegisterService(messenger.New(), "messenger")
	doWithGoRPC.RegisterService(metrics.New(registry), "metrics")
	doWithGoRPC.RegisterService(pprof.New(), "pprof")
	doWithGoRPC.RegisterService(proxy.New(), "proxy")
	doWithGoRPC.RegisterService(ratelimit.New(), "ratelimit")
	doWithGoRPC.RegisterService(shuffle.New(), "shuffle")
	doWithGoRPC.RegisterService(tictactoe.New(), "tictactoe")
	mux.Handle("/rpc/dowithgo", validateXSRF(jsonContentType(doWithGoRPC)))

	// TODO: Replace Negroni file server with a better solution, maybe CDN-based?
	stack := negroni.New()
	stack.Use(negroni.NewRecovery())
	stack.Use(negroni.NewLogger())
	stack.UseFunc(metricsMiddleware)
	// TODO: Reenable.
	// stack.UseFunc(setXSRF)
	stack.Use(negroni.NewStatic(http.Dir("public")))
	stack.UseHandler(mux)
	http.ListenAndServe(":80", stack)
}
