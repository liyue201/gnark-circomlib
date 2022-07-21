package circuits

import (
	"github.com/consensys/gnark/frontend"
)

/*
   Assume sel is binary.
   If sel == 0 then outL = L and outR=R
   If sel == 1 then outL = R and outR=L

*/

func Switcher(api frontend.API, sel, l, r frontend.Variable) (frontend.Variable, frontend.Variable) {
	aux := api.Mul(api.Sub(r, l), sel)

	outL := api.Add(aux, l)
	outR := api.Add(api.Sub(0, aux), r)

	return outL, outR
}
