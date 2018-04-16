package game

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHasWinner(t *testing.T) {
	tests := []struct {
		description string
		board       Board
		expected    bool
	}{
		{
			description: "An empty board has no winner",
			board: Board{
				[3]int8{0, 0, 0},
				[3]int8{0, 0, 0},
				[3]int8{0, 0, 0},
			},
			expected: false,
		},
		{
			description: "A full board with a column full of 'X' has a winner",
			board: Board{
				[3]int8{X, O, O},
				[3]int8{X, O, X},
				[3]int8{X, X, O},
			},
			expected: true,
		},
		{
			description: "A full board with a column full of 'O' has a winner",
			board: Board{
				[3]int8{X, O, O},
				[3]int8{O, O, X},
				[3]int8{X, O, X},
			},
			expected: true,
		},
		{
			description: "A full board with a row full of 'X' has a winner",
			board: Board{
				[3]int8{X, X, X},
				[3]int8{O, O, X},
				[3]int8{X, O, O},
			},
			expected: true,
		},
		{
			description: "A full board with a row full of 'O' has a winner",
			board: Board{
				[3]int8{X, O, X},
				[3]int8{O, X, X},
				[3]int8{O, O, O},
			},
			expected: true,
		},
		{
			description: "A full board with the left-to-right diagonal full of 'X' has a winner",
			board: Board{
				[3]int8{X, O, O},
				[3]int8{O, X, X},
				[3]int8{X, O, X},
			},
			expected: true,
		},
		{
			description: "A full board with the right-to-left diagonal full of 'X' has a winner",
			board: Board{
				[3]int8{X, O, X},
				[3]int8{O, X, X},
				[3]int8{X, O, O},
			},
			expected: true,
		},
		{
			description: "A partial board with a column full of 'O' has a winner",
			board: Board{
				[3]int8{X, O, X},
				[3]int8{0, O, 0},
				[3]int8{X, O, 0},
			},
			expected: true,
		},
		{
			description: "A partial board with a row full of 'X' has a winner",
			board: Board{
				[3]int8{X, X, X},
				[3]int8{O, 0, 0},
				[3]int8{0, O, O},
			},
			expected: true,
		},
		{
			description: "A partial board with a left-to-right diagnoal full of 'O' has a winner",
			board: Board{
				[3]int8{O, X, X},
				[3]int8{0, O, 0},
				[3]int8{0, X, O},
			},
			expected: true,
		},
		{
			description: "A full board with no columns, rows, or diagnoals with the same piece has no winner",
			board: Board{
				[3]int8{X, X, O},
				[3]int8{O, X, X},
				[3]int8{X, O, O},
			},
			expected: false,
		},
	}
	for _, test := range tests {
		test := test // Capture range variable.
		t.Run(test.description, func(t *testing.T) {
			assert.Equal(t,
				test.expected,
				hasWinner(test.board),
				"Expected and actual hasWinner results are different")
		})
	}
}

func TestIsComplete(t *testing.T) {
	tests := []struct {
		description string
		board       Board
		expected    bool
	}{
		{
			description: "An empty board is not complete",
			board: Board{
				[3]int8{0, 0, 0},
				[3]int8{0, 0, 0},
				[3]int8{0, 0, 0},
			},
			expected: false,
		},
		{
			description: "A full board is complete",
			board: Board{
				[3]int8{X, O, X},
				[3]int8{O, X, O},
				[3]int8{X, O, X},
			},
			expected: true,
		},
	}
	for _, test := range tests {
		test := test // Capture range variable.
		t.Run(test.description, func(t *testing.T) {
			assert.Equal(t,
				test.expected,
				isComplete(test.board),
				"Expected and actual isComplete results are different")
		})
	}
}
