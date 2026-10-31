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
		helpers.PrintError(true, "Database already exists. To reset it, run <kprobe db reset>")
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
	helpers.PrintInfo("Run <sudo kprobe api restart> to apply the changes")
}

func ResetDatabase() {
	if !db.DatabaseExist() {
		helpers.PrintError(true, "Database does not exist. Did you mean <kprobe db init>?")
	}

	helpers.PrintQuestion("Do you want to reset the database? \033[0;31mTHIS CANNOT BE UNDONE!\033[0m Type \"DESTROY\" to confirm")
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
	helpers.PrintInfo("Run <sudo kprobe api restart> to apply the changes")
}
