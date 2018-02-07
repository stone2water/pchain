package ethapi

import (
	"github.com/ethereum/go-ethereum/rpc"
	"golang.org/x/net/context"
	"github.com/tendermint/tendermint/rpc/core/types"
	"fmt"
	"github.com/tendermint/go-wire"
	/*
	"math/big"
	"github.com/ethereum/go-ethereum/accounts"
	ethapi "github.com/ethereum/go-ethereum/internal/ethapi"
	*/
	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/common"
)

type PublicTendermintAPI struct {
	am *accounts.Manager
	b  Backend
	Client Client
}

// NewPublicEthereumAPI creates a new Etheruem protocol API.
func NewPublicTendermintAPI(b Backend) *PublicTendermintAPI {
	return &PublicTendermintAPI{
		am: b.AccountManager(),
		b:  b,
		Client: b.Client(),
	}
}

// GasPrice returns a suggestion for a gas price.
func (s *PublicTendermintAPI) GetBlock(ctx context.Context, blockNumber rpc.BlockNumber) (string, error) {

	var result core_types.TMResult

	//fmt.Printf("GetBlock() called with startBlock: %v\n", blockNumber)
	params := map[string]interface{}{
		"height":  blockNumber,
	}

	_, err := s.Client.Call("block", params, &result)
	if err != nil {
		fmt.Println(err)
		return "", err
	}

	//fmt.Printf("tdm_getBlock: %v\n", result)
	return result.(*core_types.ResultBlock).Block.String(), nil
}

// GasPrice returns a suggestion for a gas price.
func (s *PublicTendermintAPI) GetValidator(ctx context.Context, address string) (string, error) {

	var result core_types.TMResult

	//fmt.Printf("GetValidator() called with address: %v\n", address)
	params := map[string]interface{}{
		"address":  address,
	}

	_, err := s.Client.Call("validator_epoch", params, &result)
	if err != nil {
		fmt.Println(err)
		return "", err
	}

	//fmt.Printf("tdm_getValidator: %v\n", result)
	return string(wire.JSONBytes(result.(*core_types.ResultValidatorEpoch))), nil
}


func (s *PublicTendermintAPI) SendValidatorMessage(ctx context.Context, from common.Address, epoch int, power uint64, action string) (string, error) {
	fmt.Println("in func (s *PublicTendermintAPI) SendValidatorMessage()")

	var result core_types.TMResult

	fromStr := fmt.Sprintf("%X", from.Bytes())
	data := fmt.Sprintf("%s-%X-%X-%s", fromStr, epoch, power, action)
	fmt.Printf("in func (s *PublicTendermintAPI) SendValidatorMessage(), data to sign is: %v\n", data)

	signature, err := s.Sign(ctx, from, data)
	if err != nil {
		return "", err
	}

	params := map[string]interface{}{
		"address": fromStr,
		"epoch":  epoch,
		"power":  power,
		"action":   action,
		"signature":  signature,
	}

	_, err = s.Client.Call("validator_operation", params, &result)
	if err != nil {
		fmt.Println(err)
		return "", err
	}

	return string(wire.JSONBytes(result.(*core_types.ResultValidatorOperation))), nil
}

// SendTransaction creates a transaction for the given argument, sign it and submit it to the
// transaction pool.
func (s *PublicTendermintAPI) Sign(ctx context.Context, from common.Address, dataStr string) ([]byte, error) {

	fmt.Printf("(s *PublicTendermintAPI) SignAndVerify(), from: %x, data: %v\n", from, dataStr)

	//fmt.Printf("(s *PublicTransactionPoolAPI) SendTransaction(), s.b.GetSendTxMutex() is %v\n", s.b.GetSendTxMutex())

	s.b.GetSendTxMutex().Lock()
	defer s.b.GetSendTxMutex().Unlock()

	data := []byte(dataStr)

	// Look up the wallet containing the requested signer
	account := accounts.Account{Address: from}
	//fmt.Println("(s *PublicBlockChainAPI) SendTransaction() 1")
	wallet, err := s.b.AccountManager().Find(account)
	if err != nil {
		return []byte{}, err
	}

	sig, err := wallet.SignHash(account, signHash(data))
	if err != nil {
		return []byte{}, err
	}
	sig[64] += 27

	fmt.Printf("SignAndVerify(), sig is: %X\n", sig)

	return sig, nil
}