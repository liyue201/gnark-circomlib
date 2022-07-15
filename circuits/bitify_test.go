package circuits

import (
	"fmt"
	"github.com/consensys/gnark/frontend/cs/r1cs"
	"testing"

	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark/backend"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/test"
)

type circuitNum2Bits struct {
	A    frontend.Variable   `gnark:",public"`
	Bits [8]frontend.Variable `gnark:",public"`
}

func (t *circuitNum2Bits) Define(api frontend.API) error {
	bits := Num2Bits(api, t.A, 8)
	for i := 0; i < 8; i++ {
		api.AssertIsEqual(bits[i], t.Bits[i])
	}
	return nil
}

func Test_Num2Bits(t *testing.T) {
	assert := test.NewAssert(t)
	var circuit circuitNum2Bits

	assert.ProverSucceeded(&circuit, &circuitNum2Bits{
		A:    9,
		Bits: [8]frontend.Variable{1, 0, 0, 1, 0, 0, 0, 0},
	}, test.WithCurves(ecc.BN254), test.WithBackends(backend.GROTH16))

	_r1cs, _ := frontend.Compile(ecc.BN254, r1cs.NewBuilder, &circuit)
	internal, secret, public := _r1cs.GetNbVariables()
	fmt.Printf("public, secret, internal %v, %v, %v\n", public, secret, internal)
}
