package ethapi

import (
	"encoding/binary"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core"
	st "github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
	"math/big"
	"strings"
)

const (
	CCCFuncName   = "CreateChildChain"
	JCCFuncName   = "JoinChildChain"
	DIMCFuncName  = "DepositInMainChain"
	DICCFuncName  = "DepositInChildChain"
	WFCCFuncName  = "WithdrawFromChildChain"
	WFMCFuncName  = "WithdrawFromMainChain"
	SB2MCFuncName = "SaveBlockToMainChain"

	// Create Child Chain Parameters
	CCC_ARGS_FROM                = "from"
	CCC_ARGS_CHAINID             = "chainId"
	CCC_ARGS_VALIDATOR_THRESHOLD = "validatorThreshold"
	CCC_ARGS_TOKEN_THRESHOLD     = "tokenThreshold"
	CCC_ARGS_START_BLOCK         = "startBlock"
	CCC_ARGS_END_BLOCK           = "endBlock"

	// Join Child Chain Parameters
	JCC_ARGS_FROM    = "from"
	JCC_ARGS_PUBKEY  = "pubkey"
	JCC_ARGS_CHAINID = "chainId"
	JCC_ARGS_DEPOSIT = "depositAmount"

	// Deposit In Main Chain Parameters
	DIMC_ARGS_FROM    = "from"
	DIMC_ARGS_CHAINID = "chainId"
	DIMC_ARGS_AMOUNT  = "amount"

	// Deposit In Child Chain Parameters
	DICC_ARGS_FROM    = "from"
	DICC_ARGS_TXHASH  = "txHash"
	DICC_ARGS_CHAINID = "chainId" // internal

	// Withdraw From Child Chain Parameters
	WFCC_ARGS_FROM   = "from"
	WFCC_ARGS_AMOUNT = "amount"

	// Withdraw From Main Chain Parameters
	WFMC_ARGS_FROM    = "from"
	WFMC_ARGS_CHAINID = "chainId"
	WFMC_ARGS_TXHASH  = "txHash"

	// Save Block To Main Chain Parameters
	SB2MCFuncName_ARGS_FROM  = "from"
	SB2MCFuncName_ARGS_BLOCK = "block"
)

type PublicChainAPI struct {
	am *accounts.Manager
	b  Backend
	//Client Client
}

// NewPublicChainAPI creates a new Etheruem protocol API.
func NewPublicChainAPI(b Backend) *PublicChainAPI {
	return &PublicChainAPI{
		am: b.AccountManager(),
		b:  b,
		//Client: b.Client(),
	}
}

func (s *PublicChainAPI) CreateChildChain(ctx context.Context, from common.Address, chainId string,
	minValidators uint16, minDepositAmount *big.Int, startBlock, endBlock uint64, gas *big.Int, gasPrice *big.Int) (common.Hash, error) {

	if chainId == "" || strings.Contains(chainId, ";") {
		return common.Hash{}, errors.New("chainId is nil or empty, or contains ';', should be meaningful")
	}

	params := types.MakeKeyValueSet()
	params.Set(CCC_ARGS_FROM, from)
	params.Set(CCC_ARGS_CHAINID, chainId)
	params.Set(CCC_ARGS_VALIDATOR_THRESHOLD, minValidators)
	params.Set(CCC_ARGS_TOKEN_THRESHOLD, minDepositAmount)
	params.Set(CCC_ARGS_START_BLOCK, startBlock)
	params.Set(CCC_ARGS_END_BLOCK, endBlock)

	etd := &types.ExtendTxData{
		FuncName: CCCFuncName,
		Params:   params,
	}

	args := SendTxArgs{
		From:         from,
		To:           nil,
		Gas:          (*hexutil.Big)(gas),
		GasPrice:     (*hexutil.Big)(gasPrice),
		Value:        nil,
		Data:         nil,
		Nonce:        nil,
		Type:         nil,
		ExtendTxData: etd,
	}

	return s.b.GetInnerAPIBridge().SendTransaction(ctx, args)
}

