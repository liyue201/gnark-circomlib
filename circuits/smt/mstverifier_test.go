package smt

import (
	"context"
	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/test"
	merkletree "github.com/iden3/go-merkletree-sql/v2"
	"github.com/iden3/go-merkletree-sql/v2/db/memory"
	"github.com/stretchr/testify/require"
	"math/big"
	"testing"
)

type memberShipProofCircuit struct {
	Root     frontend.Variable
	Siblings [NLevels]frontend.Variable
	Key      frontend.Variable
	Value    frontend.Variable
}

func (m *memberShipProofCircuit) Define(api frontend.API) error {
	SMTVerifier(api, 1, m.Root, m.Siblings, 0, 0, 0, m.Key, m.Value, 0)
	return nil
}

func TestSMTVerifier(t *testing.T) {
	db := memory.NewMemoryStorage()
	ctx := context.Background()
	mt, err := merkletree.NewMerkleTree(ctx, db, NLevels)
	if err != nil {
		t.Fatal(err)
	}
	err = mt.Add(ctx, big.NewInt(1), big.NewInt(2))
	require.NoError(t, err)
	err = mt.Add(ctx, big.NewInt(3), big.NewInt(8))
	require.NoError(t, err)
	proof, err := mt.GenerateCircomVerifierProof(ctx, big.NewInt(1), mt.Root())
	require.NoError(t, err)
	if proof.Fnc != 0 {
		t.Fatalf("no include: %v", proof.Fnc)
	}

	witness := memberShipProofCircuit{
		Root:  proof.Root.BigInt(),
		Key:   proof.Key.BigInt(),
		Value: proof.Value.BigInt(),
	}
	for i := 0; i < NLevels; i++ {
		witness.Siblings[i] = proof.Siblings[i].BigInt()
	}
	assert := test.NewAssert(t)
	err = test.IsSolved(&memberShipProofCircuit{}, &witness, ecc.BN254.ScalarField())
	assert.NoError(err)
}
