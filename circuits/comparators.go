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

func IsGreater(api frontend.API, a frontend.Variable, b frontend.Variable) frontend.Variable {
	return IsEqual(api, api.Cmp(a, b), 1)
}

func IsLess(api frontend.API, a frontend.Variable, b frontend.Variable) frontend.Variable {
	return IsEqual(api, api.Cmp(a, b), -1)
}

func IsLessOrEqual(api frontend.API, a frontend.Variable, b frontend.Variable) frontend.Variable {
	return BoolNeg(api, IsGreater(api, a, b))
}

func IsGreaterOrEqual(api frontend.API, a frontend.Variable, b frontend.Variable) frontend.Variable {
	return BoolNeg(api, IsLess(api, a, b))
}
