package cmd

import (
	"fmt"
	"time"
)

var Version = "dev"

func PrintHelp() {
	fmt.Println("\033[1m+------------------------+\033[0m")
	fmt.Println("\033[1m|       KProbe CLI       |\033[0m")
	fmt.Println("\033[1m+------------------------+\033[0m")

	fmt.Println("\033[1mRepository:\033[0m https://github.com/K-cermak/KProbe")
	fmt.Println("\033[1mLicense:\033[0m MIT")
	fmt.Println("\033[1mAuthor:\033[0m Karel Cermak | karlosoft.com")
	fmt.Printf("\033[1mVersion:\033[0m %s, © %d\n", Version, time.Now().Year())

	fmt.Println("\n\033[1mUsage: kprobe <command>\033[0m")
	fmt.Println("\n\033[1mCommands:\033[0m")

	fmt.Println(" \033[1m-> scan <type>\033[0m")
	fmt.Println("    \033[3mRun the selected scans once.\033[0m")
	fmt.Println("    \033[3mUse 'all' to run all scans.\033[0m")
	fmt.Println("    \033[3mUse 'all_except:<names>' to run all scans except the specified ones (comma-separated, no spaces).\033[0m")
	fmt.Println("    \033[3mUse 'only:<names>' to run only the specified scans (comma-separated, no spaces).\033[0m")

	fmt.Println("\n \033[1m-> state\033[0m")
	fmt.Println("    \033[3mShow the latest result of each scan.\033[0m")

	fmt.Println("\n \033[1m-> history <scan_name> <from> <to>\033[0m")
	fmt.Println("    \033[3mView the history of the specified scan.\033[0m")
	fmt.Println("    \033[3mFor <from> and <to>, use the format 'YYYY-MM-DD HH:MM:SS'.\033[0m")

	fmt.Println("\n \033[1m-> db init\033[0m")
	fmt.Println("    \033[3mInitialize the database.\033[0m")
	fmt.Println(" \033[1m-> db reset\033[0m")
	fmt.Println("    \033[3mReset the database. This deletes all data!\033[0m")

	fmt.Println("\n \033[1m-> config verify <path>\033[0m")
	fmt.Println("    \033[3mVerify the configuration file at the specified path.\033[0m")
	fmt.Println(" \033[1m-> config replace <path>\033[0m")
	fmt.Println("    \033[3mReplace the current configuration with the one at the specified path.\033[0m")
	fmt.Println("    \033[3mThe configuration is copied into the database, so you can delete the file afterwards.\033[0m")
	fmt.Println(" \033[1m-> config view\033[0m")
	fmt.Println("    \033[3mView the current configuration.\033[0m")

	fmt.Println("\n \033[1m-> keys view all\033[0m")
	fmt.Println("    \033[3mShow all keys and their values.\033[0m")
	fmt.Println("    \033[3mKeys with the * prefix can be changed.\033[0m")
	fmt.Println(" \033[1m-> keys view <key>\033[0m")
	fmt.Println("    \033[3mShow the value of the specified key.\033[0m")
	fmt.Println("    \033[3mIf the key has the * prefix, it can be changed.\033[0m")
	fmt.Println(" \033[1m-> keys set <key> <value>\033[0m")
	fmt.Println("    \033[3mSet the value of the specified key.\033[0m")

	fmt.Println("\n \033[1m-> test ping <address> <timeout_ms>\033[0m")
	fmt.Println("    \033[3mPing the specified address.\033[0m")
	fmt.Println("    \033[3mThe timeout is in milliseconds and applies to each packet.\033[0m")
	fmt.Println(" \033[1m-> test http <address> <timeout_ms>\033[0m")
	fmt.Println("    \033[3mSend an HTTP request to the specified address.\033[0m")
	fmt.Println("    \033[3mThe timeout is in milliseconds.\033[0m")

	fmt.Println("\n \033[1m-> api test [service|http]\033[0m")
	fmt.Println("    \033[3mCheck the API service.\033[0m")
	fmt.Println("    \033[3mUse 'service' to check the service status via systemctl.\033[0m")
	fmt.Println("    \033[3mUse 'http' to send a test HTTP request to the API.\033[0m")
	fmt.Println(" \033[1m-> api restart\033[0m")
	fmt.Println("    \033[3mRestart the API service.\033[0m")
	fmt.Println("    \033[3mThis command requires sudo privileges.\033[0m")

	fmt.Println("\n \033[1m-> help\033[0m")
	fmt.Println("    \033[3mPrint this help message.\033[0m")
}
