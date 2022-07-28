package circuits

import (
	"github.com/consensys/gnark/frontend"
	"math/big"
)

func pointAdd(x1, y1, x2, y2 *big.Int) (*big.Int, *big.Int) {
	a := big.NewInt(168700)
	d := big.NewInt(168696)

	x := BigDiv(BigAdd(BigMul(x1, y2), BigMul(y1, x2)), BigAdd(big.NewInt(1), BigMul(d, x1, x2, y1, y2)))
	y := BigDiv(BigSub(BigMul(y1, y2), BigMul(a, x1, x2)), BigSub(big.NewInt(1), BigMul(d, x1, x2, y1, y2)))
	return x, y
}

func EscalarMulW4Table(api frontend.API, base [2]*big.Int, k int) [][]*big.Int {
	out := make([][]*big.Int, 16)
	for i := 0; i < 16; i++ {
		out[i] = make([]*big.Int, 2)
	}
	var dbl [2]*big.Int
	dbl[0] = big.NewInt(0).Set(base[0])
	dbl[1] = big.NewInt(0).Set(base[1])

	for i := 0; i < 4*k; i++ {
		dbl[0], dbl[1] = pointAdd(dbl[0], dbl[1], dbl[0], dbl[1])
	}

	out[0][0] = big.NewInt(0)
	out[0][1] = big.NewInt(1)
	for i := 1; i < 16; i++ {
		out[i][0], out[i][1] = pointAdd(out[i-1][0], out[i-1][1], dbl[0], dbl[1])
	}
	return out
}
