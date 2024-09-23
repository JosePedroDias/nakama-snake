package main

import (
	"fmt"
	"math/rand"
)

const (
	CellEmpty int = iota
	CellSnake
	CellFood
)

const INITIAL_SNAKE_SIZE = 2

type Point struct {
	X int `json:"x"`
	Y int `json:"y"`
}

type IntArr2D = [][]int
type SnakeBody = []Point
type Snake struct {
	Body      SnakeBody `json:"body"`
	Direction Point     `json:"direction"`
}

type SnakeGame struct {
	W        int      `json:"w"`
	H        int      `json:"h"`
	Board    IntArr2D `json:"board"`
	Snakes   []*Snake `json:"snakes"`
	SnakeIds []string `json:"snakeIds"`
	HasFood  bool     `json:"hasFood"`
}

var moves []Point

func newSnakeGame(W, H, numSnakes int) *SnakeGame {
	moves = []Point{{0, 1}, {1, 0}, {0, -1}, {-1, 0}} // right, down, left, up

	game := &SnakeGame{
		W:        W,
		H:        H,
		Board:    make([][]int, W),
		Snakes:   make([]*Snake, 0),
		SnakeIds: make([]string, 0),
	}

	for i := range game.Board {
		game.Board[i] = make([]int, H)
	}

	for snI := 0; snI < numSnakes; snI++ {
		game.addSnake()
	}

	game.placeFood()

	return game
}

func (g *SnakeGame) addSnake() {
	freeCells := g.getFreeCells()

	var pos0 Point
	var pos1 Point
	var cellI int

	for {
		cellI = rand.Intn(len(freeCells))
		pos0 = freeCells[cellI]
		pos1 = Point{pos0.X - 1, pos0.Y}

		if pos1.X < 0 {
			continue
		}

		value1 := g.Board[pos1.X][pos1.Y]
		if value1 != CellEmpty {
			continue
		}

		break
	}

	snake := &Snake{
		Body:      make(SnakeBody, 0),
		Direction: Point{0, 1}, // right
	}
	for snBodyI := 0; snBodyI < INITIAL_SNAKE_SIZE; snBodyI++ {
		pos := Point{pos0.X - snBodyI, pos0.Y}
		g.Board[pos.X][pos.Y] = CellSnake
		snake.Body = append(snake.Body, pos)
	}

	g.Snakes = append(g.Snakes, snake)
}

func (g *SnakeGame) removeSnake(snI int) {
	snake := g.Snakes[snI]

	for _, snPos := range snake.Body {
		g.Board[snPos.X][snPos.Y] = CellEmpty
	}

	g.Snakes = append(g.Snakes[:snI], g.Snakes[snI+1:]...)
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

func (g *SnakeGame) canChangeDirection(snake *Snake, dir Point) bool {
	head_ := snake.Body[0]
	head := Point{head_.X + dir.X, head_.Y + dir.Y}
	if head.X < 0 || head.Y < 0 || head.X >= g.W || head.Y >= g.H {
		return false
	}
	cell := g.Board[head.X][head.Y]
	if cell == CellSnake {
		return false
	}
	return (snake.Direction.X+dir.X != 0) || (snake.Direction.Y+dir.Y != 0)
}

func (g *SnakeGame) move(snake *Snake) {
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

func (g *SnakeGame) getWhichSnake(pos Point) int {
	for snI, snake := range g.Snakes {
		for _, snPos := range snake.Body {
			if snPos.X == pos.X && snPos.Y == pos.Y {
				return snI
			}
		}
	}
	panic("should not happen!")
}

func (g *SnakeGame) validateDirection(snake *Snake, dir Point) bool {
	if dir.X == 0 && (dir.Y == -1 || dir.Y == 1) {
		return g.canChangeDirection(snake, dir)
	}

	if dir.Y == 0 && (dir.X == -1 || dir.X == 1) {
		return g.canChangeDirection(snake, dir)
	}

	return false
}

//lint:ignore U1000 optional method
func (g *SnakeGame) getValidDirections(snake *Snake) []Point {
	potMoves := make([]Point, 0)
	for _, dir := range moves {
		if g.canChangeDirection(snake, dir) {
			potMoves = append(potMoves, dir)
		}
	}
	return potMoves
}

func (g *SnakeGame) getSnakeIndexFromId(id string) int {
	for snI := 0; snI < len(g.SnakeIds); snI++ {
		if g.SnakeIds[snI] == id {
			return snI
		}
	}
	panic("snake having this id was not found!")
}

//lint:ignore U1000 optional method
func (g *SnakeGame) displayBoard() {
	//jsonS, _ := json.Marshal(g)
	//fmt.Printf("%s\n", jsonS)

	for y := 0; y < g.H; y++ {
		for x := 0; x < g.W; x++ {
			s := ". "
			if g.Board[x][y] == CellSnake {
				//s = "# "
				s = fmt.Sprintf("%d ", g.getWhichSnake(Point{x, y})+1)
			} else if g.Board[x][y] == CellFood {
				s = "O "
			}
			fmt.Print(s)
		}
		fmt.Println()
	}
	fmt.Println()
}

/*
func main() {
	game := newSnakeGame(10, 10, 2)
	game.displayBoard()

outer:
	for {
		//for i := 0; i < 3; i++ {

		// if i == 0 {
		// 	game.addSnake()
		// } else if i == 2 {
		// 	game.removeSnake(0)
		// }

		for snI, snake := range game.Snakes {
			potMoves := game.getValidDirections(snake)
			if len(potMoves) == 0 {
				fmt.Printf("Game over: snake #%d got stuck!\n", snI+1)
				break outer
			}
			snake.Direction = potMoves[rand.Intn(len(potMoves))]
			game.move(snake)
		}

		game.displayBoard()
		time.Sleep(100 * time.Millisecond)
	}
}
*/
