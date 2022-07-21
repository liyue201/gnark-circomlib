package circuits

import (
	"github.com/consensys/gnark/frontend"
	"math/big"
)

func BinSum(api frontend.API, in [][]frontend.Variable) []frontend.Variable {

	ops := len(in)
	n := len(in[0])

	nn := BigMul(BigSub(Lsh(1, uint(n)), big.NewInt(1)), big.NewInt(int64(ops)))
	nout := nn.BitLen()

	lin := frontend.Variable(0)
	for k := 0; k < nout; k++ {
		e2 := Lsh(1, uint(k))
		for j := 0; j < ops; j++ {
			lin = api.Add(lin, api.Mul(in[j][k], e2))
		}
	}
	out := api.ToBinary(lin)

	return out
}
