package main

import (
	"math/rand"
	"time"
)

type TransferRequest struct {
	ToAccount int64 `json:"toAccount,string"`
	Amount    int   `json:"amount,string"`
}

type CreateAccountRequest struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

type UpdateAccountRequest struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Balance   int64  `json:"balance,string"`
}

type Account struct {
	ID        int64     `json:"id"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Number    int64     `json:"number"`
	Balance   int64     `json:"balance"`
	CreatedAt time.Time `json:"created_at"`
}

func NewAccount(firstName, lastName string) *Account {
	return &Account{
		FirstName: firstName,
		LastName:  lastName,
		Number:    int64(rand.Intn(1000000)),
		CreatedAt: time.Now().UTC(),
	}
}

func UpdateAccount(id int64, firstName, lastName string, balance int64) *Account {
	return &Account{
		ID:        id,
		FirstName: firstName,
		LastName:  lastName,
		Balance:   balance,
		CreatedAt: time.Now().UTC(),
	}
}
