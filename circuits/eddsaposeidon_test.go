package circuits

import (
	crand "crypto/rand"
	"github.com/consensys/gnark-crypto/ecc"
	twistededwards2 "github.com/consensys/gnark-crypto/ecc/twistededwards"
	"github.com/consensys/gnark-crypto/signature/eddsa"
	"github.com/consensys/gnark/backend"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/std/algebra/twistededwards"
	eddsa2 "github.com/consensys/gnark/std/signature/eddsa"
	"github.com/consensys/gnark/test"
	"github.com/liyue201/gnark-circomlib/utils/poseidon"
	"math/big"
	"math/rand"
	"testing"
)

type circuitEddsaPoseidon struct {
	PublicKey eddsa2.PublicKey  `gnark:",public"`
	Signature eddsa2.Signature  `gnark:",public"`
	Message   frontend.Variable `gnark:",public"`
}

func (t *circuitEddsaPoseidon) Define(api frontend.API) error {
	EdDSAPoseidonVerifier(api, t.Signature, t.Message, t.PublicKey)
	return nil
}

func TestEdDSAPoseidonVerifier(t *testing.T) {
	assert := test.NewAssert(t)

	snarkCurve, _ := twistededwards.GetSnarkCurve(twistededwards2.BN254)

	// generate parameters for the signatures
	privKey, _ := eddsa.New(twistededwards2.BN254, crand.Reader)

	// pick a message to sign
	var msg big.Int
	msg.Rand(rand.New(rand.NewSource(0)), snarkCurve.Info().Fr.Modulus())

	t.Log("msg to sign", msg.String())
	msgData := msg.Bytes()

	hFunc := poseidon.NewPoseidon()

	// generate signature
	signature, err := privKey.Sign(msgData[:], hFunc)
	if err != nil {
		t.Errorf("Sign: %v", err)
	}

	// check if there is no problem in the signature
	pubKey := privKey.Public()
	checkSig, err := pubKey.Verify(signature, msgData[:], hFunc)
	if err != nil {
		t.Errorf("Sign: %v", err)
	}
	t.Logf("checkSig: %v", checkSig)

	// create and compile the circuit for signature verification
	var circuit circuitEddsaPoseidon

	var witness circuitEddsaPoseidon
	witness.Message = msg
	witness.PublicKey.Assign(snarkCurve, pubKey.Bytes())
	witness.Signature.Assign(snarkCurve, signature)

	assert.ProverSucceeded(&circuit, &witness, test.WithCurves(ecc.BN254), test.WithBackends(backend.GROTH16))

}
