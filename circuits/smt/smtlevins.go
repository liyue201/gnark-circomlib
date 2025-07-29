package smt

import (
	"github.com/consensys/gnark/frontend"
)

// SMTLevIns
// https://github.com/iden3/circomlib/blob/master/circuits/smt/smtlevins.circom
func SMTLevIns(api frontend.API, enabled frontend.Variable, siblings [NLevels]frontend.Variable) [NLevels]frontend.Variable {
	var levIns [NLevels]frontend.Variable

	var isZero [NLevels]frontend.Variable
	for i := 0; i < NLevels; i++ {
		isZero[i] = api.IsZero(siblings[i])
	}
	api.AssertIsEqual(api.Mul(api.Sub(isZero[NLevels-1], 1), enabled), 0)

	levIns[NLevels-1] = api.Sub(1, isZero[NLevels-2])

	done := levIns[NLevels-1]
	for i := NLevels - 2; i > 0; i-- {
		levIns[i] = api.Mul(api.Sub(1, done), api.Sub(1, isZero[i-1]))
		done = api.Add(levIns[i], done)
	}
	levIns[0] = api.Sub(1, done)
	return levIns
}
