/*
 Blockinfo is an example program using the bbrpc interface.
 It prints out the trailing 24 hour average reward from the
 blockchain (actually, 720 blocks, not strictly 24 hours).

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

	for i := (height-720); i < height; i++ {
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
	fmt.Printf("Normal: %d  Orphans:  %d  Avg Reward: %2.f\n", nonOrphans, orphans,
		float64(totalReward)/float64(nonOrphans*gobbr.Multiplier) )
}
