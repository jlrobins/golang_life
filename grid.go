package main

// go mod init life

import (
	"fmt"
	"math/rand"
	"time"
	"slices"
)

// Represent a Game of Life square/torus grid. Supports getting, setting cells by X,Y coordinates and determining neighbor count.
// Each cell value is either 0 or 1, int8. The neighbor count of a cell is returned as int8.
type Grid struct {
	// Overall size / cell count will be width * width
	width int

	// Linear array of the cells, 
	grid []int8
}

func NewGrid(width int) *Grid  {
	g := new(Grid)
	g.width = width
	g.grid = make([]int8, width*width)

	return g
}

// Am I equilavent to this other grid?
func (g *Grid) Equal(other *Grid) bool {
	return slices.Equal(g.grid, other.grid)
}

// Return 0 or 1 based on if cell at this position is active.
// Handles 'torus space' by [-1] is treated as end, and [size] is treated as [0]
func (g *Grid) get (x, y int) int8 {
	if x == -1 {
		x = g.width - 1
	} else if x == g.width {
		x = 0
	}

	if y == -1 {
		y = g.width - 1
	} else if y == g.width {
		y = 0
	}

	return g.grid[y * g.width + x]
	
}

// Set cell at [x][y] to be alive
func (g *Grid) set (x, y int) {
	g.grid[y * g.width + x] = 1
}

// Set all cells in this grid to dead
func (g *Grid) clear() {
	for i:=0; i<(g.width * g.width); i++ {
		g.grid[i] = 0
	}
}

// Randomize all cells in this grid
func (g *Grid) randomize() {
	rand.Seed(time.Now().UnixNano())

	for idx:= 0; idx < g.width * g.width; idx ++ {
		g.grid[idx] = int8(rand.Intn(2))
	}
}

func (g *Grid) Display() {
	for rowIdx:= 0; rowIdx < g.width; rowIdx++ {
		fmt.Println(g.grid[rowIdx * g.width:(rowIdx +1) * g.width])
	}
}

// Return count of neighbors of this cell which are alive
func (g *Grid) neighbors (x, y int) int8 {
	
	// whole three neighbors above, plus left and right neighbors on this row, then three neighbors below.
	return g.get(x-1, y-1) +
		g.get(x,   y-1) +
		g.get(x+1, y-1) +

		g.get(x-1, y) +
		g.get(x+1, y) +

		g.get(x-1, y+1) +
		g.get(x,   y+1) +
		g.get(x+1, y+1) 

}