package log

import "github.com/hyperledger/fabric/core/chaincode/shim"

var Logger *shim.ChaincodeLogger

func init() {
	Logger = shim.NewLogger("CTKChainCode")
	Logger.SetLevel(shim.LogDebug)
	Logger.IsEnabledFor(shim.LogCritical)
}
