package config

import (
	"fmt"
	"testing"

	"github.com/mitchellh/mapstructure"
	"github.com/stretchr/testify/assert"
)

type operations struct {
	OperationA uint64
	OperationB uint64
	OperationC uint64
	OperationD uint64
	OperationE uint64
}

func TestDecode(t *testing.T) {
	gasMap := make(map[string]uint64)
	gasMap["OperationB"] = 4
	gasMap["OperationA"] = 3
	gasMap["OperationC"] = 100
	gasMap["OperationD"] = 1000

	op := &operations{}
	err := mapstructure.Decode(gasMap, op)
	assert.Nil(t, err)

	fmt.Printf("%+v\n", op)
}

func TestDecode_ArwenGas(t *testing.T) {
	gasMap := make(map[string]uint64)
	gasMap["StorePerByte"] = 4
	gasMap["GetSCAddress"] = 4
	gasMap["GetExternalBalance"] = 4
	gasMap["BigIntByteLength"] = 4

	bigIntOp := &BigIntAPICost{}
	err := mapstructure.Decode(gasMap, bigIntOp)
	assert.Nil(t, err)

	fmt.Printf("%+v\n", bigIntOp)

	moaOp := &Kalyan3104APICost{}
	err = mapstructure.Decode(gasMap, moaOp)
	assert.Nil(t, err)

	fmt.Printf("%+v\n", moaOp)

	ethOp := &EthAPICost{}
	err = mapstructure.Decode(gasMap, ethOp)
	assert.Nil(t, err)

	fmt.Printf("%+v\n", ethOp)
}

func TestDecode_ZeroGasCostError(t *testing.T) {
	gasMap := FillGasMap_WASMOpcodeValues(1)

	wasmCosts := &WASMOpcodeCost{}
	err := mapstructure.Decode(gasMap, wasmCosts)
	assert.Nil(t, err)

	err = checkForZeroUint64Fields(*wasmCosts)
	assert.Nil(t, err)

	gasMap["BrIf"] = 0
	wasmCosts = &WASMOpcodeCost{}
	err = mapstructure.Decode(gasMap, wasmCosts)
	assert.Nil(t, err)

	err = checkForZeroUint64Fields(*wasmCosts)
	assert.Error(t, err)
}
