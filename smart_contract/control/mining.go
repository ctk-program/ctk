package control

import (
	`encoding/hex`
	"encoding/json"
	`fmt`
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
	`github.com/hyperledger/fabric/smart_contract/common`
	`github.com/hyperledger/fabric/smart_contract/decimal`
	"github.com/hyperledger/fabric/smart_contract/log"
	`github.com/hyperledger/fabric/smart_contract/module`
	`strings`
)

func (t *CTKChainCode) MiningRatio(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	allUSDT := module.GetAllUSDT(stub)
	// every 100000 reduce 1
	f , _ := common.TryGetDecimal(allUSDT).Float64()
	ratiof := f / 100000.00
	ratio := int(ratiof)
	ratioCtk := 7000 - ratio
	if ratioCtk < 100{
		ratioCtk = 100
	}
	exchangeCTK := decimal.NewFromFloat(float64(ratioCtk) / 100.00)
	
	allWithdrawCTK := module.GetAllWithdrawCTK(stub)
	// every 7000000 reduce 1
	f , _ = common.TryGetDecimal(allWithdrawCTK).Float64()
	ratiof = f / 7000000.00
	ratioCtk = 7000 - int(ratiof)
	if ratioCtk < 100{
		ratioCtk = 100
	}
	canExchangUSDT := decimal.NewFromFloat(float64(ratioCtk)).Mul(decimal.NewFromInt32(100))
	
	return shim.Success(common.FormatSuccessResult(map[string]interface{}{
		"allUSDT" : allUSDT,
		"allWithdraw" : allWithdrawCTK,
		"exchangeCTK" : exchangeCTK,
		"canExchangUSDT" : canExchangUSDT,
	}))
}

/**
* Mining Exchange
**/
func (t *CTKChainCode) Exchange(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	log.Logger.Info("##############mining exchange from usdt to ctk###############")
	if len(args) < 1 {
		return shim.Error("params error!")
	}
	var exchangeParam module.ExchangeParam
	err := json.Unmarshal([]byte(args[0]), &exchangeParam)
	if err != nil {
		return shim.Error("exchange params err:" + err.Error())
	}
	if !common.CheckAccountFormat(exchangeParam.Account){
		return shim.Error("exchange format err:" + exchangeParam.Account)
	}
	if common.CheckAccountDisabled(exchangeParam.Account){
		return shim.Error("the account is system account" )
	}
	tokenConfig := t.GetTokenInfo(stub,common.DEFAULT_TOKEN)
	messageStr := fmt.Sprintf("%s%s",exchangeParam.Amount, exchangeParam.Account)
	messageBody := hex.EncodeToString([]byte(strings.ToLower(messageStr)))
	fmt.Println(tokenConfig.OriginatorAccount,"=============",messageStr,exchangeParam.Sign)
	if !common.CheckSign(tokenConfig.OriginatorAccount,messageBody,exchangeParam.Sign){
		return shim.Error("check sign failed!" )
	}
	tokenAccountInfo := t.GetTokenAccountInfo(stub,common.DEFAULT_TOKEN,exchangeParam.Account)
	if tokenAccountInfo == nil{
		tokenAccountInfo = &module.Account{}
		tokenAccountInfo.Address = exchangeParam.Account
		tokenAccountInfo.Role = common.NODE_ROLE_USER
		tokenAccountInfo.Balance = "0"
		tokenAccountInfo.Withdraw = "0"
		tokenAccountInfo.Recharge = "0"
	}
	
	allUSDT := module.GetAllUSDT(stub)
	// every 100000 reduce 1
	f , _ := common.TryGetDecimal(allUSDT).Float64()
	ratiof := f / 100000.00
	ratio := int(ratiof)
	ratioCtk := 7000 - ratio
	if ratioCtk < 100{
		ratioCtk = 100
	}
	exchangeCTK := decimal.NewFromFloat(float64(ratioCtk) / 100.00)
	canExchangCTK := common.TryGetDecimal(exchangeParam.Amount).Mul(exchangeCTK)
	balance := common.TryGetDecimal(tokenAccountInfo.Balance).Add(canExchangCTK)
	tokenAccountInfo.Balance = balance.String()
	recharge := common.TryGetDecimal(tokenAccountInfo.Recharge).Add(canExchangCTK)
	tokenAccountInfo.Recharge = recharge.String()
	tokenAccountInfo.Save(stub,common.DEFAULT_TOKEN)
	module.AddMortUSDT(stub,exchangeParam.Amount)
	module.ReduceCTK(stub,canExchangCTK.String())
	trans := module.Transaction{}
	trans.From = module.DB_ALL_CTK
	trans.Type = common.TRANSACTION_TRANSACTION_USDT_CTK
	trans.To = exchangeParam.Account
	trans.Fee = "0"
	trans.Amount = canExchangCTK.String()
	trans.Memo = "exchange from usdt to ctk"
	trans.Save(stub,common.DEFAULT_TOKEN)
	return shim.Success(common.FormatSuccessResult(map[string]interface{}{
		"exchange_ctk" : canExchangCTK.String(),
		"account" : exchangeParam.Account,
		"mort_usdt" : exchangeParam.Amount,
	}))
}

