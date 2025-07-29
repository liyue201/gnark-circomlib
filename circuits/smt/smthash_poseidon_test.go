package smt

import (
	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/frontend/cs/r1cs"
	"github.com/consensys/gnark/test"
	"github.com/iden3/go-iden3-crypto/poseidon"
	"github.com/stretchr/testify/require"
	"math/big"
	"testing"
)

type poseidonTestCircuit struct {
	In   [2]frontend.Variable
	Hash frontend.Variable
}

func (m *poseidonTestCircuit) Define(api frontend.API) error {
	h := SMTHash2(api, m.In[0], m.In[1])
	api.AssertIsEqual(h, m.Hash)
	return nil
}

func TestHash2(t *testing.T) {

	in := []*big.Int{
		big.NewInt(1),
		big.NewInt(2),
	}
	res, err := poseidon.Hash(in)
	require.NoError(t, err)

	witness := poseidonTestCircuit{
		In:   [2]frontend.Variable{1, 2},
		Hash: res,
	}
	assert := test.NewAssert(t)
	err = test.IsSolved(&poseidonTestCircuit{},
		&witness, ecc.BN254.ScalarField())
	assert.NoError(err)

	cs, err := frontend.Compile(ecc.BN254.ScalarField(), r1cs.NewBuilder, &poseidonTestCircuit{}, frontend.IgnoreUnconstrainedInputs())
	require.NoError(t, err)
	t.Logf("NbConstraints: %v\n", cs.GetNbConstraints())
}
