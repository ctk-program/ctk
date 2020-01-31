package module

import (
	`encoding/json`
	`github.com/hyperledger/fabric/core/chaincode/shim`
	`github.com/hyperledger/fabric/smart_contract/common`
)

type Account struct {

	Address  string `json:"address"`

	Balance  string `json:"balance"`

	Role  int `json:"role"`

	PoolAddr string `json:"pool_addr"`

}

func (this *Account)Bytes() []byte {
	b,_ := json.Marshal(*this)
	return b
}

func (this *Account)Save(stub shim.ChaincodeStubInterface,token string)  {
	_ = stub.PutState(common.DB_ACCOUNT_INFO_PREFIX+token+this.Address,this.Bytes())
}
