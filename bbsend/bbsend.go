/*
 * bbsend is a sample application to transfer money using the Boolberry
 * wallet RPC package.
 */

package main

import (
	"github.com/dave-andersen/gobbr"
	"flag"
        "fmt"
	"strconv"
)

var WalletAddress = flag.String("wallet", "http://localhost:9291", "Specify address of wallet")
var PaymentID = flag.String("id", "", "Payment ID (optional)")
var Mixin = flag.Uint64("mixin", 0, "Mixin count")

func main() {

	flag.Parse()

	if (flag.NArg() < 2) {
		fmt.Println("Use:  bbrsend [-w wallet] [-id payment_id] [-m mixin] <dst> <amt>")
		return
	}

	w := gobbr.NewWallet(*WalletAddress)

	bal, err := w.GetBalance()

	if err != nil {
		fmt.Println("Could not get balance", err)
		return
	}
	fmt.Printf("Balance: %.2f\n", (float64(bal.UnlockedBalance)/gobbr.Multiplier))

	amt, err := strconv.ParseUint(flag.Arg(1), 10, 64)
	if err != nil {
		fmt.Println("Invalid amount: ", err)
		return
	}
	amt *= gobbr.Multiplier
	dst := flag.Arg(0)

	if (bal.UnlockedBalance <= (1 * gobbr.Multiplier)) {
		fmt.Println("Not enough unlocked balance to send tx")
		return
	}

	fmt.Printf("About to transfer %.2f to %s\n", float64(amt)/gobbr.Multiplier, dst);
	if (*PaymentID != "") {
		fmt.Println("with payment ID: ", *PaymentID)
	}
	fmt.Printf("\n")

	txid, err := w.Transfer(dst, amt, *Mixin, *PaymentID)
	if err != nil {
		fmt.Println("Could not do transfer: ", err)
		return
	}
	fmt.Println("Transferred BBR, txid: ", txid)
}
