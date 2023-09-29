package contexts

import (
	"fmt"
	"unsafe"

	vmcommon "github.com/kalyan3104/dme-vm-common"
	"github.com/kalyan3104/dme-vm-go/arwen"
	"github.com/kalyan3104/dme-vm-go/wasmer"
)

var _ arwen.RuntimeContext = (*runtimeContext)(nil)

type runtimeContext struct {
	host     arwen.VMHost
	instance *wasmer.Instance
	// Temporarily holding these pointers are supposed to circumvent an undesired deallocation performed by Go's GC
	instanceContextDataPointers []*int
	instanceContext             *wasmer.InstanceContext
	vmInput                     *vmcommon.VMInput
	scAddress                   []byte
	callFunction                string
	vmType                      []byte
	readOnly                    bool

	stateStack    []*runtimeContext
	instanceStack []*wasmer.Instance

	maxWasmerInstances uint64

	asyncCallInfo    *arwen.AsyncCallInfo
	asyncContextInfo *arwen.AsyncContextInfo

	validator *WASMValidator
}

// NewRuntimeContext creates a new runtimeContext
func NewRuntimeContext(host arwen.VMHost, vmType []byte) (*runtimeContext, error) {
	scAPINames := host.GetAPIMethods().Names()
	protocolBuiltinFunctions := host.GetProtocolBuiltinFunctions()

	context := &runtimeContext{
		host:                        host,
		instanceContextDataPointers: make([]*int, 0),
		vmType:                      vmType,
		stateStack:                  make([]*runtimeContext, 0),
		instanceStack:               make([]*wasmer.Instance, 0),
		validator:                   NewWASMValidator(scAPINames, protocolBuiltinFunctions),
	}

	context.InitState()

	return context, nil
}

func (context *runtimeContext) InitState() {
	context.instanceContextDataPointers = make([]*int, 0)
	context.vmInput = &vmcommon.VMInput{}
	context.scAddress = make([]byte, 0)
	context.callFunction = ""
	context.readOnly = false
	context.asyncCallInfo = nil
	context.asyncContextInfo = &arwen.AsyncContextInfo{
		AsyncContextMap: make(map[string]*arwen.AsyncContext),
	}
}

func (context *runtimeContext) StartWasmerInstance(contract []byte, gasLimit uint64) error {
	if context.RunningInstancesCount() >= context.maxWasmerInstances {
		context.instance = nil
		return arwen.ErrMaxInstancesReached
	}
	options := wasmer.CompilationOptions{
		GasLimit:           gasLimit,
		OpcodeTrace:        false,
		Metering:           true,
		RuntimeBreakpoints: true,
	}
	newInstance, err := wasmer.NewInstanceWithOptions(contract, options)
	if err != nil {
		context.instance = nil
		return err
	}

	context.instance = newInstance
	context.SetRuntimeBreakpointValue(arwen.BreakpointNone)
	return nil
}

func (context *runtimeContext) SetMaxInstanceCount(maxInstances uint64) {
	context.maxWasmerInstances = maxInstances
}

func (context *runtimeContext) InitStateFromContractCallInput(input *vmcommon.ContractCallInput) {
	context.vmInput = &input.VMInput
	context.scAddress = input.RecipientAddr
	context.callFunction = input.Function
	// Reset async map for initial state
	context.asyncContextInfo = &arwen.AsyncContextInfo{
		CallerAddr:      input.CallerAddr,
		AsyncContextMap: make(map[string]*arwen.AsyncContext),
	}
}

func (context *runtimeContext) SetCustomCallFunction(callFunction string) {
	context.callFunction = callFunction
}

func (context *runtimeContext) PushState() {
	newState := &runtimeContext{
		vmInput:          context.vmInput,
		scAddress:        context.scAddress,
		callFunction:     context.callFunction,
		readOnly:         context.readOnly,
		asyncCallInfo:    context.asyncCallInfo,
		asyncContextInfo: context.asyncContextInfo,
	}

	context.stateStack = append(context.stateStack, newState)
}

func (context *runtimeContext) PopSetActiveState() {
	stateStackLen := len(context.stateStack)
	prevState := context.stateStack[stateStackLen-1]
	context.stateStack = context.stateStack[:stateStackLen-1]

	context.vmInput = prevState.vmInput
	context.scAddress = prevState.scAddress
	context.callFunction = prevState.callFunction
	context.readOnly = prevState.readOnly
	context.asyncCallInfo = prevState.asyncCallInfo
	context.asyncContextInfo = prevState.asyncContextInfo
}

