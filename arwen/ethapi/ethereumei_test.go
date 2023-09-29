package ethapi

import (
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_convertToEthAddress(t *testing.T) {
	kalyan3104Address, _ := hex.DecodeString("aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa")
	expectedResult, _ := hex.DecodeString("aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa")

	result := convertToEthAddress(kalyan3104Address)

	assert.Equal(t, expectedResult, result)
}

func Test_convertToEthU128(t *testing.T) {
	data, _ := hex.DecodeString("aa")
	expectedResult, _ := hex.DecodeString("000000000000000000000000000000aa")

	result := convertToEthU128(data)
	assert.Equal(t, expectedResult, result)
}

func Test_convertToEthU128_whenLargeData(t *testing.T) {
	data, _ := hex.DecodeString("aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa")
	expectedResult, _ := hex.DecodeString("00000000000000000000000000000000")

	result := convertToEthU128(data)
	assert.Equal(t, expectedResult, result)
}
