package helpers

import (
	"os"
	"strconv"
)

type Variable struct {
	Name  string `json:"name"`
	Mode  string `json:"mode"`            // "bool" or "json"
	Value string `json:"value,omitempty"` // for mode="bool": the search string
	Query string `json:"query,omitempty"` // for mode="json": the JSON path (e.g. "$.data.temperature")
	As    string `json:"as,omitempty"`    // for mode="json": "bool" or "number"
}

type Scan struct {
	Name       string
	Type       string
	Address    string
	Timeout    int
	StatusCode string
	Variables  []Variable
	Expression string
}

type ScanRes struct {
	Generated string
	Passed    bool
}

func PrintInfo(info string) {
	println("\033[0;34m[*]\033[0m " + info)
}

func PrintSuccess(success string) {
	println("\033[0;32m[OK]\033[0m " + success)
}

func PrintQuestion(question string) {
	println("\033[0;36m[?]\033[0m " + question)
}

func PrintWarning(warning string) {
	println("\033[0;33m[WARN]\033[0m " + warning)
}

func PrintError(fatal bool, err string) {
	println("\033[0;31m[ERROR]\033[0m " + err)
	if fatal {
		os.Exit(1)
	}
}

func StrToInt(str string) (int, bool) {
	num, err := strconv.Atoi(str)
	if err != nil {
		return 0, false
	}

	return num, true
}

func IntToStr(num int) string {
	return strconv.Itoa(num)
}

func BoolToState(state bool) string {
	if state {
		return "\033[0;32mSuccess\033[0m"
	}

	return "\033[0;31mFailed\033[0m"
}

func ArgsMatch(args []string, expectedArgs []string) bool {
	if len(args) != len(expectedArgs) {
		return false
	}

	for i, arg := range args {
		if arg != expectedArgs[i] && expectedArgs[i] != "*" {
			return false
		}
	}

	return true
}