func (s *PublicChainAPI) JoinChildChain(ctx context.Context, from common.Address,
	pubkey string, chainId string, depositAmount *big.Int, gas *big.Int, gasPrice *big.Int) (common.Hash, error) {

	if chainId == "" || strings.Contains(chainId, ";") {
		return common.Hash{}, errors.New("chainId is nil or empty, or contains ';', should be meaningful")
	}

	params := types.MakeKeyValueSet()
	params.Set(JCC_ARGS_FROM, from)
	params.Set(JCC_ARGS_PUBKEY, pubkey)
	params.Set(JCC_ARGS_CHAINID, chainId)
	params.Set(JCC_ARGS_DEPOSIT, depositAmount)

	etd := &types.ExtendTxData{
		FuncName: JCCFuncName,
		Params:   params,
	}

	args := SendTxArgs{
		From:         from,
		To:           nil,
		Gas:          (*hexutil.Big)(gas),
		GasPrice:     (*hexutil.Big)(gasPrice),
		Value:        nil,
		Data:         nil,
		Nonce:        nil,
		Type:         nil,
		ExtendTxData: etd,
	}

	return s.b.GetInnerAPIBridge().SendTransaction(ctx, args)
}

func (s *PublicChainAPI) DepositInMainChain(ctx context.Context, from common.Address,
	chainId string, amount *big.Int, gas *big.Int, gasPrice *big.Int) (common.Hash, error) {

	if chainId == "" || strings.Contains(chainId, ";") {
		return common.Hash{}, errors.New("chainId is nil or empty, or contains ';', should be meaningful")
	}

	if chainId == "pchain" {
		return common.Hash{}, errors.New("chainId should not be \"pchain\"")
	}

	if s.b.ChainConfig().PChainId != "pchain" {
		return common.Hash{}, errors.New("this api can only be called in main chain - pchain")
	}

	params := types.MakeKeyValueSet()
	params.Set(DIMC_ARGS_FROM, from)
	params.Set(DIMC_ARGS_CHAINID, chainId)
	params.Set(DIMC_ARGS_AMOUNT, amount)

	etd := &types.ExtendTxData{
		FuncName: DIMCFuncName,
		Params:   params,
	}

	args := SendTxArgs{
		From:         from,
		To:           nil,
		Gas:          (*hexutil.Big)(gas),
		GasPrice:     (*hexutil.Big)(gasPrice),
		Value:        nil,
		Data:         nil,
		Nonce:        nil,
		Type:         nil,
		ExtendTxData: etd,
	}

	return s.b.GetInnerAPIBridge().SendTransaction(ctx, args)
}

func (s *PublicChainAPI) DepositInChildChain(ctx context.Context, from common.Address,
	txHash common.Hash, gas *big.Int, gasPrice *big.Int) (common.Hash, error) {

	chainId := s.b.ChainConfig().PChainId
	if chainId == "pchain" {
		return common.Hash{}, errors.New("this api can only be called in child chain")
	}

	params := types.MakeKeyValueSet()
	params.Set(DICC_ARGS_FROM, from)
	params.Set(DICC_ARGS_TXHASH, txHash)
	params.Set(DICC_ARGS_CHAINID, chainId)

	etd := &types.ExtendTxData{
		FuncName: DICCFuncName,
		Params:   params,
	}

	args := SendTxArgs{
		From:         from,
		To:           nil,
		Gas:          (*hexutil.Big)(gas),
		GasPrice:     (*hexutil.Big)(gasPrice),
		Value:        nil,
		Data:         nil,
		Nonce:        nil,
		Type:         nil,
		ExtendTxData: etd,
	}

	return s.b.GetInnerAPIBridge().SendTransaction(ctx, args)
}

