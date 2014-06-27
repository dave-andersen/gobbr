/*
bbstat is an example program using the bbrpc interface.
It can display different status messages from the Boolberry
daemon.
*/

package main

import (
	"github.com/dave-andersen/gobbr"
	"fmt"
	"time"
)

const DAEMON_ADDRESS = "http://localhost:10102"
const WALLET_ADDRESS = "http://localhost:9291" // pick something

func main() {
	d := gobbr.NewDaemon(DAEMON_ADDRESS)
	bh, err := d.GetBlockHeaderByHeight(1)
	if err == nil {
		fmt.Printf("Block %d has timestamp %s\n", bh.Height,
			time.Unix(int64(bh.Timestamp), 0))
	} else {
		fmt.Printf("Could not get blockheader: ", err)
	}

	w := gobbr.NewWallet(WALLET_ADDRESS)
	balance, err := w.GetBalance()
	if err == nil {
		fmt.Printf("Wallet has unlocked balance %f\n", float64(balance.UnlockedBalance)/gobbr.Multiplier)
	} else {
		fmt.Printf("Could not get balance: ", err)
	}
}


type BBBalance struct {
	Balance uint64 `json:"balance"`
	UnlockedBalance uint64 `json:"unlocked_balance"`
}

type BBJsonResponse struct {
	Id int `json:"id"`
	Jsonrpc string `json:"jsonrpc"`
}

type BBBalanceResponse struct {
	BBJsonResponse
	Result BBBalance `json:"result"`
}

type BBTransferDestination struct {
	Amount uint64 `json:"amount"`
	Address string `json:"address"`
}

type BBTransfer struct {
	Destinations []BBTransferDestination `json:"destinations"`
	Fee uint64 `json:"fee"`
	Mixin uint64 `json:"mixin"`
	UnlockTime uint64 `json:"unlock_time"`
	PaymentId string `json:"paymnet_id_hex"` // sic.  It's in the bool code.
}

type BBTransferResponse struct {
	BBJsonResponse
	Result struct {
		TxHash string `json:"tx_hash"`
	} `json:"result"`
}

type BBBlockHeader struct {
	Timestamp uint64 `json:"timestamp"`
	Height uint64 `json:"height"`
}

type BBGetBlockHeaderResponse struct {
	BBJsonResponse
	Result struct {
		Status string `json:"status"`
		BlockHeader BBBlockHeader `json:"block_header"`
	}
}

type BBQueryBlockHeader struct {
	Height uint64 `json:"height"`
}
