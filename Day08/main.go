package main

import (
	"advent/helpers"
	"container/heap"
	log "log/slog"
	"slices"
	"sort"
	"strings"
	"time"

	"github.com/samber/lo"
)

type Point struct {
	X int
	Y int
	Z int
}

type Path struct {
	A, B Point
	Dist int
}

// This needs to change to match the input
func parseInput(input string) []Point {
	return lo.Map(strings.Split(input, "\n"), func(line string, _ int) Point {
		v := lo.Map(strings.Split(line, ","), func(s string, _ int) int { return helpers.Atoi(s) })
		return Point{
			X: v[0],
			Y: v[1],
			Z: v[2],
		}
	})
}

// ----------------------------------------------------
// Min Heap and heap interface definition
// ----------------------------------------------------
type PathMinHeap []Path

func (h PathMinHeap) Len() int            { return len(h) }
func (h PathMinHeap) Less(i, j int) bool  { return h[i].Dist < h[j].Dist }
func (h PathMinHeap) Swap(i, j int)       { h[i], h[j] = h[j], h[i] }
func (h *PathMinHeap) Push(x interface{}) { *h = append(*h, x.(Path)) }
func (h *PathMinHeap) Pop() interface{} {
	old := *h
	n := len(old)
	item := old[n-1]
	*h = old[0 : n-1]
	return item
}

// ----------------------------------------------------

// Creates the min-heap from points
func ComputeDistances(points []Point) PathMinHeap {
	mph := PathMinHeap([]Path{})
	heap.Init(&mph)
	type Key struct {
		A, B Point
	}
	created := make(map[Key]int)
	for i, a := range points {
		for j, b := range points {
			if i == j {
				continue
			}
			k := Key{A: a, B: b}
			if _, ok := created[k]; ok {
				continue
			}
			if _, ok := created[Key{A: b, B: a}]; ok {
				// Duplicate!
				continue
			}
			created[k] = 1
			np := Path{
				A:    a,
				B:    b,
				Dist: Euclidean(a, b),
			}
			heap.Push(&mph, np)
		}
	}
	return mph
}

// Part 1 code
func CircuitsXConnected(mph PathMinHeap, circuits [][]Point, target int) [][]Point {
	if target == 0 {
		return circuits
	}
	shortest := heap.Pop(&mph).(Path)
	newCircuits := Connect(shortest.A, shortest.B, circuits)
	return CircuitsXConnected(mph, newCircuits, target-1)
}

// Part2 code
func LastPathConnected(mph PathMinHeap, circuits [][]Point) Path {
	// pop shortest, connect it , regen paths, decrement
	shortest := heap.Pop(&mph).(Path)

	newCircuits := Connect(shortest.A, shortest.B, circuits)
	if len(newCircuits) == 1 {
		return shortest
	}
	return LastPathConnected(mph, newCircuits)
}

// These handle the connecting and merging of circuits
func Connected(A, B Point, circuits [][]Point) bool {
	return lo.ContainsBy(circuits, func(circ []Point) bool {
		return lo.Every(circ, []Point{A, B})
	})
}
func Connect(A, B Point, circuits [][]Point) [][]Point {
	if Connected(A, B, circuits) {
		return circuits
	}
	CircA, indexA, _ := lo.FindIndexOf(circuits, func(circ []Point) bool {
		return lo.Contains(circ, A)
	})
	circuits = lo.DropByIndex(circuits, indexA)

	CircB, indexB, _ := lo.FindIndexOf(circuits, func(circ []Point) bool {
		return lo.Contains(circ, B)
	})
	circuits = lo.DropByIndex(circuits, indexB)

	newCircuit := append(CircA, CircB...)
	return append(circuits, newCircuit)
}

// you dont actually need to do the square root part!
func Euclidean(a, b Point) int {
	x2 := b.X - a.X
	y2 := b.Y - a.Y
	z2 := b.Z - a.Z
	return x2*x2 + y2*y2 + z2*z2
}

func main() {
	// Handle command line
	args := helpers.HandleCommandLine()
	data := helpers.ReadFile(args.InputFile)
	points := parseInput(data)
	for _, v := range points {
		log.Debug("", "point", v)
	}
	circuits := lo.Map(points, func(p Point, _ int) []Point {
		return []Point{p}
	})

	pre0 := time.Now()
	mph := ComputeDistances(points)
	mph2 := slices.Clone(mph)
	post0 := time.Now()
	log.Info("Setup/Parsing", "time", post0.Sub(pre0))

	pre1 := time.Now()
	res := CircuitsXConnected(mph, circuits, 1000)
	lengths := lo.Map(res, func(circ []Point, _ int) int {
		return len(circ)
	})
	sort.Sort(sort.Reverse(sort.IntSlice(lengths)))
	p1 := lo.Product(lengths[0:3])
	post1 := time.Now()
	log.Info("Part1", "answer", p1, "time", post1.Sub(pre1))

	pre2 := time.Now()
	res2 := LastPathConnected(mph2, circuits)
	p2 := res2.A.X * res2.B.X
	post2 := time.Now()
	log.Info("Part2", "answer", p2, "time", post2.Sub(pre2))
}
