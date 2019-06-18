package main

// Copied from https://github.com/tinygo-org/tinygo/blob/master/testdata/gc.go then modified

var xorshift32State uint32 = 1

func xorshift32(x uint32) uint32 {
	// Algorithm "xor" from p. 4 of Marsaglia, "Xorshift RNGs"
	x ^= x << 13
	x ^= x >> 17
	x ^= x << 5
	return x
}

func randuint32() uint32 {
	xorshift32State = xorshift32(xorshift32State)
	return xorshift32State
}

func main() {
	testNonPointerHeap()
}

func testNonPointerHeap() {
	for i := 0; i < 5; i++ {
		// Pick a random index that the optimizer can't predict.
		index := randuint32() % 4
		println(index)
	}
}