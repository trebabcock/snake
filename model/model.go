package model

import (
	"image/color"
	"math/rand"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
)

// Game holds game data
type Game struct {
	Width  float64
	Height float64
	Speed  float64
	Over   bool
	Score  int
	Player *Snake
	Foods  *Food
}

// Start starts a new Game
func Start() *Game {
	game := &Game{
		Width:  925,
		Height: 525,
		Speed:  0.25,
		Over:   false,
		Player: &Snake{},
		Foods:  &Food{},
	}

	game.Player = NewSnake(game)
	game.Foods = NewFood(game)

	return game
}

// Vector2 is a 2D vector
type Vector2 struct {
	X float64
	Y float64
}

// Add adds two vectors
func (v *Vector2) Add(v1 Vector2) {
	v.X = v1.X
	v.Y = v1.Y
}

// VectorSum returns the sum of two vectors
func VectorSum(v1, v2 Vector2) Vector2 {
	sum := Vector2{
		X: v1.X + v2.X,
		Y: v1.Y + v2.Y,
	}
	return sum
}

// VectorSumFloat returns the sum of two vectors
func VectorSumFloat(v Vector2, f float64) Vector2 {
	sum := Vector2{
		X: v.X + f,
		Y: v.Y + f,
	}
	return sum
}

// VectorProductFloat returns the product of a vector and a float
func VectorProductFloat(v Vector2, f float64) Vector2 {
	sum := Vector2{
		X: v.X * f,
		Y: v.Y * f,
	}
	return sum
}

// VectorDot returns the dot product of two vectors
func VectorDot(v1, v2 Vector2) float64 {
	x := v1.X * v2.X
	y := v1.Y * v2.Y
	return x + y
}

// Equals returns true if two vectors are equal
func (v *Vector2) Equals(v1 Vector2) bool {
	return (v.X == v1.X && v.Y == v1.Y)
}

// VectorOutOfBounds calculates if the vector is out of bounds
func VectorOutOfBounds(v Vector2, game *Game) bool {
	if v.X < 0 || v.X > game.Width-25 {
		return true
	} else if v.Y < 0 || v.Y > game.Height-25 {
		return true
	} else {
		return false
	}
}

// VectorInSnake calculates if the vector is in the snake
func VectorInSnake(v Vector2, game *Game) bool {
	vectors := []Vector2{game.Player.Head.Position, game.Player.Body[len(game.Player.Body)-1].Previous}
	for _, n := range game.Player.Body {
		vectors = append(vectors, n.Position)
	}

	for _, vec := range vectors {
		if vec.Equals(v) {
			return true
		}
	}

	return false
}

// Node is a part of the snake
type Node struct {
	Position Vector2
	Previous Vector2
	Shape    *imdraw.IMDraw
}

// NewNode creates a new node
func NewNode(position Vector2) Node {
	node := Node{
		Position: position,
		Shape:    MakeSquare(position, color.RGBA{101, 123, 131, 0}, 21),
	}
	return node
}

// Food is snake food, yum
type Food struct {
	Position Vector2
	Previous Vector2
	Shape    *imdraw.IMDraw
}

// NewFood creates a new food object
func NewFood(game *Game) *Food {
	randColor := rand.Intn(3) + 1
	col := color.RGBA{}
	switch randColor {
	case 1:
		col = color.RGBA{203, 75, 22, 0}
		break
	case 2:
		col = color.RGBA{211, 54, 130, 0}
		break
	case 3:
		col = color.RGBA{42, 161, 152, 0}
	}

	randX := rand.Intn(36) + 1
	randY := rand.Intn(20) + 1

	for (VectorInSnake(Vector2{float64(randX), float64(randY)}, game)) {
		randX = rand.Intn(36) + 1
		randY = rand.Intn(20) + 1
	}

	position := Vector2{
		X: float64(randX * 25),
		Y: float64(randY * 25),
	}

	food := Food{
		Position: position,
		Shape:    MakeSquare(position, col, 17),
	}
	return &food
}

// GetConsumed is called when food gets eaten
func (f *Food) GetConsumed(game *Game) {
	f = NewFood(game)
}

// Snake is the whole snake
type Snake struct {
	Head      Node
	Body      []Node
	Direction Vector2
}

// AddNode adds a node to the snake
func (s *Snake) AddNode() {
	node := NewNode(s.Body[len(s.Body)-1].Previous)
	s.Body = append(s.Body, node)
}

// NewSnake creates a new snake
func NewSnake(game *Game) *Snake {
	head := NewNode(Vector2{X: float64((36 / 2) * 25), Y: float64((20 / 2) * 25)})
	node1 := NewNode(Vector2{X: float64((36 / 2) * 25), Y: float64((20/2)*25 - 25)})
	node2 := NewNode(Vector2{X: float64((36 / 2) * 25), Y: float64((20/2)*25 - 50)})
	body := []Node{node1, node2}

	snake := Snake{
		Head:      head,
		Body:      body,
		Direction: Vector2{X: 0, Y: 1},
	}

	return &snake
}

// Move moves the snake
func (s *Snake) Move(game *Game) {
	for _, n := range s.Body {
		if VectorSum(s.Head.Position, VectorProductFloat(s.Direction, 25)) == game.Foods.Position {
			game.Foods = NewFood(game)
			s.AddNode()
			game.Speed = game.Speed * 0.95
			game.Score++
		}
		if VectorSum(s.Head.Position, VectorProductFloat(s.Direction, 25)) == n.Position {
			game.Over = true
		}
		if VectorOutOfBounds(VectorSum(s.Head.Position, VectorProductFloat(s.Direction, 25)), game) {
			game.Over = true
		}
	}
	s.Head.Previous = s.Head.Position
	newPosition := VectorProductFloat(s.Direction, 25)
	s.Head.Position = VectorSum(s.Head.Position, newPosition)
	s.Head.Shape = MakeSquare(s.Head.Position, color.RGBA{101, 123, 131, 0}, 21)
	for i := 0; i < len(s.Body); i++ {
		s.Body[i].Previous = s.Body[i].Position
		if i == 0 {
			s.Body[i].Position = s.Head.Previous
		} else {
			s.Body[i].Position = s.Body[i-1].Previous
		}
		s.Body[i].Shape = MakeSquare(s.Body[i].Position, color.RGBA{101, 123, 131, 0}, 21)
	}
}

// Turn turns the snake
func (s *Snake) Turn(direction Vector2, game *Game) {
	if s.Direction == direction {
		return
	} else if VectorSum(s.Head.Position, VectorProductFloat(direction, 25)) == s.Head.Previous {
		return
	} else if VectorDot(s.Direction, direction) == -1 {
		return
	} else {
		s.Direction = direction
	}
	s.Move(game)
}

// MakeSquare makes a new square
func MakeSquare(position Vector2, col color.Color, width float64) *imdraw.IMDraw {
	imd := imdraw.New(nil)
	imd.Color = col
	imd.Push(pixel.V(position.X+(25-width)/2, position.Y+(25-width)/2))
	imd.Push(pixel.V(position.X+(25-width)/2, position.Y+width+2))
	imd.Push(pixel.V(position.X+width+2, position.Y+width+2))
	imd.Push(pixel.V(position.X+width+2, position.Y+(25-width)/2))
	imd.Polygon(0)

	return imd
}
