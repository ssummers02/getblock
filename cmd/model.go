package main

type Block struct {
	Number       string `json:"number"`
	Transactions []Tx   `json:"transactions"`
}

type Tx struct {
	From  string `json:"from"`
	To    string `json:"to"`
	Value string `json:"value"`
}
