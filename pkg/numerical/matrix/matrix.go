package matrix

import "math/big"

type Matrix interface {
	// Dims returns the dimensions of a Matrix.  @TODO - use big ints?
	Dims() (r, c int)

	// At returns the value of a matrix element at row i, column j
	// panics if i or j are out of bounds for the matrix
	At(i, j int) *big.Int

	// T returns the transpose of the Matrix.
	T() Matrix
}
