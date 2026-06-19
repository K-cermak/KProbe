package utils

import (
	"fmt"
	"time"

	"KProbeCLI/db"
	"KProbeCLI/helpers"

	"github.com/prometheus-community/pro-bing"
)

func PingAddress(address string, timeout int, output bool) bool {
	pinger, err := probing.NewPinger(address)
	if err != nil {
		if output {
			helpers.PrintError(false, "Error creating ping command ("+err.Error()+")")
		}

		return false
	}

	countStr := db.GetValue("ping_retries")
	count, correct := helpers.StrToInt(countStr)
	if !correct {
		if output {
			helpers.PrintError(false, "Error parsing ping retries value")
		}

		return false
	}

	// pinger.SetPrivileged(true)
	pinger.Count = count
	pinger.Timeout = time.Duration(timeout) * time.Millisecond

	err = pinger.Run()
	if err != nil {
		if output {
			helpers.PrintError(false, "Error running ping command ("+err.Error()+")")
		}
		return false
	}

	stats := pinger.Statistics()
	if output {
		fmt.Println("\033[1mPING Stats to " + address + "\033[0m")
		fmt.Printf(" -> Sent: %d\n", stats.PacketsSent)
		fmt.Printf(" -> Received: %d\n", stats.PacketsRecv)
		fmt.Printf(" -> Lost: %f (%.2f%% loss)\n", stats.PacketLoss, stats.PacketLoss)
		fmt.Printf(" -> Min RTT: %s\n", stats.MinRtt)
		fmt.Printf(" -> Max RTT: %s\n", stats.MaxRtt)
		fmt.Printf(" -> Avg RTT: %s\n", stats.AvgRtt)
	}

	return stats.PacketLoss < 20
}
