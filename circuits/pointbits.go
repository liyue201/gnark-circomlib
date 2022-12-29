package circuits

import "github.com/consensys/gnark/frontend"

func Point2BitsStrict(api frontend.API, in []frontend.Variable) []frontend.Variable {

	n2bX := api.ToBinary(in[0], 254)
	n2bY := api.ToBinary(in[0], 254)

	sign := Sign(api, n2bX)

	out := MakeVariableArray(256)
	for i := 0; i < 254; i++ {
		out[i] = n2bY[i]
	}
	out[254] = 0
	out[255] = sign

	return out
}
