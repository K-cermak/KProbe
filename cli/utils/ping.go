package utils

import (
	"fmt"
	"time"

	"KProbeCLI/db"
	"KProbeCLI/helpers"

	"github.com/prometheus-community/pro-bing"
)

const maxPingTimeoutMs = 30000

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
	if !correct || count < 1 {
		if output {
			helpers.PrintError(false, "Error parsing ping retries value")
		}

		return false
	}

	// 0 (unlimited) would never finish for an unreachable host
	if timeout <= 0 {
		timeout = maxPingTimeoutMs
	}
	packetTimeout := time.Duration(timeout) * time.Millisecond

	// pinger.SetPrivileged(true)
	pinger.Count = count
	// Last packet is sent after (count - 1) intervals, then it still gets the full packet timeout
	pinger.Timeout = time.Duration(count-1)*pinger.Interval + packetTimeout

	receivedInTime := 0
	pinger.OnRecv = func(pkt *probing.Packet) {
		if pkt.Rtt <= packetTimeout {
			receivedInTime++
		}
	}

	err = pinger.Run()
	if err != nil {
		if output {
			helpers.PrintError(false, "Error running ping command ("+err.Error()+")")
		}
		return false
	}

	packetLoss := float64(count-receivedInTime) / float64(count) * 100

	stats := pinger.Statistics()
	if output {
		fmt.Println("\033[1mPING Stats to " + address + "\033[0m")
		fmt.Printf(" -> Sent: %d\n", stats.PacketsSent)
		fmt.Printf(" -> Received: %d\n", stats.PacketsRecv)
		fmt.Printf(" -> Received within timeout: %d of %d (%.2f%% loss)\n", receivedInTime, count, packetLoss)
		fmt.Printf(" -> Min RTT: %s\n", stats.MinRtt)
		fmt.Printf(" -> Max RTT: %s\n", stats.MaxRtt)
		fmt.Printf(" -> Avg RTT: %s\n", stats.AvgRtt)
	}

	return packetLoss < 20
}
