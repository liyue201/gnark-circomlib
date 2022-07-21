package circuits

import (
	"github.com/consensys/gnark/frontend"
)

func Xor(api frontend.API, a, b frontend.Variable) frontend.Variable {
	return api.Xor(a, b)
}

func And(api frontend.API, a, b frontend.Variable) frontend.Variable {
	return api.And(a, b)
}

func Or(api frontend.API, a, b frontend.Variable) frontend.Variable {
	return api.Or(a, b)
}

func Not(api frontend.API, a frontend.Variable) frontend.Variable {
	api.AssertIsBoolean(a)
	return api.Sub(1, a)
}

func NAND(api frontend.API, a, b frontend.Variable) frontend.Variable {
	return api.Sub(1, And(api, a, b))
}

func Nor(api frontend.API, a, b frontend.Variable) frontend.Variable {
	return api.Sub(1, Or(api, a, b))
}

func MultiAnd(api frontend.API, in[] frontend.Variable) frontend.Variable {
	out := frontend.Variable(1)
	for i := 0; i < len(in); i++ {
		out = api.And(out, in[i])
	}
	return out
}
