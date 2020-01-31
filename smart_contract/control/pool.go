package control

import (
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
	`github.com/hyperledger/fabric/smart_contract/common`
	"github.com/hyperledger/fabric/smart_contract/log"
)

/**
* Pool accounts
**/
func (t *CTKChainCode) PoolList(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	log.Logger.Info("##############pool accounts list###############")
	accounts := t.GetPoolNodeAccounts(stub)
	return shim.Success(common.FormatSuccessResult(accounts))
}