func (context *runtimeContext) PopDiscard() {
	stateStackLen := len(context.stateStack)
	context.stateStack = context.stateStack[:stateStackLen-1]
}

func (context *runtimeContext) ClearStateStack() {
	context.stateStack = make([]*runtimeContext, 0)
}

func (context *runtimeContext) PushInstance() {
	context.instanceStack = append(context.instanceStack, context.instance)
}

func (context *runtimeContext) PopInstance() {
	instanceStackLen := len(context.instanceStack)
	prevInstance := context.instanceStack[instanceStackLen-1]
	context.instanceStack = context.instanceStack[:instanceStackLen-1]

	context.CleanInstance()
	context.instance = prevInstance
}

func (context *runtimeContext) RunningInstancesCount() uint64 {
	return uint64(len(context.instanceStack))
}

func (context *runtimeContext) ClearInstanceStack() {
	for _, instance := range context.instanceStack {
		instance.Clean()
	}
	context.instanceStack = make([]*wasmer.Instance, 0)
}

func (context *runtimeContext) GetVMType() []byte {
	return context.vmType
}

func (context *runtimeContext) GetVMInput() *vmcommon.VMInput {
	return context.vmInput
}

func (context *runtimeContext) SetVMInput(vmInput *vmcommon.VMInput) {
	context.vmInput = vmInput
}

func (context *runtimeContext) GetSCAddress() []byte {
	return context.scAddress
}

func (context *runtimeContext) SetSCAddress(scAddress []byte) {
	context.scAddress = scAddress
}

func (context *runtimeContext) GetCurrentTxHash() []byte {
	return context.vmInput.CurrentTxHash
}

func (context *runtimeContext) GetOriginalTxHash() []byte {
	return context.vmInput.OriginalTxHash
}

func (context *runtimeContext) Function() string {
	return context.callFunction
}

func (context *runtimeContext) Arguments() [][]byte {
	return context.vmInput.Arguments
}

func (context *runtimeContext) ExtractCodeUpgradeFromArgs() ([]byte, []byte, error) {
	const numMinUpgradeArguments = 2

	arguments := context.vmInput.Arguments
	if len(arguments) < numMinUpgradeArguments {
		return nil, nil, arwen.ErrInvalidUpgradeArguments
	}

	code := arguments[0]
	codeMetadata := arguments[1]
	context.vmInput.Arguments = context.vmInput.Arguments[numMinUpgradeArguments:]
	return code, codeMetadata, nil
}

func (context *runtimeContext) FailExecution(err error) {
	context.host.Output().SetReturnCode(vmcommon.ExecutionFailed)

	var message string
	if err != nil {
		message = err.Error()
	} else {
		message = "execution failed"
	}

	context.host.Output().SetReturnMessage(message)
	context.SetRuntimeBreakpointValue(arwen.BreakpointExecutionFailed)
}

func (context *runtimeContext) SignalUserError(message string) {
	context.host.Output().SetReturnCode(vmcommon.UserError)
	context.host.Output().SetReturnMessage(message)
	context.SetRuntimeBreakpointValue(arwen.BreakpointSignalError)
}

func (context *runtimeContext) SetRuntimeBreakpointValue(value arwen.BreakpointValue) {
	context.instance.SetBreakpointValue(uint64(value))
}

func (context *runtimeContext) GetRuntimeBreakpointValue() arwen.BreakpointValue {
	return arwen.BreakpointValue(context.instance.GetBreakpointValue())
}

func (context *runtimeContext) VerifyContractCode() error {
	err := context.validator.verifyMemoryDeclaration(context.instance)
	if err != nil {
		return err
	}

	err = context.validator.verifyFunctions(context.instance)
	if err != nil {
		return err
	}

	return nil
}

func (context *runtimeContext) Kalyan3104APIErrorShouldFailExecution() bool {
	return true
}

func (context *runtimeContext) BigIntAPIErrorShouldFailExecution() bool {
	return true
}

func (context *runtimeContext) CryptoAPIErrorShouldFailExecution() bool {
	return true
}

func (context *runtimeContext) GetPointsUsed() uint64 {
	return context.instance.GetPointsUsed()
}

func (context *runtimeContext) SetPointsUsed(gasPoints uint64) {
	context.instance.SetPointsUsed(gasPoints)
}

func (context *runtimeContext) ReadOnly() bool {
	return context.readOnly
}

func (context *runtimeContext) SetReadOnly(readOnly bool) {
	context.readOnly = readOnly
}

func (context *runtimeContext) SetInstanceContextID(id int) {
	context.instanceContextDataPointers = append(context.instanceContextDataPointers, &id)
	context.instance.SetContextData(unsafe.Pointer(&id))
}

