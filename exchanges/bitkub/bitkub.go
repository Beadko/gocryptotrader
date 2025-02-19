package bitkub

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	exchange "github.com/thrasher-corp/gocryptotrader/exchanges"
	"github.com/thrasher-corp/gocryptotrader/exchanges/request"
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

// SendHTTPRequest sends an unauthenticated HTTP request
func (b *Bitkub) SendHTTPRequest(ctx context.Context, ep exchange.URL, path string, result interface{}) error {
	endpoint, err := b.API.Endpoints.GetURL(ep)
	if err != nil {
		return err
	}
	item := &request.Item{
		Method:        http.MethodGet,
		Path:          endpoint + path,
		Result:        result,
		Verbose:       b.Verbose,
		HTTPDebugging: b.HTTPDebugging,
		HTTPRecording: b.HTTPRecording,
	}
	return b.SendPayload(ctx, request.Unset, func() (*request.Item, error) {
		return item, nil
	}, request.UnauthenticatedRequest)
}

// GetTicker returns ticker information
// Returns only related data if symbol is specified; otherwise return all of them
func (b *Bitkub) GetTicker(ctx context.Context, symbol string) (TickerData, error) {
	params := url.Values{}
	if symbol != "" {
		params.Set("sym", strings.ToLower(symbol))
	}

	path := fmt.Sprintf("/%s/%s?%s", bitkubAPIVersion, bitkubTicker, params.Encode())
	ticker := make(TickerData)

	return ticker, b.SendHTTPRequest(ctx, exchange.RestSpot, path, &ticker)
}
