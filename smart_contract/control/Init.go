package control

import (
	`encoding/json`
	`fmt`
	`math/rand`
	`github.com/hyperledger/fabric/smart_contract/common`
	`github.com/hyperledger/fabric/smart_contract/decimal`
	`github.com/hyperledger/fabric/smart_contract/log`
	`github.com/hyperledger/fabric/smart_contract/module`
	`time`
	
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
)

type CTKChainCode struct {
	rand					 *rand.Rand
	dateZone				 *time.Location
	chargeLowLimit		   decimal.Decimal
	chargeDestroy			decimal.Decimal
	chargeDestroyUpperLimit  decimal.Decimal
	tokenCharge			  decimal.Decimal
	exchangeCharge		   decimal.Decimal
	superMinerTokenAwardAll  decimal.Decimal
	normalMinerTokenAwardAll decimal.Decimal
	tokenChargeDestroy	   decimal.Decimal
	mainToken 					 module.TokenConfig
	cid						string
	blackDestroyDateKey		string
	minersFeeDateKey		string
	addressPrefix	 string
	mortgageChaincode string
	candyChaincode	string
	maxPrecision	  int64
	maxCurrency	   string
	configKey		 string
	mineralKey		string
	blackKey		  string
	grossKey		  string
	outputKey		 string
	rateKey		   string
	dayRateKey		string
	precisionKey	  string
	tokenAwardKey	 string
	tokenTransPrefix  string
	mortgPrefix	   string
	candyPrefix	   string
	agentPayPrefix	  string
	feeAccount		float64
	storageChainCode  string
	storageFeeAccount string
	tokenMortAmount   string
	tokenMortPre	  string
	tokenMortDays	 string
	tokenMortRate	 string
	tokenSpareNodeAccount string
	tokenFoundAccount string
	minerFeeAccount	   string
}

func (self *CTKChainCode) GetTokenAccountKey(token,address string) string {
	return common.DB_ACCOUNT_INFO_PREFIX+token+address
}

func (self *CTKChainCode) GetPoolNodeAccounts(stub shim.ChaincodeStubInterface) module.NodeAccounts {
	poolAccounts , _ := stub.GetState(common.DB_POOL_ACCOUNTS)
	if len(poolAccounts) < 1{
		log.Logger.Debug("don't have any pool accounts")
		return module.NodeAccounts{}
	}
	var r module.NodeAccounts
	_ = json.Unmarshal(poolAccounts,&r)
	return r
}

func (self *CTKChainCode) RegisterNodeAccounts(stub shim.ChaincodeStubInterface,role int,data []byte) {
	db := ""
	switch role {
	case common.NODE_ROLE_SUPER:
		db = common.DB_SUPER_ACCOUNTS
	case common.NODE_ROLE_POOL:
		db = common.DB_POOL_ACCOUNTS
	case common.NODE_ROLE_NORMAL:
		db = common.DB_NORMAL_ACCOUNTS
	}
	_ = stub.PutState(db,data)
}

func (self *CTKChainCode) GetNormalNodeAccounts(stub shim.ChaincodeStubInterface) module.NodeAccounts {
	normalAccounts , _ := stub.GetState(common.DB_NORMAL_ACCOUNTS)
	if len(normalAccounts) < 1{
		log.Logger.Debug("don't have any normal accounts")
		return module.NodeAccounts{}
	}
	var r module.NodeAccounts
	_ = json.Unmarshal(normalAccounts,&r)
	return r
}

func (self *CTKChainCode) GetSuperNodeAccounts(stub shim.ChaincodeStubInterface) module.NodeAccounts {
	superAccounts , err := stub.GetState(common.DB_SUPER_ACCOUNTS)
	if len(superAccounts) < 1{
		log.Logger.Debug("Don't have any super accounts")
		return module.NodeAccounts{}
	}
	var r module.NodeAccounts
	err = json.Unmarshal(superAccounts,&r)
	if err != nil{
		return module.NodeAccounts{}
	}
	return r
}

func (self *CTKChainCode) GetTokenAccountInfo(stub shim.ChaincodeStubInterface, token,address string) *module.Account{
	
	return module.GetTokenAccountInfo(stub,token,address)
}


func (self *CTKChainCode) GetTokenInfo(stub shim.ChaincodeStubInterface, token string) *module.TokenConfig{
	return module.GetTokenInfo(stub,token)
}


