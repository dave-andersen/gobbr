/*
Package gobbr provides a wrapper interface around the Boolberry
daemon and wallet RPC interfaces.

Example use:

 bbr := gobbr.NewDaemon("http://localhost:10102")
 bh, err := d.GetBlockHeaderByHeight(1)
 if err == nil {
         fmt.Printf("Block %d has timestamp %d\n", bh.Height, bh.Timestamp)
 }

This package is still under development and many features are not yet
supported.  Interfaces may change at any time.
*/

package gobbr

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

const (
	Multiplier = (1e12)
	TransferFee = 100000000  /* in uint64 units */
)

type ErrorStruct struct {
	Code int `json:"code"`
	Message string `json:"message"`
}

type JsonResponse struct {
	Id      int    `json:"id"`
	Jsonrpc string `json:"jsonrpc"`
	Error ErrorStruct `json:"error"`
}

func (j *JsonResponse) GetError() ErrorStruct {
	return j.Error
}

type HasError interface {
	GetError() ErrorStruct
}

func ReadPostQuery(url string, data []byte) ([]byte, error) {
	rd := bytes.NewReader(data)
	resp, err := http.Post(url, "application/json", rd)
	if err != nil {
		return nil, err
	}
	response_body, err := ioutil.ReadAll(resp.Body)
	fmt.Println("responsebody: ", string(response_body))
	if err != nil {
		return nil, err
	}
	resp.Body.Close()
	return response_body, nil
}

func DoJSONQuery(url string, dest HasError, args map[string]interface{}) error {
	args["jsonrpc"] = "2.0"
	jsonbuf, err := json.Marshal(args)
	if err != nil {
		fmt.Println("error on marshal: ", err)
		os.Exit(-1)
	}
	response_body, err := ReadPostQuery(url, jsonbuf)
	if err != nil {
		return err
	}
	err = json.Unmarshal(response_body, &dest)
	if err != nil {
		return err
	}
	neterr := dest.GetError()
	if neterr.Code != 0 {
		return errors.New(neterr.Message)
	}
	return nil
}

type Daemon struct {
	address string
}

func NewDaemon(address string) *Daemon {
	/* Todo:  Connect and keep connection cached */
	d := &Daemon{address}
	return d
}

type Wallet struct {
	address string
}

func NewWallet(address string) *Wallet {
	/* Todo:  Connect and keep connection cached */
	d := &Wallet{address}
	return d
}

func (d *Daemon) DoJSONQuery(dest HasError, args map[string]interface{}) error {
	return DoJSONQuery(d.address+"/json_rpc", dest, args)
}

func (w *Wallet) DoJSONQuery(dest HasError, args map[string]interface{}) error {
	return DoJSONQuery(w.address+"/json_rpc", dest, args)
}

/*
 * Queries about the blockchain
 */

type BlockHeader struct {
	Timestamp uint64 `json:"timestamp"`
	Height    uint64 `json:"height"`
	/* INCOMPLETE */
}

type GetBlockHeaderResponse struct {
	JsonResponse
	Result struct {
		Status      string      `json:"status"`
		BlockHeader BlockHeader `json:"block_header"`
	}
}

type QueryBlockHeader struct {
	Height uint64 `json:"height"`
}

func (d *Daemon) GetBlockHeaderByHeight(height uint64) (bh BlockHeader, err error) {
	var resp GetBlockHeaderResponse
	err = d.DoJSONQuery(&resp, map[string]interface{}{
		"method": "getblockheaderbyheight",
		"params": &QueryBlockHeader{height},
	})
	return resp.Result.BlockHeader, err

}

/*
 * Wallet queries
 */

type Balance struct {
	Balance         uint64 `json:"balance"`
	UnlockedBalance uint64 `json:"unlocked_balance"`
}

type BalanceResponse struct {
	JsonResponse
	Result Balance `json:"result"`
}

func (w *Wallet) GetBalance() (balance Balance, err error) {
	var resp BalanceResponse
	err = w.DoJSONQuery(&resp, map[string]interface{}{
		"method": "getbalance",
	})
	return resp.Result, err
}

type TransferDestination struct {
	Amount uint64 `json:"amount"`
	Address string `json:"address"`
}

type Transfer struct {
	Destinations []TransferDestination `json:"destinations"`
	Fee uint64 `json:"fee"`
	Mixin uint64 `json:"mixin"`
	UnlockTime uint64 `json:"unlock_time"`
	PaymentId string `json:"payment_id_hex"`
}

type TransferResponse struct {
	JsonResponse
	Result struct {
		TxHash string `json:"tx_hash"`
	} `json:"result"`
}

func (w *Wallet) Transfer(destination string, amount uint64, mixinCount uint64, paymentId string) (txid string, err error) {
	var resp TransferResponse
	var req = &Transfer{}
	req.Destinations = []TransferDestination{ TransferDestination{amount, destination}}
	req.Fee = TransferFee
	req.Mixin = mixinCount
	req.UnlockTime = 0
	req.PaymentId = paymentId /* must be hex string */
	
	err = w.DoJSONQuery(&resp, map[string]interface{}{
		"method": "transfer",
		"params":req,
	})
	
	return resp.Result.TxHash, err
}

type PaymentDetails struct {
	TxHash string `json:"tx_hash"`
	Amount uint64 `json:"amount"`
	BlockHeight uint64 `json:"block_height"`
	UnlockTime uint64 `json:"unlock_time"`
}

type PaymentsResponse struct {
	JsonResponse
	Result struct {
		Payments []PaymentDetails `json:"payments"`
	}
}

type PaymentsRequest struct {
	PaymentId string `json:"payment_id"`
}

func (w *Wallet) GetPayments(paymentId string) (payments []PaymentDetails, err error) {
	var resp PaymentsResponse
	var req PaymentsRequest
	req.PaymentId = paymentId
	err = w.DoJSONQuery(&resp, map[string]interface{}{
		"method": "get_payments",
		"params":req,
	})
	return resp.Result.Payments, err
}

/*
 * Known bugs tracking:
 * -- GetPayments is mostly untested because the API call needs a payment ID
 *
 * Todo tracking:
 * -- Add support for some of the "info" outputs from the daemon -- those
 *    would make for nice sample apps for, e.g., showing network state.
 */