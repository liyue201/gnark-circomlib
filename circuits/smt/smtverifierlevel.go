package smt

import "github.com/consensys/gnark/frontend"

// SMTVerifierLevel
// https://github.com/iden3/circomlib/blob/master/circuits/smt/smtverifierlevel.circom
func SMTVerifierLevel(api frontend.API, stTop, stI0, stIold, stInew, stNa, sibling,
	old1Leaf, new1Leaf, lrbit, child frontend.Variable) frontend.Variable {

	// stI0 and stNa are unused but kept for consistency with the circom source.
	_ = stI0
	_ = stNa

	l := api.Select(lrbit, sibling, child)
	r := api.Select(lrbit, child, sibling)

	proofHash := SMTHash2(api, l, r)

	aux0 := api.Mul(proofHash, stTop)
	aux1 := api.Mul(old1Leaf, stIold)

	root := api.Add(aux0, aux1, api.Mul(new1Leaf, stInew))

	return root
}
