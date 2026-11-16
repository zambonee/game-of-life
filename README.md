# Conway's game of life

Implement [Conway's Game of Life](https://en.wikipedia.org/wiki/Conway%27s_Game_of_Life) in 64-bit signed integer space.

Imagine a 2D grid - each cell (coordinate) can be either "alive" or "dead". Every generation of the simulation, the system ticks forward. Each cell's value changes according to the following:

If an "alive" cell had less than 2 or more than 3 alive neighbors (in any of the 8 surrounding cells), it becomes dead.
If a "dead" cell had *exactly* 3 alive neighbors, it becomes alive.
Your input is a list of integer coordinates for live cells in the [Life 1.06 format](http://www.conwaylife.com/wiki/Life_1.06). They could be anywhere in the signed 64-bit range. This means the board could be very large!

Sample input:

#Life 1.06

0 1

1 2

2 0

2 1

2 2

-2000000000000 -2000000000000

-2000000000001 -2000000000001

-2000000000000 -2000000000001

Your program should read the state of the simulation from standard input, run 10 iterations of the Game of Life, and print the result to standard output in Life 1.06 format.

We're most interested in both the technical aspects of how you deal with very large integers and how you go about solving the problem. At the onsite, be prepared to discuss your solution, including the choices and tradeoffs you made. Though not required, you are welcome to bring a laptop with you to your interview to walk us through the code.

[epic conway's game of life](https://youtu.be/C2vgICfQawE)

## Usage

Run the simulation by piping a Life 1.06 file to the program:

```sh
go run ./cmd < input.life
```

Input file must begin with the `#Life 1.06` header. Each following nonblank line must contain an `x y` coordinate pair, with both values in the signed 64-bit integer range. The program advances the board by 10 generations and writes the live-cell coordinates to standard output in Life 1.06 format, sorted by coordinate. Input errors are reported to standard error and cause a nonzero exit status.

## Tests

Run the test suite with:

```sh
go test ./...
```