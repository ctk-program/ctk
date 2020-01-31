package module

type RegisterParam struct {
	Account  string `json:"account"`
	Role int `json:"role"`
	Sign string `json:"sign"`
	Token string `json:"token"`
}

type ExchangeParam struct {
	Account  string `json:"account"`
	Sign string `json:"sign"`
	Amount string `json:"amount"`
}


type BindPoolParam struct {
	Account  string `json:"account"`
	Pool  string `json:"pool"`
	Sign string `json:"sign"`
}

type AccountInfo struct {
	Account  string `json:"account"`
	Token string `json:"token"`
}

type AppendSuperParam struct {
	Amount  string `json:"amount"`
	Account  string `json:"account"`
	Sign string `json:"sign"`
}

type BlockAwardParam struct {
	StartHeight  int `json:"start_height"`
	EndHeight  int `json:"end_height"`
	Sign string `json:"sign"`
	Fee float64 `json:"fee"`
}

type QuitSuperParam struct {
	Account  string `json:"account"`
	Sign string `json:"sign"`
}

type Super17Param struct {
	Accounts  NodeAccounts `json:"accounts"`
	Sign string `json:"sign"`
}

type TransferParam struct {
	From string `json:"from"`
	To string `json:"to"`
	Fee string `json:"fee"`
	Amount string `json:"amount"`
	Token string `json:"token"`
	Sign string `json:"sign"`
}

type TransactionsParam struct {
	Offset int `json:"offset"`
	Limit int `json:"limit"`
	Token string   `json:"token"`
	Account string `json:"account"`
} 