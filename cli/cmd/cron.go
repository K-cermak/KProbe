package cmd

import (
	"strings"
	"sync"

	"KProbeCLI/db"
	"KProbeCLI/helpers"
	"KProbeCLI/utils"
)

func CronStart(command string) {
	scansSlice := db.GetScans()
	scansMap := make(map[string]bool)

	if len(scansSlice) == 0 {
		helpers.PrintError(true, "No scans found in database, add some via <kprobe config replace <path>>")
	}

	for _, scan := range scansSlice {
		scansMap[scan.Name] = false
	}

	if command == "all" {
		for _, scan := range scansSlice {
			scansMap[scan.Name] = true
		}
	} else if strings.HasPrefix(command, "all_except:") {
		excludedNames := strings.Split(strings.TrimPrefix(command, "all_except:"), ",")

		for _, scan := range scansSlice {
			scansMap[scan.Name] = true
		}

		for _, name := range excludedNames {
			if _, exists := scansMap[name]; exists {
				scansMap[name] = false
			} else {
				helpers.PrintWarning("Scan " + name + " not found in database, skipping")
			}
		}

	} else if strings.HasPrefix(command, "only:") {
		includedNames := strings.Split(strings.TrimPrefix(command, "only:"), ",")

		for k := range scansMap {
			scansMap[k] = false
		}

		for _, name := range includedNames {
			if _, exists := scansMap[name]; exists {
				scansMap[name] = true
			} else {
				helpers.PrintWarning("Scan " + name + " not found in database, skipping")
			}
		}

	} else {
		helpers.PrintError(true, "Invalid cron command")
	}

	scanCount := 0

	for _, scan := range scansSlice {
		if scansMap[scan.Name] {
			scanCount++
		}
	}

	if scanCount == 0 {
		helpers.PrintError(true, "No scans selected for cron job")
	}

	helpers.PrintInfo("Starting cron job for " + helpers.IntToStr(scanCount) + " scan(s)")

	ignoreSslStr := db.GetValue("ignore_ssl_errors")
	ignoreSsl := false
	if ignoreSslStr == "true" {
		ignoreSsl = true
	}

	scanResults := make(map[string]bool)
	var mu sync.Mutex
	var wg sync.WaitGroup

	for _, scan := range scansSlice {
		if !scansMap[scan.Name] {
			continue
		}

		wg.Add(1)

		go func(s helpers.Scan) {
			defer wg.Done()

			helpers.PrintInfo("Starting scan " + s.Name)

			var result bool
			if s.Type == "ping" {
				result = utils.PingAddress(s.Address, s.Timeout, false)
			} else if s.Type == "http" {
				result = utils.CheckHTTP(s.Address, s.Timeout, s.StatusCode, s.Keyword, ignoreSsl, false)
			}

			helpers.PrintSuccess("Scan " + s.Name + " finished (" + helpers.BoolToState(result) + ")")

			mu.Lock()
			scanResults[s.Name] = result
			mu.Unlock()
		}(scan)
	}

	wg.Wait()

	scanSuccessfull := 0

	for name, success := range scanResults {
		if success {
			scanSuccessfull++
			db.AddScanRes(name, true)
		} else {
			db.AddScanRes(name, false)
		}
	}

	if scanSuccessfull == scanCount {
		helpers.PrintSuccess("All " + helpers.IntToStr(scanCount) + " scan(s) finished successfully")
	} else if scanSuccessfull > 0 {
		helpers.PrintWarning("Only " + helpers.IntToStr(scanSuccessfull) + " out of " + helpers.IntToStr(scanCount) + " scan(s) finished successfully")
	} else {
		helpers.PrintError(true, "All "+helpers.IntToStr(scanCount)+" scan(s) finished with errors")
	}

	helpers.PrintInfo("Deleting old scan results")
	db.DeleteOldScanRes()
	helpers.PrintSuccess("Old scan results deleted")
}
