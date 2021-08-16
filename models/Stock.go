package models

type Stock struct {
	Id          int `pg:"pk_id"`
	Symbol      string
	Description string
	Currency    string
	TypeStock   string
	Figi        string
}
