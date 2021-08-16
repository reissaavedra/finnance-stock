package models

import (
	"fmt"

	"github.com/Finnhub-Stock-API/finnhub-go/v2"
)

type Stock struct {
	Symbol      string `pg:"symbol_id"`
	Description string
	Currency    string
	TypeStock   string
	Figi        string
}

func NewStock(stockSymbol finnhub.StockSymbol) *Stock {
	stock := Stock{
		Symbol:      stockSymbol.GetDisplaySymbol(),
		Description: stockSymbol.GetDescription(),
		Currency:    stockSymbol.GetCurrency(),
		TypeStock:   stockSymbol.GetType(),
		Figi:        stockSymbol.GetFigi(),
	}
	return &stock
}

func (s Stock) String() string {
	return fmt.Sprintf("Stock<Symbol=%s Desc=%s>\n", s.Symbol, s.Description)
}