func (s *PublicChainAPI) WithdrawFromChildChain(ctx context.Context, from common.Address,
	amount *big.Int, gas *big.Int, gasPrice *big.Int) (common.Hash, error) {

	chainId := s.b.ChainConfig().PChainId
	if chainId == "pchain" {
		return common.Hash{}, errors.New("this api can only be called in child chain")
	}

	params := types.MakeKeyValueSet()
	params.Set(WFCC_ARGS_FROM, from)
	params.Set(WFCC_ARGS_AMOUNT, amount)

	etd := &types.ExtendTxData{
		FuncName: WFCCFuncName,
		Params:   params,
	}

	args := SendTxArgs{
		From:         from,
		To:           nil,
		Gas:          (*hexutil.Big)(gas),
		GasPrice:     (*hexutil.Big)(gasPrice),
		Value:        nil,
		Data:         nil,
		Nonce:        nil,
		Type:         nil,
		ExtendTxData: etd,
	}

	return s.b.GetInnerAPIBridge().SendTransaction(ctx, args)
}

func (s *PublicChainAPI) WithdrawFromMainChain(ctx context.Context, from common.Address,
	chainId string, txHash common.Hash, gas *big.Int, gasPrice *big.Int) (common.Hash, error) {

	if chainId == "pchain" {
		return common.Hash{}, errors.New("argument can't be the main chain - pchain")
	}

	params := types.MakeKeyValueSet()
	params.Set(WFMC_ARGS_FROM, from)
	params.Set(WFMC_ARGS_CHAINID, chainId)
	params.Set(WFMC_ARGS_TXHASH, txHash)

	etd := &types.ExtendTxData{
		FuncName: WFMCFuncName,
		Params:   params,
	}

	args := SendTxArgs{
		From:         from,
		To:           nil,
		Gas:          (*hexutil.Big)(gas),
		GasPrice:     (*hexutil.Big)(gasPrice),
		Value:        nil,
		Data:         nil,
		Nonce:        nil,
		Type:         nil,
		ExtendTxData: etd,
	}

	return s.b.GetInnerAPIBridge().SendTransaction(ctx, args)
}

func (s *PublicChainAPI) SaveBlockToMainChain(ctx context.Context, from common.Address,
	block []byte) (common.Hash, error) {

	localChainId := s.b.ChainConfig().PChainId
	if localChainId != "pchain" {
		return common.Hash{}, errors.New("this api can only be called in main chain")
	}

	params := types.MakeKeyValueSet()
	params.Set(SB2MCFuncName_ARGS_FROM, from)
	params.Set(SB2MCFuncName_ARGS_BLOCK, block)

	etd := &types.ExtendTxData{
		FuncName: SB2MCFuncName,
		Params:   params,
	}

	args := SendTxArgs{
		From:         from,
		To:           nil,
		Gas:          nil,
		GasPrice:     nil,
		Value:        nil,
		Data:         nil,
		Nonce:        nil,
		Type:         nil,
		ExtendTxData: etd,
	}

	return s.b.GetInnerAPIBridge().SendTransaction(ctx, args)
}

func init() {
	//CreateChildChain
	core.RegisterValidateCb(CCCFuncName, ccc_ValidateCb)
	core.RegisterApplyCb(CCCFuncName, ccc_ApplyCb)

	//JoinChildChain
	core.RegisterValidateCb(JCCFuncName, jcc_ValidateCb)
	core.RegisterApplyCb(JCCFuncName, jcc_ApplyCb)

	//DepositInMainChain
	core.RegisterValidateCb(DIMCFuncName, dimc_ValidateCb)
	core.RegisterApplyCb(DIMCFuncName, dimc_ApplyCb)

	//DepositInChildChain
	core.RegisterValidateCb(DICCFuncName, dicc_ValidateCb)
	core.RegisterApplyCb(DICCFuncName, dicc_ApplyCb)

	//WithdrawFromChildChain
	core.RegisterValidateCb(WFCCFuncName, wfcc_ValidateCb)
	core.RegisterApplyCb(WFCCFuncName, wfcc_ApplyCb)

	//WithdrawFromMainChain
	core.RegisterValidateCb(WFMCFuncName, wfmc_ValidateCb)
	core.RegisterApplyCb(WFMCFuncName, wfmc_ApplyCb)

	//SB2MCFuncName
	core.RegisterValidateCb(SB2MCFuncName, sb2mc_ValidateCb)
	core.RegisterApplyCb(SB2MCFuncName, sb2mc_ApplyCb)
}

