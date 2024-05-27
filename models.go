package main

import "github.com/shopspring/decimal"

type Imovel struct {
	ID       int             `json:"id"`
	Endereco string          `json:"endereco"`
	Preco    decimal.Decimal `json:"preco"`
	Area     decimal.Decimal `json:"area"`
}
