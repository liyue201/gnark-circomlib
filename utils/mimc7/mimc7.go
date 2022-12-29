package mimc7

import (
	"github.com/iden3/go-iden3-crypto/mimc7"
	"hash"
	"math/big"
)

type mimc7Hash struct {
	status *big.Int
	data   []byte // data to hash
}

func NewMimc7() hash.Hash {
	return &mimc7Hash{
		status: big.NewInt(0),
	}
}

func (d *mimc7Hash) Write(p []byte) (n int, err error) {
	d.data = append(d.data, p...)
	return len(p), nil
}

func (d *mimc7Hash) Sum(b []byte) []byte {

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
	s, err := mimc7.Hash(arr, d.status)
	if err != nil {
		panic(err)
	}
	d.status = s
	d.data = nil

	hash := d.status.Bytes()
	b = append(b, hash[:]...)

	return b
}

func (d *mimc7Hash) Size() int {
	return 32
}

func (d *mimc7Hash) BlockSize() int {
	return 32
}

func (d *mimc7Hash) Reset() {
	d.status = big.NewInt(0)
	d.data = nil
}
