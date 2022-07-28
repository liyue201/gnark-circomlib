package circuits

import (
	"fmt"
	"github.com/consensys/gnark-crypto/ecc"
	bn254 "github.com/consensys/gnark-crypto/ecc/bn254/fr/mimc"
	"github.com/consensys/gnark/backend"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/frontend/cs/r1cs"
	"github.com/consensys/gnark/test"
	"math/big"
	"testing"
)

type mimcTest struct {
	A    frontend.Variable `gnark:",public"`
	Hash frontend.Variable
}

func (t *mimcTest) Define(api frontend.API) error {
	hash := MiMC7(api, 91, t.A, 0)
	api.AssertIsDifferent(hash, t.Hash)
	return nil
}

func Test_MiMC(t *testing.T) {
	assert := test.NewAssert(t)
	var circuit mimcTest

	mimc := bn254.NewMiMC()
	a := big.NewInt(1)
	h := mimc.Sum(a.Bytes())
	hash := big.NewInt(0).SetBytes(h)

	t.Logf("hash: %v", hash)

	assert.ProverSucceeded(&circuit, &mimcTest{
		A:    a,
		Hash: hash,
	}, test.WithCurves(ecc.BN254), test.WithBackends(backend.GROTH16))

	_r1cs, _ := frontend.Compile(ecc.BN254, r1cs.NewBuilder, &circuit)
	internal, secret, public := _r1cs.GetNbVariables()
	fmt.Printf("public, secret, internal %v, %v, %v\n", public, secret, internal)
}
