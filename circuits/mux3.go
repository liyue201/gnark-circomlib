package circuits

import "github.com/consensys/gnark/frontend"

func MultiMux3(api frontend.API, c [][]frontend.Variable, sel []frontend.Variable) []frontend.Variable {
	if len(c[0]) != 8 || len(sel) != 3 {
		panic("invalid param")
	}
	n := len(c)
	out := make([]frontend.Variable, n)

	a210 := make([]frontend.Variable, n)
	a21 := make([]frontend.Variable, n)
	a20 := make([]frontend.Variable, n)
	a2 := make([]frontend.Variable, n)

	a10 := make([]frontend.Variable, n)
	a1 := make([]frontend.Variable, n)
	a0 := make([]frontend.Variable, n)
	a := make([]frontend.Variable, n)

	s10 := api.Mul(sel[1], sel[0])

	for i := 0; i < n; i++ {
		a210[i] = api.Mul(s10, api.Add(c[i][7], api.Neg(c[i][6]), api.Neg(c[i][5]), c[i][4], api.Neg(c[i][3]), c[i][2], c[i][1], api.Neg(c[i][0])))
		a21[i] = api.Mul(sel[1], api.Add(c[i][6], api.Neg(c[i][4]), api.Neg(c[i][2]), c[i][0]))
		a20[i] = api.Mul(sel[0], api.Add(c[i][5], api.Neg(c[i][4]), api.Neg(c[i][1]), c[i][0]))
		a2[i] = api.Sub(c[i][4], c[i][0])

		a10[i] = api.Mul(s10, api.Add(c[i][3], api.Neg(c[i][2]), api.Neg(c[i][1]), c[i][0]))
		a1[i] = api.Mul(sel[1], api.Sub(c[i][2], c[i][0]))
		a0[i] = api.Mul(sel[0], api.Sub(c[i][1], c[i][0]))
		a[i] = c[i][0]
		out[i] = api.Add(api.Mul(api.Add(a210[i], a21[i], a20[i], a2[i]), sel[2]), a10[i], a1[i], a0[i], a[i])
	}
	return out
}

func Mux3(api frontend.API, c []frontend.Variable, sel []frontend.Variable) frontend.Variable {
	if len(c) != 8 || len(sel) != 3 {
		panic("invalid param")
	}

	a := make([][]frontend.Variable, 1)
	a[0] = make([]frontend.Variable, 8)

	for i := 0; i < 8; i++ {
		a[0][i] = c[i]
	}
	out := MultiMux3(api, a, sel)
	return out[0]
}
