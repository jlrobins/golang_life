package main

import "sync"

// A Life board. Sits atop two grids, one for this generation and one for the next.
type Board struct {
	this_grid *Grid
	next_grid *Grid

	index_channel chan int
	wg *sync.WaitGroup
}

// Construct a new board with width*width cells.
func NewBoard(width int, concurrency int) *Board {
	b := new(Board)
	b.this_grid = NewGrid(width)
	b.next_grid = NewGrid(width)

	b.index_channel = make(chan int, width)
	b.wg = new(sync.WaitGroup)

	// Go ahead and spawn `concurrency` number of goroutines to do subunits of grid evaluation.
	for i:=0; i < concurrency; i++ {
		go b.evaluation_worker() 
	}

	return b
}

// Initialize board to random values.
func (b *Board) Randomize() {
	b.this_grid.randomize()
}

// Display this board on stdout.
func (b *Board) Display() {
	b.this_grid.Display()
}

// Advance Game of Life rules by one step.
func (b *Board) Step() bool {
	// Will evaluate all cells in this_grid assigning into next_grid. Then will swap which grid
	// is considered this and next grid, moving one step forward.


	// Indicate have 'width' rows will need to wait to get evaluated by goros.
	b.wg.Add(b.this_grid.width)

	// Queue up to get each row in the board evaluated by the already running goros.
	// todo: assign index ranges, not per row.
	for rowIdx := 0; rowIdx < b.this_grid.width; rowIdx++ {
		// Message to any available evaluation_worker() goro that is now time to compute this row.
		b.index_channel <- rowIdx
	}

	// Block until each row has been completed (the count within b.wg to hit 0).
	b.wg.Wait()

	// Trivial steady state detection. If the new grid == the old grid, then we're in a steady state
	// and will remain there forevermore.
	had_changes := ! b.this_grid.Equal(b.next_grid)

	// next_grid is now correct. Swap buffers.
	b.swap()

	return had_changes

}

// Run as a goroutine. Is told what row to evaluate through the index channel. Evaluates GoL rules for the cells
// in that row, then messages the wait group that this work unit has completed.
func (b *Board) evaluation_worker() {
	for {
		// Blocks until can learn what row to compute.
		rowIdx :=  <-b.index_channel

		for colIdx := 0; colIdx < b.this_grid.width; colIdx++ {
			// default assume dead.
			neighbors := b.this_grid.neighbors(colIdx, rowIdx)
			is_alive := (b.this_grid.get(colIdx, rowIdx) == 1)

			// Game of Life: live cell lives on if has two or three neighbors. Any more is death by overcrowding, any fewer is lonliness
			// Dead cell with exactly three neighbors gets born into.
			if (is_alive && (neighbors == 2 || neighbors == 3)) ||
					(!is_alive && neighbors == 3) {
				b.next_grid.set(colIdx, rowIdx)
			}
		}

		// Signal that a row is now complete. Decrements the internal count.
		b.wg.Done()
	}
}

// Swap next_grid and this_grid. Clear the new next_grid.
func (b *Board) swap() {
	tmp := b.this_grid
	b.this_grid = b.next_grid
	b.next_grid = tmp
	tmp.clear()
}
