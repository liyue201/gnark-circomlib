package circuits

import (
	"fmt"
	"github.com/consensys/gnark/frontend/cs/r1cs"
	"testing"

	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark/backend"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/test"
)

type circuitNum2Bits struct {
	A    frontend.Variable    `gnark:",public"`
	Bits [8]frontend.Variable `gnark:",public"`
}

func (t *circuitNum2Bits) Define(api frontend.API) error {
	bits := Num2Bits(api, t.A, 8)
	for i := 0; i < 8; i++ {
		api.AssertIsEqual(bits[i], t.Bits[i])
	}
	return nil
}

func Test_Num2Bits(t *testing.T) {
	assert := test.NewAssert(t)
	var circuit circuitNum2Bits

	assert.ProverSucceeded(&circuit, &circuitNum2Bits{
		A:    9,
		Bits: [8]frontend.Variable{1, 0, 0, 1, 0, 0, 0, 0},
	}, test.WithCurves(ecc.BN254), test.WithBackends(backend.GROTH16))

	_r1cs, _ := frontend.Compile(ecc.BN254, r1cs.NewBuilder, &circuit)
	internal, secret, public := _r1cs.GetNbVariables()
	fmt.Printf("public, secret, internal %v, %v, %v\n", public, secret, internal)
}

type circuitNum2BitsStrict struct {
	A    frontend.Variable      `gnark:",public"`
	Bits [254]frontend.Variable `gnark:",public"`
}

func (t *circuitNum2BitsStrict) Define(api frontend.API) error {
	bits := Num2BitsStrict(api, t.A, 254)
	for i := 0; i < len(t.Bits); i++ {
		api.AssertIsEqual(bits[i], t.Bits[i])
	}
	return nil
}

func Test_Num2BitsStrict(t *testing.T) {

	assert := test.NewAssert(t)
	var circuit circuitNum2BitsStrict

	validAssignment := &circuitNum2BitsStrict{
		A:    9,
		Bits: [254]frontend.Variable{},
	}
	for i := 0; i < 254; i++ {
		validAssignment.Bits[i] = 0
	}
	validAssignment.Bits[0] = 1
	validAssignment.Bits[3] = 1

	assert.ProverSucceeded(&circuit, validAssignment, test.WithCurves(ecc.BN254), test.WithBackends(backend.GROTH16))

	_r1cs, _ := frontend.Compile(ecc.BN254, r1cs.NewBuilder, &circuit)
	internal, secret, public := _r1cs.GetNbVariables()
	fmt.Printf("public, secret, internal %v, %v, %v\n", public, secret, internal)
}

type circuitBits2Num struct {
	A    frontend.Variable    `gnark:",public"`
	Bits [4]frontend.Variable `gnark:",public"`
}

func (t *circuitBits2Num) Define(api frontend.API) error {
	num := Bits2Num(api, t.Bits[:])
	api.AssertIsEqual(num, t.A)
	return nil
}

func Test_Bits2Num(t *testing.T) {
	assert := test.NewAssert(t)
	var circuit circuitBits2Num

	assert.ProverSucceeded(&circuit, &circuitBits2Num{
		A:    9,
		Bits: [4]frontend.Variable{1, 0, 0, 1},
	}, test.WithCurves(ecc.BN254), test.WithBackends(backend.GROTH16))

	_r1cs, _ := frontend.Compile(ecc.BN254, r1cs.NewBuilder, &circuit)
	internal, secret, public := _r1cs.GetNbVariables()
	fmt.Printf("public, secret, internal %v, %v, %v\n", public, secret, internal)
}

type circuitBits2NumStrict struct {
	A    frontend.Variable      `gnark:",public"`
	Bits [254]frontend.Variable `gnark:",public"`
}

func (t *circuitBits2NumStrict) Define(api frontend.API) error {
	num := Bits2NumStrict(api, t.Bits[:])
	api.AssertIsEqual(num, t.A)
	return nil
}

func Test_Bits2NumStrict(t *testing.T) {
	assert := test.NewAssert(t)
	var circuit circuitBits2NumStrict

	validAssignment := &circuitNum2BitsStrict{
		A:    9,
		Bits: [254]frontend.Variable{},
	}
	for i := 0; i < 254; i++ {
		validAssignment.Bits[i] = 0
	}
	validAssignment.Bits[0] = 1
	validAssignment.Bits[3] = 1

	assert.ProverSucceeded(&circuit, validAssignment, test.WithCurves(ecc.BN254), test.WithBackends(backend.GROTH16))

	_r1cs, _ := frontend.Compile(ecc.BN254, r1cs.NewBuilder, &circuit)
	internal, secret, public := _r1cs.GetNbVariables()
	fmt.Printf("public, secret, internal %v, %v, %v\n", public, secret, internal)
}

type circuitNum2BitsNeg struct {
	A    frontend.Variable    `gnark:",public"`
	Bits [8]frontend.Variable `gnark:",public"`
}

func (t *circuitNum2BitsNeg) Define(api frontend.API) error {
	bits := Num2BitsNeg(api, t.A, 8)
	for i := 0; i < len(t.Bits); i++ {
		api.AssertIsEqual(bits[i], t.Bits[i])
	}
	return nil
}

func Test_Num2BitsNeg(t *testing.T) {
	assert := test.NewAssert(t)
	var circuit circuitNum2BitsNeg

	assert.ProverSucceeded(&circuit, &circuitNum2BitsNeg{
		A:    26, //-26
		Bits: [8]frontend.Variable{0, 1, 1, 0, 0, 1, 1, 1},
	}, test.WithCurves(ecc.BN254), test.WithBackends(backend.GROTH16))

	_r1cs, _ := frontend.Compile(ecc.BN254, r1cs.NewBuilder, &circuit)
	internal, secret, public := _r1cs.GetNbVariables()
	fmt.Printf("public, secret, internal %v, %v, %v\n", public, secret, internal)
}
