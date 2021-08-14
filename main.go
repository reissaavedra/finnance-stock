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

func getStockSymbolJson(apiKeyFiinhub string) {
	cfg := finnhub.NewConfiguration()
	cfg.AddDefaultHeader("X-Finnhub-Token", apiKeyFiinhub)
	finnhubClient := finnhub.NewAPIClient(cfg).DefaultApi
	res, _, err := finnhubClient.StockSymbols(context.Background()).Exchange("US").Execute()
	if err == nil {
		fmt.Printf("%+v\n", res)
	} else {
		fmt.Println(err)
	}
}

func main() {
	apiKeyFinnhub := initializeParameters()
	getStockSymbolJson(apiKeyFinnhub)
}
