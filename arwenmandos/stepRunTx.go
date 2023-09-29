package arwenmandos

import (
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"

	vmcommon "github.com/kalyan3104/dme-vm-common"
	vmi "github.com/kalyan3104/dme-vm-common"
	mj "github.com/kalyan3104/dme-vm-util/test-util/mandos/json/model"
)

func (ae *ArwenTestExecutor) executeTx(txIndex string, tx *mj.Transaction) (*vmi.VMOutput, error) {
	if tx.Type.HasSender() {
		beforeErr := ae.World.UpdateWorldStateBefore(
			tx.From.Value,
			tx.GasLimit.Value,
			tx.GasPrice.Value)
		if beforeErr != nil {
			return nil, fmt.Errorf("Could not set up tx %s: %w", txIndex, beforeErr)
		}
	}

	// we also use fake vm outputs for transactions that don't use the VM, just for convenience
	var output *vmi.VMOutput

	if !ae.senderHasEnoughBalance(tx) {
		// out of funds is handled by the protocol, so it needs to be mocked here
		output = outOfFundsResult()
	} else {
		switch tx.Type {
		case mj.ScDeploy:
			var err error
			output, err = ae.scCreate(tx)
			if err != nil {
				return nil, err
			}
		case mj.ScCall:
			var err error
			output, err = ae.scCall(tx)
			if err != nil {
				return nil, err
			}
		case mj.Transfer:
			var err error
			output, err = ae.simpleTransferOutput(tx)
			if err != nil {
				return nil, err
			}
		case mj.ValidatorReward:
			var err error
			output, err = ae.validatorRewardOutput(tx)
			if err != nil {
				return nil, err
			}
		default:
			return nil, errors.New("unknown transaction type")
		}

	}

	if output.ReturnCode == vmi.Ok {
		err := ae.updateStateAfterTx(tx, output)
		if err != nil {
			return nil, err
		}
	}

	return output, nil
}

func (ae *ArwenTestExecutor) senderHasEnoughBalance(tx *mj.Transaction) bool {
	if !tx.Type.HasSender() {
		return true
	}
	sender := ae.World.AcctMap.GetAccount(tx.From.Value)
	return sender.Balance.Cmp(tx.Value.Value) >= 0
}

func (ae *ArwenTestExecutor) simpleTransferOutput(tx *mj.Transaction) (*vmi.VMOutput, error) {
	outputAccounts := make(map[string]*vmcommon.OutputAccount)
	outputAccounts[string(tx.To.Value)] = &vmcommon.OutputAccount{
		Address:      tx.To.Value,
		BalanceDelta: tx.Value.Value,
	}

	return &vmcommon.VMOutput{
		ReturnData:      make([][]byte, 0),
		ReturnCode:      vmcommon.Ok,
		ReturnMessage:   "",
		GasRemaining:    0,
		GasRefund:       big.NewInt(0),
		OutputAccounts:  outputAccounts,
		DeletedAccounts: make([][]byte, 0),
		TouchedAccounts: make([][]byte, 0),
		Logs:            make([]*vmcommon.LogEntry, 0),
	}, nil
}

func (ae *ArwenTestExecutor) validatorRewardOutput(tx *mj.Transaction) (*vmi.VMOutput, error) {
	reward := tx.Value.Value
	recipient := ae.World.AcctMap.GetAccount(tx.To.Value)
	if recipient == nil {
		return nil, fmt.Errorf("Tx recipient (address: %s) does not exist", hex.EncodeToString(tx.To.Value))
	}
	recipient.BalanceDelta = reward
	storageKalyan3104Reward := recipient.StorageValue(Kalyan3104RewardKey)
	storageKalyan3104Reward = big.NewInt(0).Add(
		big.NewInt(0).SetBytes(storageKalyan3104Reward),
		reward).Bytes()

	outputAccounts := make(map[string]*vmcommon.OutputAccount)
	outputAccounts[string(tx.To.Value)] = &vmcommon.OutputAccount{
		Address:      tx.To.Value,
		BalanceDelta: tx.Value.Value,
		StorageUpdates: map[string]*vmi.StorageUpdate{
			Kalyan3104RewardKey: &vmi.StorageUpdate{
				Offset: []byte(Kalyan3104RewardKey),
				Data:   storageKalyan3104Reward,
			},
		},
	}

	return &vmcommon.VMOutput{
		ReturnData:      make([][]byte, 0),
		ReturnCode:      vmcommon.Ok,
		ReturnMessage:   "",
		GasRemaining:    0,
		GasRefund:       big.NewInt(0),
		OutputAccounts:  outputAccounts,
		DeletedAccounts: make([][]byte, 0),
		TouchedAccounts: make([][]byte, 0),
		Logs:            make([]*vmcommon.LogEntry, 0),
	}, nil
}

