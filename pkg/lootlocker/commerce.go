package lootlocker

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

type Wallet struct{}

// https://api.lootlocker.com/game/wallet/holder/<wallet_id>

func (c *Client) GetWallet(session, wallet string) (*Wallet, error) {
	raw, err := c.Request(http.MethodGet, fmt.Sprintf("game/wallet/%v", wallet), application_json, nil, map[string]string{"x-session-token": session})
	if err != nil {
		return nil, err
	}
	w := &Wallet{}
	if err := json.Unmarshal(raw, w); err != nil {
		return nil, err
	}
	return w, nil
}

type Credit struct {
	Amount     string `json:"amount"`
	WalletID   string `json:"wallet_id"`
	CurrencyID string `json:"currency_id"`
}

// https://api.lootlocker.com/game/balances/credit

func (c *Client) CreditBalance(session string, credit *Credit) error {
	raw, err := json.Marshal(credit)
	if err != nil {
		return err
	}
	_, err = c.Request(http.MethodPost, "game/balances/credit", application_json, raw, map[string]string{"x-session-token": session})
	return err
}

type Currency struct {
	ID   string `json:"id"`
	Code string `json:"code"`
	Name string `json:"name"`
}
type Balance struct {
	CreatedAt time.Time `json:"created_at"`
	Amount    float64   `json:"amount,string"`
	Currency  Currency  `json:"currency"`
}

// https://api.lootlocker.com/game/balances/wallet/<wallet_id>

func (c *Client) GetBalances(session, wallet string) (map[string]Balance, error) {
	raw, err := c.Request(http.MethodGet, fmt.Sprintf("game/balances/wallet/%v", wallet), application_json, nil, map[string]string{"x-session-token": session})
	if err != nil {
		return nil, err
	}
	pl := map[string][]Balance{}
	if err := json.Unmarshal(raw, &pl); err != nil {
		return nil, err
	}
	balances := map[string]Balance{}
	for _, b := range pl["balances"] {
		balances[b.Currency.Code] = b
	}
	log.Println(balances)
	return balances, nil
}
