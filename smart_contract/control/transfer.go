package control

import (
	`encoding/hex`
	`encoding/json`
	`fmt`
	`github.com/hyperledger/fabric/core/chaincode/shim`
	`github.com/hyperledger/fabric/protos/peer`
	`github.com/hyperledger/fabric/smart_contract/common`
	`github.com/hyperledger/fabric/smart_contract/decimal`
	`github.com/hyperledger/fabric/smart_contract/log`
	`github.com/hyperledger/fabric/smart_contract/module`
	`strings`
)

func (self *CTKChainCode) Transfer(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	log.Logger.Info("##############Transfer ###############")
	if len(args) < 1 {
		return shim.Error("params error!")
	}
	var param module.TransferParam
	err := json.Unmarshal([]byte(args[0]), &param)
	if err != nil {
		return shim.Error("transfer params err:" + err.Error())
	}
	messageStr := fmt.Sprintf("%s%s%s%s%s",param.From,param.To,param.Token,param.Fee,param.Amount)
	messageBody := hex.EncodeToString([]byte(strings.ToLower(messageStr)))
	if !common.CheckSign(param.From,messageBody,param.Sign){
		return shim.Error("check sign failed!" )
	}
	tokenConfig := self.GetTokenInfo(stub,param.Token)
	if tokenConfig == nil{
		return shim.Error(param.Token + " config not found:")
	}
	
	ctkConfig := self.GetTokenInfo(stub,common.DEFAULT_TOKEN)
	if ctkConfig == nil{
		return shim.Error("ctk config not found:")
	}
	
	fromInfo := module.GetTokenAccountInfo(stub,param.Token,param.From)
	if fromInfo == nil{
		return shim.Error("from account not exist!" )
	}
	if tokenConfig.Name != ctkConfig.Name{
		if tokenConfig.AllCostFee != "" && common.TryGetDecimal(tokenConfig.AllCostFee).Cmp(common.TryGetDecimal(common.MORT_TOKEN)) > 0 {
			if common.TryGetDecimal(param.Fee).Cmp(decimal.NewFromFloat(common.CTK_MIN_FEE)) < 0{
				return shim.Error(fmt.Sprintf("fee can not below %dctk",common.CTK_MIN_FEE) )
			}
		}
	} else{
		if common.TryGetDecimal(param.Fee).Cmp(decimal.NewFromFloat(common.CTK_MIN_FEE)) < 0{
			return shim.Error(fmt.Sprintf("fee can not below %dctk",common.CTK_MIN_FEE) )
		}
	}
	fromCtkInfo := self.GetTokenAccountInfo(stub,common.DEFAULT_TOKEN,param.From)
	if fromCtkInfo == nil{
		return shim.Error("account CTK not enough !" )
	}
	if common.TryGetDecimal(fromCtkInfo.Balance).Cmp(common.TryGetDecimal(param.Fee)) < 0{
		return shim.Error("account CTK balance not enough !" )
	}
	
	if common.TryGetDecimal(param.Amount).Cmp(decimal.NewFromInt(0)) <= 0{
		return shim.Error("tranfer count can not less than 0!" )
	}
	
	if common.TryGetDecimal(param.Amount).Add(common.TryGetDecimal(param.Fee)).Cmp(common.TryGetDecimal(fromInfo.Balance)) > 0{
		return shim.Error("from balance not enough!" )
	}
	
	trans := module.Transaction{}
	if common.TryGetDecimal(param.Fee).Cmp(decimal.NewFromFloat(0)) > 0 {
		fromCtkInfo.Balance = common.TryGetDecimal(fromCtkInfo.Balance).Sub(common.TryGetDecimal(param.Fee)).String()
		fromCtkInfo.Save(stub,common.DEFAULT_TOKEN)
	}
	
	fromInfo.Balance = common.TryGetDecimal(fromInfo.Balance).Sub(common.TryGetDecimal(param.Amount)).Sub(common.TryGetDecimal(param.Fee)).String()
	fromInfo.Save(stub,param.Token)
	
	toInfo := module.GetTokenAccountInfo(stub,param.Token,param.To)
	if toInfo == nil{
		toInfo = &module.Account{}
		toInfo.Balance = "0"
		toInfo.Address = param.To
		toInfo.Role = common.NODE_ROLE_USER
	}
	toInfo.Balance =  common.TryGetDecimal(toInfo.Balance).Add(common.TryGetDecimal(param.Amount)).String()
	toInfo.Save(stub,param.Token)
	trans.From = param.From
	trans.To = param.To
	trans.Fee = param.Fee
	trans.Memo = "transfer token"
	trans.Amount = param.Amount
	trans.Type = common.TRANSACTION_TYPE_TRANSFER_TOKEN
	trans.Save(stub,param.Token)
	return shim.Success(common.FormatSuccessResult(map[string]interface{}{}))
}