/**
* Mining Withdraw
**/
func (t *CTKChainCode) Withdraw(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	log.Logger.Info("##############mining withdraw from ctk to usdt###############")
	if len(args) < 1 {
		return shim.Error("params error!")
	}
	var exchangeParam module.ExchangeParam
	err := json.Unmarshal([]byte(args[0]), &exchangeParam)
	if err != nil {
		return shim.Error("exchange params err:" + err.Error())
	}
	if !common.CheckAccountFormat(exchangeParam.Account){
		return shim.Error("exchange format err:" + exchangeParam.Account)
	}
	if common.CheckAccountDisabled(exchangeParam.Account){
		return shim.Error("the account is system account" )
	}
	messageStr := fmt.Sprintf("%s%s",exchangeParam.Amount,exchangeParam.Account)
	tokenConfig := t.GetTokenInfo(stub,common.DEFAULT_TOKEN)
	messageBody := hex.EncodeToString([]byte(strings.ToLower(messageStr)))
	if !common.CheckSign(tokenConfig.OriginatorAccount,messageBody,exchangeParam.Sign){
		return shim.Error("check sign failed!" )
	}
	tokenAccountInfo := t.GetTokenAccountInfo(stub,common.DEFAULT_TOKEN,exchangeParam.Account)
	if tokenAccountInfo == nil{
		return shim.Error("account not exist!")
	}
	
	allWithdrawCTK := module.GetAllWithdrawCTK(stub)
	// every 7000000 reduce 1
	f , _ := common.TryGetDecimal(allWithdrawCTK).Float64()
	ratiof := f / 7000000.00
	ratioCtk := 7000 - int(ratiof)
	if ratioCtk < 100{
		ratioCtk = 100
	}
	canExchangUSDT := common.TryGetDecimal(exchangeParam.Amount).Div(decimal.NewFromFloat(float64(ratioCtk))).Mul(decimal.NewFromInt32(100))
	
	balance := common.TryGetDecimal(tokenAccountInfo.Balance).Sub(common.TryGetDecimal(exchangeParam.Amount))
	withdraw := common.TryGetDecimal(tokenAccountInfo.Withdraw).Sub(common.TryGetDecimal(exchangeParam.Amount))
	tokenAccountInfo.Balance = balance.String()
	tokenAccountInfo.Withdraw = withdraw.String()
	tokenAccountInfo.Save(stub,common.DEFAULT_TOKEN)
	module.ReduceMortUSDT(stub,canExchangUSDT.String())
	module.AddWithdrawCTK(stub,exchangeParam.Amount)
	trans := module.Transaction{}
	trans.From = module.DB_ALL_MORT_USDT
	trans.To = exchangeParam.Account
	trans.Fee = "0"
	trans.Amount = exchangeParam.Amount
	trans.Memo = "exchange from usdt to ctk"
	trans.Save(stub,common.DEFAULT_TOKEN)
	return shim.Success(common.FormatSuccessResult(map[string]interface{}{
		"exchange_usdt" : canExchangUSDT.String(),
		"account" : exchangeParam.Account,
		"withdraw_ctk" : exchangeParam.Amount,
	}))
}

/**
* Mining Bind Pool
**/
func (t *CTKChainCode) BindPool(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	log.Logger.Info("##############mining bind pool###############")
	if len(args) < 1 {
		return shim.Error("params error!")
	}
	var param module.BindPoolParam
	err := json.Unmarshal([]byte(args[0]), &param)
	if err != nil {
		return shim.Error(" params err:" + err.Error())
	}
	if !common.CheckAccountFormat(param.Account){
		return shim.Error("account format err:" + param.Account)
	}
	if common.CheckAccountDisabled(param.Account){
		return shim.Error("the account is system account" )
	}
	if !common.CheckAccountFormat(param.Pool){
		return shim.Error("pool format err:" + param.Account)
	}
	if common.CheckAccountDisabled(param.Pool){
		return shim.Error("the pool is system account" )
	}
	messageStr := fmt.Sprintf("%s%s",param.Pool, param.Account)
	messageBody := hex.EncodeToString([]byte(strings.ToLower(messageStr)))
	if !common.CheckSign(param.Account,messageBody,param.Sign){
		return shim.Error("check sign failed!" )
	}
	tokenAccountInfo := t.GetTokenAccountInfo(stub,common.DEFAULT_TOKEN,param.Account)
	if tokenAccountInfo == nil{
		return shim.Error("account not exist!")
	}
	poolAccountInfo := t.GetTokenAccountInfo(stub,common.DEFAULT_TOKEN,param.Pool)
	
	if poolAccountInfo == nil{
		return shim.Error("pool not exist!")
	}
	if poolAccountInfo.Role != common.NODE_ROLE_POOL{
		return shim.Error("pool address not found!")
	}
	if tokenAccountInfo.Address == poolAccountInfo.Address{
		return shim.Error("error pool cannot self !")
	}
	
	tokenAccountInfo.PoolAddr = param.Pool
	tokenAccountInfo.Save(stub,common.DEFAULT_TOKEN)
	return shim.Success(common.FormatSuccessResult(map[string]interface{}{
		"account":tokenAccountInfo,
	}))
}
