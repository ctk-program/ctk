package module

import (
	`encoding/json`
	`github.com/hyperledger/fabric/core/chaincode/shim`
	`github.com/hyperledger/fabric/smart_contract/common`
	`github.com/hyperledger/fabric/smart_contract/log`
	`strings`
)

type NodeAccounts []string

func (this *NodeAccounts) Has( addr string) bool {
	for _,v := range *this{
		if v == addr{
			return true
		}
	}
	return false
}

func (this *NodeAccounts) Add( addr string) {
	*this = append(*this,addr)
}

func (this *NodeAccounts) Remove( addr string) {
	for i :=0;i< len(*this);i++{
		if (*this)[i] == addr{
			if i == 0{
				*this = make([]string,0)
				return
			}
			if i == len(*this) - 1{
				*this = (*this)[:i]
				return
			} else{
				*this = append((*this)[:i], (*this)[i+1:]...)
				return
			}
		}
	}
}

func (this *NodeAccounts) Bytes() []byte {
	b,_ := json.Marshal(*this)
	return b
}

func GetSuperNodeAccounts(stub shim.ChaincodeStubInterface) NodeAccounts {
	superAccounts , _ := stub.GetState(common.DB_SUPER_ACCOUNTS)
	if len(superAccounts) < 1{
		log.Logger.Debug("don't have any super accounts")
		return NodeAccounts{}
	}
	var r NodeAccounts
	_ = json.Unmarshal(superAccounts,&r)
	return r
}

func GetNormalNodeAccounts(stub shim.ChaincodeStubInterface) NodeAccounts {
	normalAccounts , _ := stub.GetState(common.DB_NORMAL_ACCOUNTS)
	if len(normalAccounts) < 1{
		log.Logger.Debug("don't have any normal accounts")
		return NodeAccounts{}
	}
	var r NodeAccounts
	_ = json.Unmarshal(normalAccounts,&r)
	return r
}

func GetTokenAccountKey(token,address string) string {
	return common.DB_ACCOUNT_INFO_PREFIX+token+address
}

func GetTokenAccountInfo(stub shim.ChaincodeStubInterface, token,address string) *Account{
	address = strings.ToLower(address)
	accountBytes, err := stub.GetState(GetTokenAccountKey(token,address))
	
	if len(accountBytes) < 1{
		log.Logger.Error("account not exist!")
		return nil
	}
	var account Account
	err = json.Unmarshal(accountBytes,&account)
	if err != nil{
		log.Logger.Error(err.Error())
		return nil
	}
	return &account
}