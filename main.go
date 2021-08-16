package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"project/config"
	"project/models"

	finnhub "github.com/Finnhub-Stock-API/finnhub-go/v2"
	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
)

func initializeParametersFinnhub() string {
	var c config.Conf
	finnhubApiKey := c.GetConf().FinnhubConf.FinnhubApiKey
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

func initializeDatabase() {
	db := initializeParametersPgDatabase()
	defer db.Close()
	errorCreateSchema := createSchema(db)
	if errorCreateSchema != nil {
		f, err := os.OpenFile("testlogfile", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			log.Fatalf("error opening file: %v", err)
		}
		defer f.Close()
		log.SetOutput(f)
		log.Println(errorCreateSchema)
	}
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

func loadIntoDatabase(stockSymbolsArray []finnhub.StockSymbol) {
	stocks := []models.Stock{}
	for _, stockSymbol := range stockSymbolsArray {
		stock := models.NewStock(stockSymbol)
		stocks = append(stocks, *stock)
	}

	db := initializeParametersPgDatabase()
	defer db.Close()

	_, err := db.Model(&stocks).Insert()
	if err != nil {
		panic(err)
	}
}

func main() {
	initializeDatabase()
	finnhubApiKey := initializeParametersFinnhub()
	stockSymbolsArray := getStockSymbolJson(finnhubApiKey)
	loadIntoDatabase(stockSymbolsArray)
}