func ccc_ValidateCb(tx *types.Transaction, state *st.StateDB, cch core.CrossChainHelper) error {
	etd := tx.ExtendTxData()

	fromVar, _ := etd.Params.Get(CCC_ARGS_FROM)
	from := fromVar.(common.Address)
	chainIdVar, _ := etd.Params.Get(CCC_ARGS_CHAINID)
	chainId := chainIdVar.(string)
	minValidatorsVar, _ := etd.Params.Get(CCC_ARGS_VALIDATOR_THRESHOLD)
	minValidators := minValidatorsVar.(uint16)
	minDepositAmountVar, _ := etd.Params.Get(CCC_ARGS_TOKEN_THRESHOLD)
	minDepositAmount := minDepositAmountVar.(*big.Int)
	startBlockVar, _ := etd.Params.Get(CCC_ARGS_START_BLOCK)
	startBlock := startBlockVar.(uint64)
	endBlockVar, _ := etd.Params.Get(CCC_ARGS_END_BLOCK)
	endBlock := endBlockVar.(uint64)

	err := cch.CanCreateChildChain(from, chainId, minValidators, minDepositAmount, startBlock, endBlock)
	if err != nil {
		return err
	}

	return nil
}

func ccc_ApplyCb(tx *types.Transaction, state *st.StateDB, cch core.CrossChainHelper) error {
	etd := tx.ExtendTxData()

	fromVar, _ := etd.Params.Get(CCC_ARGS_FROM)
	from := common.BytesToAddress(fromVar.([]byte))
	chainIdVar, _ := etd.Params.Get(CCC_ARGS_CHAINID)
	chainId := string(chainIdVar.([]byte))

	minValidatorsVar, _ := etd.Params.Get(CCC_ARGS_VALIDATOR_THRESHOLD)
	minValidators := convertByteSliceToUint16(minValidatorsVar.([]byte))

	minDepositAmountVar, _ := etd.Params.Get(CCC_ARGS_TOKEN_THRESHOLD)
	minDepositAmount := new(big.Int).SetBytes(minDepositAmountVar.([]byte))

	startBlockVar, _ := etd.Params.Get(CCC_ARGS_START_BLOCK)
	startBlock := convertByteSliceToUint64(startBlockVar.([]byte))

	endBlockVar, _ := etd.Params.Get(CCC_ARGS_END_BLOCK)
	endBlock := convertByteSliceToUint64(endBlockVar.([]byte))

	err := cch.CreateChildChain(from, chainId, minValidators, minDepositAmount, startBlock, endBlock)
	if err != nil {
		return err
	}

	return nil
}

func jcc_ValidateCb(tx *types.Transaction, state *st.StateDB, cch core.CrossChainHelper) error {
	etd := tx.ExtendTxData()

	fromVar, _ := etd.Params.Get(JCC_ARGS_FROM)
	from := fromVar.(common.Address)
	pubkeyVar, _ := etd.Params.Get(JCC_ARGS_PUBKEY)
	pubkey := pubkeyVar.(string)
	chainIdVar, _ := etd.Params.Get(JCC_ARGS_CHAINID)
	chainId := chainIdVar.(string)
	depositAmountVar, _ := etd.Params.Get(JCC_ARGS_DEPOSIT)
	depositAmount := depositAmountVar.(*big.Int)

	// Check Balance
	if state.GetBalance(from).Cmp(depositAmount) == -1 {
		return core.ErrBalance
	}

	if err := cch.ValidateJoinChildChain(from, pubkey, chainId, depositAmount); err != nil {
		return err
	}

	return nil
}

