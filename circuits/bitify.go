package circuits

import (
	"github.com/consensys/gnark/frontend"
)

func Num2Bits(api frontend.API, v frontend.Variable, n int) []frontend.Variable {
	return api.ToBinary(v, n)
}

func Num2BitsStrict(api frontend.API, v frontend.Variable, n int)  []frontend.Variable {
	return api.ToBinary(v, n)
}