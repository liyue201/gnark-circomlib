package smt

import (
	"github.com/consensys/gnark/frontend"
)

// SMTProcessorSM
// https://github.com/iden3/circomlib/blob/master/circuits/smt/smtprocessorsm.circom
func SMTProcessorSM(api frontend.API, xor, is0, levIns frontend.Variable, fnc [2]frontend.Variable,
	prevTop, prevOld0, prevBot, prevNew1, prevNa, prevUpd frontend.Variable) (
	frontend.Variable, frontend.Variable, frontend.Variable, frontend.Variable, frontend.Variable, frontend.Variable) {

	aux1 := api.Mul(prevTop, levIns)
	aux2 := api.Mul(aux1, fnc[0])

	stStop := api.Sub(prevTop, aux1)
	stOld0 := api.Mul(aux2, is0)
	stNew1 := api.Mul(api.Add(aux2, api.Neg(stOld0), prevBot), xor)
	stBot := api.Mul(api.Sub(1, xor), api.Add(aux2, api.Neg(stOld0), prevBot))
	stUpd := api.Sub(aux1, aux2)
	stNa := api.Add(prevNew1, prevOld0, prevNa, prevUpd)
	return stStop, stOld0, stBot, stNew1, stNa, stUpd

}
