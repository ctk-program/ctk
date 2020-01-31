package module

import (
	`encoding/json`
	`fmt`
	`github.com/hyperledger/fabric/core/chaincode/shim`
	`github.com/hyperledger/fabric/smart_contract/common`
)

// block award
type BlockAwardLog struct {
	StartHeight int 	`json:"start_height"`
	EndHeight int 	`json:"end_height"`
	Transactions []Transaction `json:"transactions"`
	CreateTime int64 `json:"create_time"`
}

func (this *BlockAwardLog) Bytes() []byte {
	b,_ := json.Marshal(*this)
	return b
}

func (this *BlockAwardLog) Save(stub shim.ChaincodeStubInterface,token string) {
	fmt.Println("==============save================",this.StartHeight,this.EndHeight,string(this.Bytes()))
	_ = stub.PutState(common.DB_BLOCK_REWARD_HISTORY,this.Bytes())
}