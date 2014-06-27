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
        "net/http"
        "fmt"
        "io/ioutil"
	"os"
        "encoding/json"
)

const (
	Multiplier = 1000000000000
)

type JsonResponse struct {
	Id int `json:"id"`
	Jsonrpc string `json:"jsonrpc"`
}

func ReadPostQuery(url string, data []byte) ([]byte, error) {
	rd := bytes.NewReader(data)
        resp, err := http.Post(url, "application/json", rd)
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

func DoJSONQuery(url string, dest interface{}, args map[string]interface{}) error {
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

func (d *Daemon) DoJSONQuery(dest interface{}, args map[string]interface{}) error {
	return DoJSONQuery(d.address + "/json_rpc", dest, args)
}

func (w *Wallet) DoJSONQuery(dest interface{}, args map[string]interface{}) error {
	return DoJSONQuery(w.address + "/json_rpc", dest, args)
}

/*
 * Queries about the blockchain
 */

type BlockHeader struct {
	Timestamp uint64 `json:"timestamp"`
	Height uint64 `json:"height"`
	/* INCOMPLETE */
}

type GetBlockHeaderResponse struct {
	JsonResponse
	Result struct {
		Status string `json:"status"`
		BlockHeader BlockHeader `json:"block_header"`
	}
}

type QueryBlockHeader struct {
	Height uint64 `json:"height"`
}


func (d *Daemon) GetBlockHeaderByHeight(height uint64) (bh BlockHeader, err error) {
	var resp GetBlockHeaderResponse
	err = d.DoJSONQuery(&resp, map[string]interface{} {
		"method" : "getblockheaderbyheight",
		"params" : &QueryBlockHeader{height},
	})
	return resp.Result.BlockHeader, err

}

/*
 * Wallet queries
 */

type Balance struct {
	Balance uint64 `json:"balance"`
	UnlockedBalance uint64 `json:"unlocked_balance"`
}

type BalanceResponse struct {
	JsonResponse
	Result Balance `json:"result"`
}

func (w *Wallet) GetBalance() (balance Balance, err error) {
	var resp BalanceResponse
	err = w.DoJSONQuery(&resp, map[string]interface{} {
		"method" :  "getbalance",
	})
	return resp.Result, err
}