func jcc_ApplyCb(tx *types.Transaction, state *st.StateDB, cch core.CrossChainHelper) error {
	etd := tx.ExtendTxData()

	fromVar, _ := etd.Params.Get(JCC_ARGS_FROM)
	from := common.BytesToAddress(fromVar.([]byte))
	pubkeyVar, _ := etd.Params.Get(JCC_ARGS_PUBKEY)
	pubkey := string(pubkeyVar.([]byte))
	chainIdVar, _ := etd.Params.Get(JCC_ARGS_CHAINID)
	chainId := string(chainIdVar.([]byte))
	depositAmountVar, _ := etd.Params.Get(JCC_ARGS_DEPOSIT)
	depositAmount := new(big.Int).SetBytes(depositAmountVar.([]byte))

	// Check Balance
	if state.GetBalance(from).Cmp(depositAmount) == -1 {
		return core.ErrBalance
	}
	// Add the validator into Chain DB
	err := cch.JoinChildChain(from, pubkey, chainId, depositAmount)
	if err != nil {
		return err
	} else {
		// Everything fine, Lock the Balance for this account
		state.SubBalance(from, depositAmount)
		state.AddLockedBalance(from, depositAmount)
	}

	return nil
}

func dimc_ValidateCb(tx *types.Transaction, state *st.StateDB, cch core.CrossChainHelper) error {
	etd := tx.ExtendTxData()

	fromInt, _ := etd.Params.Get(DIMC_ARGS_FROM)
	from := fromInt.(common.Address)
	amountInt, _ := etd.Params.Get(DIMC_ARGS_AMOUNT)
	amount := amountInt.(*big.Int)

	if state.GetBalance(from).Cmp(amount) < 0 {
		return errors.New(fmt.Sprintf("%x has no enough balance for deposit", from))
	}

	return nil
}

func dimc_ApplyCb(tx *types.Transaction, state *st.StateDB, cch core.CrossChainHelper) error {
	etd := tx.ExtendTxData()

	fromInt, _ := etd.Params.Get(DIMC_ARGS_FROM)
	from := common.BytesToAddress(fromInt.([]byte))
	chainIdInt, _ := etd.Params.Get(DIMC_ARGS_CHAINID)
	chainId := string(chainIdInt.([]byte))
	amountInt, _ := etd.Params.Get(DIMC_ARGS_AMOUNT)
	amount := new(big.Int).SetBytes(amountInt.([]byte))

	chainInfo := core.GetChainInfo(cch.GetChainInfoDB(), chainId)
	state.SubBalance(from, amount)
	state.AddChainBalance(chainInfo.Owner, amount)

	// record this cross chain tx
	cch.RecordCrossChainTx(from, tx.Hash())

	return nil
}

func dicc_ValidateCb(tx *types.Transaction, state *st.StateDB, cch core.CrossChainHelper) error {
	etd := tx.ExtendTxData()

	fromInt, _ := etd.Params.Get(DICC_ARGS_FROM)
	from := fromInt.(common.Address)
	chainIdInt, _ := etd.Params.Get(DICC_ARGS_CHAINID)
	chainId := chainIdInt.(string)
	txHashInt, _ := etd.Params.Get(DICC_ARGS_TXHASH)
	txHash := txHashInt.(common.Hash)

	mainTx := cch.GetTxFromMainChain(txHash)
	if mainTx == nil {
		return errors.New(fmt.Sprintf("tx %x does not exist in main chain", txHash))
	}

	if !cch.VerifyCrossChainTx(txHash) {
		return errors.New(fmt.Sprintf("tx %x already been used", txHash))
	}

	mainEtd := mainTx.ExtendTxData()
	if mainEtd == nil || mainEtd.FuncName != DIMCFuncName {
		return errors.New(fmt.Sprintf("not expected tx %s", mainEtd))
	}

	mainFromInt, _ := mainEtd.Params.Get(DIMC_ARGS_FROM)
	mainFrom := common.BytesToAddress(mainFromInt.([]byte))
	mainChainIdInt, _ := mainEtd.Params.Get(DIMC_ARGS_CHAINID)
	mainChainId := string(mainChainIdInt.([]byte))

	if mainFrom != from || mainChainId != chainId {
		return errors.New("params are not consistent with tx in main chain")
	}

	return nil
}

