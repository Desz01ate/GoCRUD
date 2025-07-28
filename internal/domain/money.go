package domain

import (
	"errors"
	"fmt"
)

type Currency string

const (
	THB Currency = "THB"
	USD Currency = "USD"
)

type Money struct {
	Amount   int64    `json:"amount"` // intentional, avoiding floating point issues.
	Currency Currency `json:"currency"`
}

func NewMoney(amount int64, currency Currency) Money {
	return Money{
		Amount:   amount,
		Currency: currency,
	}
}

func (m Money) IsZero() bool {
	return m.Amount == 0
}

func (m Money) IsNegative() bool {
	return m.Amount < 0
}

func (m Money) IsPositive() bool {
	return m.Amount > 0
}

func (m Money) Add(other Money) (Money, error) {
	if m.Currency != other.Currency {
		return Money{}, errors.New("cannot add money with different currencies")
	}

	return Money{
		Amount:   m.Amount + other.Amount,
		Currency: m.Currency,
	}, nil
}

func (m Money) Subtract(other Money) (Money, error) {
	if m.Currency != other.Currency {
		return Money{}, errors.New("cannot subtract money with different currencies")
	}

	return Money{
		Amount:   m.Amount - other.Amount,
		Currency: m.Currency,
	}, nil
}

func (m Money) String() string {
	return fmt.Sprintf("%.2f %s", float64(m.Amount)/100, m.Currency)
}

func (m Money) ToFloat() float64 {
	return float64(m.Amount) / 100
}
