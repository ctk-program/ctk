package control

import (
	`encoding/hex`
	`encoding/json`
	`fmt`
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
	`github.com/hyperledger/fabric/smart_contract/common`
	"github.com/hyperledger/fabric/smart_contract/log"
	`github.com/hyperledger/fabric/smart_contract/module`
	`sort`
	`strings`
)

/**
* Vote Super account
**/
func (t *CTKChainCode) Vote(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	log.Logger.Info("##############vote for 31 node account###############")
	//auto vote the 31 super supernodes from super node 17 main 14 second

	var mortsAccounts module.MortageToken

	b,err := stub.GetState(common.DB_MORT_RECORD)
	if err != nil{
		return shim.Error("not have enough super nodes !")
	}
	if len(b) <= 0 {
		return shim.Error("mort accounts not enough!")
	}
	_ = json.Unmarshal(b,&mortsAccounts)
	if len(mortsAccounts) < 1{
		return shim.Error("not have enough super nodes !")
	}
	accountsMap := make([]module.Account,0)
	for account,balance := range mortsAccounts{
		accountInfo := t.GetTokenAccountInfo(stub,common.DEFAULT_TOKEN,account)
		if accountInfo == nil{
			continue
		}
		accountInfo.Balance = balance
		accountsMap = append(accountsMap,*accountInfo)
	}
	if len(accountsMap) > 1{
		// sort by balance
		sort.Slice(accountsMap, func(i, j int) bool {
			return t.TryGetDecimal(accountsMap[i].Balance).Cmp(t.TryGetDecimal(accountsMap[j].Balance)) > 0
		})
	}
	supersuperAccounts := module.NodeAccounts{}
	b,_ = stub.GetState(common.DB_17_SUPER_NODES)
	var super17Accounts module.NodeAccounts
	if len(b) > 0 {
		_ = json.Unmarshal(b,&super17Accounts)
		supersuperAccounts = append(super17Accounts,supersuperAccounts...)
	}
	// select 31 super super nodes
	for i := 0;i<len(accountsMap);i++{
		if supersuperAccounts.Has(accountsMap[i].Address){
			continue
		}
		supersuperAccounts = append(supersuperAccounts,accountsMap[i].Address)
		if len(supersuperAccounts) >= 31{
			break
		}
	}
	
	fmt.Println("============supersuperAccounts===========",supersuperAccounts)
	t.RegisterNodeAccounts(stub,common.NODE_ROLE_SUPER,supersuperAccounts.Bytes())

	return shim.Success(common.FormatSuccessResult(supersuperAccounts))
}

/**
* Vote Super 17 account
**/
func (t *CTKChainCode) VoteSuperNodes(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	log.Logger.Info("##############register 17 super node account###############")
	if len(args) < 1 {
		return shim.Error("params error!")
	}
	tokenConfig := t.GetTokenInfo(stub,common.DEFAULT_TOKEN)
	
	var super17accounts module.Super17Param
	err := json.Unmarshal([]byte(args[0]), &super17accounts)
	if err != nil {
		return shim.Error("params err:" + err.Error())
	}
	if len(super17accounts.Accounts) != 17{
		return shim.Error("super nodes not enough")
	}
	b,_ := json.Marshal(super17accounts.Accounts)
	messageStr := fmt.Sprintf("%s",hex.EncodeToString(b) )
	messageBody := hex.EncodeToString([]byte(strings.ToLower(messageStr)))
	fmt.Println(messageBody,"tokenConfig.OriginatorAccount",tokenConfig.OriginatorAccount,super17accounts.Sign)
	if !common.CheckSign(tokenConfig.OriginatorAccount,messageBody,super17accounts.Sign){
		return shim.Error("check sign failed!" )
	}
	for _,addr := range super17accounts.Accounts{
		info := module.GetTokenAccountInfo(stub,common.DEFAULT_TOKEN,addr)
		if info == nil{
			info = &module.Account{}
			info.Balance = "0"
			info.Address = addr
			info.Role = common.NODE_ROLE_SUPER
			info.Save(stub,common.DEFAULT_TOKEN)
		}
	}
	_ = stub.PutState(common.DB_17_SUPER_NODES,super17accounts.Accounts.Bytes())
	return shim.Success([]byte("success"))
}
