package circuits

import (
	"github.com/consensys/gnark/frontend"
	"math/big"
)

func AliasCheck(api frontend.API, in []frontend.Variable) {
	api.AssertIsEqual(CompConstant(api, in, big.NewInt(0).Sub(Lsh(1, 255), big.NewInt(1))), 0)
}
