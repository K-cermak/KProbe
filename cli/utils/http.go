package utils

import (
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"KProbeCLI/helpers"
)

func CheckHTTP(url string, timeout int, acceptCodes string, keyword string, ignoreSslErrors bool, output bool) bool {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: ignoreSslErrors},
	}

	client := http.Client{
		Timeout:   time.Duration(timeout) * time.Millisecond,
		Transport: tr,
	}

	resp, err := client.Get(url)
	if err != nil {
		if output {
			helpers.PrintError(false, "Error performing HTTP request ("+err.Error()+")")
		}
		return false
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		if output {
			helpers.PrintError(false, "Error reading response body ("+err.Error()+")")
		}
		return false
	}

	if acceptCodes != "" {
		foundCode := false
		acceptCodesArray := strings.Split(acceptCodes, ",")
		for _, code := range acceptCodesArray {
			codeInt, correct := helpers.StrToInt(code)
			if !correct {
				helpers.PrintError(true, "Invalid timeout value")
			}

			if resp.StatusCode == codeInt {
				foundCode = true
				break
			}
		}

		if !foundCode {
			if output {
				helpers.PrintError(false, "Invalid status code received ("+helpers.IntToStr(resp.StatusCode)+")")
			}
			return false
		}
	}

	bodyStr := string(body)
	displayBody := bodyStr
	truncated := 0

	if keyword != "" {
		if !strings.Contains(bodyStr, keyword) {
			if output {
				helpers.PrintError(false, "Keyword not found in response body")
			}
			return false
		}
	}

	if len(bodyStr) > 100 {
		displayBody = bodyStr[:100]
		truncated = len(bodyStr) - 100
	}

	if output {
		fmt.Println("\033[1mHTTP Response from " + url + "\033[0m")
		fmt.Printf(" -> Status Code: %d\n", resp.StatusCode)
		fmt.Println(" -> Response Body:")
		fmt.Print("    " + displayBody)
		if truncated > 0 {
			fmt.Printf("... (truncated %d characters)\n", truncated)
		} else {
			fmt.Println()
		}
	}

	return true
}