func outOfFundsResult() *vmi.VMOutput {
	return &vmcommon.VMOutput{
		ReturnData:      make([][]byte, 0),
		ReturnCode:      vmcommon.OutOfFunds,
		ReturnMessage:   "",
		GasRemaining:    0,
		GasRefund:       big.NewInt(0),
		OutputAccounts:  make(map[string]*vmcommon.OutputAccount),
		DeletedAccounts: make([][]byte, 0),
		TouchedAccounts: make([][]byte, 0),
		Logs:            make([]*vmcommon.LogEntry, 0),
	}
}

func (ae *ArwenTestExecutor) scCreate(tx *mj.Transaction) (*vmi.VMOutput, error) {
	input := &vmi.ContractCreateInput{
		ContractCode: tx.Code.Value,
		VMInput: vmi.VMInput{
			CallerAddr:  tx.From.Value,
			Arguments:   mj.JSONBytesValues(tx.Arguments),
			CallValue:   tx.Value.Value,
			GasPrice:    tx.GasPrice.Value,
			GasProvided: tx.GasLimit.Value,
		},
	}

	return ae.vm.RunSmartContractCreate(input)
}

func (ae *ArwenTestExecutor) scCall(tx *mj.Transaction) (*vmi.VMOutput, error) {
	recipient := ae.World.AcctMap.GetAccount(tx.To.Value)
	if recipient == nil {
		return nil, fmt.Errorf("Tx recipient (address: %s) does not exist", hex.EncodeToString(tx.To.Value))
	}
	if len(recipient.Code) == 0 {
		return nil, fmt.Errorf("Tx recipient (address: %s) is not a smart contract", hex.EncodeToString(tx.To.Value))
	}
	input := &vmi.ContractCallInput{
		RecipientAddr: tx.To.Value,
		Function:      tx.Function,
		VMInput: vmi.VMInput{
			CallerAddr:  tx.From.Value,
			Arguments:   mj.JSONBytesValues(tx.Arguments),
			CallValue:   tx.Value.Value,
			GasPrice:    tx.GasPrice.Value,
			GasProvided: tx.GasLimit.Value,
		},
	}

	return ae.vm.RunSmartContractCall(input)
}

func (ae *ArwenTestExecutor) updateStateAfterTx(
	tx *mj.Transaction,
	output *vmi.VMOutput) error {

	// subtract call value from sender (this is not reflected in the delta)
	// except for validatorReward, there is no sender there
	if tx.Type.HasSender() {
		_ = ae.World.UpdateBalanceWithDelta(tx.From.Value, big.NewInt(0).Neg(tx.Value.Value))
	}

	accountsSlice := make([]*vmi.OutputAccount, len(output.OutputAccounts))
	i := 0
	for _, account := range output.OutputAccounts {
		accountsSlice[i] = account
		i++
	}

	// update accounts based on deltas
	updErr := ae.World.UpdateAccounts(accountsSlice, output.DeletedAccounts, tx.From.Value)
	if updErr != nil {
		return updErr
	}

	// sum of all balance deltas should equal call value (unless we got an error)
	// (unless it is validatorReward, when funds just pop into existence)
	if tx.Type.HasSender() {
		sumOfBalanceDeltas := big.NewInt(0)
		for _, oa := range output.OutputAccounts {
			sumOfBalanceDeltas = sumOfBalanceDeltas.Add(sumOfBalanceDeltas, oa.BalanceDelta)
		}
		if sumOfBalanceDeltas.Cmp(tx.Value.Value) != 0 {
			return fmt.Errorf("sum of balance deltas should equal call value. Sum of balance deltas: %d (0x%x). Call value: %d (0x%x)",
				sumOfBalanceDeltas, sumOfBalanceDeltas, tx.Value.Value, tx.Value.Value)
		}
	}

	return nil
}
