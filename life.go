package main

import (
	"fmt"
)

func main() {

	// Construct a 10x10 board whose cells will be evaluated by 5 goros in parallel.
	b := NewBoard(10, 5)

	// Start out with a random cell configuration.
	b.Randomize()

	// Keep advancing the board until we're in a trivial to detect cycle. May be infinite loop based
	// on how we randomized.
	for step:= 0; true; step++ {
		b.Display()

		fmt.Println()

		// Advance one Life tick.
		had_changes := b.Step()

		if ! had_changes {
			fmt.Printf("No more changes, breaking on step %d\n", step)
			break
		}
	}
}