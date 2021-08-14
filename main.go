package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"

	finnhub "github.com/Finnhub-Stock-API/finnhub-go/v2"
	"gopkg.in/yaml.v2"
)

type conf struct {
	ApiKeyFinnhub string `yaml:"apiKeyFinnhub"`
}

func (c *conf) getConf() *conf {
	ymlFile, err := ioutil.ReadFile("./api-key-finnhub.yml")
	if err != nil {
		log.Printf("yamlFile.Get err   #%v ", err)
	}
	err = yaml.Unmarshal(ymlFile, c)
	if err != nil {
		log.Fatalf("Unmarshall: %v", err)
	}
	return c
}

func initializeParameters() string {
	var c conf
	c.getConf()
	apiKeyFinnhub := c.ApiKeyFinnhub
	return apiKeyFinnhub
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
		status = 200
	}
	return status
}

func main() {
	apiKeyFinnhub := initializeParameters()
	stockSymbolsArray := getStockSymbolJson(apiKeyFinnhub)
	loadIntoDatabase(stockSymbolsArray)
}
