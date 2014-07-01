/*
 Blockinfo is an example program using the bbrpc interface.
 It retrieves the header for every block in the blockchain and
 will print out which ones were orphans.

 The secondary purpose of this program is to benchmark the time
 it takes to retrieve a lot of blockchain data.
 */

package main

import (
	"fmt"
	"github.com/dave-andersen/gobbr"
)

const DAEMON_ADDRESS = "http://localhost:10102"

func main() {
	d := gobbr.NewDaemon(DAEMON_ADDRESS)

	height, err := d.GetHeight()
	if err != nil {
		fmt.Println("Error getting height: ", err)
		return
	}

	orphans := 0
	nonOrphans := 0
	totalReward := uint64(0)

	// For benchmarking, set
	//height = 1000
	for i := uint64(0); i < height; i++ {
		bh, err := d.GetBlockHeaderByHeight(i)
		if err != nil {
			fmt.Println("Error getting blockheader: ", err)
			return
		}

		if bh.OrphanStatus {
			orphans++
		} else {
			nonOrphans++
			totalReward += bh.Reward
		}
	}
	fmt.Printf("Normal: %d  Orphans:  %d  Reward: %2.f\n", nonOrphans, orphans,
		float64(totalReward)/gobbr.Multiplier )
}
