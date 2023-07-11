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

type circuitIsEqual struct {
	A frontend.Variable `gnark:",public"`
	B frontend.Variable `gnark:",public"`
}

func (t *circuitIsEqual) Define(api frontend.API) error {
	//q := api.Compiler().Field()
	//fmt.Printf("%s\n", q.String())

	//q.Sub(q, big.NewInt(1))

	//fmt.Printf("%s\n", q.String())

	//api.AssertIsEqual(IsEqual(api, t.A, q), 1)
	api.AssertIsEqual(IsEqual(api, t.A, -1), 1)

	api.AssertIsEqual(IsEqual(api, t.B, 123), 1)
	return nil
}

func TestIsEqual(t *testing.T) {
	assert := test.NewAssert(t)
	var circuit circuitIsEqual

	assert.ProverSucceeded(&circuit, &circuitIsEqual{
		A: -1,
		B: 123,
	}, test.WithCurves(ecc.BN254), test.WithBackends(backend.GROTH16), test.WithCompileOpts(frontend.IgnoreUnconstrainedInputs()))

	_r1cs, _ := frontend.Compile(ecc.BN254.ScalarField(), r1cs.NewBuilder, &circuit)
	internal, secret, public := _r1cs.GetNbVariables()
	fmt.Printf("public, secret, internal %v, %v, %v\n", public, secret, internal)
}
