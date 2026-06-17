package cmd

import (
	"fmt"

	"KProbeCLI/db"
	"KProbeCLI/helpers"
)

func ViewKeys(key string) {
	found := false

	if key == "all" || key == "probe_name" {
		found = true
		fmt.Println("\033[1m*probe_name\033[0m")
		fmt.Println(" -> " + db.GetValue("probe_name"))
		fmt.Println("    \033[3mProbe name is used to identify the probe in the API requests.\033[0m")
	}

	if key == "all" || key == "db_version" {
		found = true
		fmt.Println("\033[1mdb_version\033[0m")
		fmt.Println(" -> " + db.GetValue("db_version"))
		fmt.Println("    \033[3mDatabase version.\033[0m")
	}

	if key == "all" || key == "db_init_time" {
		found = true
		fmt.Println("\033[1mdb_init_time\033[0m")
		fmt.Println(" -> " + db.GetValue("db_init_time"))
		fmt.Println("    \033[3mDatabase initialization time.\033[0m")
	}

	if key == "all" || key == "config_set" {
		found = true
		fmt.Println("\033[1mconfig_set\033[0m")
		fmt.Println(" -> " + db.GetValue("config_set"))
		fmt.Println("    \033[3mConfig set status, if set to true, you already set the config file.\033[0m")
	}

	if key == "all" || key == "delete_after" {
		found = true
		fmt.Println("\033[1m*delete_after\033[0m")
		fmt.Println(" -> " + db.GetValue("delete_after"))
		fmt.Println("    \033[3mDelete after is used to set the time to delete the scan history data after. Updating this value will not affect the existing data.\033[0m")
	}

	if key == "all" || key == "api_port" {
		found = true
		fmt.Println("\033[1m*api_port\033[0m")
		fmt.Println(" -> " + db.GetValue("api_port"))
		fmt.Println("    \033[3mAPI port is used to set the port for the API server.\033[0m")
	}

	if key == "all" || key == "editor_endpoint" {
		found = true
		fmt.Println("\033[1m*editor_endpoint\033[0m")
		fmt.Println(" -> " + db.GetValue("editor_endpoint"))
		fmt.Println("    \033[3mEditor endpoint is used to enable/disable the /editor endpoint.\033[0m")
	}

	if key == "all" || key == "ping_retries" {
		found = true
		fmt.Println("\033[1m*ping_retries\033[0m")
		fmt.Println(" -> " + db.GetValue("ping_retries"))
		fmt.Println("    \033[3mNumber of retries for ping requests.\033[0m")
	}

	if key == "all" || key == "ignore_ssl_errors" {
		found = true
		fmt.Println("\033[1m*ignore_ssl_errors\033[0m")
		fmt.Println(" -> " + db.GetValue("ignore_ssl_errors"))
		fmt.Println("    \033[3mIgnore SSL errors is used to ignore TLS (SSL) errors in the HTTP requests.\033[0m")
	}

	if found {
		fmt.Println("\n\033[3m(values with\033[0m \033[1m*\033[0m \033[3mcan be changed using <kprobe keys set <key> <value>> command)\033[0m")
	} else {
		helpers.PrintError(true, "Key "+key+" not found")
	}
}

func SetKeys(key string, value string) {
	switch key {
	case "probe_name":
		if len(value) < 3 || len(value) > 32 {
			helpers.PrintError(true, "Probe name must be between 3 and 32 characters")
		}

		db.InsertValue("probe_name", value)
		helpers.PrintInfo("You should now run <sudo kprobe api restart> to apply some changes")

	case "delete_after":
		if len(value) < 1 || len(value) > 36500 {
			helpers.PrintError(true, "Delete after must be between 1 and 36500 days")
		}

		db.InsertValue("delete_after", value)

	case "api_port":
		if len(value) < 1 || len(value) > 65535 {
			helpers.PrintError(true, "API port must be between 1 and 65535")
		}

		db.InsertValue("api_port", value)
		helpers.PrintInfo("You should now run <sudo kprobe api restart> to apply some changes")

	case "editor_endpoint":
		if value != "true" && value != "false" {
			helpers.PrintError(true, "Editor endpoint must be true or false")
		}

		db.InsertValue("editor_endpoint", value)
		helpers.PrintInfo("You should now run <sudo kprobe api restart> to apply some changes")

	case "ping_retries":
		if len(value) < 1 || len(value) > 100 {
			helpers.PrintError(true, "HTTP retries must be between 1 and 100")
		}

		db.InsertValue("ping_retries", value)

	case "ignore_ssl_errors":
		if value != "true" && value != "false" {
			helpers.PrintError(true, "Ignore SSL errors must be true or false")
		}

		db.InsertValue("ignore_ssl_errors", value)

	default:
		helpers.PrintError(true, "Key "+key+" not found or cannot be changed")
	}

	helpers.PrintSuccess("Key " + key + " updated to " + value)
}
