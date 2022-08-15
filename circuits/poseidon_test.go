package circuits

import (
	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark/backend"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/frontend/cs/r1cs"
	"github.com/consensys/gnark/test"
	"github.com/iden3/go-iden3-crypto/poseidon"
	"math/big"
	"testing"
)

type circuitPoseidon struct {
	A    [3]frontend.Variable `gnark:",public"`
	Hash frontend.Variable    `gnark:",public"`
}

func (t *circuitPoseidon) Define(api frontend.API) error {
	hash := Poseidon(api, t.A[:])
	api.Println(t.Hash)
	api.Println(hash)
	api.AssertIsEqual(hash, t.Hash)
	return nil
}

func TestPoseidon(t *testing.T) {
	assert := test.NewAssert(t)
	var circuit circuitPoseidon

	input := []*big.Int{big.NewInt(1), big.NewInt(2), big.NewInt(3)}

	h, _ := poseidon.Hash(input)

	t.Logf("h: %s", h.String())

	assert.ProverSucceeded(&circuit, &circuitPoseidon{
		A:    [3]frontend.Variable{1, 2, 3},
		Hash: h,
	}, test.WithCurves(ecc.BN254), test.WithBackends(backend.GROTH16))

	_r1cs, _ := frontend.Compile(ecc.BN254, r1cs.NewBuilder, &circuit)
	internal, secret, public := _r1cs.GetNbVariables()
	t.Logf("public, secret, internal %v, %v, %v\n", public, secret, internal)
}
