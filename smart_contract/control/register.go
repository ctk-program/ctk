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
	`time`
)

/**
* Register Account
**/
func (t *CTKChainCode) Register(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	log.Logger.Info("##############register node account###############")
	if len(args) < 1 {
		return shim.Error("params error!")
	}
	var registerParam module.RegisterParam
	err := json.Unmarshal([]byte(args[0]), &registerParam)
	if err != nil {
		return shim.Error("register params err:" + err.Error())
	}
	if !common.CheckAccountFormat(registerParam.Account){
		return shim.Error("account format err:" + registerParam.Account)
	}
	if common.CheckAccountDisabled(registerParam.Account){
		return shim.Error("the account is system account" )
	}
	messageStr := fmt.Sprintf("%s%s%d",registerParam.Token,registerParam.Account , registerParam.Role)
	messageBody := hex.EncodeToString([]byte(strings.ToLower(messageStr)))
	if !common.CheckSign(registerParam.Account,messageBody,registerParam.Sign){
		return shim.Error("check sign failed!" )
	}
	tokenConfig := t.GetTokenInfo(stub,registerParam.Token)
	if tokenConfig == nil{
		return shim.Error("token config failed" )
	}
	tokenAccountInfo := t.GetTokenAccountInfo(stub,registerParam.Token,registerParam.Account)
	fmt.Println(tokenAccountInfo,"==========",registerParam.Token,registerParam.Account)
	if tokenAccountInfo == nil{
		tokenAccountInfo = &module.Account{}
		tokenAccountInfo.Balance = "0"
		tokenAccountInfo.Role = registerParam.Role
		tokenAccountInfo.Address = registerParam.Account
	}

	freezeMoney := t.GetNeedFreezeCTK(stub,tokenConfig,registerParam.Role)
	if freezeMoney <= 0{
		return shim.Error("freeze token config failed" )
	}
	fmt.Println("tokenAccountInfo.Balance=====",tokenAccountInfo.Balance,"need amount",freezeMoney)
	if t.TryGetDecimal(tokenAccountInfo.Balance).Cmp(decimal.NewFromInt(int64(freezeMoney))) < 0{
		return shim.Error("account CTK not enough!" )
	}
	var accounts module.NodeAccounts
	typ := 0
	switch registerParam.Role {
	case common.NODE_ROLE_POOL:
		accounts = t.GetPoolNodeAccounts(stub)
		typ = common.TRANSACTION_TYPE_BUY_POOL
		module.AddPoolFee(stub,decimal.NewFromInt32(int32(freezeMoney)).String())
	case common.NODE_ROLE_NORMAL:
		accounts = t.GetPoolNodeAccounts(stub)
		typ = common.TRANSACTION_TYPE_BUY_NODE
		module.AddNodesFee(stub,decimal.NewFromInt32(int32(freezeMoney)).String())
	default:
		return shim.Error("unknown register role!" )
	}
	if accounts.Has(registerParam.Account){
		return shim.Success(common.FormatSuccessResult(tokenAccountInfo))
	}
	balance := t.TryGetDecimal(tokenAccountInfo.Balance).Sub(decimal.NewFromInt(int64(freezeMoney)))
	//spend ctk
	tokenAccountInfo.Role = registerParam.Role
	tokenAccountInfo.Balance = balance.StringFixed(int32(tokenConfig.Precision))
	tokenAccountInfo.Save(stub,tokenConfig.Name)
	// tx record
	accounts.Add(registerParam.Account)
	t.RegisterNodeAccounts(stub,registerParam.Role,accounts.Bytes())
	trans := module.Transaction{}
	trans.CreateTime = time.Now().Unix()
	trans.To = tokenConfig.OriginatorAccount
	trans.From = registerParam.Account
	trans.Type = typ
	trans.Amount = decimal.NewFromInt(int64(freezeMoney)).String()
	trans.Fee = "0"
	trans.Memo = fmt.Sprintf("spend %d for node",freezeMoney)
	trans.Save(stub,tokenConfig.Name)
	// assign ctk bonus to nodes
	module.NodeBonus(stub,decimal.NewFromInt32(int32(freezeMoney)),tokenConfig.Name,registerParam.Account,true)
	return shim.Success(common.FormatSuccessResult(tokenAccountInfo))
}

