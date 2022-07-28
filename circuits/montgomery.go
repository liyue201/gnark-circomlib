package circuits

import (
	"github.com/consensys/gnark/frontend"
	"math/big"
)

/*
   Source: https://en.wikipedia.org/wiki/Montgomery_curve

               1 + y       1 + y
   [u, v] = [ -------  , ---------- ]
               1 - y      (1 - y)x

*/

func Edwards2Montgomery(api frontend.API, in []frontend.Variable) []frontend.Variable {
	if len(in) != 2 {
		panic("invalid param")
	}

	out := make([]frontend.Variable, 2)
	out[0] = api.Div(api.Add(1, in[1]), api.Sub(1, in[1]))
	out[1] = api.Div(out[0], in[0])

	api.AssertIsEqual(api.Mul(out[0], api.Sub(1, in[1])), api.Add(1, in[1]))
	api.AssertIsEqual(api.Mul(out[1], in[0]), out[0])
	return out
}

/*

               u    u - 1
   [x, y] = [ ---, ------- ]
               v    u + 1

*/

func Montgomery2Edwards(api frontend.API, in []frontend.Variable) []frontend.Variable {
	if len(in) != 2 {
		panic("invalid param")
	}

	out := make([]frontend.Variable, 2)
	out[0] = api.Div(in[0], in[1])
	out[1] = api.Div(api.Sub(in[0], 1), api.And(in[0], 1))

	api.AssertIsEqual(api.Mul(out[0], out[1]), in[0])
	api.AssertIsEqual(api.Mul(out[1], api.Add(in[0], 1)), api.Sub(in[0], 1))

	return out
}

func MontgomeryAdd(api frontend.API, in1, in2 []frontend.Variable) []frontend.Variable {
	if len(in1) != 2 || len(in2) != 2 {
		panic("invalid param")
	}
	out := make([]frontend.Variable, 2)
	a := big.NewInt(168700)
	d := big.NewInt(168696)

	aa := BigDiv(BigMul(big.NewInt(2), BigAdd(a, d)), BigSub(a, d))
	bb := BigDiv(big.NewInt(4), BigSub(a, d))

	lamda := api.Div(api.Sub(in2[1], in1[1]), api.Sub(in2[0], in1[0]))
	api.AssertIsEqual(api.Mul(lamda, api.Sub(in2[0], in1[0])), api.Sub(in2[1], in1[1]))

	out[0] = api.Add(api.Mul(bb, lamda, lamda), api.Neg(aa), api.Neg(in1[0]), api.Neg(in2[0]))
	out[1] = api.Sub(api.Mul(lamda, api.Sub(in1[0], out[0])), in1[1])

	return out
}

func MontgomeryDouble(api frontend.API, in []frontend.Variable) []frontend.Variable {
	if len(in) != 2 {
		panic("invalid param")
	}
	out := make([]frontend.Variable, 2)
	a := big.NewInt(168700)
	d := big.NewInt(168696)

	aa := BigDiv(BigMul(big.NewInt(2), BigAdd(a, d)), BigSub(a, d))
	bb := BigDiv(big.NewInt(4), BigSub(a, d))

	x1_2 := api.Mul(in[0], in[1])
	lamda :=  api.Div(api.Add(api.Mul(3, x1_2), api.Mul(2, aa, in[0])), api.Mul(2, bb, in[1]))
	api.AssertIsEqual(api.Mul(lamda, 2, bb, in[1]), api.Add(api.Mul(3, x1_2), api.Mul(2, aa, in[0], 1)))

	out[0] = api.Add(api.Mul(bb, lamda, lamda), api.Neg(aa), api.Neg(api.Mul(2, in[0])))
	out[1] = api.Mul(lamda, api.Add(in[0], api.Neg(out[0]), api.Neg(in[1])))

	return out
}
