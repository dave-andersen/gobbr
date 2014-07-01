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

func ReadPostQuery(cli *http.Client, url string, data []byte) ([]byte, error) {
	rd := bytes.NewReader(data)
	resp, err := cli.Post(url, "application/json", rd)
	if err != nil {
		return nil, err
	}
	response_body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	resp.Body.Close()
	return response_body, nil
}

func DoJSONQuery(cli *http.Client, url string, dest HasError, args map[string]interface{}) error {
	args["jsonrpc"] = "2.0"
	jsonbuf, err := json.Marshal(args)
	if err != nil {
		fmt.Println("error on marshal: ", err)
		os.Exit(-1)
	}
	response_body, err := ReadPostQuery(cli, url, jsonbuf)
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

func DoGetJSON(url string, dest interface{}) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	response_body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	resp.Body.Close()

	err = json.Unmarshal(response_body, &dest)
	if err != nil {
		return err
	}
	return nil
}


type Daemon struct {
	address string
	cli *http.Client
}

func NewDaemon(address string) *Daemon {
	d := &Daemon{address, &http.Client{}}
	return d
}

type Wallet struct {
	address string
	cli *http.Client
}

func NewWallet(address string) *Wallet {
	d := &Wallet{address, &http.Client{}}
	return d
}

func (d *Daemon) DoJSONQuery(dest HasError, args map[string]interface{}) error {
	return DoJSONQuery(d.cli, d.address+"/json_rpc", dest, args)
}

func (w *Wallet) DoJSONQuery(dest HasError, args map[string]interface{}) error {
	return DoJSONQuery(w.cli, w.address+"/json_rpc", dest, args)
}

/*
 * Daemon queries - not JSON_RPC encoded
 */

type GetHeightResponse struct {
	Height uint64 `json:"height"`
	Status string `json:"status"`
}

func (d *Daemon) GetHeight() (uint64, error) {
	var resp GetHeightResponse
	err := DoGetJSON(d.address + "/getheight", &resp)
	if err != nil {
		return 0, err
	}
	if (resp.Status != "OK") {
		return 0, errors.New(resp.Status)
	}
	return resp.Height, nil
}

type GetInfoResponse struct {
	AliasCount uint64 `json:"alias_count"`
	AltBlocksCount uint64 `json:"alt_blocks_count"`
	CurrentBlocksMedian uint64 `json:"current_blocks_median"`
	CurrentNetworkHashrate350 uint64 `json:"current_network_hashrate_350"`
	CurrentNetworkHashrate50 uint64 `json:"current_network_hashrate_50"`
	Difficulty uint64 `json:"difficulty"`
	GreyPeerlistSize uint64 `json:"grey_peerlist_size"`
	Height uint64 `json:"height"`
	IncomingConnectionsCount uint64 `json:"incoming_connections_count"`
	OutgoingConnectionsCount uint64 `json:"outgoing_connections_count"`
	ScratchpadSize uint64 `json:"scratchpad_size"`
	TxCount uint64 `json:"tx_count"`
	TxPoolSize uint64 `json:"tx_pool_size"`
	WhitePeerlistSize uint64 `json:"white_peerlist_size"`
	Status string `json:"status"`
}

func (d *Daemon) GetInfo() (GetInfoResponse, error) {
	var resp GetInfoResponse
	err := DoGetJSON(d.address + "/getinfo", &resp)
	if err != nil {
		return resp, err
	}
	if (resp.Status != "OK") {
		return resp, errors.New(resp.Status)
	}
	return resp, nil
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
 * -- Build a more interesting sample app.
 * -- Add an RPC method to the daemon to get more payments
 */