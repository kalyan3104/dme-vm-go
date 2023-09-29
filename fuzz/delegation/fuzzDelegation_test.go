package delegation

import (
	"flag"
	"fmt"
	"math/big"
	"math/rand"
	"os"
	"path/filepath"
	"testing"
	"time"

	mc "github.com/kalyan3104/dme-vm-util/test-util/mandos/controller"
	"github.com/stretchr/testify/require"
)

var fuzz = flag.Bool("fuzz", false, "fuzz")

func getTestRoot() string {
	exePath, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	arwenTestRoot := filepath.Join(exePath, "../../test")
	return arwenTestRoot
}

func newExecutorWithPaths() *fuzzDelegationExecutor {
	fileResolver := mc.NewDefaultFileResolver().
		ReplacePath(
			"delegation.wasm",
			filepath.Join(getTestRoot(), "delegation_v0.3/delegation.wasm")).
		ReplacePath(
			"auction-mock.wasm",
			filepath.Join(getTestRoot(), "delegation_v0.3/auction-mock.wasm"))

	pfe, err := newFuzzDelegationExecutor(fileResolver)
	if err != nil {
		panic(err)
	}
	return pfe
}

func TestFuzzDelegation(t *testing.T) {
	if !*fuzz {
		t.Skip("skipping test; only run with --fuzz argument")
	}

	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	pfe := newExecutorWithPaths()
	defer pfe.saveGeneratedScenario()

	err := pfe.init(&fuzzDelegationExecutorInitArgs{
		serviceFee:                  r.Intn(10000),
		numBlocksBeforeForceUnstake: r.Intn(1000),
		numBlocksBeforeUnbond:       r.Intn(1000),
		numDelegators:               10,
		stakePerNode:                big.NewInt(1000000000),
	})
	require.Nil(t, err)
	pfe.enableAutoActivation()
	pfe.increaseBlockNonce(r.Intn(10000))

	maxStake := big.NewInt(0).Mul(pfe.stakePerNode, big.NewInt(2))
	maxSystemReward := big.NewInt(1000000000)

	re := newRandomEventProvider()
	for stepIndex := 0; stepIndex < 1000; stepIndex++ {
		re.reset()
		switch {
		case re.withProbability(0.05):
			// increment block nonce
			err = pfe.increaseBlockNonce(r.Intn(1000))
			require.Nil(t, err)
		case re.withProbability(0.05):
			// add nodes
			err = pfe.addNodes(r.Intn(3))
			require.Nil(t, err)
		case re.withProbability(0.05):
			// stake
			delegatorIdx := r.Intn(pfe.numDelegators + 1)
			stake := big.NewInt(0).Rand(r, maxStake)
			err = pfe.stake(delegatorIdx, stake)
			require.Nil(t, err)
		case re.withProbability(0.05):
			// withdraw inactive stake
			delegatorIdx := r.Intn(pfe.numDelegators + 1)
			stake := big.NewInt(0).Rand(r, maxStake)
			err = pfe.withdrawInactiveStake(delegatorIdx, stake)
			require.Nil(t, err)
		case re.withProbability(0.05):
			// add system rewards
			rewards := big.NewInt(0).Rand(r, maxSystemReward)
			err = pfe.addRewards(rewards)
			require.Nil(t, err)
		case re.withProbability(0.2):
			// claim rewards
			delegatorIdx := r.Intn(pfe.numDelegators + 1)
			err = pfe.claimRewards(delegatorIdx)
			require.Nil(t, err)
		case re.withProbability(0.05):
			// computeAllRewards
			err = pfe.computeAllRewards()
			require.Nil(t, err)
		case re.withProbability(0.05):
			// announceUnStake
			delegatorIdx := r.Intn(pfe.numDelegators + 1)
			amount := big.NewInt(0).Rand(r, maxStake)
			err = pfe.announceUnStake(delegatorIdx, amount)
			require.Nil(t, err)
		case re.withProbability(0.05):
			// purchaseStake
			sellerIdx := r.Intn(pfe.numDelegators + 1)
			buyerIdx := r.Intn(pfe.numDelegators + 1)
			amount := big.NewInt(0).Rand(r, maxStake)
			err = pfe.purchaseStake(sellerIdx, buyerIdx, amount)
			require.Nil(t, err)
		case re.withProbability(0.05):
			// unStake
			delegatorIdx := r.Intn(pfe.numDelegators + 1)
			err = pfe.unStake(delegatorIdx)
			require.Nil(t, err)
		case re.withProbability(0.05):
			// unBondAllAvailable
			err = pfe.unBondAllAvailable()
			require.Nil(t, err)
		default:
		}
	}

	err = pfe.checkContractBalanceVsState()
	if err != nil {
		fmt.Println(err)
		return
	}

	// all delegators (incl. owner) withdraw all inactive stake
	for delegatorIdx := 0; delegatorIdx <= pfe.numDelegators; delegatorIdx++ {
		err = pfe.withdrawAllInactiveStake(delegatorIdx)
		require.Nil(t, err)
	}

	// all delegators (incl. owner) claim all rewards
	err = pfe.computeAllRewards()
	for delegatorIdx := 0; delegatorIdx <= pfe.numDelegators; delegatorIdx++ {
		err = pfe.claimRewards(delegatorIdx)
		require.Nil(t, err)

		err = pfe.checkContractBalanceVsState()
		if err != nil {
			fmt.Println(err)
			return
		}
	}

	// check that delegators got all rewards out
	totalDelegatorBalance := pfe.getAllDelegatorsBalance()
	require.True(t, pfe.totalRewards.Cmp(totalDelegatorBalance) == 0,
		"Rewards don't match. Total rewards: %d. Total delegator balance: %d.",
		pfe.totalRewards, totalDelegatorBalance)

	// all delegators (incl. owner) announce unstake
	for delegatorIdx := 0; delegatorIdx <= pfe.numDelegators; delegatorIdx++ {
		err = pfe.announceUnStakeAll(delegatorIdx)
		require.Nil(t, err)
	}

	pfe.increaseBlockNonce(pfe.numBlocksBeforeForceUnstake + 1)

	// all delegators (incl. owner) unstake
	for delegatorIdx := 0; delegatorIdx <= pfe.numDelegators; delegatorIdx++ {
		err = pfe.unStake(delegatorIdx)
		require.Nil(t, err)
	}

	pfe.increaseBlockNonce(pfe.numBlocksBeforeUnbond + 1)

	// unBondAllAvailable
	err = pfe.unBondAllAvailable()
	require.Nil(t, err)

	// auction SC should have no more funds
	auctionBalanceAfterUnbond := pfe.getAuctionBalance()
	require.True(t, auctionBalanceAfterUnbond.Sign() == 0,
		"Auction still has balance after full unbond: %d",
		auctionBalanceAfterUnbond)

	// all delegators (incl. owner) withdraw all inactive stake
	for delegatorIdx := 0; delegatorIdx <= pfe.numDelegators; delegatorIdx++ {
		err = pfe.withdrawAllInactiveStake(delegatorIdx)
		require.Nil(t, err)
	}

	withdrawnAtTheEnd := pfe.getWithdrawTargetBalance()
	require.True(t, withdrawnAtTheEnd.Cmp(pfe.totalStakeAdded) == 0,
		"Stake added and withdrawn doesn't match. Staked: %d. Withdrawn: %d. Off by: %d",
		pfe.totalStakeAdded, withdrawnAtTheEnd,
		big.NewInt(0).Sub(pfe.totalStakeAdded, withdrawnAtTheEnd))

}
