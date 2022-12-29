package poseidon

import (
	"github.com/iden3/go-iden3-crypto/poseidon"
	"hash"
	"math/big"
)

type poseidonHash struct {
	status *big.Int
	data   []byte // data to hash
}

func NewPoseidon() hash.Hash {
	return &poseidonHash{
		status: big.NewInt(0),
	}
}

func (d *poseidonHash) Write(p []byte) (n int, err error) {
	d.data = append(d.data, p...)
	return len(p), nil
}

func (d *poseidonHash) Sum(b []byte) []byte {

	blockSize := d.BlockSize()
	if len(d.data)%blockSize != 0 {
		q := len(d.data) / blockSize
		r := len(d.data) % blockSize
		sliceq := make([]byte, q*blockSize)
		copy(sliceq, d.data)
		slicer := make([]byte, r)
		copy(slicer, d.data[q*blockSize:])
		sliceremainder := make([]byte, blockSize-r)
		d.data = append(sliceq, sliceremainder...)
		d.data = append(d.data, slicer...)
	}
	nbChunks := len(d.data) / blockSize

	arr := make([]*big.Int, nbChunks)
	for i := 0; i < nbChunks; i++ {
		arr[i] = new(big.Int).SetBytes(d.data[i*blockSize : i*blockSize+blockSize])
	}
	s, err := poseidon.Hash(arr)
	if err != nil {
		panic(err)
	}
	d.status = s
	d.data = nil

	hash := d.status.Bytes()
	b = append(b, hash[:]...)

	return b
}

func (d *poseidonHash) Size() int {
	return 32
}

func (d *poseidonHash) BlockSize() int {
	return 32
}

func (d *poseidonHash) Reset() {
	d.status = big.NewInt(0)
	d.data = nil
}
