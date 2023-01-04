package StrideLiquidStake

import (
	"StrideLiquidStake/module"
	"StrideLiquidStake/tool"
	"context"
	"fmt"
	"github.com/cosmos/cosmos-sdk/client"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"strconv"
)

type LiquidStakeBot struct {
	ctx       context.Context
	TxContext client.Context
	//Accounts is a map from address to account object
	Accounts map[string]module.Account
	//LiquidStakeDenom is the Denom name that you want to liquid stake
	//Now the stride has support: uatom,ujuno,uosmo,ustars
	LiquidStakeDenom string
	//SwapRate is a map from hostDenom to it's swap ratio.
	//e.g. uatom's ratio is 1.05 means 1 statom can swap 1.05 atom.
	SwapRate map[string]float64
}

// NewLiquidStakeBot new a liquid stake bot to help user lquidstake, claim and auto delegate rewards and query ratio
// params:
// denom is the token you want to liquid stake
// accounts is the users account list
// returns:
// A LiquidStakeBot object
func NewLiquidStakeBot(denom string, accounts []module.Account) *LiquidStakeBot {
	accountsMap := make(map[string]module.Account)
	for _, v := range accounts {
		accountsMap[v.DelAddr] = v
	}
	swapRate := make(map[string]float64)
	return &LiquidStakeBot{
		ctx:              context.Background(),
		TxContext:        tool.InitTxContext(),
		Accounts:         accountsMap,
		LiquidStakeDenom: denom,
		SwapRate:         swapRate,
	}
}

// LiquidStake stake users native token and receive st token.
// params:
// account is the account you want to liquid stake.
// amount: how much token you want to stake.tips: amount is the number of "utoken",1 atom=1,000,000 uatom
// returns:
// true: success, false: some error occured
func (l *LiquidStakeBot) LiquidStake(account module.Account, amount uint64) bool {
	addr, err := sdk.AccAddressFromBech32(account.DelAddr)
	if err != nil {
		fmt.Println("err:", err)
		return false
	}
	_, sequence, accnum, err := tool.QueryAccountInfo(l.ctx, addr, l.TxContext)
	msg, err := tool.NewLiquidStakeMsg(addr.String(), amount, l.LiquidStakeDenom)
	resp, err := tool.SendStrideMsgWithMnemonic(l.ctx, l.TxContext, account.Mnemonic, "", sequence, accnum, msg)
	if err != nil {
		fmt.Println("err:", err)
		return false
	}
	fmt.Println("send tx success,tx hash is:", resp.TxHash)
	return true
}

// ReinvestAllRewards claim all accounts' stake reward and delegate them.
// params:
// none
// returns:
// true: success, false: some error occured
func (l *LiquidStakeBot) ReinvestAllRewards() bool {
	for _, acc := range l.Accounts {
		addr, err := sdk.AccAddressFromBech32(acc.DelAddr)
		reward, err := tool.QueryTotalReward(l.ctx, addr, l.TxContext)
		_, sequence, accnum, err := tool.QueryAccountInfo(l.ctx, addr, l.TxContext)
		if err != nil {
			fmt.Println("err:", err)
			return false
		}
		fmt.Println("addr:", acc.DelAddr, "total reward is:", reward.Total[0].Amount, reward.Total[0].Denom)
		validators, err := tool.QueryStakeValidators(l.ctx, addr)
		if err != nil {
			fmt.Println("err:", err)
			return false
		}
		valAddrs := make([]sdk.ValAddress, 0, len(validators))
		for _, v := range validators {
			valAddr, err := sdk.ValAddressFromBech32(v)
			if err != nil {
				fmt.Println("err:", err)
				return false
			}
			valAddrs = append(valAddrs, valAddr)
		}
		msgs, err := tool.NewClaimRewardMsgs(addr, valAddrs)
		for _, msg := range msgs {
			resp, err := tool.SendStrideMsgWithMnemonic(l.ctx, l.TxContext, acc.Mnemonic, "", sequence, accnum, msg)
			if err != nil {
				fmt.Println("err:", err)
				return false
			}
			sequence++
			fmt.Println(resp.TxHash)
		}
		delAmount, err := strconv.ParseFloat(reward.Total[0].Amount, 64)
		if err != nil {
			fmt.Println("err:", err)
			return false
		}
		delIntAmount := int64(delAmount)
		amount := sdk.NewInt64Coin(reward.Total[0].Denom, delIntAmount)
		msg, err := tool.NewDelegateMsg(addr, valAddrs[0], amount)
		if err != nil {
			fmt.Println("err:", err)
			return false
		}
		resp, err := tool.SendStrideMsgWithMnemonic(l.ctx, l.TxContext, acc.Mnemonic, "", sequence, accnum, msg)
		if err != nil {
			fmt.Println("err:", err)
			return false
		}
		sequence++
		fmt.Println(resp.TxHash)
	}
	return true
}

// RebaseTokens query the ratio of the native token and st token and save the ratio in SwapRate.
// params:
// none
// returns:
// true: success, false: some error occured
func (l *LiquidStakeBot) RebaseTokens() bool {
	h, err := tool.QuerySwapRate(l.ctx, l.TxContext)
	if err != nil {
		fmt.Println("err:", err)
		return false
	}
	for _, v := range h.HostZone {
		fmt.Println(v.HostDenom, v.LastRedemptionRate)
		rate, err := strconv.ParseFloat(v.LastRedemptionRate, 64)
		if err != nil {
			fmt.Println("err:", err)
			return false
		}
		l.SwapRate[v.HostDenom] = rate
	}
	return true
}
