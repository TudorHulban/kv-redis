package redis

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type tstruct struct {
	F1 int
	F2 []byte
}

func TestGobers(t *testing.T) {
	v := tstruct{
		F1: 1,
		F2: []byte("a"),
	}

	// testing encoder
	resEnc, errEncode := Encoder(v)
	assert.NoError(t, errEncode)
	assert.NotNil(t, resEnc)

	// testing decoder
	var p tstruct

	errDecode := Decoder(resEnc, &p)
	assert.NoError(t, errDecode)
	assert.Equal(t, v, p)
}