func (context *runtimeContext) SetInstanceContext(instCtx *wasmer.InstanceContext) {
	context.instanceContext = instCtx
}

func (context *runtimeContext) GetInstanceContext() *wasmer.InstanceContext {
	return context.instanceContext
}

func (context *runtimeContext) GetInstanceExports() wasmer.ExportsMap {
	return context.instance.Exports
}

func (context *runtimeContext) CleanInstance() {
	if context.instance == nil {
		return
	}
	context.instance.Clean()
	context.instance = nil
}

func (context *runtimeContext) GetFunctionToCall() (wasmer.ExportedFunctionCallback, error) {
	exports := context.instance.Exports
	if function, ok := exports[context.callFunction]; ok {
		return function, nil
	}

	if function, ok := exports["main"]; ok {
		return function, nil
	}

	return nil, arwen.ErrFuncNotFound
}

func (context *runtimeContext) GetInitFunction() wasmer.ExportedFunctionCallback {
	exports := context.instance.Exports

	if init, ok := exports[arwen.InitFunctionName]; ok {
		return init
	}

	if init, ok := exports[arwen.InitFunctionNameEth]; ok {
		return init
	}

	return nil
}

func (context *runtimeContext) SetAsyncCallInfo(asyncCallInfo *arwen.AsyncCallInfo) {
	context.asyncCallInfo = asyncCallInfo
}

func (context *runtimeContext) AddAsyncContextCall(contextIdentifier []byte, asyncCall *arwen.AsyncGeneratedCall) error {
	_, ok := context.asyncContextInfo.AsyncContextMap[string(contextIdentifier)]
	currentContextMap := context.asyncContextInfo.AsyncContextMap
	if !ok {
		currentContextMap[string(contextIdentifier)] = &arwen.AsyncContext{
			AsyncCalls: make([]*arwen.AsyncGeneratedCall, 0),
		}
	}

	currentContextMap[string(contextIdentifier)].AsyncCalls =
		append(currentContextMap[string(contextIdentifier)].AsyncCalls, asyncCall)

	return nil
}

func (context *runtimeContext) GetAsyncContextInfo() *arwen.AsyncContextInfo {
	return context.asyncContextInfo
}

func (context *runtimeContext) GetAsyncContext(contextIdentifier []byte) (*arwen.AsyncContext, error) {
	asyncContext, ok := context.asyncContextInfo.AsyncContextMap[string(contextIdentifier)]
	if !ok {
		return nil, arwen.ErrAsyncContextDoesNotExist
	}

	return asyncContext, nil
}

func (context *runtimeContext) GetAsyncCallInfo() *arwen.AsyncCallInfo {
	return context.asyncCallInfo
}

func (context *runtimeContext) MemLoad(offset int32, length int32) ([]byte, error) {
	if length == 0 {
		return []byte{}, nil
	}

	memory := context.instanceContext.Memory()
	memoryView := memory.Data()
	memoryLength := memory.Length()
	requestedEnd := uint32(offset + length)
	isOffsetTooSmall := offset < 0
	isOffsetTooLarge := uint32(offset) > memoryLength
	isRequestedEndTooLarge := requestedEnd > memoryLength
	isLengthNegative := length < 0

	if isOffsetTooSmall || isOffsetTooLarge {
		return nil, fmt.Errorf("mem load: %w", arwen.ErrBadBounds)
	}
	if isLengthNegative {
		return nil, fmt.Errorf("mem load: %w", arwen.ErrNegativeLength)
	}

	result := make([]byte, length)
	if isRequestedEndTooLarge {
		copy(result, memoryView[offset:])
	} else {
		copy(result, memoryView[offset:requestedEnd])
	}

	return result, nil
}

func (context *runtimeContext) MemStore(offset int32, data []byte) error {
	dataLength := int32(len(data))
	if dataLength == 0 {
		return nil
	}

	memory := context.instanceContext.Memory()
	memoryView := memory.Data()
	memoryLength := memory.Length()
	requestedEnd := uint32(offset + dataLength)
	isOffsetTooSmall := offset < 0
	isNewPageNecessary := requestedEnd > memoryLength

	if isOffsetTooSmall {
		return arwen.ErrBadLowerBounds
	}
	if isNewPageNecessary {
		err := memory.Grow(1)
		if err != nil {
			return err
		}

		memoryView = memory.Data()
		memoryLength = memory.Length()
	}

	isRequestedEndTooLarge := requestedEnd > memoryLength
	if isRequestedEndTooLarge {
		return arwen.ErrBadUpperBounds
	}

	copy(memoryView[offset:requestedEnd], data)
	return nil
}
