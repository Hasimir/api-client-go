package ice3x

import (
	"crypto/hmac"
	"crypto/sha512"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

var _key, _secret string

var _url string = "https://api.ice3x.com"

type Trade struct {
	Id           int64  `json:"id"`
	CreationTime int64  `json:"creationTime"`
	Description  string `json:"description"`
	Price        int64  `json:"price"`
	Volume       int64  `json:"volume"`
	Fee          int64  `json:"fee"`
}

type TradeRequest struct {
	Currency   string `json:"currency"`
	Instrument string `json:"instrument"`
	Limit      int32  `json:"limit"`
	Since      int64  `json:"since"`
}

type TradeHistoryResponse struct {
	Success      bool    `json:"success"`
	ErrorCode    int32   `json:"errorCode"`
	ErrorMessage string  `json:"errorMessage"`
	Trades       []Trade `json:"trades, []Trade"`
}

func SetAuth(key, secret string) {
	_key = key
	_secret = secret
}

func signMessage(message string, secret string) (string, error) {
	key, err := base64.StdEncoding.DecodeString(secret)
	if err != nil {
		return "", err
	}
	h := hmac.New(sha512.New, key)
	h.Write([]byte(message))
	return base64.StdEncoding.EncodeToString(h.Sum(nil)), nil
}

func requestHttp(path string, values url.Values, postData string, v interface{}) error {

	endpoint, err := url.Parse(_url)
	if err != nil {
		return err
	}
	endpoint.Path += path

	nowInMilisecond := int64(time.Now().UnixNano()) / int64(time.Millisecond)
	timestamp := strconv.FormatInt(nowInMilisecond, 10)

	// create string to sign
	stringToSign := path + "\n" + timestamp + "\n" + postData

	signature, err := signMessage(stringToSign, _secret)
	if err != nil {
		return err
	}

	reqBody := strings.NewReader(postData)

	// create the request
	req, err := http.NewRequest("POST", endpoint.String(), reqBody)
	if err != nil {
		return err
	}

	req.Header.Add("User-Agent", "Mozilla/4.0 (compatible; Ice3x GO client)")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("accept", "application/json")
	req.Header.Add("accept-charset", "utf-8")
	req.Header.Add("signature", signature)
	req.Header.Add("apikey", _key)
	req.Header.Add("timestamp", timestamp)

	// submit the http request
	r, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	// if no result interface, return
	if v == nil {
		return nil
	}

	// read the body of the http message into a byte array
	body, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		return err
	}

	if len(body) == 0 {
		return fmt.Errorf("Response body length is 0")
	}

	e := make(map[string]interface{})
	err = json.Unmarshal(body, &e)
	if bsEr, ok := e["error"]; ok {
		return fmt.Errorf("%v", bsEr)
	}
	var datastr = string(body[:])

	//parse the JSON response into the response object
	return json.Unmarshal(body, v)
}

func TradeHistory(currency string, instrument string, limit int32, since int64) (*TradeHistoryResponse, error) {

	request := TradeRequest{Currency: currency, Instrument: instrument, Limit: limit, Since: since}

	postData, err := json.Marshal(request)

	if err != nil {
		return nil, err
	}

	postDataStr := string(postData[:])

	res := &TradeHistoryResponse{}
	err = requestHttp("/order/trade/history", url.Values{}, postDataStr, res)
	if err != nil {
		return nil, err
	}
	return res, nil
}
