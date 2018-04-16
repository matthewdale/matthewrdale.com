package flipflop

import (
	"net/http"
	"sync/atomic"
)

type Service struct {
	count uint64
}

func New() *Service {
	return &Service{}
}

type FlipFlopArgs struct{}

type FlipFlopReply struct {
	On    bool
	Count uint64
}

func (svc *Service) FlipFlop(r *http.Request, args *FlipFlopArgs, reply *FlipFlopReply) error {
	newCount := atomic.AddUint64(&svc.count, 1)
	reply.On = newCount%2 == 0
	reply.Count = newCount
	return nil
}
