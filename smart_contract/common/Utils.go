package common

import (
	`encoding/hex`
	`encoding/json`
	`fmt`
	`github.com/hyperledger/fabric/smart_contract/crypto`
	`github.com/hyperledger/fabric/smart_contract/decimal`
	`github.com/hyperledger/fabric/smart_contract/log`
	`regexp`
	"strings"
)

/**
 * ios_bbbbbbbb -> IosBbbbbbbbb
 */
func StrFirstToUpper(str string) string {
	temp := strings.Split(str, "_")
	var upperStr string
	for y := 0; y < len(temp); y++ {
		vv := []rune(temp[y])
		for i := 0; i < len(vv); i++ {
			if i == 0 {
				vv[i] -= 32
				upperStr += string(vv[i]) // + string(vv[i+1])
			} else {
				upperStr += string(vv[i])
			}
		}
	}
	return upperStr
}

//check account format
func CheckAccountFormat(account string) bool {
	match, err := regexp.MatchString(`^`+ADDR_PREFIX+`[A-Za-z0-9]{40}$`, account)
	if err != nil {
		return false
	}
	return match
}


//check account black
func CheckAccountDisabled(account string) bool {
	if _,ok := BlackAddrs[account];ok{
		return true
	}
	if len(account) != 42 {
		return true
	}
	accountPrefix := string(account[0:10])
	if _,ok := BlackAddrPres[accountPrefix];ok {
		return true
	}
	return false
}

//check sign
func CheckSign(pubaddress, message, signedAddress string) bool {
	if signedAddress == "" || pubaddress =="" || message == ""{
		return false
	}
	var err error
	var addr = strings.ToLower(pubaddress)
	
	addr = string([]byte(addr)[2:])
	var msg = crypto.Keccak256([]byte(message))
	sign, err := hex.DecodeString(strings.ToLower(signedAddress))
	if err != nil {
		log.Logger.Errorf("\n sign decode error: %s", err)
		return false
	}
	recoveredPub, err := crypto.Ecrecover(msg, sign)
	if err != nil {
		log.Logger.Errorf("\n ECRecover error: %s", err)
		return false
	}
	pubKey,err := crypto.UnmarshalPubkey(recoveredPub)
	if err != nil {
		log.Logger.Errorf("\n UnmarshalPubkey error: %s", err)
		return false
	}
	recoveredAddr := crypto.PubkeyToAddress(*pubKey)
	recoveredAddrstr := hex.EncodeToString(recoveredAddr[:])
	if addr != recoveredAddrstr {
		log.Logger.Errorf("\n sign of account addres is", "0x"+recoveredAddrstr, "not is", "0x"+addr)
		return false
	} else {
		return true
	}
}


//====================standard response================
type CTKResponse struct {
	Code int `json:"code"`
	Msg string `json:"msg"`
	Data interface{} `json:"data"`
}
var StandardResponse = CTKResponse{200,"ok", map[string]interface{}{}}

//format shim result
func FormatResult(r interface{}) []byte {
	b,err := json.Marshal(r)
	if err != nil{
		return []byte(err.Error())
	}
	return b
}
//format error result
func FormatErrorResult(msg string) []byte {
	StandardResponse.Code = 500
	StandardResponse.Msg = msg
	StandardResponse.Data = map[string]interface{}{}
	return FormatResult(StandardResponse)
}
//format success result
func FormatSuccessResult(data interface{}) []byte {
	StandardResponse.Code = 200
	StandardResponse.Msg = "ok"
	StandardResponse.Data = data
	return FormatResult(StandardResponse)
}
//====================standard response================
//get decimal from string
func TryGetDecimal(number string) decimal.Decimal {
	newDecimal, err := decimal.NewDecimalFromString(number)
	if err != nil {
		fmt.Println("==================== Error NewDecimalFromString:" + number,"================")
		return decimal.NewFromFloat(0)
	}
	return newDecimal
}
