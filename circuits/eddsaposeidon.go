package circuits

import (
	tedwards "github.com/consensys/gnark-crypto/ecc/twistededwards"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/std/algebra/native/twistededwards"
	"github.com/consensys/gnark/std/signature/eddsa"
)

func EdDSAPoseidonVerifier(api frontend.API, signature eddsa.Signature, msg frontend.Variable, pubKey eddsa.PublicKey) {
	curve, err := twistededwards.NewEdCurve(api, tedwards.BN254)
	if err != nil {
		panic(err)
	}
	hashFun := NewPoseidonHash(api)
	err = eddsa.Verify(curve, signature, msg, pubKey, hashFun)
	if err != nil {
		panic(err)
	}
}
