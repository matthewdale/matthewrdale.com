package messenger

import "net/http"

type Service struct {
	messages chan int64
}

func New() *Service {
	return &Service{
		messages: make(chan int64, 100),
	}
}

type SendArgs struct {
	Number int64
}

type SendReply struct {
	Sent bool
}

// TODO: Are there privacy concerns here, even with numbers? You can store anything
// in strings of numbers...
func (svc *Service) Send(r *http.Request, args *SendArgs, reply *SendReply) error {
	select {
	case svc.messages <- args.Number:
		reply.Sent = true
	default:
	}
	return nil
}

type ReceiveArgs struct{}

type ReceiveReply struct {
	Number   int64
	Received bool
}

func (svc *Service) Receive(r *http.Request, args *ReceiveArgs, reply *ReceiveReply) error {
	select {
	case reply.Number = <-svc.messages:
		reply.Received = true
	default:
	}
	return nil
}
