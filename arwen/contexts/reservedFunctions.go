package contexts

import (
	vmcommon "github.com/kalyan3104/dme-vm-common"
	"github.com/kalyan3104/dme-vm-go/arwen"
)

// ReservedFunctions holds the reserved function names
type ReservedFunctions struct {
	functionNames vmcommon.FunctionNames
}

// NewReservedFunctions creates a new ReservedFunctions
func NewReservedFunctions(scAPINames vmcommon.FunctionNames, protocolBuiltinFunctions vmcommon.FunctionNames) *ReservedFunctions {
	result := &ReservedFunctions{
		functionNames: make(vmcommon.FunctionNames),
	}

	for name, value := range protocolBuiltinFunctions {
		result.functionNames[name] = value
	}

	for name, value := range scAPINames {
		result.functionNames[name] = value
	}

	var empty struct{}
	result.functionNames[arwen.UpgradeFunctionName] = empty

	return result
}

// IsReserved returns whether a function is reserved
func (reservedFunctions *ReservedFunctions) IsReserved(functionName string) bool {
	if _, ok := reservedFunctions.functionNames[functionName]; ok {
		return true
	}

	return false
}

// GetReserved gets the reserved functions as a slice of strings
func (reservedFunctions *ReservedFunctions) GetReserved() []string {
	keys := make([]string, len(reservedFunctions.functionNames))

	i := 0
	for key := range reservedFunctions.functionNames {
		keys[i] = key
		i++
	}

	return keys
}
