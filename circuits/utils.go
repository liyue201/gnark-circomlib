package circuits

import (
	"github.com/consensys/gnark/frontend"
	"math/big"
)

func Lsh(k int64, n uint) *big.Int {
	z := big.NewInt(k)
	return z.Lsh(z, n)
}

func BigLsh(k *big.Int, n uint) *big.Int {
	return big.NewInt(0).Lsh(k, n)
}

func BigRsh(k *big.Int, n uint) *big.Int {
	return big.NewInt(0).Rsh(k, n)
}

func BigAnd(x *big.Int, y *big.Int) *big.Int {
	return big.NewInt(0).And(x, y)
}

func BigAdd(x *big.Int, y *big.Int) *big.Int {
	return big.NewInt(0).Add(x, y)
}

func BigSub(x *big.Int, y *big.Int) *big.Int {
	z := big.NewInt(0)
	return z.Sub(x, y)
}

func BigMul(x *big.Int, y *big.Int, c ...*big.Int) *big.Int {
	z := big.NewInt(0).Mul(x, y)
	for _, v := range c {
		z = z.Mul(z, v)
	}
	return z
}

func BigDiv(x *big.Int, y *big.Int) *big.Int {
	return big.NewInt(0).Div(x, y)
}

func MakeVariableArray(n int) []frontend.Variable {
	a := make([]frontend.Variable, n)
	for i := 0; i < n; i++ {
		a[i] = 0
	}
	return a
}

func Make2DVariableArray(n, m int) [][]frontend.Variable {
	a := make([][]frontend.Variable, n)
	for i := 0; i < n; i++ {
		a[i] = MakeVariableArray(m)
	}
	return a
}

func Make3DVariableArray(n, m, o int) [][][]frontend.Variable {
	a := make([][][]frontend.Variable, n)
	for i := 0; i < n; i++ {
		a[i] = Make2DVariableArray(m, o)
	}
	return a
}
