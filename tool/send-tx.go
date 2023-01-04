package tool

import (
	"StrideLiquidStake/module"
	"context"
	"encoding/json"
	"fmt"
	"github.com/Stride-Labs/stride/v4/app"
	stakeibctypes "github.com/Stride-Labs/stride/v4/x/stakeibc/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	ctx "github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/cosmos/cosmos-sdk/crypto/hd"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	ctypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/cosmos/cosmos-sdk/types/tx"
	"github.com/cosmos/cosmos-sdk/types/tx/signing"
	asigning "github.com/cosmos/cosmos-sdk/x/auth/signing"
	"github.com/cosmos/cosmos-sdk/x/auth/types"
	atypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	dtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	stypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"google.golang.org/grpc"
	"strconv"

	"os"
)

const (
	CHAIN_ID            = "stride-1"
	GRPC_SERVER_ADDRESS = "65.108.141.109:11090"
	GAS_LIMIT           = 200000
	GAS_FEE             = 2000
)

func InitTxContext() client.Context {
	encodingConfig := app.MakeEncodingConfig()
	TxContext := client.Context{}.WithChainID(CHAIN_ID).
		WithCodec(encodingConfig.Marshaler).
		WithInterfaceRegistry(encodingConfig.InterfaceRegistry).
		WithTxConfig(encodingConfig.TxConfig).
		WithLegacyAmino(encodingConfig.Amino).
		WithInput(os.Stdin).
		WithAccountRetriever(types.AccountRetriever{}).
		WithBroadcastMode(flags.BroadcastBlock).
		WithHomeDir(app.DefaultNodeHome).
		WithViper("ICA")
	return TxContext
}

func NewLiquidStakeMsg(creator string, amount uint64, hostDenom string) (sdk.Msg, error) {
	msg := stakeibctypes.NewMsgLiquidStake(creator, amount, hostDenom)
	if err := msg.ValidateBasic(); err != nil {
		return nil, err
	}

	return msg, nil
}

func NewGrpcConnection() *grpc.ClientConn {
	grpcConn, _ := grpc.Dial(
		GRPC_SERVER_ADDRESS, // your gRPC server address.
		grpc.WithInsecure(), // The SDK doesn't support any transport security mechanism.
	)
	return grpcConn
}
func CloseGrpcConnection(grpcConn *grpc.ClientConn) {
	err := grpcConn.Close()
	if err != nil {
		fmt.Println("grpcConn.Close() err:", err)
	}
}

func QueryAccountInfo(ctx context.Context, address sdk.AccAddress, TxContext client.Context) (acc *module.AccountInfo, seq uint64, accnum uint64, err error) {
	// Create a connection to the gRPC server.
	grpcConn := NewGrpcConnection()
	defer CloseGrpcConnection(grpcConn)

	// This creates a gRPC client to query the x/bank service.
	authClient := atypes.NewQueryClient(grpcConn)
	authRes, err := authClient.Account(ctx, &atypes.QueryAccountRequest{address.String()})
	if err != nil {
		return nil, 0, 0, err
	}
	res, err := TxContext.Codec.MarshalJSON(authRes)
	if err != nil {
		return nil, 0, 0, err
	}
	acc = &module.AccountInfo{}
	err = json.Unmarshal(res, acc)
	if err != nil {
		return nil, 0, 0, err
	}
	accnum, err = strconv.ParseUint(acc.Account.BaseVestingAccount.BaseAccount.AccountNumber, 10, 64)
	if err != nil {
		return nil, 0, 0, err
	}
	seq, err = strconv.ParseUint(acc.Account.BaseVestingAccount.BaseAccount.Sequence, 10, 64)
	if err != nil {
		return nil, 0, 0, err
	}
	return acc, seq, accnum, nil
}

func QueryStakeValidators(ctx context.Context, address sdk.AccAddress) ([]string, error) {
	grpcConn := NewGrpcConnection()
	defer CloseGrpcConnection(grpcConn)
	queryClient := dtypes.NewQueryClient(grpcConn)
	delValsRes, err := queryClient.DelegatorValidators(ctx, &dtypes.QueryDelegatorValidatorsRequest{DelegatorAddress: address.String()})
	if err != nil {
		return nil, err
	}
	validators := delValsRes.Validators
	return validators, nil
}

func QueryTotalReward(ctx context.Context, delAddr sdk.AccAddress, TxContext client.Context) (*module.TotalReward, error) {
	grpcConn := NewGrpcConnection()
	defer CloseGrpcConnection(grpcConn)
	queryClient := dtypes.NewQueryClient(grpcConn)
	rewardRes, err := queryClient.DelegationTotalRewards(ctx, &dtypes.QueryDelegationTotalRewardsRequest{DelegatorAddress: delAddr.String()})
	if err != nil {
		return nil, err
	}
	res, err := TxContext.Codec.MarshalJSON(rewardRes)
	if err != nil {
		return nil, err
	}
	reward := &module.TotalReward{}
	err = json.Unmarshal(res, reward)
	if err != nil {
		return nil, err
	}
	return reward, nil
}

