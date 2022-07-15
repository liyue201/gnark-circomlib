package circuits

import "math/big"

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

func BigMul(x *big.Int, y *big.Int) *big.Int {
	return big.NewInt(0).Mul(x, y)
}

func BigDiv(x *big.Int, y *big.Int) *big.Int {
	return big.NewInt(0).Div(x, y)
}