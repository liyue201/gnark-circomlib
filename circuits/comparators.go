package circuits

import "github.com/consensys/gnark/frontend"

func IsZero(api frontend.API, a frontend.Variable) frontend.Variable {
	return api.IsZero(a)
}

func BoolNeg(api frontend.API, a frontend.Variable) frontend.Variable {
	api.AssertIsBoolean(a)
	return api.Xor(a, 1)
}

func IsEqual(api frontend.API, a frontend.Variable, b frontend.Variable) frontend.Variable {
	return api.IsZero(api.Sub(a, b))
}

func ForceEqualIfEnabled(api frontend.API, a, b, enabled frontend.Variable) {
	c := api.IsZero(api.Sub(a, b))
	api.AssertIsEqual(api.Mul(api.Sub(1, c), enabled), 0)
}

func LessThan(api frontend.API, a frontend.Variable, b frontend.Variable) frontend.Variable {
	return IsEqual(api, api.Cmp(a, b), -1)
}

func LessEqThan(api frontend.API, a frontend.Variable, b frontend.Variable) frontend.Variable {
	return BoolNeg(api, GreaterThan(api, a, b))
}

func GreaterThan(api frontend.API, a frontend.Variable, b frontend.Variable) frontend.Variable {
	return IsEqual(api, api.Cmp(a, b), 1)
}

func GreaterEqThan(api frontend.API, a frontend.Variable, b frontend.Variable) frontend.Variable {
	return BoolNeg(api, LessThan(api, a, b))
}
