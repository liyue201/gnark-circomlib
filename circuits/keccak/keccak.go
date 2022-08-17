package keccak

import (
	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark/backend/hint"
	"github.com/consensys/gnark/frontend"
	"math/big"
)


func init() {
	hint.Register(divFunc)
}

func u64Xor(api frontend.API, a frontend.Variable, b frontend.Variable, cs ...frontend.Variable) frontend.Variable {
	bitsA := api.ToBinary(a, 64)
	bitsB := api.ToBinary(b, 64)
	bitsRes := make([]frontend.Variable, 64)
	for i := 0; i < 64; i++ {
		bitsRes[i] = api.Xor(bitsA[i], bitsB[i])
	}

	for _, c := range cs {
		bitsC := api.ToBinary(c)
		for i := 0; i < 64; i++ {
			bitsRes[i] = api.Xor(bitsRes[i], bitsC[i])
		}
	}
	return api.FromBinary(bitsRes...)
}

func u64And(api frontend.API, a frontend.Variable, b frontend.Variable, cs ...frontend.Variable) frontend.Variable {
	bitsA := api.ToBinary(a, 64)
	bitsB := api.ToBinary(b, 64)
	bitsRes := make([]frontend.Variable, 64)
	for i := 0; i < 64; i++ {
		bitsRes[i] = api.And(bitsA[i], bitsB[i])
	}

	for _, c := range cs {
		bitsC := api.ToBinary(c)
		for i := 0; i < 64; i++ {
			bitsRes[i] = api.And(bitsRes[i], bitsC[i])
		}
	}

	return api.FromBinary(bitsRes...)
}

func u64Or(api frontend.API, a frontend.Variable, b frontend.Variable, cs ...frontend.Variable) frontend.Variable {
	bitsA := api.ToBinary(a, 64)
	bitsB := api.ToBinary(b, 64)
	bitsRes := make([]frontend.Variable, 64)
	for i := 0; i < 64; i++ {
		bitsRes[i] = api.Or(bitsA[i], bitsB[i])
	}

	for _, c := range cs {
		bitsC := api.ToBinary(c)
		for i := 0; i < 64; i++ {
			bitsRes[i] = api.Or(bitsRes[i], bitsC[i])
		}
	}

	return api.FromBinary(bitsRes...)
}

func u8Xor(api frontend.API, a frontend.Variable, b frontend.Variable, cs ...frontend.Variable) frontend.Variable {
	bitsA := api.ToBinary(a, 8)
	bitsB := api.ToBinary(b, 8)
	bitsRes := make([]frontend.Variable, 8)
	for i := 0; i < 8; i++ {
		bitsRes[i] = api.Xor(bitsA[i], bitsB[i])
	}

	for _, c := range cs {
		bitsC := api.ToBinary(c)
		for i := 0; i < 8; i++ {
			bitsRes[i] = api.Xor(bitsRes[i], bitsC[i])
		}
	}
	return api.FromBinary(bitsRes...)
}

func Keccak512(api frontend.API, mBytes []frontend.Variable) []frontend.Variable {

	rate := 72
	p := mBytes
	s := make([]frontend.Variable, 25)
	for i := 0; i < len(s); i++ {
		s[i] = 0
	}

	for len(p) >= rate {
		permute(api, s, p, rate)
		p = p[rate:]
	}
	if len(p) > 0 {
		p = append(p, 1)
		padNum := rate - len(p)
		if padNum > 0 {
			pad := make([]frontend.Variable, padNum)
			for i := 0; i < padNum; i++ {
				pad[i] = 0
			}
			p = append(p, pad...)
		}
		p[rate-1] = u8Xor(api, p[rate-1], 0x80)
		s = permute(api, s, p, rate)
		p = p[rate:]
	}
	return copyOutUnaligned(api, s)
}

func permute(api frontend.API, s, p []frontend.Variable, rate int) []frontend.Variable {
	api.Println(p...)
	buf := make([]frontend.Variable, rate/8)
	for i := 0; i < len(buf); i++ {
		var bits []frontend.Variable
		for j := 0; j < 8; j++ {
			bits = append(bits, api.ToBinary(p[i*8+j], 8)...)
		}
		buf[i] = api.FromBinary(bits...)
	}
	//api.Println(buf...)
	s = xorIn(api, s, buf)
	s = keccakF(api, s)
	return s
}

func copyOutUnaligned(api frontend.API, s []frontend.Variable) []frontend.Variable {
	out := make([]frontend.Variable, 64)
	r := 72
	w := 8
	for b := 0; b < 64; {
		for y := 0; y < 5; y++ {
			for x := 0; x < 5; x++ {
				if x+5*y < (r/w) && (b < 64) {
					bits := api.ToBinary(s[5*x+y], 64)
					for i := 0; i < 8; i++ {
						out[b+i] = api.FromBinary(bits[i*8 : (i+1)*8]...)
					}
					b += 8
				}
			}
		}
	}
	return out
}

