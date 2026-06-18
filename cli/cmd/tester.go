package cmd

import (
	"os/exec"
	"runtime"

	"KProbeCLI/db"
	"KProbeCLI/helpers"
	"KProbeCLI/utils"
)

func PingTest(address string, timeout string) {
	count, correct := helpers.StrToInt(timeout)
	if !correct {
		helpers.PrintError(true, "Invalid timeout value")
	}

	helpers.PrintInfo("Pinging " + address + " with timeout " + timeout + " ms")
	ret := utils.PingAddress(address, count, true)
	if ret {
		helpers.PrintSuccess("Ping successful")
	} else {
		helpers.PrintWarning("Ping failed")
	}
}

func HttpTest(address string, timeout string) {
	count, correct := helpers.StrToInt(timeout)
	if !correct {
		helpers.PrintError(true, "Invalid timeout value")
	}

	ignoreSslStr := db.GetValue("ignore_ssl_errors")
	ignoreSsl := false
	if ignoreSslStr == "true" {
		ignoreSsl = true
	}

	helpers.PrintInfo("Performing HTTP request to " + address + " with timeout " + timeout + " ms")
	ret := utils.CheckHTTP(address, count, "", nil, "", ignoreSsl, true)
	if ret {
		helpers.PrintSuccess("HTTP request successful")
	} else {
		helpers.PrintWarning("HTTP request failed")
	}
}

func ApiTest(testType string) {
	if testType != "service" && testType != "http" {
		helpers.PrintError(true, "Invalid type, expected <service> or <http>")
	}

	switch testType {
	case "http":
		apiPort := db.GetValue("api_port")
		HttpTest("http://127.0.0.1:" + apiPort, "5000")

	case "service":
		if runtime.GOOS != "linux" {
			helpers.PrintError(true, "This action is only available on Linux")
		}

		cmd := exec.Command("systemctl", "is-active", "kprobe")
		output, err := cmd.CombinedOutput()

		if string(output) == "active\n" {
			helpers.PrintSuccess("Service is active")
		} else if err != nil {
			helpers.PrintWarning("Service is not active or systemctl failed (" + err.Error() + ")")
		} else {
			helpers.PrintWarning("Service is not active")
		}
	}
}

func ApiRestart() {
	if runtime.GOOS != "linux" {
		helpers.PrintError(true, "This action is only available on Linux")
	}

	cmd := exec.Command("systemctl", "restart", "kprobe")
	output, err := cmd.CombinedOutput()
	if err != nil {
		errMsg := "Failed to restart service (" + err.Error() + "). Tip: Use 'sudo' to run the command."
		if len(output) > 0 {
			errMsg += "\nOutput: " + string(output)
		}
		helpers.PrintError(true, errMsg)
	}

	helpers.PrintSuccess("Service restarted")
}
