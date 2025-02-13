package main

import (
	"context"
	"errors"
	"log"
	"math/rand"
	"sync"
	"time"

	"github.com/thrasher-corp/gocryptotrader/common"
	"github.com/thrasher-corp/gocryptotrader/currency"
	"github.com/thrasher-corp/gocryptotrader/engine"
	exchange "github.com/thrasher-corp/gocryptotrader/exchanges"
	"github.com/thrasher-corp/gocryptotrader/exchanges/asset"
	"github.com/thrasher-corp/gocryptotrader/exchanges/kline"
)

const (
	totalWrappers = 25
)

func main() {
	var err error
	engine.Bot, err = engine.New()
	if err != nil {
		log.Fatalf("Failed to initialise engine. Err: %s", err)
	}

	engine.Bot.Settings = engine.Settings{
		DisableExchangeAutoPairUpdates: true,
	}

	log.Printf("Loading exchanges..")
	var wg sync.WaitGroup
	for x := range exchange.Exchanges {
		err := engine.Bot.LoadExchange(exchange.Exchanges[x], &wg)
		if err != nil {
			log.Printf("Failed to load exchange %s. Err: %s",
				exchange.Exchanges[x],
				err)
			continue
		}
	}
	wg.Wait()
	log.Println("Done.")

	log.Printf("Testing exchange wrappers..")
	results := make(map[string][]string)
	wg = sync.WaitGroup{}
	exchanges := engine.Bot.GetExchanges()
	for x := range exchanges {
		exch := exchanges[x]
		wg.Add(1)
		go func(e exchange.IBotExchange) {
			results[e.GetName()] = testWrappers(e)
			wg.Done()
		}(exch)
	}
	wg.Wait()
	log.Println("Done.")

	log.Println()
	for name, funcs := range results {
		pct := float64(totalWrappers-len(funcs)) / float64(totalWrappers) * 100
		log.Printf("Exchange %s wrapper coverage [%d/%d - %.2f%%] | Total missing: %d", name, totalWrappers-len(funcs), totalWrappers, pct, len(funcs))
		log.Printf("\t Wrappers not implemented:")

		for x := range funcs {
			log.Printf("\t - %s", funcs[x])
		}
		log.Println()
	}
}

