package cmd

import (
	"bufio"
	"os"
	"strings"

	"KProbeCLI/db"
	"KProbeCLI/helpers"
)

func InitDatabase() {
	if db.DatabaseExist() {
		helpers.PrintError(true, "Database already exists, if you want to reset it, run <kprobe db reset>")
	}

	helpers.PrintQuestion("Do you want to initialize the database? (y/n)")
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)

	if input != "y" && input != "Y" {
		helpers.PrintWarning("Database initialization aborted")
		return
	}

	helpers.PrintInfo("Initializing database")
	db.InitDatabase()
	helpers.PrintSuccess("Database initialized successfully")
	helpers.PrintInfo("You should now run <sudo kprobe api restart> to apply some changes")
}

func ResetDatabase() {
	if !db.DatabaseExist() {
		helpers.PrintError(true, "Database does not exist, did you want to run <kprobe db init>?")
	}

	helpers.PrintQuestion("Do you want to reset the database? \033[0;31mTHIS ACTION IS IRREVERSIBLE!\033[0m TYPE \"DESTROY\" TO CONFIRM")
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)

	if input != "DESTROY" {
		helpers.PrintWarning("Database reset aborted")
		return
	}

	helpers.PrintInfo("Resetting database")
	db.InitDatabase()
	helpers.PrintSuccess("Database reset successfully")
	helpers.PrintInfo("You should now run <sudo kprobe api restart> to apply some changes")
}
