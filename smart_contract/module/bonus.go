package module

import (
	`github.com/hyperledger/fabric/core/chaincode/shim`
	`github.com/hyperledger/fabric/smart_contract/common`
	`github.com/hyperledger/fabric/smart_contract/decimal`
	`time`
)

// node bonus when normal node register or pool node register
func NodeBonus(stub shim.ChaincodeStubInterface,freezeMoney decimal.Decimal,token string,from string,isRegister bool) (txs []Transaction){
	supers := GetSuperNodeAccounts(stub)
	nodes := GetNormalNodeAccounts(stub)
	everySuperFee := freezeMoney.Mul(decimal.NewFromFloat(0.4)).Div(decimal.NewFromInt(31))
	for i := 0;i<len(supers);i++{
		a := supers[i]
		accountInfo := GetTokenAccountInfo(stub,token,a)
		if accountInfo == nil{
			continue
		}
		if !isRegister && i > 16{
			if accountInfo.PoolAddr == ""{
				continue
			}
			poolInfo := GetTokenAccountInfo(stub,token,accountInfo.PoolAddr)
			if poolInfo == nil || poolInfo.Role != common.NODE_ROLE_POOL{
				continue
			}
			//pool award ***************
			b := common.TryGetDecimal(poolInfo.Balance).Add(everySuperFee.Mul(decimal.NewFromFloat(0.1)))
			poolInfo.Balance = b.String()
			poolInfo.Save(stub,token)
		}
		
		
		b := common.TryGetDecimal(accountInfo.Balance).Add(everySuperFee.Mul(decimal.NewFromFloat(0.9)))
		accountInfo.Balance = b.String()
		accountInfo.Save(stub,token)
		
		trans := Transaction{}
		trans.To = accountInfo.Address
		trans.From = from
		trans.Type = common.TRANSACTION_TYPE_BONUS_SUPER
		trans.Amount = everySuperFee.String()
		trans.Memo = "super node award"
		trans.Fee = "0"
		trans.CreateTime = time.Now().Unix()
		trans.Save(stub,token)
		txs = append(txs,trans)

	}
	if len(nodes) > 14{
		nodesCount := decimal.NewFromInt32(int32(len(nodes)-14))
		everyNodesFee := freezeMoney.Mul(decimal.NewFromFloat(0.5)).Div(nodesCount)
		for _,a := range nodes{
			if supers.Has(a){
				continue
			}
			accountInfo := GetTokenAccountInfo(stub,token,a)
			if accountInfo == nil{
				continue
			}
			if !isRegister{
				if accountInfo.PoolAddr == ""{
					continue
				}
				poolInfo := GetTokenAccountInfo(stub,token,accountInfo.PoolAddr)
				if poolInfo == nil || poolInfo.Role != common.NODE_ROLE_POOL{
					continue
				}
				//pool award ***************
				b := common.TryGetDecimal(poolInfo.Balance).Add(everySuperFee.Mul(decimal.NewFromFloat(0.1)))
				poolInfo.Balance = b.String()
				poolInfo.Save(stub,token)
			}
			
			b := common.TryGetDecimal(accountInfo.Balance).Add(everyNodesFee.Mul(decimal.NewFromFloat(0.9)))
			accountInfo.Balance = b.String()
			accountInfo.Save(stub,token)
			
			trans := Transaction{}
			trans.To = accountInfo.Address
			trans.From = from
			trans.Type = common.TRANSACTION_TYPE_BONUS_NORMAL
			trans.Amount = everyNodesFee.String()
			trans.Memo = "normal node award"
			trans.Fee = "0"
			trans.CreateTime = time.Now().Unix()
			trans.Save(stub,token)
			txs = append(txs,trans)
		}
	}
	
	//black account
	blackAccountCTK := common.BLACK_KEY
	blackBalance := freezeMoney.Mul(decimal.NewFromFloat(0.1))
	accountInfo := GetTokenAccountInfo(stub,token,blackAccountCTK)
	if accountInfo == nil{
		accountInfo = &Account{}
		accountInfo.Balance = "0"
		accountInfo.Role = common.NODE_ROLE_USER
		accountInfo.Address = blackAccountCTK
	}
	accountInfo.Balance = common.TryGetDecimal(accountInfo.Balance).Add(blackBalance).String()
	accountInfo.Save(stub,token)
	return
}