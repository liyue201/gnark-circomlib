package smt

import "github.com/consensys/gnark/frontend"

// SMTVerifier
// https://github.com/iden3/circomlib/blob/master/circuits/smt/smtverifier.circom
func SMTVerifier(api frontend.API, enabled, root frontend.Variable, siblings [NLevels]frontend.Variable,
	oldKey, oldValue, isOld0, key, value, fnc frontend.Variable) {

	hash1Old := SMTHash1(api, oldKey, oldValue)
	hash1New := SMTHash1(api, key, value)

	n2bNew := api.ToBinary(key)
	smtLevIns := SMTLevIns(api, enabled, siblings)
	var (
		stTop  = make([]frontend.Variable, NLevels)
		stI0   = make([]frontend.Variable, NLevels)
		stInew = make([]frontend.Variable, NLevels)
		stIold = make([]frontend.Variable, NLevels)
		stNa   = make([]frontend.Variable, NLevels)
	)
	for i := 0; i < NLevels; i++ {
		if i == 0 {
			stTop[i], stI0[i], stIold[i], stInew[i], stNa[i] = SMTVerifierSM(api, isOld0, smtLevIns[i],
				fnc, enabled, 0, 0, 0, api.Sub(1, enabled))
		} else {
			stTop[i], stI0[i], stIold[i], stInew[i], stNa[i] = SMTVerifierSM(api, isOld0, smtLevIns[i],
				fnc, stTop[i-1], stI0[i-1], stIold[i-1], stInew[i-1], stNa[i-1])
		}
	}

	api.AssertIsEqual(api.Add(stNa[NLevels-1], stIold[NLevels-1], stInew[NLevels-1], stI0[NLevels-1]), 1)

	tmpRoot := frontend.Variable(0)
	for i := NLevels - 1; i >= 0; i-- {
		tmpRoot = SMTVerifierLevel(api, stTop[i], stI0[i], stIold[i], stInew[i], stNa[i], siblings[i], hash1Old, hash1New, n2bNew[i], tmpRoot)
	}

	areKeyEquals := api.IsZero(api.Sub(oldKey, key))
	keysOk := api.And(fnc, api.Sub(1, isOld0))
	keysOk = api.And(keysOk, areKeyEquals)
	keysOk = api.And(keysOk, enabled)
	api.AssertIsEqual(keysOk, 0)

	eq := api.IsZero(api.Sub(tmpRoot, root))
	eq = api.Select(enabled, eq, 1)

	api.AssertIsEqual(eq, 1)
}