func dicc_ApplyCb(tx *types.Transaction, state *st.StateDB, cch core.CrossChainHelper) error {
	etd := tx.ExtendTxData()

	txHashInt, _ := etd.Params.Get(DICC_ARGS_TXHASH)
	txHash := common.BytesToHash(txHashInt.([]byte))

	mainTx := cch.GetTxFromMainChain(txHash)
	if mainTx == nil {
		return errors.New(fmt.Sprintf("tx %x does not exist in main chain", txHash))
	}

	mainEtd := mainTx.ExtendTxData()
	if mainEtd == nil || mainEtd.FuncName != DIMCFuncName {
		return errors.New(fmt.Sprintf("not expected tx %s", mainEtd))
	}

	mainFromInt, _ := mainEtd.Params.Get(DIMC_ARGS_FROM)
	mainFrom := common.BytesToAddress(mainFromInt.([]byte))
	mainAmountInt, _ := mainEtd.Params.Get(DIMC_ARGS_AMOUNT)
	mainAmount := new(big.Int).SetBytes(mainAmountInt.([]byte))

	// delete this cross chain tx
	cch.DeleteCrossChainTx(txHash)

	state.AddBalance(mainFrom, mainAmount)

	return nil
}

func wfcc_ValidateCb(tx *types.Transaction, state *st.StateDB, cch core.CrossChainHelper) error {
	etd := tx.ExtendTxData()

	fromInt, _ := etd.Params.Get(WFCC_ARGS_FROM)
	from := fromInt.(common.Address)
	amountInt, _ := etd.Params.Get(WFCC_ARGS_AMOUNT)
	amount := amountInt.(*big.Int)

	if state.GetBalance(from).Cmp(amount) < 0 {
		return errors.New("no enough balance to withdraw")
	}

	return nil
}

func wfcc_ApplyCb(tx *types.Transaction, state *st.StateDB, cch core.CrossChainHelper) error {
	etd := tx.ExtendTxData()

	fromInt, _ := etd.Params.Get(WFCC_ARGS_FROM)
	from := common.BytesToAddress(fromInt.([]byte))
	amountInt, _ := etd.Params.Get(WFCC_ARGS_AMOUNT)
	amount := new(big.Int).SetBytes(amountInt.([]byte))

	if state.GetBalance(from).Cmp(amount) < 0 {
		return errors.New("no enough balance to withdraw")
	}

	state.SubBalance(from, amount)

	// record this cross chain tx
	cch.RecordCrossChainTx(from, tx.Hash())

	return nil
}

func wfmc_ValidateCb(tx *types.Transaction, state *st.StateDB, cch core.CrossChainHelper) error {
	etd := tx.ExtendTxData()

	fromInt, _ := etd.Params.Get(WFMC_ARGS_FROM)
	from := fromInt.(common.Address)
	chainIdInt, _ := etd.Params.Get(WFMC_ARGS_CHAINID)
	chainId := chainIdInt.(string)
	txHashInt, _ := etd.Params.Get(WFMC_ARGS_TXHASH)
	txHash := txHashInt.(common.Hash)

	childTx := cch.GetTxFromChildChain(txHash, chainId)
	if childTx == nil {
		return errors.New(fmt.Sprintf("tx %x does not exist in child chain %s", txHash, chainId))
	}

	if !cch.VerifyCrossChainTx(txHash) {
		return errors.New(fmt.Sprintf("tx %x already been used", txHash))
	}

	childEtd := childTx.ExtendTxData()
	if childEtd == nil || childEtd.FuncName != WFCCFuncName {
		return errors.New(fmt.Sprintf("not expected tx %s", childEtd))
	}

	childFromInt, _ := childEtd.Params.Get(WFCC_ARGS_FROM)
	childFrom := common.BytesToAddress(childFromInt.([]byte))
	childAmountInt, _ := childEtd.Params.Get(WFCC_ARGS_AMOUNT)
	childAmount := new(big.Int).SetBytes(childAmountInt.([]byte))

	if childFrom != from {
		return errors.New("params are not consistent with tx in child chain")
	}

	chainInfo := core.GetChainInfo(cch.GetChainInfoDB(), chainId)
	if state.GetChainBalance(chainInfo.Owner).Cmp(childAmount) < 0 {
		return errors.New("no enough balance to withdraw")
	}

	return nil
}

