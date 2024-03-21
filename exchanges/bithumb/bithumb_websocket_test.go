package bithumb

import (
	"errors"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/thrasher-corp/gocryptotrader/currency"
	exchange "github.com/thrasher-corp/gocryptotrader/exchanges"
	"github.com/thrasher-corp/gocryptotrader/exchanges/asset"
	"github.com/thrasher-corp/gocryptotrader/exchanges/stream"
	"github.com/thrasher-corp/gocryptotrader/exchanges/subscription"
	"github.com/thrasher-corp/gocryptotrader/exchanges/ticker"
	testexch "github.com/thrasher-corp/gocryptotrader/internal/testing/exchange"
)

var (
	wsTickerResp    = []byte(`{"type":"ticker","content":{"tickType":"24H","date":"20210811","time":"132017","openPrice":"33400","closePrice":"34010","lowPrice":"32660","highPrice":"34510","value":"45741663716.89916828275244531","volume":"1359398.496892086826189907","sellVolume":"198021.237915860451480504","buyVolume":"1161377.258976226374709403","prevClosePrice":"33530","chgRate":"1.83","chgAmt":"610","volumePower":"500","symbol":"UNI_KRW"}}`)
	wsTransResp     = []byte(`{"type":"transaction","content":{"list":[{"buySellGb":"1","contPrice":"1166","contQty":"125.2400","contAmt":"146029.8400","contDtm":"2021-08-13 15:23:42.911273","updn":"dn","symbol":"DAI_KRW"}]}}`)
	wsOrderbookResp = []byte(`{"type":"orderbookdepth","content":{"list":[{"symbol":"XLM_KRW","orderType":"ask","price":"401.2","quantity":"0","total":"0"},{"symbol":"XLM_KRW","orderType":"ask","price":"401.6","quantity":"21277.735","total":"1"},{"symbol":"XLM_KRW","orderType":"ask","price":"403.3","quantity":"4000","total":"1"},{"symbol":"XLM_KRW","orderType":"bid","price":"399.5","quantity":"0","total":"0"},{"symbol":"XLM_KRW","orderType":"bid","price":"398.2","quantity":"0","total":"0"},{"symbol":"XLM_KRW","orderType":"bid","price":"399.8","quantity":"31416.8779","total":"1"},{"symbol":"XLM_KRW","orderType":"bid","price":"398.5","quantity":"34328.387","total":"1"}],"datetime":"1628835823604483"}}`)
	// this pair is used to ensure that endpoints match it correctly
	testPairMapping = currency.NewPair(currency.BTC, currency.USDT)
)

func TestWsHandleData(t *testing.T) {
	t.Parallel()
	pairs := currency.Pairs{
		currency.Pair{
			Base:  currency.BTC,
			Quote: currency.USDT,
		},
	}

	dummy := Bithumb{
		Base: exchange.Base{
			Name: "dummy",
			Features: exchange.Features{
				Enabled: exchange.FeaturesEnabled{SaveTradeData: true},
			},
			CurrencyPairs: currency.PairsManager{
				Pairs: map[asset.Item]*currency.PairStore{
					asset.Spot: {
						Available: pairs,
						Enabled:   pairs,
						ConfigFormat: &currency.PairFormat{
							Uppercase: true,
							Delimiter: currency.DashDelimiter,
						},
					},
				},
			},
			Websocket: &stream.Websocket{
				Wg:          *new(sync.WaitGroup),
				DataHandler: make(chan interface{}, 1),
			},
		},
	}

	dummy.setupOrderbookManager()
	dummy.API.Endpoints = b.NewEndpoints()

	welcomeMsg := []byte(`{"status":"0000","resmsg":"Connected Successfully"}`)
	err := dummy.wsHandleData(welcomeMsg)
	if err != nil {
		t.Fatal(err)
	}

	err = dummy.wsHandleData([]byte(`{"status":"1336","resmsg":"Failed"}`))
	if !errors.Is(err, stream.ErrSubscriptionFailure) {
		t.Fatalf("received: %v but expected: %v",
			err,
			stream.ErrSubscriptionFailure)
	}

	err = dummy.wsHandleData(wsTickerResp)
	if !errors.Is(err, nil) {
		t.Fatalf("received: %v but expected: %v", err, nil)
	}

	handled := <-dummy.Websocket.DataHandler
	if _, ok := handled.(*ticker.Price); !ok {
		t.Fatal("unexpected value")
	}

	err = dummy.wsHandleData(wsTransResp) // This doesn't pipe to datahandler
	if !errors.Is(err, nil) {
		t.Fatalf("received: %v but expected: %v", err, nil)
	}

	err = dummy.wsHandleData(wsOrderbookResp) // This doesn't pipe to datahandler
	if !errors.Is(err, nil) {
		t.Fatalf("received: %v but expected: %v", err, nil)
	}
}

func TestGenerateSubscriptions(t *testing.T) {
	t.Parallel()

	expected := subscription.List{}
	pairs, err := b.GetEnabledPairs(asset.Spot)
	assert.NoError(t, err, "GetEnabledPairs should not error")
	for _, baseSub := range b.Features.Subscriptions {
		baseSub.Channel = channelName(baseSub.Channel)
		s := baseSub.Clone()
		s.Pairs = pairs
		expected = append(expected, s)
	}
	subs, err := b.GenerateSubscriptions()
	assert.NoError(t, err, "GenerateSubscriptions should not error")
	if assert.Len(t, subs, len(expected), "Should have the correct number of subs") {
		assert.ElementsMatch(t, subs, expected, "Should get the correct subscriptions")
	}
}

func TestSubscribe(t *testing.T) {
	t.Parallel()

	b := new(Bithumb)
	require.NoError(t, testexch.TestInstance(b), "TestInstance must not error")
	testexch.SetupWs(t, b)
	err := b.Subscribe(subscription.List{
		{Channel: "private", Pairs: currency.Pairs{testPairMapping}},
		{Channel: subscription.OrderbookChannel, Pairs: currency.Pairs{testPairMapping}},
	},
	)
	require.NoError(t, err, "Subscribe must not error")
}
