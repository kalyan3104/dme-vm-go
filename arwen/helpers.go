package arwen

import (
	"fmt"
	"io/ioutil"
	"math/big"
	"path/filepath"
	"time"
	"unsafe"

	logger "github.com/kalyan3104/dme-logger-go"
)

var logDuration = logger.GetOrCreate("arwen/duration")

// Zero is the big integer 0
var Zero = big.NewInt(0)

func CustomStorageKey(keyType string, associatedKey []byte) []byte {
	return append(associatedKey, []byte(keyType)...)
}

func BooleanToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}

func GuardedMakeByteSlice2D(length int32) ([][]byte, error) {
	if length < 0 {
		return nil, fmt.Errorf("GuardedMakeByteSlice2D: negative length (%d)", length)
	}

	result := make([][]byte, length)
	return result, nil
}

func GuardedGetBytesSlice(data []byte, offset int32, length int32) ([]byte, error) {
	dataLength := uint32(len(data))
	isOffsetTooSmall := offset < 0
	isOffsetTooLarge := uint32(offset) > dataLength
	requestedEnd := uint32(offset + length)
	isRequestedEndTooLarge := requestedEnd > dataLength
	isLengthNegative := length < 0

	if isOffsetTooSmall || isOffsetTooLarge || isRequestedEndTooLarge {
		return nil, fmt.Errorf("GuardedGetBytesSlice: bad bounds")
	}

	if isLengthNegative {
		return nil, fmt.Errorf("GuardedGetBytesSlice: negative length")
	}

	result := data[offset : offset+length]
	return result, nil
}

func PadBytesLeft(data []byte, size int) []byte {
	if data == nil {
		return nil
	}
	if len(data) == 0 {
		return []byte{}
	}
	padSize := size - len(data)
	if padSize <= 0 {
		return data
	}

	paddedBytes := make([]byte, padSize)
	paddedBytes = append(paddedBytes, data...)
	return paddedBytes
}

func InverseBytes(data []byte) []byte {
	length := len(data)
	invBytes := make([]byte, length)
	for i := 0; i < length; i++ {
		invBytes[length-i-1] = data[i]
	}
	return invBytes
}

func WithFault(err error, context unsafe.Pointer, failExecution bool) bool {
	if err == nil {
		return false
	}

	if failExecution {
		runtime := GetRuntimeContext(context)
		metering := GetMeteringContext(context)

		metering.UseGas(metering.GasLeft())
		runtime.FailExecution(err)
	}

	return true
}

func GetSCCode(fileName string) []byte {
	code, _ := ioutil.ReadFile(filepath.Clean(fileName))
	return code
}

func TimeTrack(start time.Time, message string) {
	elapsed := time.Since(start)
	logDuration.Trace(message, "duration", elapsed)
}

func U64MulToBigInt(x, y uint64) *big.Int {
	bx := big.NewInt(0).SetUint64(x)
	by := big.NewInt(0).SetUint64(y)

	return big.NewInt(0).Mul(bx, by)
}

// IfNil tests if the provided interface pointer or underlying object is nil
func IfNil(checker nilInterfaceChecker) bool {
	if checker == nil {
		return true
	}
	return checker.IsInterfaceNil()
}

type nilInterfaceChecker interface {
	IsInterfaceNil() bool
}
