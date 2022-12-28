package circuits

import "github.com/consensys/gnark/frontend"

func BabyAdd(api frontend.API, x1, y1, x2, y2 frontend.Variable) (frontend.Variable, frontend.Variable) {

	a := 168700
	d := 168696

	beta := api.Mul(x1, y2)
	gamma := api.Mul(y1, x2)
	delta := api.Mul(api.Sub(y1, api.Mul(a, x1)), api.Add(x2, y2))

	tau := api.Mul(beta, gamma)

	xout := api.Div(api.Add(beta, gamma), api.Add(1, api.Mul(d, tau)))

	yout := api.Div(api.Add(delta, api.Mul(a, beta), api.Neg(gamma)),
		api.Sub(1, api.Mul(d, tau)))

	return xout, yout
}
