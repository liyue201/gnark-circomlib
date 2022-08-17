package keccak

import (
	"fmt"
	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark/backend"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/frontend/cs/r1cs"
	"github.com/consensys/gnark/test"
	"github.com/ethereum/go-ethereum/crypto"
	"math/rand"
	"testing"
	"time"
)

type keccak512Circuit struct {
	M    [72*2-1]frontend.Variable `gnark:",public"`
	Hash [64]frontend.Variable `gnark:",public"`
}

func (t *keccak512Circuit) Define(api frontend.API) error {

	hash := Keccak512(api, t.M[:])
	for i := 0; i < len(hash); i++ {
		api.AssertIsEqual(hash[i], t.Hash[i])
	}
	//api.Println(hash...)

	return nil
}

func Test_Keccak512(t *testing.T) {
	assert := test.NewAssert(t)
	var circuit keccak512Circuit
	rand.Seed(time.Now().Unix())
	m := make([]byte, len(circuit.M))
	for i := 0; i< len(circuit.M); i++ {
		m[i] = byte(rand.Int() % 256)
		circuit.M[i] = m[i]
	}
	hash := crypto.Keccak512(m)
	var assignment keccak512Circuit
	for i := 0; i < len(assignment.M); i++ {
		assignment.M[i] = m[i]
	}
	for i := 0; i < len(assignment.Hash); i++ {
		assignment.Hash[i] = hash[i]
	}

	assert.ProverSucceeded(&circuit, &assignment, test.WithCurves(ecc.BN254), test.WithBackends(backend.GROTH16))

	_r1cs, err := frontend.Compile(ecc.BN254, r1cs.NewBuilder, &circuit, frontend.IgnoreUnconstrainedInputs())
	if err != nil {
		t.Errorf("Compile: %v", err)
	}
	internal, secret, public := _r1cs.GetNbVariables()
	fmt.Printf("public, secret, internal %v, %v, %v\n", public, secret, internal)
}
