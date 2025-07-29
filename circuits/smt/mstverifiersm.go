package smt

import "github.com/consensys/gnark/frontend"

// SMTVerifierSM
// https://github.com/iden3/circomlib/blob/master/circuits/smt/smtverifiersm.circom
func SMTVerifierSM(api frontend.API, is0, levIns, fnc, prevTop, prevI0, prevIold, prevInew, prevNa frontend.Variable) (
	frontend.Variable, frontend.Variable, frontend.Variable, frontend.Variable, frontend.Variable) {

	prevTopLevIns := api.Mul(prevTop, levIns)
	prevTopLevInsFnc := api.Mul(prevTopLevIns, fnc)

	stTop := api.Sub(prevTop, prevTopLevIns)
	stInew := api.Sub(prevTopLevIns, prevTopLevInsFnc)
	stIold := api.Mul(prevTopLevInsFnc, api.Sub(1, is0))
	stI0 := api.Mul(prevTopLevIns, is0)
	stNa := api.Add(prevNa, prevInew, prevIold, prevI0)

	return stTop, stI0, stIold, stInew, stNa
}
