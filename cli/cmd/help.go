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

	fmt.Println(" \033[1m-> cron <type>\033[0m")
	fmt.Println("    \033[3mStart the cron job with the specified type.\033[0m")
	fmt.Println("    \033[3mUse 'all' to start all cron jobs.\033[0m")
	fmt.Println("    \033[3mUse 'all_except:<names>' to start all cron jobs except the specified ones (seperate names with comma without space).\033[0m")
	fmt.Println("    \033[3mUse 'only:<names>' to start only the specified cron jobs (seperate names with comma without space).\033[0m")

	fmt.Println("\n \033[1m-> state\033[0m")
	fmt.Println("    \033[3mView the current state of the scans.\033[0m")

	fmt.Println("\n \033[1m-> history <scan_name> <from> <to>\033[0m")
	fmt.Println("    \033[3mView the history of the specified scan.\033[0m")
	fmt.Println("    \033[3mFor <from> and <to> use the format 'YYYY-MM-DD HH:MM:SS'.\033[0m")

	fmt.Println("\n \033[1m-> db init\033[0m")
	fmt.Println("    \033[3mInitialize the database.\033[0m")
	fmt.Println(" \033[1m-> db reset\033[0m")
	fmt.Println("    \033[3mReset the database, this will delete all the data!\033[0m")

	fmt.Println("\n \033[1m-> config verify <path>\033[0m")
	fmt.Println("    \033[3mVerify the configuration file at the specified path.\033[0m")
	fmt.Println(" \033[1m-> config replace <path>\033[0m")
	fmt.Println("    \033[3mReplace the current configuration with the one at the specified path.\033[0m")
	fmt.Println("    \033[3mFile is copied to the database, so you can delete the original file afterwards.\033[0m")
	fmt.Println(" \033[1m-> config view\033[0m")
	fmt.Println("    \033[3mView the current configuration.\033[0m")

	fmt.Println("\n \033[1m-> keys view all\033[0m")
	fmt.Println("    \033[3mView all the keys with their values in the database.\033[0m")
	fmt.Println("    \033[3mKeys with the * prefix can be changed.\033[0m")
	fmt.Println(" \033[1m-> keys view <key>\033[0m")
	fmt.Println("    \033[3mView the value of the specified key.\033[0m")
	fmt.Println("    \033[3mIf the key has the * prefix, it can be changed.\033[0m")
	fmt.Println(" \033[1m-> keys set <key> <value>\033[0m")
	fmt.Println("    \033[3mSet the value of the specified key.\033[0m")

	fmt.Println("\n \033[1m-> test ping <address> <timeout_ms>\033[0m")
	fmt.Println("    \033[3mTest the ping to the specified address with the specified timeout.\033[0m")
	fmt.Println("    \033[3mTimeout is in milliseconds.\033[0m")
	fmt.Println(" \033[1m-> test http <address> <timeout_ms>\033[0m")
	fmt.Println("    \033[3mTest the http request to the specified address with the specified timeout.\033[0m")
	fmt.Println("    \033[3mTimeout is in milliseconds.\033[0m")

	fmt.Println("\n \033[1m-> api test [service|http]\033[0m")
	fmt.Println("    \033[3mTest the api service or the http service.\033[0m")
	fmt.Println("    \033[3mUse 'service' to test the api service via systemctl.\033[0m")
	fmt.Println("    \033[3mUse 'http' to test the api service via http request.\033[0m")
	fmt.Println(" \033[1m-> api restart\033[0m")
	fmt.Println("    \033[3mRestart the api service.\033[0m")
	fmt.Println("    \033[3mThis command requires sudo privileges.\033[0m")

	fmt.Println("\n \033[1m-> help\033[0m")
	fmt.Println("    \033[3mPrint this help message.\033[0m")
}
