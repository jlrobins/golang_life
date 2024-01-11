# golang_life
Basic demonstration golang implementing [Conway's Game of Life](https://en.wikipedia.org/wiki/Conway%27s_Game_of_Life)


## Architecture
* Main in `life.go` constructing a 10x10 life `Board` which will get evaluated in parallel with 5 goros.
* A Board struct (`board.go`) contains two grids, one for the current board state, and one for then 'next' state.
  - Each evaluation step swaps which is the current vs next board
  - Cell evaluation is done by goroutines which loop forever draining row index numbers whose row needs evaluating (via an integer pushed onto a channel). When done evaluating a row, will then message that the work unit is complete (via a `sync.WaitGroup`).
* A Grid (`grid.go`) manages mapping row/col coordinates to indexes within a single slice of int8s holding the life cell grid.
  - int8 value 1 represents 'live cell'. Int8s are used to simplify .neighbors() being able to add up the neighbor count.

* Current configuration uses 5 long-lived goros to evaluate the 10 rows within the board's grids. So, each clock step, each
goro will consume on average 2 of the row index messages pushed onto the work queue channel. A more interesting and viable use of this parallelism would be to have a _huge_ Life board, say, 10k x 10k cells, then rearrange slightly to have the goros get messaged with a low / high cell index value pair (or perhaps just a pair of row indexes), where then completing a work unit will take some less stupidly trivial computation time and be worth the inter-goro communication overhead. But larger Life grids are harder to see and interpret when dumped to stdout.

## To build
* `go build`
* `./life`
  -- if notices that the board is in a steady state, will terminate. If board never reaches a simple steady state, will loop forever, needing control-c to terminate.




