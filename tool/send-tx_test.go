package tool

import (
	"context"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"testing"
)

var mnemonic = ""
var bench32Addr = ""

func TestQuerySwapRate(t *testing.T) {
	fmt.Println("QuerySwapRate testing!")
	ctx := context.Background()
	TxContext := InitTxContext()
	h, err := QuerySwapRate(ctx, TxContext)
	if err != nil {
		t.Errorf("%s", err)
	}
	for _, v := range h.HostZone {
		fmt.Println(v.HostDenom, v.LastRedemptionRate)
	}
}

func TestQueryTotalReward(t *testing.T) {
	fmt.Println("QueryTotalReward testing!")
	ctx := context.Background()
	TxContext := InitTxContext()
	acc, err := sdk.AccAddressFromBech32(bench32Addr)
	if err != nil {
		t.Errorf("%s", err)
	}
	reward, err := QueryTotalReward(ctx, acc, TxContext)
	if err != nil {
		t.Errorf("%s", err)
	}
	fmt.Println(reward.Total[0].Denom, reward.Total[0].Amount)
}

func TestQueryAccountInfo(t *testing.T) {
	fmt.Println("QueryAccountInfo testing!")
	ctx := context.Background()
	TxContext := InitTxContext()
	addr, err := sdk.AccAddressFromBech32(bench32Addr)
	if err != nil {
		t.Errorf("%s", err)
	}
	acc, _, _, err := QueryAccountInfo(ctx, addr, TxContext)
	if err != nil {
		t.Errorf("%s", err)
	}
	fmt.Println(acc.Account.BaseVestingAccount.BaseAccount.AccountNumber,
		acc.Account.BaseVestingAccount.BaseAccount.Sequence)
}

func TestQueryStakeValidators(t *testing.T) {
	fmt.Println("QueryStakeValidators testing!")
	ctx := context.Background()
	addr, err := sdk.AccAddressFromBech32(bench32Addr)
	if err != nil {
		t.Errorf("%s", err)
	}
	valdators, err := QueryStakeValidators(ctx, addr)
	if err != nil {
		t.Errorf("%s", err)
	}
	fmt.Println(valdators[0])
}

func TestClaimReward(t *testing.T) {
	fmt.Println("ClaimReward testing!")
	ctx := context.Background()
	TxContext := InitTxContext()
	addr, err := sdk.AccAddressFromBech32(bench32Addr)
	if err != nil {
		t.Errorf("%s", err)
	}
	validators, err := QueryStakeValidators(ctx, addr)
	if err != nil {
		t.Errorf("%s", err)
	}
	valAddrs := make([]sdk.ValAddress, 0, len(validators))
	for _, v := range validators {
		valAddr, err := sdk.ValAddressFromBech32(v)
		if err != nil {
			t.Errorf("%s", err)
		}
		valAddrs = append(valAddrs, valAddr)
	}
	msgs, err := NewClaimRewardMsgs(addr, valAddrs)
	_, sequence, accnum, err := QueryAccountInfo(ctx, addr, TxContext)
	for _, msg := range msgs {
		resp, err := SendStrideMsgWithMnemonic(ctx, TxContext, mnemonic, "", sequence, accnum, msg)
		if err != nil {
			t.Errorf("%s", err)
		}
		fmt.Println(resp.TxHash)
	}
}

func TestLiquidStake(t *testing.T) {
	fmt.Println("ClaimReward testing!")
	ctx := context.Background()
	TxContext := InitTxContext()
	addr, err := sdk.AccAddressFromBech32(bench32Addr)
	if err != nil {
		t.Errorf("%s", err)
	}
	validators, err := QueryStakeValidators(ctx, addr)
	if err != nil {
		t.Errorf("%s", err)
	}
	valAddrs := make([]sdk.ValAddress, 0, len(validators))
	for _, v := range validators {
		valAddr, err := sdk.ValAddressFromBech32(v)
		if err != nil {
			t.Errorf("%s", err)
		}
		valAddrs = append(valAddrs, valAddr)
	}
	_, sequence, accnum, err := QueryAccountInfo(ctx, addr, TxContext)
	msg, err := NewLiquidStakeMsg(addr.String(), 10, "uatom")
	resp, err := SendStrideMsgWithMnemonic(ctx, TxContext, mnemonic, "", sequence, accnum, msg)
	if err != nil {
		t.Errorf("%s", err)
	}
	fmt.Println(resp.TxHash)
}

func TestDelegate(t *testing.T) {
	fmt.Println("ClaimReward testing!")
	ctx := context.Background()
	TxContext := InitTxContext()
	addr, err := sdk.AccAddressFromBech32(bench32Addr)
	if err != nil {
		t.Errorf("%s", err)
	}
	validators, err := QueryStakeValidators(ctx, addr)
	if err != nil {
		t.Errorf("%s", err)
	}
	valAddrs := make([]sdk.ValAddress, 0, len(validators))
	for _, v := range validators {
		valAddr, err := sdk.ValAddressFromBech32(v)
		if err != nil {
			t.Errorf("%s", err)
		}
		valAddrs = append(valAddrs, valAddr)
	}
	_, accnum, sequence, err := QueryAccountInfo(ctx, addr, TxContext)
	coin := sdk.NewInt64Coin("ustrd", 10)
	msg, err := NewDelegateMsg(addr, valAddrs[0], coin)
	resp, err := SendStrideMsgWithMnemonic(ctx, TxContext, mnemonic, "", sequence, accnum, msg)
	if err != nil {
		t.Errorf("%s", err)
	}
	fmt.Println(resp.TxHash)

}
