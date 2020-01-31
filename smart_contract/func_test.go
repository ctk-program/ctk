package main

import (
	`fmt`
	`github.com/stretchr/testify/assert`
	`github.com/hyperledger/fabric/smart_contract/common`
	`github.com/hyperledger/fabric/smart_contract/module`
	`testing`
)

func TestDecimal(t *testing.T) {
	allFee := common.TryGetDecimal(common.TryGetDecimal("24498550000").String()).Add(common.TryGetDecimal("0")).String()
	assert.Equal(t,"24498550000",allFee)
}


func TestMortAccounts(t *testing.T) {
	accounts := module.NodeAccounts{}
	accounts.Add("abc")
	accounts.Add("bcd")
	fmt.Println(accounts)
	accounts.Remove("bcd")
	accounts.Remove("abc")
	fmt.Println(accounts)
}
