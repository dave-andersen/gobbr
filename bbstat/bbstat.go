/*
bbstat is an example program using the bbrpc interface.
It can display different status messages from the Boolberry
daemon.
*/

package main

import (
	"fmt"
	"github.com/dave-andersen/gobbr"
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
