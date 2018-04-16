package proxy

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	"github.com/pkg/errors"
)

var domainWhitelist = map[string]bool{
	"golang.org":       true,
	"matthewrdale.com": true,
}

type Service struct{}

func New() *Service {
	return &Service{}
}

type ProxyArgs struct {
	URL string
}

type ProxyReply struct {
	StatusCode int
	Body       string
}

var client = &http.Client{Timeout: 10 * time.Second}

func (svc *Service) Get(r *http.Request, args *ProxyArgs, reply *ProxyReply) error {
	u, err := url.Parse(args.URL)
	if err != nil {
		return errors.WithMessage(err, "error parsing URL")
	}
	if !domainWhitelist[u.Host] {
		return fmt.Errorf("domain %q is not in whitelist: %#v", u.Host, domainWhitelist)
	}
	res, err := client.Get(u.String())
	if err != nil {
		return errors.WithMessage(err, "error getting result")
	}
	reply.StatusCode = res.StatusCode
	body, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		return errors.WithMessage(err, "error reading body")
	}
	reply.Body = string(body)
	return nil
}