func xorIn(api frontend.API, s []frontend.Variable, buf []frontend.Variable) []frontend.Variable {
	for y := 0; y < 5; y++ {
		for x := 0; x < 5; x++ {
			if x+5*y < len(buf) {
				s[5*x+y] = u64Xor(api, s[5*x+y], buf[x+5*y])
			}
		}
	}
	return s
}

func Sponge(api frontend.API, m []frontend.Variable) []frontend.Variable {
	//m[5] = 1
	//m[8], _ = new(big.Int).SetString("8000000000000000", 16)
	//api.Println(m...)

	r := 72
	w := 8
	size := len(m) * 8
	s := make([]frontend.Variable, 25)
	for i := 0; i < len(s); i++ {
		s[i] = frontend.Variable(0)
	}
	for i := 0; i < size/r; i++ {
		for y := 0; y < 5; y++ {
			for x := 0; x < 5; x++ {
				if x+5*y < r/w {
					s[5*x+y] = u64Xor(api, s[5*x+y], m[i*9+x+5*y])
				}
			}
		}
		s = keccakF(api, s)
	}

	//api.Println(s...)

	seedBytes := make([]frontend.Variable, 64)
	for b := 0; b < 64; {
		for y := 0; y < 5; y++ {
			for x := 0; x < 5; x++ {
				if x+5*y < (r/w) && (b < 64) {
					bits := api.ToBinary(s[5*x+y], 64)
					for i := 0; i < 8; i++ {
						seedBytes[b+i] = api.FromBinary(bits[i*8 : (i+1)*8]...)
					}
					b += 8
				}
			}
		}
	}
	return seedBytes
}

func pow2(n uint) *big.Int {
	return new(big.Int).Lsh(big.NewInt(1), n)
}

func fixedToU64(api frontend.API, a frontend.Variable) frontend.Variable {
	return api.FromBinary(api.ToBinary(a, 128)[:64]...)
}

func divFunc(curveID ecc.ID, inputs []*big.Int, outputs []*big.Int) error {
	a := inputs[0]
	b := inputs[1]
	outputs[0], outputs[1] = new(big.Int).QuoRem(a, b, new(big.Int))
	return nil
}

func div(api frontend.API, a, b frontend.Variable) frontend.Variable {
	outputs, _ := api.Compiler().NewHint(divFunc, 2, a, b)
	api.AssertIsEqual(a, api.Add(api.Mul(b, outputs[0]), outputs[1]))
	return outputs[0]
}

