package StrideLiquidStake

import (
	"StrideLiquidStake/module"
	"fmt"
)
import "testing"

// mnemonic used to store the mnemonic phrase for the test
// bench32Addr are used to store the  bench32 address. It must generate by the above mnemonic
var mnemonic = ""
var bench32Addr = ""

func TestLiquidStake(t *testing.T) {
	fmt.Println("LiquidStakeBot testing!")
	liquidBot := NewLiquidStakeBot("uatom", []module.Account{module.Account{
		Mnemonic: mnemonic,
		DelAddr:  bench32Addr,
	}})
	ok := liquidBot.LiquidStake(module.Account{
		Mnemonic: mnemonic,
		DelAddr:  bench32Addr,
	}, 10)
	if !ok {
		t.Errorf("LiquidStake failed!")
	}
}

func TestReinvestAllRewards(t *testing.T) {
	fmt.Println("ReinvestAllRewards testing!")
	liquidBot := NewLiquidStakeBot("uatom", []module.Account{module.Account{
		Mnemonic: mnemonic,
		DelAddr:  bench32Addr,
	}})
	ok := liquidBot.ReinvestAllRewards()
	if !ok {
		t.Errorf("LiquidStake failed!")
	}

}

func TestRebaseToken(t *testing.T) {
	fmt.Println("RebaseTokens testing!")
	liquidBot := NewLiquidStakeBot("uatom", []module.Account{module.Account{
		Mnemonic: mnemonic,
		DelAddr:  bench32Addr,
	}})
	ok := liquidBot.RebaseTokens()
	if !ok {
		t.Errorf("LiquidStake failed!")
	}
	for k, v := range liquidBot.SwapRate {
		fmt.Println(k, v)
	}
}
