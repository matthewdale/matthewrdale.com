package main

import (
	"net"
	"net/http"
	"strings"
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
	"github.com/pkg/errors"
	gometrics "github.com/rcrowley/go-metrics"
	"github.com/urfave/negroni"
	"golang.org/x/net/xsrftoken"
)

const xsrfCookieName = "xsrf"

var xsrfKey = random.String(60)

func xsrfMiddleware(writer http.ResponseWriter, request *http.Request, next http.HandlerFunc) {
	host, _, err := net.SplitHostPort(request.RemoteAddr)
	if err != nil {
		http.Error(
			writer,
			errors.WithMessage(err, "error parsing host from remote addr").Error(),
			http.StatusInternalServerError)
		return
	}
	user := host + request.UserAgent()
	cookie, _ := request.Cookie(xsrfCookieName)
	if cookie == nil {
		http.SetCookie(writer, &http.Cookie{
			Name:    xsrfCookieName,
			Value:   xsrftoken.Generate(xsrfKey, user, ""),
			Expires: time.Now().Add(xsrftoken.Timeout),
		})
	}

	// All RPC requests require an XSRF token cookie. Enforce that the cookie
	// exists and that the token is valid. If it's not, return unauthorized.
	if strings.HasPrefix(request.URL.Path, "/rpc") {
		if cookie == nil || !xsrftoken.Valid(cookie.Value, xsrfKey, user, "") {
			http.Error(writer, "missing XSRF token", http.StatusUnauthorized)
			return
		}
	}
	next(writer, request)
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

	// RPC service for "10 Things You Can Do With Go"
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
	mux.Handle("/rpc/dowithgo", jsonContentType(doWithGoRPC))

	// TODO: Replace Negroni file server with a better solution, maybe CDN-based?
	stack := negroni.New()
	stack.Use(negroni.NewRecovery())
	stack.Use(negroni.NewLogger())
	stack.UseFunc(metricsMiddleware)
	// TODO: Reenable.
	// stack.UseFunc(xsrfMiddleware)
	stack.Use(negroni.NewStatic(http.Dir("public")))
	stack.UseHandler(mux)
	http.ListenAndServe(":8080", stack)
}
