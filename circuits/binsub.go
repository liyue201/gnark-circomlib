package circuits

import "github.com/consensys/gnark/frontend"

/*
This component creates a binary substraction.


Main Constraint:
   (a[0]     * 2^0  +  a[1]     * 2^1  + ..... + a[n-1]    * 2^(n-1))  +
 +  2^n
 - (b[0]     * 2^0  +  b[1]     * 2^1  + ..... + b[n-1]    * 2^(n-1))
 ===
   out[0] * 2^0  + out[1] * 2^1 +   + out[n-1] *2^(n-1) + aux


    out[0]     * (out[0] - 1) === 0
    out[1]     * (out[0] - 1) === 0
    .
    .
    .
    out[n-1]   * (out[n-1] - 1) === 0
    aux * (aux-1) == 0

*/

func BinSub(api frontend.API, a, b []frontend.Variable) []frontend.Variable {
	if len(a) != len(b) {
		panic("invalid params")
	}
	n := len(a)

	lin := frontend.Variable(Lsh(1, uint(n)))
	for i := 0; i < n; i++ {
		lin = api.Add(lin, api.Mul(a[i], Lsh(1, uint(i))))
		lin = api.Sub(lin, b[i], Lsh(1, uint(i)))
	}
	out := api.ToBinary(lin, n)

	return out
}
