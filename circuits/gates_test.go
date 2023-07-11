package circuits

import (
	"fmt"
	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark/backend"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/frontend/cs/r1cs"
	"github.com/consensys/gnark/test"
	"testing"
)

type circuitTest struct {
	A frontend.Variable `gnark:",public"`
	B frontend.Variable
}

func (t *circuitTest) Define(api frontend.API) error {
	api.AssertIsEqual(Or(api, t.A, t.B), 1)
	api.AssertIsEqual(Or(api, 0, t.A), 1)
	api.AssertIsEqual(Or(api, 0, 0), 0)
	api.AssertIsEqual(Or(api, 1, 1), 1)
	return nil
}

func Test_Or(t *testing.T) {
	assert := test.NewAssert(t)
	var circuit circuitTest

	assert.ProverSucceeded(&circuit, &circuitTest{
		A: 1,
		B: 0,
	}, test.WithCurves(ecc.BN254), test.WithBackends(backend.GROTH16))

	_r1cs, _ := frontend.Compile(ecc.BN254.ScalarField(), r1cs.NewBuilder, &circuit)
	internal, secret, public := _r1cs.GetNbVariables()
	fmt.Printf("public, secret, internal %v, %v, %v\n", public, secret, internal)
}