func testWrappers(e exchange.IBotExchange) []string {
	p := currency.NewPair(currency.BTC, currency.USD)
	assetType := asset.Spot
	if !e.SupportsAsset(assetType) {
		assets := e.GetAssetTypes(false)
		rand.Seed(time.Now().Unix())
		assetType = assets[rand.Intn(len(assets))] // nolint:gosec // basic number generation required, no need for crypo/rand
	}

	var funcs []string

	_, err := e.FetchTicker(context.TODO(), p, assetType)
	if errors.Is(err, common.ErrNotYetImplemented) {
		funcs = append(funcs, "FetchTicker")
	}

	_, err = e.UpdateTicker(context.TODO(), p, assetType)
	if errors.Is(err, common.ErrNotYetImplemented) {
		funcs = append(funcs, "UpdateTicker")
	}

	_, err = e.FetchOrderbook(context.TODO(), p, assetType)
	if errors.Is(err, common.ErrNotYetImplemented) {
		funcs = append(funcs, "FetchOrderbook")
	}

	_, err = e.UpdateOrderbook(context.TODO(), p, assetType)
	if errors.Is(err, common.ErrNotYetImplemented) {
		funcs = append(funcs, "UpdateOrderbook")
	}

	_, err = e.FetchTradablePairs(context.TODO(), asset.Spot)
	if errors.Is(err, common.ErrNotYetImplemented) {
		funcs = append(funcs, "FetchTradablePairs")
	}

	err = e.UpdateTradablePairs(context.TODO(), false)
	if errors.Is(err, common.ErrNotYetImplemented) {
		funcs = append(funcs, "UpdateTradablePairs")
	}

	_, err = e.FetchAccountInfo(context.TODO(), assetType)
	if errors.Is(err, common.ErrNotYetImplemented) {
		funcs = append(funcs, "GetAccountInfo")
	}

	_, err = e.GetRecentTrades(context.TODO(), p, assetType)
	if errors.Is(err, common.ErrNotYetImplemented) {
		funcs = append(funcs, "GetRecentTrades")
	}

	_, err = e.GetHistoricTrades(context.TODO(), p, assetType, time.Time{}, time.Time{})
	if errors.Is(err, common.ErrNotYetImplemented) {
		funcs = append(funcs, "GetHistoricTrades")
	}

	_, err = e.GetFundingHistory(context.TODO())
	if errors.Is(err, common.ErrNotYetImplemented) {
		funcs = append(funcs, "GetFundingHistory")
	}

	_, err = e.SubmitOrder(context.TODO(), nil)
	if errors.Is(err, common.ErrNotYetImplemented) {
		funcs = append(funcs, "SubmitOrder")
	}

	_, err = e.ModifyOrder(context.TODO(), nil)
	if errors.Is(err, common.ErrNotYetImplemented) {
		funcs = append(funcs, "ModifyOrder")
	}

	err = e.CancelOrder(context.TODO(), nil)
	if errors.Is(err, common.ErrNotYetImplemented) {
		funcs = append(funcs, "CancelOrder")
	}

	_, err = e.CancelBatchOrders(context.TODO(), nil)
	if errors.Is(err, common.ErrNotYetImplemented) {
		funcs = append(funcs, "CancelBatchOrders")
	}

	_, err = e.CancelAllOrders(context.TODO(), nil)
	if errors.Is(err, common.ErrNotYetImplemented) {
		funcs = append(funcs, "CancelAllOrders")
	}

	_, err = e.GetOrderInfo(context.TODO(), "1", p, assetType)
	if errors.Is(err, common.ErrNotYetImplemented) {
		funcs = append(funcs, "GetOrderInfo")
	}

	_, err = e.GetOrderHistory(context.TODO(), nil)
	if errors.Is(err, common.ErrNotYetImplemented) {
		funcs = append(funcs, "GetOrderHistory")
	}

	_, err = e.GetActiveOrders(context.TODO(), nil)
	if errors.Is(err, common.ErrNotYetImplemented) {
		funcs = append(funcs, "GetActiveOrders")
	}

	_, err = e.GetDepositAddress(context.TODO(), currency.BTC, "")
	if errors.Is(err, common.ErrNotYetImplemented) {
		funcs = append(funcs, "GetDepositAddress")
	}

	_, err = e.WithdrawCryptocurrencyFunds(context.TODO(), nil)
	if errors.Is(err, common.ErrNotYetImplemented) {
		funcs = append(funcs, "WithdrawCryptocurrencyFunds")
	}

	_, err = e.WithdrawFiatFunds(context.TODO(), nil)
	if errors.Is(err, common.ErrNotYetImplemented) {
		funcs = append(funcs, "WithdrawFiatFunds")
	}
	_, err = e.WithdrawFiatFundsToInternationalBank(context.TODO(), nil)
	if errors.Is(err, common.ErrNotYetImplemented) {
		funcs = append(funcs, "WithdrawFiatFundsToInternationalBank")
	}

	_, err = e.GetHistoricCandles(context.TODO(), currency.Pair{}, asset.Spot, time.Unix(0, 0), time.Unix(0, 0), kline.OneDay)
	if errors.Is(err, common.ErrNotYetImplemented) {
		funcs = append(funcs, "GetHistoricCandles")
	}

	_, err = e.GetHistoricCandlesExtended(context.TODO(), currency.Pair{}, asset.Spot, time.Unix(0, 0), time.Unix(0, 0), kline.OneDay)
	if errors.Is(err, common.ErrNotYetImplemented) {
		funcs = append(funcs, "GetHistoricCandlesExtended")
	}

	_, err = e.UpdateAccountInfo(context.TODO(), assetType)
	if errors.Is(err, common.ErrNotYetImplemented) {
		funcs = append(funcs, "UpdateAccountInfo")
	}

	_, err = e.GetFeeByType(context.TODO(), &exchange.FeeBuilder{})
	if errors.Is(err, common.ErrNotYetImplemented) {
		funcs = append(funcs, "GetFeeByType")
	}

	err = e.UpdateOrderExecutionLimits(context.TODO(), asset.DownsideProfitContract)
	if errors.Is(err, common.ErrNotYetImplemented) {
		funcs = append(funcs, "UpdateOrderExecutionLimits")
	}
	return funcs
}