func (self *CTKChainCode) GetNeedFreezeCTK(stub shim.ChaincodeStubInterface, tokenConfig *module.TokenConfig,role int) int{
	if tokenConfig == nil{
		return 0
	}
	switch role {
	case common.NODE_ROLE_NORMAL:
		return tokenConfig.NormalNodeFreezeCTK
	case common.NODE_ROLE_SUPER:
		return tokenConfig.SuperNodeFreezeCTK
	case common.NODE_ROLE_POOL:
		return tokenConfig.PoolNodeFreezeCTK
	}
	return 0
}

func (self *CTKChainCode) Init(stub shim.ChaincodeStubInterface) peer.Response {
	token := self.mainToken
	accountInfo := &module.Account{}
	accountInfo.Address = token.OriginatorAccount
	accountInfo.Balance = token.AllCount
	accountInfo.Role = common.NODE_ROLE_SUPER_SUPER
	accountInfo.Save(stub,token.Name)
	module.AddCTK(stub,token.AllCount)
	self.mainToken.CreateTime = self.GetCurrentDateString()
	self.mainToken.Save(stub)
	return shim.Success(nil)
}

func (self *CTKChainCode) GetCurrentDateString() string {
	return time.Now().In(time.FixedZone("CST", 28800)).Format("2006-01-02")
}


func (self *CTKChainCode) FormatNumber(number decimal.Decimal, precision int64) string {
	return number.StringFixed(int32(precision))
}


func (self *CTKChainCode) TryGetDecimal(number string) decimal.Decimal {
	return common.TryGetDecimal(number)
}


func (token *CTKChainCode) InitGenesis() {
	fmt.Println("==============init genesis==============")
	token.rand = rand.New(rand.NewSource(time.Now().Unix()))
	mainConfig := module.TokenConfig{}
	mainConfig.Name = common.DEFAULT_TOKEN
	mainConfig.Title = common.DEFAULT_TOKEN
	mainConfig.OriginatorAccount = "0xe53617F9e0f76E8B63B9741F5f670C6d81c86e96"
	mainConfig.AllCount = "24498550000"
	mainConfig.AllCostFee = "0"
	mainConfig.Precision = 6
	mainConfig.Site = ""
	mainConfig.Logo = ""
	mainConfig.Email = ""
	mainConfig.NormalNodeFreezeCTK = 300000
	mainConfig.SuperNodeFreezeCTK = 10000000
	mainConfig.PoolNodeFreezeCTK = 3000000
	token.mainToken = mainConfig
	token.addressPrefix = "0x"
	token.cid = "sxlchannel"
	token.blackDestroyDateKey = "black_destroy_info"
	token.minersFeeDateKey = "miners_fee_info"
	token.mortgageChaincode = "ctk-mortg"
	token.maxPrecision = 10
	token.maxCurrency = "1000000000000"
	token.configKey = common.TOKEN_CONFIG_PREFIX
	token.mineralKey = token.addressPrefix+common.MINER_FEE_ACCOUNT
	token.blackKey = token.addressPrefix+common.BLACK_KEY
	token.minerFeeAccount =  token.addressPrefix + common.MINER_FEE_ACCOUNT
	token.tokenChargeDestroy = token.TryGetDecimal("30")
	token.chargeDestroyUpperLimit = token.TryGetDecimal("24477550000")
	token.chargeLowLimit = token.TryGetDecimal("1")
	token.tokenCharge = token.TryGetDecimal("3000")
	token.tokenMortAmount = "10000"
	token.tokenMortPre = common.MORT_PREFIX
	token.tokenMortDays = "30"
	token.tokenMortRate = "0.1"
	token.normalMinerTokenAwardAll = token.TryGetDecimal("150")
	token.superMinerTokenAwardAll = token.TryGetDecimal("150")
	token.grossKey = "_gross"
	token.outputKey = "_output"
	token.rateKey = "_rate"
	token.dayRateKey = "_dayrate"
	token.precisionKey = "_precision"
	token.tokenAwardKey = "_award_"
	token.tokenTransPrefix = token.addressPrefix+common.TOKEN_TRANSFER_PREFIX
	token.mortgPrefix = token.addressPrefix+common.MORTG_PREFIX
	token.candyPrefix = token.addressPrefix+common.CANDY_PREFIX
	token.agentPayPrefix = token.addressPrefix+common.AGENT_PAY_PREFIX
	token.feeAccount = 0.01
	token.dateZone = time.FixedZone("CST", 28800)
}
