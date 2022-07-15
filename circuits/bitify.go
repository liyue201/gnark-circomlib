package circuits

import (
	"github.com/consensys/gnark/frontend"
)

func Num2Bits(api frontend.API, in frontend.Variable, n int) []frontend.Variable {
	return api.ToBinary(in, n)
}

func Num2BitsStrict(api frontend.API, in frontend.Variable, n int) []frontend.Variable {
	bits := api.ToBinary(in, n)
	AliasCheck(api, bits)
	return bits
}

func Bits2Num(api frontend.API, in []frontend.Variable) frontend.Variable {
	return api.FromBinary(in...)
}

func Bits2NumStrict(api frontend.API, in []frontend.Variable) frontend.Variable {
	AliasCheck(api, in)
	return api.FromBinary(in...)
}

func Num2BitsNeg(api frontend.API, in frontend.Variable, n int) []frontend.Variable {
	var neg frontend.Variable
	if n == 0 {
		neg = frontend.Variable(0)
	} else {
		neg = api.Sub(Lsh(1, uint(n)), in)
	}
	out := api.ToBinary(neg, n)
	return out
}
