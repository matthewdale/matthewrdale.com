package ratelimit

import (
	"net/http"
	"time"
)

type Service struct {
	limiter chan struct{}
}

func New() *Service {
	svc := &Service{
		limiter: make(chan struct{}, 10),
	}
	go svc.drain()
	return svc
}

func (svc *Service) drain() {
	for range time.Tick(2 * time.Second) {
		select {
		case <-svc.limiter:
		default:
		}
	}
}

type IncrArgs struct{}

type IncrReply struct {
	Limited bool
}

func (svc *Service) Incr(r *http.Request, args *IncrArgs, reply *IncrReply) error {
	select {
	case svc.limiter <- struct{}{}:
	default:
		reply.Limited = true
	}
	return nil
}
