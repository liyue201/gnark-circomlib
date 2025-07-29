package smt

import "github.com/consensys/gnark/frontend"

// SMTProcessorLevel
// https://github.com/iden3/circomlib/blob/master/circuits/smt/smtprocessorlevel.circom
func SMTProcessorLevel(api frontend.API, stStop, stOld0, stBot, stNew1, stNa, stdUpd, sibling,
	old1Leaf, new1Leaf, newlrbit, oldChild, newChild frontend.Variable) (frontend.Variable, frontend.Variable) {

	var aux [4]frontend.Variable

	// Old side
	l := api.Select(newlrbit, sibling, oldChild)
	r := api.Select(newlrbit, oldChild, sibling)
	oldProofHash := SMTHash2(api, l, r)
	aux[0] = api.Mul(old1Leaf, api.Add(stBot, stNew1, stdUpd))
	oldRoot := api.Add(aux[0], api.Mul(oldProofHash, stStop))

	//New side
	aux[1] = api.Mul(newChild, api.Add(stStop, stBot))
	aux[2] = api.Mul(sibling, stStop)

	newSwitcherL := api.Add(aux[1], api.Mul(new1Leaf, stNew1))
	newSwitcherR := api.Add(aux[2], api.Mul(old1Leaf, stNew1))

	l = api.Select(newlrbit, newSwitcherR, newSwitcherL)
	r = api.Select(newlrbit, newSwitcherL, newSwitcherR)

	newProofHash := SMTHash2(api, l, r)
	aux[3] = api.Mul(newProofHash, api.Add(stStop, stBot, stNew1))

	newRoot := api.Add(aux[3], api.Mul(new1Leaf, api.Add(stOld0, stdUpd)))
	return oldRoot, newRoot
}
