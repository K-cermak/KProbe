package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"KProbeCLI/db"
	"KProbeCLI/helpers"
)

const (
	INVALID_FORMAT       string = "Invalid config file format"
	DUPL_SCAN_NAME       string = "Duplicate scan name detected"
	INVALID_SCAN_NAME    string = "Invalid scan name"
	TIMEOUT_NOT_INT      string = "Timeout is not an integer"
	TIMEOUT_OUT_OF_RANGE string = "Timeout out of range"
	INVALID_CODE         string = "Invalid status code"
	CODE_TOO_LONG        string = "Status code too long"
	INVALID_KEYWORD      string = "Invalid keyword"
	KEYWORD_UNSUPPORTED  string = "Keyword is not supported for ping scans"
)

func VerifyConfig(path string) {
	helpers.PrintInfo("Verifying config file")

	file, err := os.Open(path)
	if err != nil {
		helpers.PrintError(true, "Failed to open file ("+err.Error()+")")
	}
	defer file.Close()

	scanNames := make(map[string]bool)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		fields := strings.Fields(line)
		if len(fields) == 0 {
			continue
		}
		if len(fields) < 4 {
			helpers.PrintError(true, INVALID_FORMAT)
		}

		scanName := fields[0]
		if err := validateScanName(scanName, scanNames); err != "" {
			helpers.PrintError(true, "Invalid scan name: "+err)
		}

		scanType := fields[1]
		if scanType != "http" && scanType != "ping" {
			helpers.PrintError(true, "Invalid scan type: "+scanType)
		}

		scanAddress := fields[2]
		if len(scanAddress) > 256 {
			helpers.PrintError(true, "Invalid scan address")
		}

		scanTimeout := fields[3]
		scanTimeout = strings.TrimPrefix(scanTimeout, "timeout=")
		if err := validateTimeout(scanTimeout); err != "" {
			helpers.PrintError(true, "Invalid scan interval: "+err)
		}

		if len(fields) > 4 {
			if strings.HasPrefix(fields[4], "status_code=") {
				if scanType == "ping" {
					helpers.PrintError(true, "Status code is not supported for ping scans")
				}

				if err := validateStatusCode(fields[4]); err != "" {
					helpers.PrintError(true, "Invalid status code: "+err)
				}
			}
		}

		if err := validateKeyword(line, scanType); err != "" {
			helpers.PrintError(true, "Invalid keyword: "+err)
		}

		scanNames[scanName] = true
	}

	if err := scanner.Err(); err != nil {
		helpers.PrintError(true, "Failed to read file ("+err.Error()+")")
	}

	helpers.PrintSuccess("Config file verified successfully")
}

func SetConfig(path string) {
	helpers.PrintInfo("Replacing config file")

	helpers.PrintQuestion("Do you want to replace the config file? (y/n)")
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)

	if input != "y" && input != "Y" {
		helpers.PrintWarning("Config file replacement aborted")
		return
	}

	helpers.PrintInfo("Deleting old config file")
	db.DeleteScans()
	helpers.PrintSuccess("Old config file deleted successfully")
	helpers.PrintInfo("Adding new scans")

	file, err := os.Open(path)
	if err != nil {
		helpers.PrintError(true, "Failed to open file ("+err.Error()+")")
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		fields := strings.Fields(line)

		scanName := fields[0]
		scanType := fields[1]
		scanAddress := fields[2]

		scanTimeout := fields[3]
		scanTimeout = strings.TrimPrefix(scanTimeout, "timeout=")
		scanTimeoutInt, correct := helpers.StrToInt(scanTimeout)
		if !correct {
			helpers.PrintError(true, "Failed to convert timeout to integer")
		}

		statusCode := ""
		keyword := ""

		if len(fields) > 4 {
			if strings.HasPrefix(fields[4], "status_code=") {
				statusCode = fields[4]
				statusCode = statusCode[13 : len(statusCode)-1]
			}
		}

		startIdx := strings.Index(line, "keyword=\"")
		if startIdx != -1 {
			endIdx := strings.LastIndex(line, "\"")
			keyword = line[startIdx+9 : endIdx]
		}

		helpers.PrintInfo("Adding scan: " + scanName)
		scan := helpers.Scan{
			Name:       scanName,
			Type:       scanType,
			Address:    scanAddress,
			Timeout:    scanTimeoutInt,
			StatusCode: statusCode,
			Keyword:    keyword,
		}

		db.AddScan(scan)

		db.InsertValue("config_set", "true")
		helpers.PrintSuccess("Scan added successfully")
	}

	if err := scanner.Err(); err != nil {
		helpers.PrintError(true, "Failed to read file ("+err.Error()+")")
	}

	helpers.PrintSuccess("Config file replaced successfully")
}

func ViewConfig() {
	scans := db.GetScans()
	if len(scans) == 0 {
		helpers.PrintWarning("No scans found")
		return
	}

	for _, scan := range scans {
		fmt.Println("\033[1m" + scan.Name + "\033[0m")
		fmt.Println(" -> Type: " + scan.Type)
		fmt.Println(" -> Address: " + scan.Address)
		fmt.Println(" -> Timeout: " + helpers.IntToStr(scan.Timeout) + "ms")
		if scan.StatusCode != "" {
			fmt.Println(" -> Status code(s): " + scan.StatusCode)
		}
		if scan.Keyword != "" {
			fmt.Println(" -> Keyword: \"" + scan.Keyword + "\"")
		}
		fmt.Println()
	}
}

func validateScanName(scanName string, scanNames map[string]bool) string {
	if _, exists := scanNames[scanName]; exists {
		return DUPL_SCAN_NAME
	}
	if scanName == "" || len(scanName) > 32 {
		return INVALID_SCAN_NAME
	}
	for _, c := range scanName {
		if !('a' <= c && c <= 'z') && !('0' <= c && c <= '9') && c != '_' {
			return INVALID_SCAN_NAME
		}
	}
	return ""
}

func validateTimeout(scanTimeout string) string {
	num, correct := helpers.StrToInt(scanTimeout)
	if !correct {
		return TIMEOUT_NOT_INT
	}

	if num < 0 || num > 30000 {
		return TIMEOUT_OUT_OF_RANGE
	}

	return ""
}

func validateStatusCode(statusCode string) string {
	if len(statusCode) > 256 {
		return CODE_TOO_LONG
	}

	statusCode = strings.TrimPrefix(statusCode, "status_code=\"")
	statusCode = statusCode[:len(statusCode)-1]

	codes := strings.Split(statusCode, ",")
	if len(codes) == 0 {
		return INVALID_CODE
	}

	for _, code := range codes {
		if num, ok := helpers.StrToInt(code); !ok || num < 100 || num > 599 {
			return INVALID_CODE
		}
	}

	return ""
}

func validateKeyword(line string, scanType string) string {
	startIdx := strings.Index(line, "keyword=\"")
	if startIdx == -1 {
		return ""
	}

	if scanType == "ping" {
		return KEYWORD_UNSUPPORTED
	}

	endIdx := strings.LastIndex(line, "\"")
	if endIdx == startIdx {
		return INVALID_KEYWORD
	}

	return ""
}
