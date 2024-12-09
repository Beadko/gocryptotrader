package exchange

import (
	"time"

	"github.com/thrasher-corp/gocryptotrader/currency"
	"github.com/thrasher-corp/gocryptotrader/exchanges/order"
)

type Liquidations struct {
	Exchange  string
	Pair      currency.Pair
	Amount    float64
	Price     float64
	Side      order.Side
	Timestamp time.Time
}
