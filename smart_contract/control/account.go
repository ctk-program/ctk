package control

import (
	`encoding/json`
	`fmt`
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
	`github.com/hyperledger/fabric/smart_contract/common`
	"github.com/hyperledger/fabric/smart_contract/log"
	`github.com/hyperledger/fabric/smart_contract/module`
	`strings`
)

/**
* Account Info
**/
func (t *CTKChainCode) AccountInfo(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	log.Logger.Info("##############account info###############")
	if len(args) < 1 {
		return shim.Error("params error!")
	}
	var param module.AccountInfo
	err := json.Unmarshal([]byte(args[0]), &param)
	if err != nil {
		return shim.Error(" params err:" + err.Error())
	}
	account := t.GetTokenAccountInfo(stub,param.Token,param.Account)
	if account == nil{
		return shim.Success(common.FormatSuccessResult(map[string]interface{}{
			"account":map[string]interface{}{
				"address":param.Account,
				"balance":"0",
			},
		}))
	}
	return shim.Success(common.FormatSuccessResult(map[string]interface{}{
		"account":account,
	}))
}

/**
* Account Txs
**/
func (t *CTKChainCode) Transactions(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	log.Logger.Info("################# search transaction list  ############################")
	if len(args) < 1 {
		log.Logger.Error("QueryTransferByType params error。")
		return shim.Error("params error")
	}
	var param module.TransactionsParam
	err := json.Unmarshal([]byte(args[0]), &param)
	if err != nil {
		log.Logger.Error("args json unmarshal err")
		return shim.Error("unmarshal error")
	}
	param.Account = strings.ToLower(param.Account)
	
	key := module.DB_TRANS_TO+module.DB_TRANSACTIONS+param.Token+param.Account
	fmt.Println(key)
	historyIter, err := stub.GetHistoryForKey(key)
	var list []map[string]interface{}
	i := 0

	if err == nil {
		fmt.Println("=======",historyIter.HasNext(),err)
		for historyIter.HasNext(){
			modification, err := historyIter.Next()
			if err != nil {
				return shim.Error(err.Error())
			}
			fmt.Println("history receive tx", string(modification.Value))
			var v map[string]interface{}
			_ = json.Unmarshal(modification.Value,&v)
			list = append(list,v)
			i++
			if i > 1000{
				break
			}
		}
		
	}
	i = 0
	key = module.DB_TRANS_FROM+module.DB_TRANSACTIONS+param.Token+param.Account
	fmt.Println(key)
	historyIter, err = stub.GetHistoryForKey(key)

	if err != nil {
		return shim.Success(common.FormatSuccessResult(list))
	}
	fmt.Println("=======",historyIter.HasNext(),err)
	for historyIter.HasNext(){
		modification, err := historyIter.Next()
		fmt.Println("=======",modification)
		if err != nil {
			return shim.Error(err.Error())
		}
		fmt.Println("history spend tx", string(modification.Value))
		var v map[string]interface{}
		_ = json.Unmarshal(modification.Value,&v)
		list = append(list,v)
		i++
		if i > 1000{
			break
		}
	}
	return shim.Success(common.FormatSuccessResult(list))
}
