package game

import (
	"errors"
	"sync"

	"github.com/matthewdale/matthewrdale.com/random"
)

// Piece is a tic-tac-toe piece, represented as a 1 or -1 for X and O, respectively.
type Piece = int8

const (
	// X is the X tic-tac-toe game piece.
	X Piece = 1
	// O is the O tic-tac-toe game piece.
	O Piece = -1
)

// Board is a 3x3 tic-tac-toe game board. Each space in the game board is either
// 0 (no piece), 1 (X), or -1 (O).
type Board [3][3]Piece

// Game is a tic-tac-toe game.
type Game struct {
	board     Board
	over      bool
	winner    string
	players   map[string]Piece
	nextPiece Piece
	mu        sync.RWMutex
}

// New returns a new tic-tac-toe game.
func New() *Game {
	return &Game{
		players:   make(map[string]Piece, 2),
		nextPiece: X,
	}
}

// Board returns a copy of the game board. Each space in the game board is either
// 0 (no piece), 1 (X), or -1 (O).
func (g *Game) Board() Board {
	g.mu.RLock()
	defer g.mu.RUnlock()

	var b Board
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			b[i][j] = g.board[i][j]
		}
	}
	return b
}

// Over returns true if the game is over. The game is over when there is a winner
// or when there is a piece placed on every position in the board.
func (g *Game) Over() bool {
	g.mu.RLock()
	defer g.mu.RUnlock()

	return g.over
}

// Winner returns the player ID of the game winner or 0 if there is no winner.
// If the game is not over and there is no winner, there may be a winner in the
// future. If the game is over and there is no winner, the game is a draw.
func (g *Game) Winner() string {
	g.mu.RLock()
	defer g.mu.RUnlock()

	return g.winner
}

// Full returns true if the game has two players.
func (g *Game) Full() bool {
	g.mu.RLock()
	defer g.mu.RUnlock()

	return len(g.players) == 2
}

// NewPlayer adds a new player to the game and returns the player ID, the game
// piece assigned to the player (either an 'X' or an 'O', represented by 1 and
// -1, respectively). Once there are two players, NewPlayer will always return
// an error.
func (g *Game) NewPlayer() (string, Piece, error) {
	if g.Full() {
		return "", 0, errors.New("game is full")
	}

	g.mu.Lock()
	defer g.mu.Unlock()

	// Player 0 always gets "X", player 1 always gets "O".
	piece := X
	if len(g.players) == 1 {
		piece = O
	}
	playerID := random.String(10)
	g.players[playerID] = piece
	return playerID, piece, nil
}

// Place adds a player's game piece to the given position. The position is a tuple
// of ints where the first int indexes the row and the second int indexes the column.
// An error is returned if:
//   - the game is over
//   - the give position is not inside the game board
//   - the player does not exist
//   - the same player attempts to place a piece twice in a row
//   - the position is already occupied by a game piece
// The game over and winner states are checked after every game piece is placed.
func (g *Game) Place(playerID string, position [2]int) error {
	if g.Over() {
		return errors.New("game is over")
	}
	if position[0] < 0 || position[0] >= 3 || position[1] < 0 || position[1] >= 3 {
		return errors.New("position is outside game board")
	}

	g.mu.Lock()
	defer g.mu.Unlock()

	piece, found := g.players[playerID]
	if !found {
		return errors.New("player does not exist")
	}
	if piece != g.nextPiece {
		return errors.New("you cannot place twice in a row")
	}
	if g.board[position[0]][position[1]] != 0 {
		return errors.New("position is already occupied")
	}
	g.board[position[0]][position[1]] = piece

	if hasWinner(g.board) {
		g.winner = playerID
		g.over = true
		return nil
	}
	if isComplete(g.board) {
		g.over = true
		return nil
	}

	if piece == X {
		g.nextPiece = O
	} else {
		g.nextPiece = X
	}
	return nil
}

// hasWinner returns true if the provided game board has a winner.
func hasWinner(b Board) bool {
	var colSums [3]int8
	var diagSums [2]int8
	for i := 0; i < 3; i++ {
		var rowSum int8
		for j := 0; j < 3; j++ {
			rowSum += b[i][j]
			colSums[j] += b[i][j]
			if i == j {
				diagSums[0] += b[i][j]
			}
			if i == 2-j {
				diagSums[1] += b[i][j]
			}
		}
		if rowSum == 3 || rowSum == -3 {
			return true
		}
	}
	for i := 0; i < len(colSums); i++ {
		if colSums[i] == 3 || colSums[i] == -3 {
			return true
		}
	}
	for i := 0; i < len(diagSums); i++ {
		if diagSums[i] == 3 || diagSums[i] == -3 {
			return true
		}
	}
	return false
}

// isComplete returns true if the provided game board is completely full.
func isComplete(b Board) bool {
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			if b[i][j] == 0 {
				return false
			}
		}
	}
	return true
}
