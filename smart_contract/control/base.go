/**
 *	Author Fanxu(746439274@qq.com)
 */

package control

import (
	`fmt`
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
	`github.com/hyperledger/fabric/smart_contract/log`
	"reflect"
	`strings`
)

//define controller maps
type ControllerMapsType map[string]reflect.Value

var ControllerMaps = make(ControllerMapsType, 0)

func Init() {
	var routers CTKChainCode
	vf := reflect.ValueOf(&routers)
	vft := vf.Type()
	mNum := vf.NumMethod()
	for i := 0; i < mNum; i++ {
		mName := vft.Method(i).Name
		log.Logger.Debug("index:", i, " MethodName:", mName)
		ControllerMaps[mName] = vf.Method(i) //<<<
	}
}

func Execute(fun string, stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if _, ok := ControllerMaps[fun]; !ok {
		return shim.Error("Invalid invoke function method : " + fun)
	}
	params := make([]reflect.Value, 2)
	params[0] = reflect.ValueOf(stub)
	if len(args) > 0 {
		args[0] = strings.ReplaceAll(args[0],"`",`"`)
		fmt.Println("===========",fun,"=========",args[0])
	}
	params[1] = reflect.ValueOf(args)
	resp := ControllerMaps[fun].Call(params)
	return resp[0].Interface().(peer.Response)
}