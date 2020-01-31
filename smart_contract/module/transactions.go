package module

import (
	`encoding/json`
	`github.com/hyperledger/fabric/core/chaincode/shim`
	`strings`
	`time`
)

//transaction
type Transaction struct {
	Type int `json:"type"`
	From string `json:"from"`
	To string `json:"to"`
	Fee string `json:"fee"`
	Amount string `json:"amount"`
	CreateTime int64 `json:"create_time"`
	Memo string `json:"memo"`
	Token string `json:"token"`
}
const DB_TRANSACTIONS = "transactions"
const DB_TRANS_FROM  = "from"
const DB_TRANS_TO  = "to"

func (this *Transaction) Bytes() []byte {
	b,_ := json.Marshal(*this)
	return b
}

func (this *Transaction) Save(stub shim.ChaincodeStubInterface,token string) {
	this.Token = token
	this.From = strings.ToLower(this.From)
	this.To = strings.ToLower(this.To)
	this.CreateTime = time.Now().Unix()
	// fmt.Println(DB_TRANS_FROM+DB_TRANSACTIONS+token+this.From,this.Memo)
	_ = stub.PutState(DB_TRANS_FROM+DB_TRANSACTIONS+token+this.From,this.Bytes())
	// fmt.Println(DB_TRANS_TO+DB_TRANSACTIONS+token+this.To,this.Memo)
	_ = stub.PutState(DB_TRANS_TO+DB_TRANSACTIONS+token+this.To,this.Bytes())
}