func keccakF(api frontend.API, a []frontend.Variable) []frontend.Variable {
	var b [25]frontend.Variable
	for i := 0; i < len(b); i++ {
		b[i] = 0
	}
	var c [5]frontend.Variable
	for i := 0; i < len(c); i++ {
		c[i] = 0
	}
	var d [5]frontend.Variable
	for i := 0; i < len(d); i++ {
		d[i] = 0
	}
	var rc [24]frontend.Variable
	rc[0], _ = new(big.Int).SetString("0000000000000001", 16)
	rc[1], _ = new(big.Int).SetString("0000000000008082", 16)
	rc[2], _ = new(big.Int).SetString("800000000000808A", 16)
	rc[3], _ = new(big.Int).SetString("8000000080008000", 16)
	rc[4], _ = new(big.Int).SetString("000000000000808B", 16)
	rc[5], _ = new(big.Int).SetString("0000000080000001", 16)
	rc[6], _ = new(big.Int).SetString("8000000080008081", 16)
	rc[7], _ = new(big.Int).SetString("8000000000008009", 16)
	rc[8], _ = new(big.Int).SetString("000000000000008A", 16)
	rc[9], _ = new(big.Int).SetString("0000000000000088", 16)
	rc[10], _ = new(big.Int).SetString("0000000080008009", 16)
	rc[11], _ = new(big.Int).SetString("000000008000000A", 16)
	rc[12], _ = new(big.Int).SetString("000000008000808B", 16)
	rc[13], _ = new(big.Int).SetString("800000000000008B", 16)
	rc[14], _ = new(big.Int).SetString("8000000000008089", 16)
	rc[15], _ = new(big.Int).SetString("8000000000008003", 16)
	rc[16], _ = new(big.Int).SetString("8000000000008002", 16)
	rc[17], _ = new(big.Int).SetString("8000000000000080", 16)
	rc[18], _ = new(big.Int).SetString("000000000000800A", 16)
	rc[19], _ = new(big.Int).SetString("800000008000000A", 16)
	rc[20], _ = new(big.Int).SetString("8000000080008081", 16)
	rc[21], _ = new(big.Int).SetString("8000000000008080", 16)
	rc[22], _ = new(big.Int).SetString("0000000080000001", 16)
	rc[23], _ = new(big.Int).SetString("8000000080008008", 16)

	mask := new(big.Int).Sub(pow2(64), big.NewInt(1))

	for i := 0; i < 24; i++ {
		c[0] = u64Xor(api, a[0], a[1], a[2], a[3], a[4])
		c[1] = u64Xor(api, a[5], a[6], a[7], a[8], a[9])
		c[2] = u64Xor(api, a[10], a[11], a[12], a[13], a[14])
		c[3] = u64Xor(api, a[15], a[16], a[17], a[18], a[19])
		c[4] = u64Xor(api, a[20], a[21], a[22], a[23], a[24])

		pow63 := pow2(63)

		//api.Println(c[1])
		//api.Println(api.Div(c[1], pow63))

		d[0] = u64Xor(api, c[4], u64Or(api, fixedToU64(api, api.Mul(c[1], 2)), div(api, c[1], pow63)))
		d[1] = u64Xor(api, c[0], u64Or(api, fixedToU64(api, api.Mul(c[2], 2)), div(api, c[2], pow63)))
		d[2] = u64Xor(api, c[1], u64Or(api, fixedToU64(api, api.Mul(c[3], 2)), div(api, c[3], pow63)))
		d[3] = u64Xor(api, c[2], u64Or(api, fixedToU64(api, api.Mul(c[4], 2)), div(api, c[4], pow63)))
		d[4] = u64Xor(api, c[3], u64Or(api, fixedToU64(api, api.Mul(c[0], 2)), div(api, c[0], pow63)))

		a[0] = u64Xor(api, a[0], d[0])
		a[1] = u64Xor(api, a[1], d[0])
		a[2] = u64Xor(api, a[2], d[0])
		a[3] = u64Xor(api, a[3], d[0])
		a[4] = u64Xor(api, a[4], d[0])

		a[5] = u64Xor(api, a[5], d[1])
		a[6] = u64Xor(api, a[6], d[1])
		a[7] = u64Xor(api, a[7], d[1])
		a[8] = u64Xor(api, a[8], d[1])
		a[9] = u64Xor(api, a[9], d[1])

		a[10] = u64Xor(api, a[10], d[2])
		a[11] = u64Xor(api, a[11], d[2])
		a[12] = u64Xor(api, a[12], d[2])
		a[13] = u64Xor(api, a[13], d[2])
		a[14] = u64Xor(api, a[14], d[2])

		a[15] = u64Xor(api, a[15], d[3])
		a[16] = u64Xor(api, a[16], d[3])
		a[17] = u64Xor(api, a[17], d[3])
		a[18] = u64Xor(api, a[18], d[3])
		a[19] = u64Xor(api, a[19], d[3])

		a[20] = u64Xor(api, a[20], d[4])
		a[21] = u64Xor(api, a[21], d[4])
		a[22] = u64Xor(api, a[22], d[4])
		a[23] = u64Xor(api, a[23], d[4])
		a[24] = u64Xor(api, a[24], d[4])

		/*Rho and pi steps*/
		b[0] = a[0]

		b[8] = u64Or(api, fixedToU64(api, api.Mul(a[1], pow2(36))), div(api, a[1], pow2(28)))
		b[11] = u64Or(api, fixedToU64(api, api.Mul(a[2], pow2(3))), div(api, a[2], pow2(61)))
		b[19] = u64Or(api, fixedToU64(api, api.Mul(a[3], pow2(41))), div(api, a[3], pow2(23)))
		b[22] = u64Or(api, fixedToU64(api, api.Mul(a[4], pow2(18))), div(api, a[4], pow2(46)))

		b[2] = u64Or(api, fixedToU64(api, api.Mul(a[5], pow2(1))), div(api, a[5], pow2(63)))
		b[5] = u64Or(api, fixedToU64(api, api.Mul(a[6], pow2(44))), div(api, a[6], pow2(20)))
		b[13] = u64Or(api, fixedToU64(api, api.Mul(a[7], pow2(10))), div(api, a[7], pow2(54)))
		b[16] = u64Or(api, fixedToU64(api, api.Mul(a[8], pow2(45))), div(api, a[8], pow2(19)))
		b[24] = u64Or(api, fixedToU64(api, api.Mul(a[9], pow2(2))), div(api, a[9], pow2(62)))

		b[4] = u64Or(api, fixedToU64(api, api.Mul(a[10], pow2(62))), div(api, a[10], pow2(2)))
		b[7] = u64Or(api, fixedToU64(api, api.Mul(a[11], pow2(6))), div(api, a[11], pow2(58)))
		b[10] = u64Or(api, fixedToU64(api, api.Mul(a[12], pow2(43))), div(api, a[12], pow2(21)))
		b[18] = u64Or(api, fixedToU64(api, api.Mul(a[13], pow2(15))), div(api, a[13], pow2(49)))
		b[21] = u64Or(api, fixedToU64(api, api.Mul(a[14], pow2(61))), div(api, a[14], pow2(3)))

		b[1] = u64Or(api, fixedToU64(api, api.Mul(a[15], pow2(28))), div(api, a[15], pow2(36)))
		b[9] = u64Or(api, fixedToU64(api, api.Mul(a[16], pow2(55))), div(api, a[16], pow2(9)))
		b[12] = u64Or(api, fixedToU64(api, api.Mul(a[17], pow2(25))), div(api, a[17], pow2(39)))
		b[15] = u64Or(api, fixedToU64(api, api.Mul(a[18], pow2(21))), div(api, a[18], pow2(43)))
		b[23] = u64Or(api, fixedToU64(api, api.Mul(a[19], pow2(56))), div(api, a[19], pow2(8)))

		b[3] = u64Or(api, fixedToU64(api, api.Mul(a[20], pow2(27))), div(api, a[20], pow2(37)))
		b[6] = u64Or(api, fixedToU64(api, api.Mul(a[21], pow2(20))), div(api, a[21], pow2(44)))
		b[14] = u64Or(api, fixedToU64(api, api.Mul(a[22], pow2(39))), div(api, a[22], pow2(25)))
		b[17] = u64Or(api, fixedToU64(api, api.Mul(a[23], pow2(8))), div(api, a[23], pow2(56)))
		b[20] = u64Or(api, fixedToU64(api, api.Mul(a[24], pow2(14))), div(api, a[24], pow2(50)))

		/*Xi state*/

		a[0] = u64Xor(api, b[0], u64And(api, u64Xor(api, b[5], mask), b[10]))
		a[1] = u64Xor(api, b[1], u64And(api, u64Xor(api, b[6], mask), b[11]))
		a[2] = u64Xor(api, b[2], u64And(api, u64Xor(api, b[7], mask), b[12]))
		a[3] = u64Xor(api, b[3], u64And(api, u64Xor(api, b[8], mask), b[13]))
		a[4] = u64Xor(api, b[4], u64And(api, u64Xor(api, b[9], mask), b[14]))

		a[5] = u64Xor(api, b[5], u64And(api, u64Xor(api, b[10], mask), b[15]))
		a[6] = u64Xor(api, b[6], u64And(api, u64Xor(api, b[11], mask), b[16]))
		a[7] = u64Xor(api, b[7], u64And(api, u64Xor(api, b[12], mask), b[17]))
		a[8] = u64Xor(api, b[8], u64And(api, u64Xor(api, b[13], mask), b[18]))
		a[9] = u64Xor(api, b[9], u64And(api, u64Xor(api, b[14], mask), b[19]))

		a[10] = u64Xor(api, b[10], u64And(api, u64Xor(api, b[15], mask), b[20]))
		a[11] = u64Xor(api, b[11], u64And(api, u64Xor(api, b[16], mask), b[21]))
		a[12] = u64Xor(api, b[12], u64And(api, u64Xor(api, b[17], mask), b[22]))
		a[13] = u64Xor(api, b[13], u64And(api, u64Xor(api, b[18], mask), b[23]))
		a[14] = u64Xor(api, b[14], u64And(api, u64Xor(api, b[19], mask), b[24]))

		a[15] = u64Xor(api, b[15], u64And(api, u64Xor(api, b[20], mask), b[0]))
		a[16] = u64Xor(api, b[16], u64And(api, u64Xor(api, b[21], mask), b[1]))
		a[17] = u64Xor(api, b[17], u64And(api, u64Xor(api, b[22], mask), b[2]))
		a[18] = u64Xor(api, b[18], u64And(api, u64Xor(api, b[23], mask), b[3]))
		a[19] = u64Xor(api, b[19], u64And(api, u64Xor(api, b[24], mask), b[4]))

		a[20] = u64Xor(api, b[20], u64And(api, u64Xor(api, b[0], mask), b[5]))
		a[21] = u64Xor(api, b[21], u64And(api, u64Xor(api, b[1], mask), b[6]))
		a[22] = u64Xor(api, b[22], u64And(api, u64Xor(api, b[2], mask), b[7]))
		a[23] = u64Xor(api, b[23], u64And(api, u64Xor(api, b[3], mask), b[8]))
		a[24] = u64Xor(api, b[24], u64And(api, u64Xor(api, b[4], mask), b[9]))

		///*Last step*/

		a[0] = u64Xor(api, a[0], rc[i])
	}

	return a
}
