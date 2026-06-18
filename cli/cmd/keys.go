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
		fmt.Println("    \033[3mName used to identify this probe instance in API requests.\033[0m")
	}

	if key == "all" || key == "db_version" {
		found = true
		fmt.Println("\033[1mdb_version\033[0m")
		fmt.Println(" -> " + db.GetValue("db_version"))
		fmt.Println("    \033[3mThe schema version of the local SQLite database.\033[0m")
	}

	if key == "all" || key == "db_init_time" {
		found = true
		fmt.Println("\033[1mdb_init_time\033[0m")
		fmt.Println(" -> " + db.GetValue("db_init_time"))
		fmt.Println("    \033[3mThe timestamp when the database was originally initialized.\033[0m")
	}

	if key == "all" || key == "config_set" {
		found = true
		fmt.Println("\033[1mconfig_set\033[0m")
		fmt.Println(" -> " + db.GetValue("config_set"))
		fmt.Println("    \033[3mIndicates whether the initial configuration has been completed.\033[0m")
	}

	if key == "all" || key == "delete_after" {
		found = true
		fmt.Println("\033[1m*delete_after\033[0m")
		fmt.Println(" -> " + db.GetValue("delete_after"))
		fmt.Println("    \033[3mThe retention period (in days) for scan history. Older records are automatically deleted. Changes apply only to new scans.\033[0m")
	}

	if key == "all" || key == "api_port" {
		found = true
		fmt.Println("\033[1m*api_port\033[0m")
		fmt.Println(" -> " + db.GetValue("api_port"))
		fmt.Println("    \033[3mThe listening port for the local API server.\033[0m")
	}

	if key == "all" || key == "editor_endpoint" {
		found = true
		fmt.Println("\033[1m*editor_endpoint\033[0m")
		fmt.Println(" -> " + db.GetValue("editor_endpoint"))
		fmt.Println("    \033[3mToggles the availability of the /editor web interface endpoint.\033[0m")
	}

	if key == "all" || key == "ping_retries" {
		found = true
		fmt.Println("\033[1m*ping_retries\033[0m")
		fmt.Println(" -> " + db.GetValue("ping_retries"))
		fmt.Println("    \033[3mThe number of packets to send (retries) during a single Ping scan.\033[0m")
	}

	if key == "all" || key == "ignore_ssl_errors" {
		found = true
		fmt.Println("\033[1m*ignore_ssl_errors\033[0m")
		fmt.Println(" -> " + db.GetValue("ignore_ssl_errors"))
		fmt.Println("    \033[3mToggles whether invalid or self-signed TLS/SSL certificates should be ignored during HTTP scans.\033[0m")
	}

	if key == "all" || key == "max_http_body_size" {
		found = true
		fmt.Println("\033[1m*max_http_body_size\033[0m")
		fmt.Println(" -> " + db.GetValue("max_http_body_size"))
		fmt.Println("    \033[3mThe maximum payload size (in megabytes) to read during HTTP scans.\033[0m")
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
			helpers.PrintError(true, "Ping retries must be between 1 and 100")
		}

		db.InsertValue("ping_retries", value)

	case "ignore_ssl_errors":
		if value != "true" && value != "false" {
			helpers.PrintError(true, "Ignore SSL errors must be true or false")
		}

		db.InsertValue("ignore_ssl_errors", value)

	case "max_http_body_size":
		valInt, correct := helpers.StrToInt(value)
		if !correct || valInt < 1 || valInt > 1024 {
			helpers.PrintError(true, "Max HTTP body size must be between 1 and 1024 MB")
		}

		db.InsertValue("max_http_body_size", value)

	default:
		helpers.PrintError(true, "Key "+key+" not found or cannot be changed")
	}

	helpers.PrintSuccess("Key " + key + " updated to " + value)
}
