package quickbooks

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/url"
)

type Payment struct {
	ID             string      `json:"Id,omitempty"`
	TxnDate   Date   `json:",omitempty"`
	TotalAmt       json.Number `json:",omitempty"`
	ProcessPayment bool        `json:",omitempty"`
	CustomerRef    ReferenceType
	Line           []PaymentLine
}

// PaymentLine ...
type PaymentLine struct {
	ID        string `json:"Id,omitempty"`
	LineNum   int    `json:",omitempty"`
	Amount    json.Number
	LinkedTxn []TxnLine
}

type TxnLine struct {
	TxnId   json.Number `json:",omitempty"`
	TxnType string      `json:",omitempty"`
}

// CreatePayment creates the given Payment on the QuickBooks server, returning
// the resulting Payment object.
func (c *Client) CreatePayment(inv *Payment) (*Payment, error) {
	var u, err = url.Parse(string(c.Endpoint))
	if err != nil {
		return nil, err
	}
	u.Path = "/v3/company/" + c.RealmID + "/payment"
	var v = url.Values{}
	v.Add("minorversion", minorVersion)
	u.RawQuery = v.Encode()
	var j []byte
	j, err = json.Marshal(inv)
	if err != nil {
		return nil, err
	}
	var req *http.Request
	req, err = http.NewRequest("POST", u.String(), bytes.NewBuffer(j))
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
	var res *http.Response
	res, err = c.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, parseFailure(res)
	}

	var r struct {
		Payment Payment
		Time    Date
	}
	err = json.NewDecoder(res.Body).Decode(&r)
	return &r.Payment, err
}
