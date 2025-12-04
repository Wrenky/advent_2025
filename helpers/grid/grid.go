package grid

import (
	"fmt"
	"math"
	"slices"

	"github.com/samber/lo"
)

type Grid[T comparable] [][]T

// Iterator, to get all points.
// `for point := range g.AllPointsIt` etc
func (g Grid[T]) AllPointsIt(yield func(Coord) bool) {
	for i, rows := range g {
		for j := range rows {
			if !yield(Coord{X: i, Y: j}) {
				return
			}
		}
	}
}

func AllPoints[T comparable](g Grid[T]) []Coord {
	return lo.Flatten(lo.Map(g, func(row []T, i int) []Coord {
		return lo.Map(row, func(c T, j int) Coord {
			return Coord{X: i, Y: j}
		})
	}))
}

// ---------------------------------------------------------------
// Grid/2d array helpers
// ---------------------------------------------------------------
func Transpose[T comparable](grid Grid[T]) Grid[T] {
	xl := len(grid[0])
	yl := len(grid)
	result := make([][]T, xl)
	for i := range result {
		result[i] = make([]T, yl)
	}
	for i := 0; i < xl; i++ {
		for j := 0; j < yl; j++ {
			result[i][j] = grid[j][i]
		}
	}
	return result
}

func Rotate90[T comparable](grid Grid[T]) Grid[T] {
	transposed := Transpose(grid)
	result := [][]T{}
	for _, v := range transposed {
		result = append(result, lo.Reverse(v))
	}
	return result
}

func RotateN90[T comparable](grid Grid[T]) Grid[T] {
	result := [][]T{}
	for _, v := range grid {
		result = append(result, lo.Reverse(v))
	}
	return Transpose(result)
}

// Helper for get operations to avoid
func Get[T comparable](grid Grid[T], a Coord) T {
	return grid[a.X][a.Y]
}

func Print[T comparable](g Grid[T]) {
	for _, v := range g {
		fmt.Printf("%v\n", v)
	}
}

func Copy[T comparable](g Grid[T]) Grid[T] {
	newG := make(Grid[T], len(g))
	for i := range g {
		newG[i] = make([]T, len(g[i]))
		copy(newG[i], g[i])
	}
	return newG
}

func InBounds[T comparable](g Grid[T], c Coord) bool {
	if c.X < 0 || c.Y < 0 {
		return false
	}
	if (c.X > len(g[0])-1) || (c.Y > len(g)-1) {
		return false
	}
	return true
}

func Clone[T comparable](g Grid[T]) Grid[T] {
	return lo.Map(g, func(row []T, _ int) []T {
		return slices.Clone(row)
	})
}

// These are annoying because in math its x,y, but in head is [col (y)][row (x)]
// y is up/down, x is left right
// They also tend to be counter intuitive to a graph- 0,0 is the top left corner!
type Coord struct {
	X, Y int
}

// Helpful grid coordinate variables for movement
var (
	UPLEFT        = Coord{X: -1, Y: -1}
	UP            = Coord{X: -1, Y: 0}
	UPRIGHT       = Coord{X: -1, Y: 1}
	LEFT          = Coord{X: 0, Y: -1}
	RIGHT         = Coord{X: 0, Y: 1}
	DOWNRIGHT     = Coord{X: 1, Y: -1}
	DOWN          = Coord{X: 1, Y: 0}
	DOWNLEFT      = Coord{X: 1, Y: 1}
	ALLDIRECTIONS = []Coord{
		UPLEFT, UP, UPRIGHT,
		LEFT, RIGHT,
		DOWNLEFT, DOWN, DOWNRIGHT,
	}
)

func (c Coord) String() string {
	return fmt.Sprintf("(%d, %d)", c.X, c.Y)
}
func ManhattanDist(a, b Coord) int {
	distance := math.Abs(float64(a.X-b.X)) + math.Abs(float64(a.Y-b.Y))
	return int(distance)
}

func Add(a, b Coord) Coord {
	return Coord{
		X: a.X + b.X,
		Y: a.Y + b.Y,
	}
}
