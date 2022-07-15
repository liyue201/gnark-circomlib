package circuits

import (
	"fmt"
	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark/backend"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/frontend/cs/r1cs"
	"github.com/consensys/gnark/test"
	"math/big"
	"testing"
)

type circuitCompConstant struct {
	A frontend.Variable `gnark:",public"`
}

func (t *circuitCompConstant) Define(api frontend.API) error {
	in := api.ToBinary(t.A, 254)
	api.AssertIsEqual(CompConstant(api, in, big.NewInt(0)), 1)
	api.AssertIsEqual(CompConstant(api, in, big.NewInt(1)), 1)
	api.AssertIsEqual(CompConstant(api, in, big.NewInt(744)), 1)
	api.AssertIsEqual(CompConstant(api, in, big.NewInt(999)), 1)
	api.AssertIsEqual(CompConstant(api, in, big.NewInt(1000)), 0)
	api.AssertIsEqual(CompConstant(api, in, big.NewInt(1001)), 0)
	api.AssertIsEqual(CompConstant(api, in, big.NewInt(4155)), 0)
	return nil
}

func TestCompConstant(t *testing.T) {
	assert := test.NewAssert(t)
	var circuit circuitCompConstant

	assert.ProverSucceeded(&circuit, &circuitCompConstant{
		A: 1000,
	}, test.WithCurves(ecc.BN254), test.WithBackends(backend.PLONK))

	_r1cs, _ := frontend.Compile(ecc.BN254, r1cs.NewBuilder, &circuit)
	internal, secret, public := _r1cs.GetNbVariables()
	fmt.Printf("public, secret, internal %v, %v, %v\n", public, secret, internal)
}
