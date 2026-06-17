package main

import (
	"os"

	"KProbeCLI/cmd"
	"KProbeCLI/helpers"
)

func main() {
	args := os.Args

	// cron <type>
	if helpers.ArgsMatch(args, []string{"*", "cron", "*"}) {
		cmd.CronStart(args[2])
		return
	}

	// state
	if helpers.ArgsMatch(args, []string{"*", "state"}) {
		cmd.ViewCurrentState()
		return
	}

	// history <scan_name> <from> <to>
	if helpers.ArgsMatch(args, []string{"*", "history", "*", "*", "*"}) {
		cmd.ViewScanInfo(args[2], args[3], args[4])
		return
	}

	// db init
	if helpers.ArgsMatch(args, []string{"*", "db", "init"}) {
		cmd.InitDatabase()
		return
	}

	// db reset
	if helpers.ArgsMatch(args, []string{"*", "db", "reset"}) {
		cmd.ResetDatabase()
		return
	}

	// config verify <path>
	if helpers.ArgsMatch(args, []string{"*", "config", "verify", "*"}) {
		cmd.VerifyConfig(args[3])
		return
	}

	// config replace <path>
	if helpers.ArgsMatch(args, []string{"*", "config", "replace", "*"}) {
		cmd.VerifyConfig(args[3])
		cmd.SetConfig(args[3])
		return
	}

	// config view
	if helpers.ArgsMatch(args, []string{"*", "config", "view"}) {
		cmd.ViewConfig()
		return
	}

	// keys view all / keys view <key>
	if helpers.ArgsMatch(args, []string{"*", "keys", "view", "*"}) {
		cmd.ViewKeys(args[3])
		return
	}

	// keys set <key> <value>
	if helpers.ArgsMatch(args, []string{"*", "keys", "set", "*", "*"}) {
		cmd.SetKeys(args[3], args[4])
		return
	}

	// test ping <address> <timeout_ms>
	if helpers.ArgsMatch(args, []string{"*", "test", "ping", "*", "*"}) {
		cmd.PingTest(args[3], args[4])
		return
	}

	// test http <address> <timeout_ms>
	if helpers.ArgsMatch(args, []string{"*", "test", "http", "*", "*"}) {
		cmd.HttpTest(args[3], args[4])
		return
	}

	// api test [service|http]
	if helpers.ArgsMatch(args, []string{"*", "api", "test", "*"}) {
		cmd.ApiTest(args[3])
		return
	}

	// api restart
	if helpers.ArgsMatch(args, []string{"*", "api", "restart"}) {
		cmd.ApiRestart()
		return
	}

	// help
	if helpers.ArgsMatch(args, []string{"*", "help"}) {
		cmd.PrintHelp()
		return
	}

	helpers.PrintError(true, "Invalid command, rerun with <kprobe help> for help")
}
