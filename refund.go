package quickbooks

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/url"
)

type Refund struct {
	ID        string   `json:"Id,omitempty"`
	TxnDate   Date   `json:",omitempty"`
	Line []RefundLine
	DepositToAccountRef ReferenceType       `json:",omitempty"`
	CustomerRef         ReferenceType       `json:",omitempty"`
}

type RefundLine struct {
	ID        string   `json:"Id,omitempty"`

	DetailType          string              `json:",omitempty"`
	Amount              json.Number         `json:",omitempty"`
	SalesItemLineDetail SalesItemLineDetail `json:",omitempty"`
}

// CreateRefund creates the given Refund on the QuickBooks server, returning
// the resulting Refund object.
func (c *Client) CreateRefund(inv *Refund) (*Refund, error) {
	var u, err = url.Parse(string(c.Endpoint))
	if err != nil {
		return nil, err
	}
	u.Path = "/v3/company/" + c.RealmID + "/refundreceipt"
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
		Refund Refund
		Time    Date
	}
	err = json.NewDecoder(res.Body).Decode(&r)
	return &r.Refund, err
}
