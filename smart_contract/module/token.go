package module

import (
	`encoding/json`
	`github.com/hyperledger/fabric/core/chaincode/shim`
	`github.com/hyperledger/fabric/smart_contract/common`
	`github.com/hyperledger/fabric/smart_contract/log`
)

type TokenConfig struct {
	OriginatorAccount string `json:"originator_account"`
	Name			  string `json:"name"`
	Title			 string `json:"title"`
	Logo			  string `json:"logo"`
	Precision		 int64 `json:"precision"`
	Site			  string `json:"site"`
	Email			 string `json:"email"`
	CreateTime string `json:"create_time"`
	TokenCharge int `json:"token_charge"` //publish token need the ctk count
	NormalNodeFreezeCTK int `json:"normal_node_freeze_ctk"`
	SuperNodeFreezeCTK int `json:"super_node_freeze_ctk"`
	PoolNodeFreezeCTK int `json:"pool_node_freeze_ctk"`
	AllCostFee string `json:"all_cost_fee"`
	AllCount string `json:"all_count"`
}

func GetTokenInfo(stub shim.ChaincodeStubInterface, token string) *TokenConfig{
	tokenBytes, err := stub.GetState(token+common.TOKEN_CONFIG_PREFIX)
	
	if len(tokenBytes) < 1{
		log.Logger.Error("token not exist!")
		return nil
	}
	var config TokenConfig
	err = json.Unmarshal(tokenBytes,&config)
	if err != nil{
		log.Logger.Error(err.Error())
		return nil
	}
	return &config
}

func (this *TokenConfig) Byte() []byte{
	b,_ := json.Marshal(*this)
	return b
}

func (this *TokenConfig) Save(stub shim.ChaincodeStubInterface) {
	_ = stub.PutState(this.Name+common.TOKEN_CONFIG_PREFIX,this.Byte())
}