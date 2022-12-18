package handler

// Handler manage http contract
type Handler struct {
	WalletHandler     WalletHandler
	AccountHandler    AccountHandler
	DepositHandler    DepositHandler
	WithdrawalHandler WithdrawalHandler
}
