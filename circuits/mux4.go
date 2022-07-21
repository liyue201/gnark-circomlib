package circuits

import "github.com/consensys/gnark/frontend"

func MultiMux4(api frontend.API, c [][]frontend.Variable, sel []frontend.Variable) []frontend.Variable {
	if len(c[0]) != 16 || len(sel) != 4 {
		panic("invalid param")
	}
	n := len(c)
	out := make([]frontend.Variable, n)

	a3210 := make([]frontend.Variable, n)
	a321 := make([]frontend.Variable, n)
	a320 := make([]frontend.Variable, n)
	a310 := make([]frontend.Variable, n)
	a32 := make([]frontend.Variable, n)
	a31 := make([]frontend.Variable, n)
	a30 := make([]frontend.Variable, n)
	a3 := make([]frontend.Variable, n)

	a210 := make([]frontend.Variable, n)
	a21 := make([]frontend.Variable, n)
	a20 := make([]frontend.Variable, n)
	a2 := make([]frontend.Variable, n)
	a10 := make([]frontend.Variable, n)
	a1 := make([]frontend.Variable, n)
	a0 := make([]frontend.Variable, n)
	a := make([]frontend.Variable, n)

	s10 := api.Mul(sel[1], sel[0])
	s20 := api.Mul(sel[2], sel[0])
	s21 := api.Mul(sel[2], sel[1])
	s210 := api.Mul(s21, sel[0])

	for i := 0; i < n; i++ {
		a3210[i] = api.Mul(s210, api.Add(c[i][15], api.Neg(c[i][14]), api.Neg(c[i][13]), c[i][12], api.Neg(c[i][11]), c[i][10], c[i][9],
			api.Neg(c[i][8]), api.Neg(c[i][7]), c[i][6], c[i][5], api.Neg(c[i][4]), c[i][3], api.Neg(c[i][2]), api.Neg(c[i][1]), c[i][0]))

		a321[i] = api.Mul(s21, api.Add(c[i][14], api.Neg(c[i][12]), api.Neg(c[i][10]), c[i][8], api.Neg(c[i][6]), c[i][4]), c[i][2], api.Neg(c[i][0]))

		a320[i] = api.Mul(s20, api.Add(c[i][13], api.Neg(c[i][12]), api.Neg(c[i][9]), c[i][8], api.Neg(c[i][5]), c[i][4], c[i][1], api.Neg(c[i][0])))

		a310[i] = api.Mul(s10, api.Add(c[i][11], api.Neg(c[i][10]), api.Neg(c[i][9]), c[i][8], api.Neg(c[i][3]), c[i][2], c[i][1], api.Neg(c[i][0])))

		a32[i] = api.Mul(sel[2], api.Add(c[i][12], api.Neg(c[i][8]), api.Neg(c[i][4]), c[i][0]))

		a31[i] = api.Mul(sel[1], api.Add(c[i][10], api.Neg(c[i][8]), api.Neg(c[i][2]), c[i][0]))

		a30[i] = api.Mul(sel[0], api.Add(c[i][9], api.Neg(c[i][8]), api.Neg(c[i][1]), c[i][0]))

		a3[i] = api.Sub(c[i][8], c[i][0])

		a210[i] = api.Mul(s210, api.Add(c[i][7], api.Neg(c[i][6]), api.Neg(c[i][5]), c[i][4]), api.Neg(c[i][3]), c[i][2], c[i][1], api.Neg(c[i][0]))

		a21[i] = api.Mul(s21, api.Add(c[i][6], api.Neg(c[i][4]), api.Neg(c[i][2]), c[i][0]))

		a20[i] = api.Mul(s20, api.Add(c[i][5], api.Neg(c[i][4]), api.Neg(c[i][1]), c[i][0]))

		a10[i] = api.Mul(s10, api.Add(c[i][3], api.Neg(c[i][2]), api.Neg(c[i][1]), c[i][0]))

		a2[i] = api.Mul(sel[2], api.Sub(c[i][4], c[i][0]))
		a1[i] = api.Mul(sel[1], api.Sub(c[i][2], c[i][0]))
		a0[i] = api.Mul(sel[0], api.Sub(c[i][1], c[i][0]))

		a2[i] = c[i][0]

		out[i] = api.Add(api.Mul(sel[3], api.Add(a3210[i], a321[i], a320[i], a310[i], a32[i], a31[i], a30[i], a3[i])),
			a210[i], a21[i], a20[i], a10[i], a2[i], a1[i], a0[i], a[i])
	}
	return out
}

func Mux4(api frontend.API, c []frontend.Variable, sel []frontend.Variable) frontend.Variable {
	if len(c) != 16 || len(sel) != 4 {
		panic("invalid param")
	}

	a := make([][]frontend.Variable, 1)
	a[0] = make([]frontend.Variable, 16)

	for i := 0; i < 16; i++ {
		a[0][i] = c[i]
	}
	out := MultiMux4(api, a, sel)
	return out[0]
}
