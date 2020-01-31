package module

import (
	`encoding/json`
	`github.com/hyperledger/fabric/core/chaincode/shim`
	`github.com/hyperledger/fabric/smart_contract/common`
)

type MortageToken map[string]string

func (this *MortageToken) Bytes() []byte {
	b,_ := json.Marshal(*this)
	return b
}


type TokenInfo struct {
	Creator string `json:"creator"`
	CC                string `json:"cc"`
	Name              string `json:"name"`
	Desc              string `json:"desc"`
	Logo              string `json:"logo"`
	Total             string `json:"total"`
	Award             string `json:"award"`
	Balance           string `json:"balance"`
	Decimal           string `json:"decimal"`
	Mineral           string `json:"mineral"`
	Url               string `json:"url"`
	Email             string `json:"email"`
	PublishTime       string `json:"publish_time"`
}

const DB_ALL_NODES_FEE = "all_nodes_fee"

// get all nodes fee
func GetAllNodesFee( stub shim.ChaincodeStubInterface) string {
	b,err := stub.GetState(DB_ALL_NODES_FEE)
	allFee := "0"
	if err == nil{
		allFee = string(b)
		if allFee == ""{
			return "0"
		}
	}
	return allFee
}

//add nodes fee
func AddNodesFee( stub shim.ChaincodeStubInterface,fee string)  {
	allFee := GetAllNodesFee(stub)
	allFee = common.TryGetDecimal(fee).Add(common.TryGetDecimal(allFee)).String()
	_ = stub.PutState(DB_ALL_NODES_FEE,[]byte(allFee))
}

const DB_ALL_POOL_FEE = "all_pool_fee"

// get all pool fee
func GetAllPoolFee( stub shim.ChaincodeStubInterface) string {
	b,err := stub.GetState(DB_ALL_POOL_FEE)
	allFee := "0"
	if err == nil{
		allFee = string(b)
		if allFee == ""{
			return "0"
		}
	}
	return allFee
}

//add pool fee
func AddPoolFee( stub shim.ChaincodeStubInterface,fee string)  {
	allFee := GetAllNodesFee(stub)
	allFee = common.TryGetDecimal(fee).Add(common.TryGetDecimal(allFee)).String()
	_ = stub.PutState(DB_ALL_POOL_FEE,[]byte(allFee))
}

const DB_ALL_SUPER_FEE = "all_super_fee"

// get all pool fee
func GetAllSuperFee( stub shim.ChaincodeStubInterface) string {
	b,err := stub.GetState(DB_ALL_SUPER_FEE)
	allFee := "0"
	if err == nil{
		allFee = string(b)
		if allFee == ""{
			allFee = "0"
		}
	}
	return allFee
}

//add pool fee
func AddSuperFee( stub shim.ChaincodeStubInterface,fee string)  {
	allFee := GetAllNodesFee(stub)
	allFee = common.TryGetDecimal(fee).Add(common.TryGetDecimal(allFee)).String()
	_ = stub.PutState(DB_ALL_SUPER_FEE,[]byte(allFee))
}

const DB_ALL_MORT_USDT = "all_mrot_usdt"

// get all usdt fee
func GetAllUSDT( stub shim.ChaincodeStubInterface) string {
	b,err := stub.GetState(DB_ALL_MORT_USDT)
	allFee := "0"
	if err == nil{
		allFee = string(b)
		if allFee == ""{
			allFee = "0"
		}
	}
	return allFee
}

//add usdt fee
func AddMortUSDT( stub shim.ChaincodeStubInterface,fee string)  {
	allFee := GetAllUSDT(stub)
	allFee = common.TryGetDecimal(fee).Add(common.TryGetDecimal(allFee)).String()
	_ = stub.PutState(DB_ALL_MORT_USDT,[]byte(allFee))
}

//reduce usdt fee
func ReduceMortUSDT( stub shim.ChaincodeStubInterface,fee string)  {
	allFee := GetAllUSDT(stub)
	allFee = common.TryGetDecimal(fee).Sub(common.TryGetDecimal(allFee)).String()
	_ = stub.PutState(DB_ALL_MORT_USDT,[]byte(allFee))
}

const DB_ALL_CTK = "all_ctk"

// get all usdt fee
func GetAllCTK( stub shim.ChaincodeStubInterface) string {
	b,err := stub.GetState(DB_ALL_MORT_USDT)
	allFee := "0"
	if err == nil{
		allFee = string(b)
		if allFee == ""{
			allFee = "0"
		}
	}
	return allFee
}

//add ctk fee
func AddCTK( stub shim.ChaincodeStubInterface,fee string)  {
	allFee := GetAllUSDT(stub)
	allFee = common.TryGetDecimal(fee).Add(common.TryGetDecimal(allFee)).String()
	_ = stub.PutState(DB_ALL_CTK,[]byte(allFee))
}

//reduce usdt fee
func ReduceCTK( stub shim.ChaincodeStubInterface,fee string)  {
	allFee := GetAllUSDT(stub)
	allFee = common.TryGetDecimal(fee).Sub(common.TryGetDecimal(allFee)).String()
	_ = stub.PutState(DB_ALL_CTK,[]byte(allFee))
}

const DB_ALL_WITHDRAW_CTK = "all_withdraw_ctk"

// get all usdt fee
func GetAllWithdrawCTK( stub shim.ChaincodeStubInterface) string {
	b,err := stub.GetState(DB_ALL_WITHDRAW_CTK)
	allFee := "0"
	if err == nil{
		allFee = string(b)
		if allFee == ""{
			allFee = "0"
		}
	}
	return allFee
}

//add withdraw ctk fee
func AddWithdrawCTK( stub shim.ChaincodeStubInterface,fee string)  {
	allFee := GetAllUSDT(stub)
	allFee = common.TryGetDecimal(fee).Add(common.TryGetDecimal(allFee)).String()
	_ = stub.PutState(DB_ALL_WITHDRAW_CTK,[]byte(allFee))
}