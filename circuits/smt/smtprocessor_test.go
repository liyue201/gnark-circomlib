package smt

import (
	"context"
	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/test"
	merkletree "github.com/iden3/go-merkletree-sql/v2"
	"github.com/iden3/go-merkletree-sql/v2/db/memory"
	"math/big"
	"testing"
)

type insertProofCircuit struct {
	OldRoot  frontend.Variable
	NewRoot  frontend.Variable
	Siblings [NLevels]frontend.Variable
	OldKey   frontend.Variable
	OldValue frontend.Variable
	IsOld0   frontend.Variable
	NewKey   frontend.Variable
	NewValue frontend.Variable
}

func (m *insertProofCircuit) Define(api frontend.API) error {
	SMTVerifyInsertAnElement(api, m.OldRoot, m.NewRoot, m.Siblings, m.OldKey, m.OldValue, m.IsOld0, m.NewKey, m.NewValue)
	return nil
}

func TestInsertProofCircuit(t *testing.T) {
	db := memory.NewMemoryStorage()
	ctx := context.Background()
	mt, err := merkletree.NewMerkleTree(ctx, db, NLevels)
	if err != nil {
		t.Fatal(err)
	}
	err = mt.Add(ctx, big.NewInt(1), big.NewInt(2))
	if err != nil {
		t.Fatalf("%v", err)
	}
	err = mt.Add(ctx, big.NewInt(3), big.NewInt(6))
	if err != nil {
		t.Fatalf("%v", err)
	}
	err = mt.Add(ctx, big.NewInt(6), big.NewInt(3))
	if err != nil {
		t.Fatalf("%v", err)
	}
	err = mt.Add(ctx, big.NewInt(8), big.NewInt(6))
	if err != nil {
		t.Fatalf("%v", err)
	}

	proof, err := mt.AddAndGetCircomProof(ctx, big.NewInt(99), big.NewInt(8))
	if err != nil {
		t.Fatal(err)
	}
	if proof.Fnc != 2 {
		t.Fatalf("no insert: %v", proof.Fnc)
	}

	witness := insertProofCircuit{
		OldRoot:  proof.OldRoot.BigInt(),
		NewRoot:  proof.NewRoot.BigInt(),
		OldKey:   proof.OldKey.BigInt(),
		OldValue: proof.OldValue.BigInt(),
		NewKey:   proof.NewKey.BigInt(),
		NewValue: proof.NewValue.BigInt(),
	}
	if proof.IsOld0 {
		witness.IsOld0 = 1
	} else {
		witness.IsOld0 = 0
	}
	for i := 0; i < NLevels; i++ {
		witness.Siblings[i] = proof.Siblings[i].BigInt()
	}
	assert := test.NewAssert(t)
	err = test.IsSolved(&insertProofCircuit{}, &witness, ecc.BN254.ScalarField())
	assert.NoError(err)
}
