package main

import (
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/smart_contract/control"
	`github.com/hyperledger/fabric/smart_contract/log`
)

func main() {
	control.Init()
	token := new(control.CTKChainCode)
	token.InitGenesis()
	err := shim.Start(token)
	if err != nil {
		log.Logger.Errorf("Error starting CTKChainCode: %s",err.Error())
	}
}
