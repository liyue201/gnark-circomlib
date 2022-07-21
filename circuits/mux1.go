package circuits

import (
	"github.com/consensys/gnark/frontend"
)

func MultiMux1(api frontend.API, c [][]frontend.Variable, sel frontend.Variable) []frontend.Variable {
	if len(c[0]) != 2 {
		panic("invalid param")
	}

	n := len(c)
	out := make([]frontend.Variable, n)
	for i := 0; i < n; i++ {
		out[i] = api.Add(api.Mul(api.Sub(c[i][1], c[i][0]), sel), c[i][0])
	}
	return out
}

func Mux1(api frontend.API, c []frontend.Variable, sel frontend.Variable) frontend.Variable {
	if len(c) != 2 {
		panic("invalid parm")
	}

	a := make([][]frontend.Variable, 1)
	a[0] = make([]frontend.Variable, 2)

	for i := 0; i < 2; i++ {
		a[0][i] = c[i]
	}
	out := MultiMux1(api, a, sel)
	return out[0]
}
