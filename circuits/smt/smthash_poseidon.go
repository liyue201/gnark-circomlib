package smt

import (
	"github.com/consensys/gnark/frontend"
	"github.com/liyue201/gnark-circomlib/circuits"
)

// SMTHash1
// https://github.com/iden3/circomlib/blob/master/circuits/smt/smthash_poseidon.circom
// Hash1 = H(key | value | 1)
func SMTHash1(api frontend.API, key, value frontend.Variable) frontend.Variable {
	return circuits.Poseidon(api, []frontend.Variable{key, value, 1})
}

// SMTHash2
// This component is used to create the 2 nodes.
// Hash2 = H(Hl | Hr)
func SMTHash2(api frontend.API, L, R frontend.Variable) frontend.Variable {
	return circuits.Poseidon(api, []frontend.Variable{L, R})
}
