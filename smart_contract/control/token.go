package control

import (
	`encoding/hex`
	`encoding/json`
	`fmt`
	`github.com/hyperledger/fabric/core/chaincode/shim`
	`github.com/hyperledger/fabric/protos/peer`
	`github.com/hyperledger/fabric/smart_contract/common`
	`github.com/hyperledger/fabric/smart_contract/log`
	`github.com/hyperledger/fabric/smart_contract/module`
	`strings`
	`time`
)

type PublishParams struct {
	Token module.TokenConfig `json:"token"`
	Sign string `json:"sign"`
} 
func (self *CTKChainCode) Publish(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	log.Logger.Info("##############Token Publish###############")
	if len(args) < 1 {
		return shim.Error("params error!")
	}
	var param PublishParams
	err := json.Unmarshal([]byte(args[0]), &param)
	if err != nil {
		return shim.Error("token publish params err:" + err.Error())
	}
	if param.Token.Name == "" || param.Token.Title == "" || param.Token.AllCount == "" {
		return shim.Error("token base info err:")
	}
	accountInfo := module.GetTokenAccountInfo(stub,common.DEFAULT_TOKEN,param.Token.OriginatorAccount)
	if accountInfo == nil{
		return shim.Error("account not exist!")
	}
	if common.TryGetDecimal(accountInfo.Balance).Cmp(common.TryGetDecimal(common.MORT_TOKEN)) < 0 {
		return shim.Error("ctk balance not enough!")
	}
	if self.GetTokenInfo(stub,param.Token.Name) != nil{
		return shim.Error("token name has exist!")
	}
	b , _ := json.Marshal(param.Token)
	messageStr := fmt.Sprintf("%s",hex.EncodeToString(b))
	messageBody := hex.EncodeToString([]byte(strings.ToLower(messageStr)))
	if !common.CheckSign(param.Token.OriginatorAccount,messageBody,param.Sign){
		return shim.Error("check sign failed!" )
	}
	param.Token.AllCostFee = "0"
	param.Token.CreateTime = self.GetCurrentDateString()
	param.Token.Save(stub)
	accountInfo.Balance = common.TryGetDecimal(accountInfo.Balance).Sub(common.TryGetDecimal(common.MORT_TOKEN)).String()
	accountInfo.Save(stub,common.DEFAULT_TOKEN)
	trans := module.Transaction{}
	trans.From = param.Token.OriginatorAccount
	trans.To = common.BLACK_KEY
	trans.Fee = "0"
	trans.CreateTime = time.Now().Unix()
	trans.Type = common.TRANSACTION_TRANSACTION_PUBLISH_TOKEN
	trans.Save(stub,common.DEFAULT_TOKEN)
	
	tokenInfo := module.GetTokenAccountInfo(stub,param.Token.Name,param.Token.OriginatorAccount)
	if tokenInfo == nil{
		tokenInfo = &module.Account{}
		tokenInfo.Balance = param.Token.AllCount
		tokenInfo.Role = common.NODE_ROLE_USER
		tokenInfo.Address = param.Token.OriginatorAccount
		tokenInfo.Save(stub,param.Token.Name)
	} else{
		tokenInfo.Balance = common.TryGetDecimal(tokenInfo.Balance).Add(common.TryGetDecimal(param.Token.AllCount)).String()
		tokenInfo.Save(stub,param.Token.Name)
	}
	trans.From = "0"
	trans.To = param.Token.OriginatorAccount
	trans.Fee = "0"
	trans.Amount = param.Token.AllCount
	trans.CreateTime = time.Now().Unix()
	trans.Type = common.TRANSACTION_TRANSACTION_PUBLISH_TOKEN
	trans.Save(stub,param.Token.Name)
	return shim.Success(common.FormatSuccessResult(map[string]interface{}{}))
}
func (self *CTKChainCode) UpdateToken(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	log.Logger.Info("##############Update Token ###############")
	if len(args) < 1 {
		return shim.Error("params error!")
	}
	var param PublishParams
	err := json.Unmarshal([]byte(args[0]), &param)
	if err != nil {
		return shim.Error("token update params err:" + err.Error())
	}
	accountInfo := module.GetTokenAccountInfo(stub,common.DEFAULT_TOKEN,param.Token.OriginatorAccount)
	if accountInfo == nil{
		return shim.Error("account not exist!")
	}
	tokenConfig := self.GetTokenInfo(stub,param.Token.Name)
	if tokenConfig == nil{
		return shim.Error("token name has not exist!")
	}
	b , _ := json.Marshal(param.Token)
	messageStr := fmt.Sprintf("%s",hex.EncodeToString(b))
	messageBody := hex.EncodeToString([]byte(strings.ToLower(messageStr)))
	if !common.CheckSign(param.Token.OriginatorAccount,messageBody,param.Sign){
		return shim.Error("check sign failed!" )
	}
	tokenConfig.PoolNodeFreezeCTK = param.Token.PoolNodeFreezeCTK
	tokenConfig.SuperNodeFreezeCTK = param.Token.SuperNodeFreezeCTK
	tokenConfig.NormalNodeFreezeCTK = param.Token.NormalNodeFreezeCTK
	tokenConfig.Save(stub)
	return shim.Success(common.FormatSuccessResult(map[string]interface{}{}))
}
