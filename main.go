package main

import (
	"context"
	"fmt"
	"project/config"
	"project/models"

	finnhub "github.com/Finnhub-Stock-API/finnhub-go/v2"
	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
)

func initializeParametersFinnhub() string {
	var c config.Conf
	finnhubApiKey := c.GetConf().FinnhubConf.FinnhubApiKey
	fmt.Printf("%#v\n", finnhubApiKey)
	return finnhubApiKey
}

func initializeParametersPgDatabase() *pg.DB {
	var c config.Conf
	pgConfig := c.GetConf().DatabaseConf
	db := pg.Connect(&pg.Options{
		Addr:     pgConfig.PgHost + ":" + pgConfig.PgPort,
		User:     pgConfig.PgUser,
		Password: pgConfig.PgPassword,
		Database: pgConfig.PgDatabase,
	})
	return db
}

func testPgConnection() {
	db := initializeParametersPgDatabase()
	defer db.Close()

	createSchema(db)
}

func createSchema(db *pg.DB) error {
	models := []interface{}{
		(*models.Stock)(nil),
	}
	for _, model := range models {
		err := db.Model(model).DropTable(&orm.DropTableOptions{
			IfExists: true,
			Cascade:  true,
		})
		if err != nil {
			return err
		}

		err = db.Model(model).CreateTable(nil)
		if err != nil {
			return err
		}
	}
	return nil
}

func getStockSymbolJson(apiKeyFiinhub string) []finnhub.StockSymbol {
	cfg := finnhub.NewConfiguration()
	cfg.AddDefaultHeader("X-Finnhub-Token", apiKeyFiinhub)
	finnhubClient := finnhub.NewAPIClient(cfg).DefaultApi
	res, _, err := finnhubClient.StockSymbols(context.Background()).Exchange("US").Execute()
	if err != nil {
		fmt.Println(err)
	}
	return res
}

func loadIntoDatabase(stockSymbolsArray []finnhub.StockSymbol) int {
	status := 200
	for _, stockSymbol := range stockSymbolsArray {
		fmt.Printf("%+v\n", stockSymbol.GetDescription())
		stock := newStock(stockSymbol)
		fmt.Printf("%+v\n", *stock)
	}
	return status
}

func newStock(stockSymbol finnhub.StockSymbol) *models.Stock {
	stock := models.Stock{
		Symbol:      stockSymbol.GetDisplaySymbol(),
		Description: stockSymbol.GetDescription(),
		Currency:    stockSymbol.GetCurrency(),
		TypeStock:   stockSymbol.GetType(),
		Figi:        stockSymbol.GetFigi(),
	}
	return &stock
}

func main() {
	testPgConnection()
}