//append ctk to become super node
func (t *CTKChainCode) AppendSuper(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	log.Logger.Info("##############append super node ###############")
	if len(args) < 1 {
		return shim.Error("params error!")
	}
	var appendParam module.AppendSuperParam
	err := json.Unmarshal([]byte(args[0]), &appendParam)
	if err != nil {
		return shim.Error("append params err:" + err.Error())
	}

	if !common.CheckAccountFormat(appendParam.Account){
		return shim.Error("account format err:" + appendParam.Account)
	}
	if common.CheckAccountDisabled(appendParam.Account){
		return shim.Error("the account is system account" )
	}
	messageStr := fmt.Sprintf("%s%s",appendParam.Amount, appendParam.Account )
	messageBody := hex.EncodeToString([]byte(strings.ToLower(messageStr)))
	if !common.CheckSign(appendParam.Account,messageBody,appendParam.Sign){
		return shim.Error("check sign failed!" )
	}
	tokenConfig := t.GetTokenInfo(stub,common.DEFAULT_TOKEN)
	if tokenConfig == nil{
		return shim.Error("token config failed" )
	}
	tokenAccountInfo := t.GetTokenAccountInfo(stub,common.DEFAULT_TOKEN,appendParam.Account)
	if tokenAccountInfo == nil{
		return shim.Error("account not exist!" )
	}

	freezeMoney := t.GetNeedFreezeCTK(stub,tokenConfig,common.NODE_ROLE_SUPER)
	if freezeMoney <= 0{
		return shim.Error("freeze token config failed" )
	}
	appendAmount,err := decimal.NewDecimalFromString(appendParam.Amount)
	if err != nil{
		return shim.Error("append amount error" )
	}
	if appendAmount.Cmp(decimal.NewFromInt(int64(100))) < 0{
		return shim.Error("append money at least 100 CTK!" )
	}
	if t.TryGetDecimal(tokenAccountInfo.Balance).Cmp(appendAmount) < 0{
		return shim.Error("account CTK not enough!" )
	}
	accounts := t.GetSuperNodeAccounts(stub)
	fmt.Println("=========super accounts=========",accounts)
	if !accounts.Has(appendParam.Account){
		if appendAmount.Cmp(decimal.NewFromInt(int64(freezeMoney+100))) < 0{
			return shim.Error(fmt.Sprintf("account CTK not enough! at least %d",int64(freezeMoney+100)) )
		}
		accounts.Add(appendParam.Account)
		t.RegisterNodeAccounts(stub,common.NODE_ROLE_SUPER,accounts.Bytes())
	}
	balance := t.TryGetDecimal(tokenAccountInfo.Balance).Sub(appendAmount)
	tokenAccountInfo.Balance = balance.StringFixed(int32(tokenConfig.Precision))
	tokenAccountInfo.Save(stub,tokenConfig.Name)
	
	trans := module.Transaction{}
	trans.CreateTime = time.Now().Unix()
	trans.To = common.MORTG_PREFIX
	trans.From = appendParam.Account
	trans.Type = common.TRANSACTION_TYPE_MORT_SUPER
	trans.Amount = appendAmount.String()
	trans.Fee = "100"
	trans.Memo = fmt.Sprintf("freeze %d for node contain 100 CTK fee",freezeMoney)
	trans.Save(stub,common.DEFAULT_TOKEN)
	mortsRecords := module.MortageToken{}
	b,_ := stub.GetState(common.DB_MORT_RECORD)
	if len(b) > 0{
		_ = json.Unmarshal(b,&mortsRecords)
		fmt.Println("===============mortsRecords init ==============",mortsRecords)
	}

	if _,ok := mortsRecords[appendParam.Account];ok{
		allMortAmount := mortsRecords[appendParam.Account]
		newAmount := t.TryGetDecimal(allMortAmount).Add(appendAmount).Sub(decimal.NewFromInt32(int32(100)))
		mortsRecords[appendParam.Account] = newAmount.String()
		fmt.Println("===============mortsRecords exist====",mortsRecords)
	} else{
		fmt.Println("===============mortsRecords not exist====",mortsRecords)
		newAmount := appendAmount.Sub(decimal.NewFromInt32(int32(100)))
		mortsRecords[appendParam.Account] = newAmount.String()
	}

	err = stub.PutState(common.DB_MORT_RECORD,mortsRecords.Bytes())
	if err != nil{
		fmt.Println(err.Error())
		return shim.Error(err.Error())
	}
	return shim.Success([]byte("success"))
}

