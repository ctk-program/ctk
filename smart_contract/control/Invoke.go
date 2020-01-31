package control

import (
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
	`github.com/hyperledger/fabric/smart_contract/common`
	`github.com/hyperledger/fabric/smart_contract/log`
)

func (t *CTKChainCode) Invoke(stub shim.ChaincodeStubInterface) peer.Response {
	log.Logger.Info("Invoke CTK")
	function, args := stub.GetFunctionAndParameters()
	return Execute(common.StrFirstToUpper(function), stub, args)
}
