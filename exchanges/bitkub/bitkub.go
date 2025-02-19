package bitkub

import (
	exchange "github.com/thrasher-corp/gocryptotrader/exchanges"
)

// Bitkub is the overarching type across this package
type Bitkub struct {
	exchange.Base
}

const (
	bitkubAPIURL     = "https://api.bitkub.com/api"
	bitkubAPIVersion = "3"

	// Public endpoints
	bitkubTicker     = "market/ticker"
	bitkubBids       = "market/bids"
	bitkubAsks       = "market/asks"
	bitkubDepth      = "market/depth"
	bitkubTrades     = "market/trades"
	bitkubServerTime = "market/servertime"

	// Authenticated endpoints

	// User endpoints
	bitkubTradingCredits = "user/trading-credits" // #nosec G101
	bitkubLimits         = "user/limits"

	// Market endpoints
	bitkubBalances         = "market/balances"
	bitkubPlaceBid         = "market/place-bid"
	bitkubPlaceAsk         = "market/place-ask"
	bitkubPCancelOrder     = "market/cancel-order"
	bitkubWsToken          = "market/wstoken"
	bitkubMyOpenOrders     = "market/my-open-orders"
	bitkubWsMyOrderHistory = "market/my-order-history"
	bitkubOrderInfo        = "market/order-info"

	// Crypto endpoints
	bitkubInternalWithdraw      = "crypto/internal-withdraw"
	bitkubAddresses             = "crypto/addresses"
	bitkubCryptoWithdraw        = "crypto/withdraw"
	bitkubCryptoDepositHistory  = "crypto/deposit-history"
	bitkubCryptoWithdrawHistory = "crypto/withdraw-history"
	bitkubGenerateAddress       = "crypto/generate-address"

	// Fiat endpoints
	bitkubFiatAccounts        = "fiat/accounts"
	bitkubFiatWithdraw        = "fiat/withdraw"
	bitkubFiatDepositHistory  = "fiat/deposit-history"
	bitkubFiatWithdrawHistory = "fiat/withdraw-history"
)

// Start implementing public and private exchange API funcs below
