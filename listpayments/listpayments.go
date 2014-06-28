/*
bbstat is an example program using the bbrpc interface.
It can display different status messages from the Boolberry
daemon.
*/

package main

import (
	"flag"
	"fmt"
	"github.com/dave-andersen/gobbr"
)

const WALLET_ADDRESS = "http://localhost:9291" // pick something

func main() {
	flag.Parse()
	if flag.NArg() < 1 {
		fmt.Println("Use:  listpayments <payment id>")
		return
	}

	w := gobbr.NewWallet(WALLET_ADDRESS)

	payments, err := w.GetPayments(flag.Arg(0))
	if err != nil {
		fmt.Println("Could not get payments:", err)
		return
	}

	for p, _ := range payments {
		fmt.Println("Payment: ", p)
	}
}
