package tictactoe

import (
	"net/http"
	"sync"

	"github.com/karlseguin/ccache"
	"github.com/matthewdale/matthewrdale.com/dowithgo/tictactoe/game"
	"github.com/matthewdale/matthewrdale.com/random"
	"github.com/pkg/errors"
)

type Service struct {
	games        *ccache.Cache
	latestGame   *game.Game
	latestGameID string
	joinMu       sync.Mutex
}

func New() *Service {
	return &Service{
		games: ccache.New(ccache.Configure().MaxSize(1000)),
	}
}

type JoinArgs struct{}

type JoinReply struct {
	GameID   string
	Board    game.Board
	PlayerID string
	Piece    game.Piece
}

func (svc *Service) Join(r *http.Request, args *JoinArgs, reply *JoinReply) error {
	// Only one player can join a new game at a time, so synchronize the entire
	// join function.
	svc.joinMu.Lock()
	defer svc.joinMu.Unlock()

	if svc.latestGame == nil || svc.latestGame.Full() {
		gameID := random.String(20)
		g := game.New()
		svc.games.Set(gameID, g, 0)
		svc.latestGame = g
		svc.latestGameID = gameID
	}

	playerID, piece, err := svc.latestGame.NewPlayer()
	if err != nil {
		return errors.WithMessage(err, "error adding new player to game")
	}

	// Return the game state to the player who just joined the game.
	reply.GameID = svc.latestGameID
	reply.Board = svc.latestGame.Board()
	reply.PlayerID = playerID
	reply.Piece = piece
	return nil
}

type PlaceArgs struct {
	GameID   string
	PlayerID string
	Position [2]int
}

type PlaceReply struct {
	Board  game.Board
	Winner string
	Over   bool
}

func (svc *Service) Place(r *http.Request, args *PlaceArgs, reply *PlaceReply) error {
	item := svc.games.Get(args.GameID)
	if item == nil {
		return errors.New("game does not exist")
	}

	g := item.Value().(*game.Game)
	err := g.Place(args.PlayerID, args.Position)
	if err != nil {
		return errors.WithMessage(err, "error placing piece in game")
	}
	reply.Board = g.Board()
	reply.Winner = g.Winner()
	reply.Over = g.Over()
	return nil
}

type GetGameArgs struct {
	GameID string
}

type GetGameReply struct {
	Board  game.Board
	Winner string
	Over   bool
}

func (svc *Service) GetGame(r *http.Request, args *GetGameArgs, reply *GetGameReply) error {
	item := svc.games.Get(args.GameID)
	if item == nil {
		return errors.New("game does not exist")
	}

	g := item.Value().(*game.Game)
	reply.Board = g.Board()
	reply.Winner = g.Winner()
	reply.Over = g.Over()
	return nil
}
