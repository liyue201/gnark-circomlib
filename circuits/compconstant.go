package circuits

import (
	"fmt"
	"github.com/consensys/gnark/frontend"
	"math/big"
)

// CompConstant returns 1 if in (in binary) > ct
func CompConstant(api frontend.API, in []frontend.Variable, ct *big.Int) frontend.Variable {
	if len(in) != 254 {
		panic(fmt.Sprintf("CompConstant: invalid len: %v", len(in)))
	}
	var (
		parts [127]frontend.Variable
		clsb  *big.Int
		cmsb  *big.Int
	)

	sum := frontend.Variable(0)
	b := frontend.Variable(BigSub(Lsh(1, 128), big.NewInt(1)))
	a := frontend.Variable(1)
	e := frontend.Variable(1)

	for i := 0; i < 127; i++ {
		clsb = BigAnd(BigRsh(ct, uint(i<<1)), big.NewInt(1))
		cmsb = BigAnd(BigRsh(ct, uint((i<<1)+1)), big.NewInt(1))
		slsb := in[i<<1]
		smsb := in[(i<<1)+1]
		if cmsb.Int64() == 0 && clsb.Int64() == 0 {
			parts[i] = api.Add(api.Mul(-1, b, smsb, slsb), api.Mul(b, smsb), api.Mul(b, slsb))
		} else if cmsb.Int64() == 0 && clsb.Int64() == 1 {
			parts[i] = api.Add(api.Mul(a, smsb, slsb), api.Mul(-1, a, slsb), api.Mul(b, smsb), api.Mul(-1, a, smsb), a)
		} else if cmsb.Int64() == 1 && clsb.Int64() == 0 {
			parts[i] = api.Add(api.Mul(b, smsb, slsb), api.Mul(-1, a, smsb), a)
		} else {
			parts[i] = api.Add(api.Mul(-1, a, smsb, slsb), a)
		}
		sum = api.Add(sum, parts[i])
		b = api.Sub(b, e)
		a = api.Add(a, e)
		e = api.Mul(e, 2)
	}
	bits := api.ToBinary(sum, 135)

	return bits[127]
}
