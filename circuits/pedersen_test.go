package circuits

import (
	"fmt"
	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark/backend"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/test"
	"github.com/minya-konka/pedersen"
	"math/rand"
	"testing"
)

type circuitPedersen struct {
	MsgBytes [77]frontend.Variable
	X        frontend.Variable
	Y        frontend.Variable
}

func (t *circuitPedersen) Define(api frontend.API) error {

	var bits []frontend.Variable
	for i := 0; i < len(t.MsgBytes); i++ {
		bits = append(bits, api.ToBinary(t.MsgBytes[i], 8)...)
	}
	out := Pedersen(api, bits)

	api.AssertIsEqual(out[0], t.X)
	api.AssertIsEqual(out[1], t.Y)

	return nil
}

func TestPedersen(t *testing.T) {
	assert := test.NewAssert(t)
	var circuit circuitPedersen
	var assign circuitPedersen

	msg := make([]byte, len(circuit.MsgBytes))
	for i := 0; i < len(circuit.MsgBytes); i++ {
		msg[i] = byte(rand.Int31n(256))
		assign.MsgBytes[i] = msg[i]
	}
	p := pedersen.Hash(msg)

	fmt.Printf("%s\n", p.X)
	fmt.Printf("%s\n", p.Y)
	assign.X = p.X
	assign.Y = p.Y

	assert.ProverSucceeded(&circuit, &assign, test.WithCurves(ecc.BN254), test.WithBackends(backend.GROTH16))
}
