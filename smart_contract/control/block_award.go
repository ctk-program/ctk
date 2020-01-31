package control

import (
	`encoding/hex`
	`encoding/json`
	`fmt`
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
	`github.com/hyperledger/fabric/smart_contract/common`
	`github.com/hyperledger/fabric/smart_contract/decimal`
	`github.com/hyperledger/fabric/smart_contract/log`
	`github.com/hyperledger/fabric/smart_contract/module`
	`strings`
	`time`
)

/**
* BlockAward
**/
func (t *CTKChainCode) BlockAward(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	log.Logger.Info("##############BlockAward ###############")
	if len(args) < 1 {
		return shim.Error("params error!")
	}
	var awardParam module.BlockAwardParam
	err := json.Unmarshal([]byte(args[0]), &awardParam)
	if err != nil {
		return shim.Error("award params err:" + err.Error())
	}
	if awardParam.StartHeight <= 0{
		awardParam.StartHeight = 1
	}
	if awardParam.StartHeight > awardParam.EndHeight{
		return shim.Error("end height can not more than start height!")
	}
	ctkToken := t.GetTokenInfo(stub,common.DEFAULT_TOKEN)
	messageStr := fmt.Sprintf("%f%d%d",awardParam.Fee,awardParam.StartHeight,awardParam.EndHeight)
	messageBody := hex.EncodeToString([]byte(strings.ToLower(messageStr)))
	if !common.CheckSign(ctkToken.OriginatorAccount,messageBody,awardParam.Sign){
		return shim.Error("check sign failed!" )
	}
	var blockRewardLog module.BlockAwardLog
	//get last block award log
	b,err := stub.GetState(common.DB_BLOCK_REWARD_HISTORY)
	if err != nil{
		fmt.Println("==========err=========",err.Error())
	}
	if len(b) > 0{
		_ = json.Unmarshal(b,&blockRewardLog)
		fmt.Println("===============blockRecord===============",blockRewardLog.StartHeight,blockRewardLog.EndHeight)
		if awardParam.StartHeight <= blockRewardLog.EndHeight || awardParam.StartHeight<= blockRewardLog.StartHeight {
			return shim.Error(fmt.Sprintf("start height wrong!need start more than %d , current is %d",(blockRewardLog.EndHeight+1),awardParam.StartHeight))
		}
	}
	
	txs := module.NodeBonus(stub,decimal.NewFromFloat(awardParam.Fee),common.DEFAULT_TOKEN,common.MINER_FEE_ACCOUNT,false)
	
	// block transaction fee
	// for i := awardParam.StartHeight;i<awardParam.EndHeight;i++{
	// 	myargs := [][]byte{[]byte("GetBlockByNumber"), []byte("sxlchannel"), []byte(strconv.FormatUint(uint64(i),10))}
	// 	response := stub.InvokeChaincode("qscc", myargs, "sxlchannel")
	//
	// 	var block fcommon.Block
	// 	_ = json.Unmarshal(response.Payload,&block)
	// 	fmt.Println("============response=============",block)
	// }
	heightGap := awardParam.StartHeight - blockRewardLog.EndHeight
	newBlockRewardLog := &module.BlockAwardLog{}
	newBlockRewardLog.StartHeight = blockRewardLog.EndHeight+1
	newBlockRewardLog.EndHeight = awardParam.EndHeight
	newBlockRewardLog.Transactions = txs
	newBlockRewardLog.CreateTime = time.Now().Unix()
	newBlockRewardLog.Save(stub,common.DEFAULT_TOKEN)
	needMaintenceFee := len(txs) * common.CTK_MAINTENANCE_FEE * heightGap
	start := time.Now().Unix()
	for needMaintenceFee > 0 {
		//Node Maintenance Fee
		for _,tx := range txs{
			accountInfo := module.GetTokenAccountInfo(stub,common.DEFAULT_TOKEN,tx.To)
			if common.TryGetDecimal(accountInfo.Balance).Cmp(decimal.NewFromInt32(common.CTK_MAINTENANCE_FEE)) > 0{
				accountInfo.Balance = common.TryGetDecimal(accountInfo.Balance).Sub(decimal.NewFromInt32(common.CTK_MAINTENANCE_FEE)).String()
				accountInfo.Save(stub,common.DEFAULT_TOKEN)
				needMaintenceFee--
				trans := module.Transaction{}
				trans.To = common.BLACK_KEY
				trans.From = tx.To
				trans.Fee = "0"
				trans.Amount = fmt.Sprintf("%d",common.CTK_MAINTENANCE_FEE)
				trans.Type = common.TRANSACTION_TRANSACTION_NODE_MAINTENCE_FEE
				trans.Memo = "node maintenance fee"
				trans.Save(stub,common.DEFAULT_TOKEN)
			}
		}
		end := time.Now().Unix()
		if end - start >= 3{
			break
		}
	}
	
	
	return shim.Success(common.FormatSuccessResult(newBlockRewardLog))
}