//quit super node
func (t *CTKChainCode) QuitSuper(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	log.Logger.Info("##############quit super node ###############")
	if len(args) < 1 {
		return shim.Error("params error!")
	}
	var quitParam module.QuitSuperParam
	err := json.Unmarshal([]byte(args[0]), &quitParam)
	if err != nil {
		return shim.Error("register params err:" + err.Error())
	}

	if !common.CheckAccountFormat(quitParam.Account){
		return shim.Error("account format err:" + quitParam.Account)
	}
	if common.CheckAccountDisabled(quitParam.Account){
		return shim.Error("the account is system account" )
	}
	messageStr := fmt.Sprintf("%s", quitParam.Account )
	messageBody := hex.EncodeToString([]byte(strings.ToLower(messageStr)))
	if !common.CheckSign(quitParam.Account,messageBody,quitParam.Sign){
		return shim.Error("check sign failed!" )
	}
	tokenConfig := t.GetTokenInfo(stub,common.DEFAULT_TOKEN)
	if tokenConfig == nil{
		return shim.Error("token config failed" )
	}
	tokenAccountInfo := t.GetTokenAccountInfo(stub,common.DEFAULT_TOKEN,quitParam.Account)
	if tokenAccountInfo == nil{
		return shim.Error("account not found!" )
	}
	accounts := t.GetSuperNodeAccounts(stub)
	if !accounts.Has(quitParam.Account){
		return shim.Error("token is not super node" )
	}
	fmt.Println("=================accounts================",accounts)
	accounts.Remove(quitParam.Account)
	mortsRecords := module.MortageToken{}
	b,err := stub.GetState(common.DB_MORT_RECORD)
	if len(b)>0{
		_ = json.Unmarshal(b,&mortsRecords)
	}
	fmt.Println("==================mortsRecords================",mortsRecords)
	needQuitAmount := "0"
	if _,ok := mortsRecords[quitParam.Account];ok{
		needQuitAmount = mortsRecords[quitParam.Account]
		delete(mortsRecords,quitParam.Account)
	}

	balance := t.TryGetDecimal(tokenAccountInfo.Balance).Add(t.TryGetDecimal(needQuitAmount))
	tokenAccountInfo.Balance = balance.StringFixed(int32(tokenConfig.Precision))
	tokenAccountInfo.Save(stub,tokenConfig.Name)
	t.RegisterNodeAccounts(stub,common.NODE_ROLE_SUPER,accounts.Bytes())
	trans := module.Transaction{}
	trans.CreateTime = time.Now().Unix()
	trans.To = quitParam.Account
	trans.From = t.storageFeeAccount
	trans.Type = common.TRANSACTION_TYPE_RETURN_SUPER
	trans.Amount = needQuitAmount
	trans.Fee = "0"
	trans.Memo = fmt.Sprintf("return %s for quit super node",needQuitAmount)
	trans.Save(stub,common.DEFAULT_TOKEN)

	_ = stub.PutState(common.DB_MORT_RECORD,mortsRecords.Bytes())
	return shim.Success([]byte("exit super node success!"))
}