func wfmc_ApplyCb(tx *types.Transaction, state *st.StateDB, cch core.CrossChainHelper) error {
	etd := tx.ExtendTxData()

	fromInt, _ := etd.Params.Get(WFMC_ARGS_FROM)
	from := common.BytesToAddress(fromInt.([]byte))
	chainIdInt, _ := etd.Params.Get(WFMC_ARGS_CHAINID)
	chainId := string(chainIdInt.([]byte))
	txHashInt, _ := etd.Params.Get(WFMC_ARGS_TXHASH)
	txHash := common.BytesToHash(txHashInt.([]byte))

	childTx := cch.GetTxFromChildChain(txHash, chainId)
	if childTx == nil {
		return errors.New(fmt.Sprintf("tx %x does not exist in child chain %s", txHash, chainId))
	}

	childEtd := childTx.ExtendTxData()
	if childEtd == nil || childEtd.FuncName != WFCCFuncName {
		return errors.New(fmt.Sprintf("not expected tx %s", childEtd))
	}

	childFromInt, _ := childEtd.Params.Get(WFCC_ARGS_FROM)
	childFrom := common.BytesToAddress(childFromInt.([]byte))
	childAmountInt, _ := childEtd.Params.Get(WFCC_ARGS_AMOUNT)
	childAmount := new(big.Int).SetBytes(childAmountInt.([]byte))

	if childFrom != from {
		return errors.New("params are not consistent with tx in child chain")
	}

	chainInfo := core.GetChainInfo(cch.GetChainInfoDB(), chainId)
	chainOwner := chainInfo.Owner
	if state.GetChainBalance(chainOwner).Cmp(childAmount) < 0 {
		return errors.New("no enough balance to withdraw")
	}

	// delete this cross chain tx
	cch.DeleteCrossChainTx(txHash)

	state.SubChainBalance(chainOwner, childAmount)
	state.AddBalance(from, childAmount)

	return nil
}

func sb2mc_ValidateCb(tx *types.Transaction, state *st.StateDB, cch core.CrossChainHelper) error {
	etd := tx.ExtendTxData()

	fromInt, _ := etd.Params.Get(SB2MCFuncName_ARGS_FROM)
	from := fromInt.(common.Address)
	blockInt, _ := etd.Params.Get(SB2MCFuncName_ARGS_BLOCK)
	block := blockInt.([]byte)

	err := cch.VerifyTdmBlock(from, block)
	if err != nil {
		return errors.New("block does not pass verification")
	}

	return nil
}

func sb2mc_ApplyCb(tx *types.Transaction, state *st.StateDB, cch core.CrossChainHelper) error {
	etd := tx.ExtendTxData()

	//fromInt, _ := etd.Params.Get(SB2MCFuncName_ARGS_FROM)
	//from := common.BytesToAddress(fromInt.([]byte))
	blockInt, _ := etd.Params.Get(SB2MCFuncName_ARGS_BLOCK)
	block := blockInt.([]byte)

	return cch.SaveTdmBlock2MainBlock(block)
}

// ---------------------------------------------
// Utility Func
const (
	UINT64_BYTE_SIZE = 8
	UINT16_BYTE_SIZE = 2
)

// Convert the Byte Slice to uint64
func convertByteSliceToUint64(input []byte) uint64 {
	result := make([]byte, UINT64_BYTE_SIZE)

	l := UINT64_BYTE_SIZE - len(input)
	copy(result[l:], input)

	return binary.BigEndian.Uint64(result)
}

// Convert the Byte Slice to uint16
func convertByteSliceToUint16(input []byte) uint16 {
	result := make([]byte, UINT16_BYTE_SIZE)

	l := UINT16_BYTE_SIZE - len(input)
	copy(result[l:], input)

	return binary.BigEndian.Uint16(result)
}