func QuerySwapRate(ctx context.Context, TxContext client.Context) (*module.HostZone, error) {
	grpcConn := NewGrpcConnection()
	defer CloseGrpcConnection(grpcConn)
	queryClient := stakeibctypes.NewQueryClient(grpcConn)
	hostZoneRes, err := queryClient.HostZoneAll(ctx, &stakeibctypes.QueryAllHostZoneRequest{Pagination: &query.PageRequest{
		Key:        nil,
		Offset:     0,
		Limit:      100,
		CountTotal: false,
		Reverse:    false,
	}})
	if err != nil {
		return nil, err
	}
	res, err := TxContext.Codec.MarshalJSON(hostZoneRes)
	hostZone := &module.HostZone{}
	err = json.Unmarshal(res, hostZone)
	return hostZone, err
}

func NewClaimRewardMsgs(delAddr sdk.AccAddress, valAddrs []sdk.ValAddress) ([]sdk.Msg, error) {
	msgs := make([]sdk.Msg, 0, len(valAddrs))
	for _, valAddr := range valAddrs {
		msg := dtypes.NewMsgWithdrawDelegatorReward(delAddr, valAddr)
		if err := msg.ValidateBasic(); err != nil {
			return nil, err
		}
		msgs = append(msgs, msg)
	}
	return msgs, nil
}

func NewDelegateMsg(delAddr sdk.AccAddress, valAddrs sdk.ValAddress, amount sdk.Coin) (sdk.Msg, error) {
	msg := stypes.NewMsgDelegate(delAddr, valAddrs, amount)
	return msg, nil
}
func SendStrideMsgWithMnemonic(ctx context.Context, TxContext client.Context, mnemonic, memo string, sequence, accnum uint64, msg sdk.Msg) (*sdk.TxResponse, error) {
	txBuilder := TxContext.TxConfig.NewTxBuilder()
	txBuilder.SetGasLimit(GAS_LIMIT)
	txBuilder.SetFeeAmount(sdk.NewCoins(sdk.NewInt64Coin("ustrd", GAS_FEE)))
	txBuilder.SetMemo(memo)
	priv, err := NewPrivateKeyByMnemonic(mnemonic)
	if err != nil {
		return nil, err
	}
	if err = SignTx(txBuilder, TxContext, priv, sequence, accnum, msg); err != nil {
		return nil, err
	}
	txBytes, err := TxContext.TxConfig.TxEncoder()(txBuilder.GetTx())
	if err != nil {
		return nil, err
	}
	res, err := BrocastTransaction(ctx, GRPC_SERVER_ADDRESS, txBytes)
	return res, nil

}

func NewPrivateKeyByMnemonic(mnemonic string) (ctypes.PrivKey, error) {
	algo, err := keyring.NewSigningAlgoFromString(string(hd.Secp256k1Type), keyring.SigningAlgoList{hd.Secp256k1})
	if err != nil {
		return nil, err
	}

	// create master key and derive first key for keyring
	prvbz, err := algo.Derive()(mnemonic, "", hd.CreateHDPath(118, 0, 0).String())
	if err != nil {
		return nil, err
	}
	prv := algo.Generate()(prvbz)
	return prv, nil
}

func SignTx(txBuilder client.TxBuilder, TxContext client.Context, priv ctypes.PrivKey, sequence uint64, accountNum uint64, msg sdk.Msg) error {
	var err error
	if err = txBuilder.SetMsgs(msg); err != nil {
		return err
	}
	sigV2 := signing.SignatureV2{
		PubKey: priv.PubKey(),
		Data: &signing.SingleSignatureData{
			SignMode:  TxContext.TxConfig.SignModeHandler().DefaultMode(),
			Signature: nil,
		},
		Sequence: sequence,
	}
	err = txBuilder.SetSignatures(sigV2)
	if err != nil {
		return err
	}
	signerData := asigning.SignerData{
		ChainID:       CHAIN_ID,
		AccountNumber: accountNum,
		Sequence:      sequence,
	}
	sigV2, err = ctx.SignWithPrivKey(TxContext.TxConfig.SignModeHandler().DefaultMode(), signerData, txBuilder, priv, TxContext.TxConfig, sequence)
	if err != nil {
		return err
	}
	if err = txBuilder.SetSignatures(sigV2); err != nil {
		return err
	}
	return nil
}

func BrocastTransaction(ctx context.Context, grpcAdrr string, txBytes []byte) (*sdk.TxResponse, error) {
	grpcConn, err := grpc.Dial(
		grpcAdrr,            // Or your gRPC server address.
		grpc.WithInsecure(), // The SDK doesn't support any transport security mechanism.
	)
	if err != nil {
		return nil, err
	}
	txClient := tx.NewServiceClient(grpcConn)
	grpcRes, err := txClient.BroadcastTx(
		ctx,
		&tx.BroadcastTxRequest{
			Mode:    tx.BroadcastMode_BROADCAST_MODE_SYNC,
			TxBytes: txBytes, // Proto-binary of the signed transaction, see previous step.
		},
	)
	if err != nil {
		return nil, err
	}

	fmt.Println(grpcRes.TxResponse)
	defer grpcConn.Close()
	return grpcRes.TxResponse, nil
}
