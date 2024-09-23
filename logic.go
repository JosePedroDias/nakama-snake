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

const INITIAL_SNAKE_SIZE = 2

type Point struct {
	X, Y int
}

type IntArr2D = [][]int
type SnakeBody = []Point
type Snake struct {
	Body      SnakeBody
	Direction Point
}

type SnakeGame struct {
	Board   IntArr2D
	Snakes  []*Snake
	HasFood bool
	W, H    int
}

func NewSnakeGame(W, H, numSnakes int) *SnakeGame {
	game := &SnakeGame{
		W:      W,
		H:      H,
		Board:  make([][]int, W),
		Snakes: make([]*Snake, 0),
	}

	for i := range game.Board {
		game.Board[i] = make([]int, H)
	}

	freeCells := game.getFreeCells()

	for snI := 0; snI < numSnakes; snI++ {
		// use one freeCell, get its value and remove it from the slice
		cellI := rand.Intn(len(freeCells))
		pos0 := freeCells[cellI]
		freeCells = append(freeCells[:cellI], freeCells[cellI+1:]...)

		snake := &Snake{
			Body:      make(SnakeBody, 0),
			Direction: Point{0, 1}, // right
		}
		for snBodyI := 0; snBodyI < INITIAL_SNAKE_SIZE; snBodyI++ {
			pos := Point{pos0.X - snBodyI, pos0.Y}
			// avoid placing body out of bounds (hacky as it can be non-empty)
			if pos0.X == 0 {
				pos.X = pos0.X + snBodyI
			}
			game.Board[pos.X][pos.Y] = CellSnake
			snake.Body = append(snake.Body, pos)
		}

		game.Snakes = append(game.Snakes, snake)
	}

	game.placeFood()

	return game
}

func (g *SnakeGame) getFreeCells() []Point {
	freeCells := []Point{}

	for x := 0; x < g.W; x++ {
		for y := 0; y < g.H; y++ {
			if g.Board[x][y] == CellEmpty {
				freeCells = append(freeCells, Point{x, y})
			}
		}
	}

	return freeCells
}

func (g *SnakeGame) getFreeCell() Point {
	freeCells := g.getFreeCells()
	return freeCells[rand.Intn(len(freeCells))]
}

func (g *SnakeGame) placeFood() {
	foodPosition := g.getFreeCell()
	g.Board[foodPosition.X][foodPosition.Y] = CellFood
	g.HasFood = true
}

func (g *SnakeGame) CanChangeDirection(snake *Snake, newDir Point) bool {
	head_ := snake.Body[0]
	head := Point{head_.X + newDir.X, head_.Y + newDir.Y}
	if head.X < 0 || head.Y < 0 || head.X >= g.W || head.Y >= g.H {
		return false
	}
	cell := g.Board[head.X][head.Y]
	if cell == CellSnake {
		return false
	}
	return (snake.Direction.X+newDir.X != 0) || (snake.Direction.Y+newDir.Y != 0)
}

func (g *SnakeGame) Move(snake *Snake) {
	head := snake.Body[0]
	newHead := Point{head.X + snake.Direction.X, head.Y + snake.Direction.Y}

	// Check if food is found
	if g.Board[newHead.X][newHead.Y] == CellFood {
		// Eat the food and grow the snake
		snake.Body = append([]Point{newHead}, snake.Body...)
		g.Board[newHead.X][newHead.Y] = CellSnake
		g.HasFood = false
		// Place new food
		g.placeFood()
	} else {
		// Remove tail from the board
		tail := snake.Body[len(snake.Body)-1]
		g.Board[tail.X][tail.Y] = CellEmpty

		// Move the snake normally (no food)
		snake.Body = append([]Point{newHead}, snake.Body[:len(snake.Body)-1]...)
		g.Board[newHead.X][newHead.Y] = CellSnake
	}
}

func (g *SnakeGame) GetWhichSnake(pos Point) int {
	for snI, snake := range g.Snakes {
		for _, snPos := range snake.Body {
			if snPos.X == pos.X && snPos.Y == pos.Y {
				return snI
			}
		}
	}
	panic("should not happen!")
}

func (g *SnakeGame) DisplayBoard() {
	for y := 0; y < g.H; y++ {
		for x := 0; x < g.W; x++ {
			s := ". "
			if g.Board[x][y] == CellSnake {
				//s = "# "
				s = fmt.Sprintf("%d ", g.GetWhichSnake(Point{x, y})+1)
			} else if g.Board[x][y] == CellFood {
				s = "O "
			}
			fmt.Print(s)
		}
		fmt.Println()
	}
	fmt.Println()
}

func main() {
	game := NewSnakeGame(10, 10, 2)

	game.DisplayBoard()

	moves := []Point{{0, 1}, {1, 0}, {0, -1}, {-1, 0}} // right, down, left, up

outer:
	//for i := 0; i < 3; i++ {
	for {
		for snI, snake := range game.Snakes {
			potMoves := make([]Point, 0)
			for _, dir := range moves {
				if game.CanChangeDirection(snake, dir) {
					potMoves = append(potMoves, dir)
				}
			}

			if len(potMoves) == 0 {
				fmt.Printf("Game over: snake #%d got stuck!\n", snI+1)
				break outer
			}

			snake.Direction = potMoves[rand.Intn(len(potMoves))]
			game.Move(snake)
		}

		game.DisplayBoard()
		time.Sleep(100 * time.Millisecond)
	}
}
