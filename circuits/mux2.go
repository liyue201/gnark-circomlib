package circuits

import "github.com/consensys/gnark/frontend"

func MultiMux2(api frontend.API, c [][]frontend.Variable, sel []frontend.Variable) []frontend.Variable {
	if len(c[0]) != 4 {
		panic("invalid param")
	}
	n := len(c)
	out := make([]frontend.Variable, n)

	a10 := make([]frontend.Variable, n)
	a0 := make([]frontend.Variable, n)
	a1 := make([]frontend.Variable, n)
	a := make([]frontend.Variable, n)
	s10 := api.Mul(sel[1], sel[0])

	for i := 0; i < n; i++ {
		a10[i] = api.Mul(api.Add(api.Sub(c[i][3], c[i][2]), c[i][0]), s10)
		a1[i] = api.Mul(api.Sub(c[i][2], c[i][0]), sel[1])
		a0[i] = api.Mul(api.Sub(c[i][1], c[i][0]), sel[0])
		a[i] = c[i][0]
		out[i] = api.Add(a10[i], a1[i], a0[i], a[i])
	}
	return out
}

func Mux2(api frontend.API, c []frontend.Variable, sel []frontend.Variable) frontend.Variable {
	if len(c) != 4 || len(sel) != 2 {
		panic("invalid parm")
	}
	a := make([][]frontend.Variable, 1)
	a[0] = make([]frontend.Variable, 4)

	for i := 0; i < 4; i++ {
		a[0][i] = c[i]
	}

	out := MultiMux2(api, a, sel)
	return out[0]
}
