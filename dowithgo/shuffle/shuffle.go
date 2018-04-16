package shuffle

import (
	"math/rand"
	"net/http"

	// Import the random package to make sure the math/rand package-level
	// pseudo-random number generator is seeded.
	_ "github.com/matthewdale/matthewrdale.com/random"
)

type Service struct{}

func New() *Service {
	return &Service{}
}

type ShuffleArgs struct {
	List []string
}

type ShuffleReply struct {
	List []string
}

func (svc *Service) Shuffle(r *http.Request, args *ShuffleArgs, reply *ShuffleReply) error {
	if len(args.List) > 100 {
		args.List = args.List[:100]
	}
	rand.Shuffle(len(args.List), func(i, j int) {
		args.List[i], args.List[j] = args.List[j], args.List[i]
	})
	reply.List = args.List
	return nil
}
