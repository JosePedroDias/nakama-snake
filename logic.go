package main

import (
	"fmt"
	"math/rand"
	"time"
)

const (
	CellEmpty int = iota
	CellSnake
	CellFood
)

type Point struct {
	X, Y int
}

type SnakeGame struct {
	Board     [][]int // 0 for empty, 1 for snake, 2 for food
	Snake     []Point // The body of the snake, where the first element is the head
	Direction Point   // Movement direction (e.g., {-1, 0} for left)
	HasFood   bool
	W, H      int
}

// Initialize the game with an MxN board, a snake of size 2, and one food cell
func NewSnakeGame(W, H int) *SnakeGame {
	game := &SnakeGame{
		Board:     make([][]int, W),
		Snake:     []Point{{W / 2, H / 2}, {W / 2, H/2 - 1}}, // Initial snake with size 2
		Direction: Point{0, 1},                               // Moving right initially
		W:         W,
		H:         H,
	}

	for i := range game.Board {
		game.Board[i] = make([]int, H)
	}

	// Place the initial snake on the board
	game.Board[game.Snake[0].X][game.Snake[0].Y] = CellSnake
	game.Board[game.Snake[1].X][game.Snake[1].Y] = CellSnake

	// Place the first food
	game.placeFood()

	return game
}

// Places food at a random free location
func (g *SnakeGame) placeFood() {
	freeCells := []Point{}

	// Collect all free cells
	for x := 0; x < g.W; x++ {
		for y := 0; y < g.H; y++ {
			if g.Board[x][y] == 0 { // Empty cell
				freeCells = append(freeCells, Point{x, y})
			}
		}
	}

	if len(freeCells) == 0 {
		return
	}

	// Choose a random free cell for food
	foodPosition := freeCells[rand.Intn(len(freeCells))]
	g.Board[foodPosition.X][foodPosition.Y] = CellFood
	g.HasFood = true
}

// Move the snake in the current direction
func (g *SnakeGame) Move() bool {
	head := g.Snake[0]
	newHead := Point{head.X + g.Direction.X, head.Y + g.Direction.Y}

	// Check for boundary collisions
	if newHead.X < 0 || newHead.X >= g.W || newHead.Y < 0 || newHead.Y >= g.H {
		return false // Game over
	}

	// Check for self-collision
	if g.Board[newHead.X][newHead.Y] == CellSnake {
		return false // Game over
	}

	// Check if food is found
	if g.Board[newHead.X][newHead.Y] == CellFood {
		// Eat the food and grow the snake
		g.Snake = append([]Point{newHead}, g.Snake...)
		g.Board[newHead.X][newHead.Y] = CellSnake
		g.HasFood = false
		// Place new food
		g.placeFood()
	} else {
		// Remove tail from the board
		tail := g.Snake[len(g.Snake)-1]
		g.Board[tail.X][tail.Y] = CellEmpty

		// Move the snake normally (no food)
		g.Snake = append([]Point{newHead}, g.Snake[:len(g.Snake)-1]...)
		g.Board[newHead.X][newHead.Y] = CellSnake
	}

	return true
}

func (g *SnakeGame) CanChangeDirection(newDir Point) bool {
	head_ := g.Snake[0]
	head := Point{head_.X + newDir.X, head_.Y + newDir.Y}
	if head.X < 0 || head.Y < 0 || head.X >= g.W || head.Y >= g.H {
		return false
	}
	return (g.Direction.X+newDir.X != 0) || (g.Direction.Y+newDir.Y != 0)
}

func (g *SnakeGame) DisplayBoard() {
	fmt.Printf("%#v\n", g.Direction)
	for y := 0; y < g.H; y++ {
		for x := 0; x < g.W; x++ {
			if g.Board[x][y] == CellSnake {
				fmt.Print("# ")
			} else if g.Board[x][y] == CellFood {
				fmt.Print("O ")
			} else {
				fmt.Print(". ")
			}
		}
		fmt.Println()
	}
	fmt.Println()
}

func main() {
	game := NewSnakeGame(10, 10)

	game.DisplayBoard()

	moves := []Point{{0, 1}, {1, 0}, {0, -1}, {-1, 0}} // Right, down, left, up

	for i := 0; i < 500; i++ {
		potMoves := make([]Point, 0)
		for _, dir := range moves {
			if game.CanChangeDirection(dir) {
				potMoves = append(potMoves, dir)
			}
		}
		game.Direction = potMoves[rand.Intn(len(potMoves))]
		if !game.Move() {
			fmt.Println("Game Over!")
			break
		}
		game.DisplayBoard()
		time.Sleep(100 * time.Millisecond) // Delay for visualization (optional)
	}
}
