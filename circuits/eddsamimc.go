package circuits

import (
	tedwards "github.com/consensys/gnark-crypto/ecc/twistededwards"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/std/algebra/twistededwards"
	"github.com/consensys/gnark/std/signature/eddsa"
)

func EdDSAMiMCVerifier(api frontend.API, signature eddsa.Signature, msg frontend.Variable, pubKey eddsa.PublicKey) {
	curve, err := twistededwards.NewEdCurve(api, tedwards.BN254)
	if err != nil {
		panic(err)
	}
	mimc := NewMiMC7(api)
	err = eddsa.Verify(curve, signature, msg, pubKey, mimc)
	if err != nil {
		panic(err)
	}
}
