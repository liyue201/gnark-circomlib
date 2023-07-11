package circuits

import (
	"github.com/consensys/gnark/frontend"
	"math/big"
)

// Sign returns 1 if in is positive, or 0 otherwise
func Sign(api frontend.API, in []frontend.Variable) frontend.Variable {
	// The order of bn254 is
	// FrModulus = 21888242871839275222246405745257275088548364400416034343698204186575808495617
	// (FrModulus-1)/2 =  10944121435919637611123202872628637544274182200208017171849102093287904247808
	//v, _ := big.NewInt(0).SetString("10944121435919637611123202872628637544274182200208017171849102093287904247808", 10)

	v := api.Compiler().Field()
	v.Sub(v, big.NewInt(1))
	v.Rsh(v, 1)
	return CompConstant(api, in, v)
}
