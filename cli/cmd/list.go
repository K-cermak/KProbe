package cmd

import (
	"fmt"

	"KProbeCLI/db"
	"KProbeCLI/helpers"
)

func ViewCurrentState() {
	scans := db.GetScans()
	if len(scans) == 0 {
		helpers.PrintWarning("No scans found")
		return
	}

	for _, scan := range scans {
		fmt.Println("\033[1m" + scan.Name + "\033[0m")

		state, success := db.GetScanNewest(scan.Name)
		if !success {
			fmt.Println(" -> State: Unknown")
		} else {
			fmt.Println(" -> " + helpers.BoolToState(state.Passed) + " (" + state.Generated + ")")
		}
	}
}

func ViewScanInfo(scanName string, start string, end string) {
	data := db.GetScanRes(scanName, start, end)
	if len(data) == 0 {
		helpers.PrintWarning("No data found for scan " + scanName + " from " + start + " to " + end)
		helpers.PrintQuestion("You sure your date format is the YYYY-MM-DD HH:MM:SS?")
		return
	}

	fmt.Println("\033[1mShowing history for scan " + scanName + " from " + start + " to " + end + ":\033[0m")
	for _, res := range data {
		fmt.Println(" -> " + helpers.BoolToState(res.Passed) + " (" + res.Generated + ")")
	}
}