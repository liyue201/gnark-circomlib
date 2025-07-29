package smt

import (
	"github.com/consensys/gnark/frontend"
)

const NLevels = 254

// SMTVerifyInsertAnElement verifies the correctness of the new root hash after inserting an element into the Sparse Merkle Tree (SMT).
// This function calls SMTProcessor with the appropriate parameters to compute the new root and asserts that it matches the expected newRoot.
// Parameters:
// - api: The API interface used for constraint operations.
// - oldRoot: The original root of the SMT before the insertion.
// - newRoot: The expected root of the SMT after the insertion.
// - siblings: An array of sibling nodes in the path from the leaf to the root.
// - oldKey, oldValue: The key-value pair of the node being replaced (if any).
// - isOld0: A flag indicating whether the old node is a zero node.
// - newKey, newValue: The key-value pair of the new node being inserted.
func SMTVerifyInsertAnElement(api frontend.API, oldRoot, newRoot frontend.Variable, siblings [NLevels]frontend.Variable,
	oldKey, oldValue, isOld0, newKey, newValue frontend.Variable) {
	resultRoot := SMTProcessor(api, oldRoot, siblings, oldKey, oldValue, isOld0, newKey, newValue, [2]frontend.Variable{1, 0}, 0)
	api.AssertIsEqual(resultRoot, newRoot)
}

// SMTProcessor
// https://github.com/iden3/circomlib/blob/master/circuits/smt/smtprocessor.circom
func SMTProcessor(api frontend.API, oldRoot frontend.Variable, siblings [NLevels]frontend.Variable,
	oldKey, oldValue, isOld0, newKey, newValue frontend.Variable, fnc [2]frontend.Variable, isDummy frontend.Variable) frontend.Variable {

	enabled := api.Add(fnc[0], fnc[1], api.Neg(api.Mul(fnc[0], fnc[1])))
	enabled = api.Select(isDummy, 0, enabled)

	hash1Old := SMTHash1(api, oldKey, oldValue)
	hash1New := SMTHash1(api, newKey, newValue)

	n2bOld := api.ToBinary(oldKey)
	n2bNew := api.ToBinary(newKey)

	smtLevIns := SMTLevIns(api, enabled, siblings)
	var xors [NLevels]frontend.Variable
	for i := 0; i < NLevels; i++ {
		xors[i] = api.Xor(n2bOld[i], n2bNew[i])
	}

	var (
		stStop = make([]frontend.Variable, NLevels)
		stOld0 = make([]frontend.Variable, NLevels)
		stBot  = make([]frontend.Variable, NLevels)
		stNew1 = make([]frontend.Variable, NLevels)
		stNa   = make([]frontend.Variable, NLevels)
		stUpd  = make([]frontend.Variable, NLevels)
	)

	for i := 0; i < NLevels; i++ {
		if i == 0 {
			stStop[i], stOld0[i], stBot[i], stNew1[i], stNa[i], stUpd[i] = SMTProcessorSM(api, xors[i], isOld0, smtLevIns[i], fnc,
				enabled, 0, 0, 0, api.Sub(1, enabled), 0)
		} else {
			stStop[i], stOld0[i], stBot[i], stNew1[i], stNa[i], stUpd[i] = SMTProcessorSM(api, xors[i], isOld0, smtLevIns[i], fnc,
				stStop[i-1], stOld0[i-1], stBot[i-1], stNew1[i-1], stNa[i-1], stUpd[i-1])
		}
	}
	api.AssertIsEqual(api.Add(stNa[NLevels-1], stNew1[NLevels-1], stOld0[NLevels-1], stUpd[NLevels-1]), 1)

	tmpOldRoot := frontend.Variable(0)
	tmpNewRoot := frontend.Variable(0)
	for i := NLevels - 1; i >= 0; i-- {
		tmpOldRoot, tmpNewRoot = SMTProcessorLevel(api, stStop[i], stOld0[i], stBot[i], stNew1[i], stNa[i], stUpd[i], siblings[i],
			hash1Old, hash1New, n2bNew[i], tmpOldRoot, tmpNewRoot)
	}

	s := api.Mul(fnc[0], fnc[1])
	topSwitcherL := api.Select(s, tmpNewRoot, tmpOldRoot)
	topSwitcherR := api.Select(s, tmpOldRoot, tmpNewRoot)

	eq := api.IsZero(api.Sub(oldRoot, topSwitcherL))
	eq = api.Select(enabled, eq, 1)
	api.AssertIsEqual(eq, 1)

	newRoot := api.Add(api.Mul(enabled, api.Sub(topSwitcherR, oldRoot)), oldRoot)

	areKeyEquals := api.IsZero(api.Sub(oldKey, newKey))
	keysOk := api.And(api.Sub(1, fnc[0]), fnc[1])
	keysOk = api.And(keysOk, api.Sub(1, areKeyEquals))

	api.AssertIsEqual(keysOk, 0)

	return newRoot
}
