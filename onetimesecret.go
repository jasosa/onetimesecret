package onetimesecret

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/pkg/errors"
)

// Onetimesecret represents the information of a one time secret
type Onetimesecret struct {
	SecretKey   string `json:"secret_key"`
	SecretValue string `json:"value"`
}

//Client onetimesecret api client
type Client struct {
	user  string
	token string
	url   string
}

// NewClient initializes a new instance of Client
func NewClient(user, token, url string) *Client {
	c := &Client{
		user:  user,
		token: token,
		url:   url,
	}

	return c
}

// Generate generates a short, unique secret and returns the key to share it and the value
func (c *Client) Generate(ttl int) (string, string, error) {
	url := fmt.Sprintf("%s/generate?ttl=%d", c.url, ttl)
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return "", "", errors.Wrap(err, "error creating request")
	}
	req.SetBasicAuth(c.user, c.token)

	var res *http.Response
	if res, err = http.DefaultClient.Do(req); err != nil {
		return "", "", errors.Wrap(err, "error executing request")
	}

	bytess, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", "", errors.Wrap(err, "error reading response")
	}

	ots := &Onetimesecret{}
	err = json.Unmarshal(bytess, ots)
	if err != nil {
		return "", "", errors.Wrap(err, "error unmarshaling response")
	}
	return ots.SecretKey, ots.SecretValue, nil
}
