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
		fmt.Println("    \033[3mName that identifies this probe instance in API responses.\033[0m")
	}

	if key == "all" || key == "db_version" {
		found = true
		fmt.Println("\033[1mdb_version\033[0m")
		fmt.Println(" -> " + db.GetValue("db_version"))
		fmt.Println("    \033[3mSchema version of the local SQLite database.\033[0m")
	}

	if key == "all" || key == "db_init_time" {
		found = true
		fmt.Println("\033[1mdb_init_time\033[0m")
		fmt.Println(" -> " + db.GetValue("db_init_time"))
		fmt.Println("    \033[3mTime when the database was initialized.\033[0m")
	}

	if key == "all" || key == "config_set" {
		found = true
		fmt.Println("\033[1mconfig_set\033[0m")
		fmt.Println(" -> " + db.GetValue("config_set"))
		fmt.Println("    \033[3mWhether a scan configuration has been loaded.\033[0m")
	}

	if key == "all" || key == "delete_after" {
		found = true
		fmt.Println("\033[1m*delete_after\033[0m")
		fmt.Println(" -> " + db.GetValue("delete_after"))
		fmt.Println("    \033[3mRetention period (in days) for scan history. Older records are deleted automatically. Changes apply only to new scan results.\033[0m")
	}

	if key == "all" || key == "api_port" {
		found = true
		fmt.Println("\033[1m*api_port\033[0m")
		fmt.Println(" -> " + db.GetValue("api_port"))
		fmt.Println("    \033[3mPort the API server listens on.\033[0m")
	}

	if key == "all" || key == "editor_endpoint" {
		found = true
		fmt.Println("\033[1m*editor_endpoint\033[0m")
		fmt.Println(" -> " + db.GetValue("editor_endpoint"))
		fmt.Println("    \033[3mEnables or disables the /editor web endpoint.\033[0m")
	}

	if key == "all" || key == "ping_retries" {
		found = true
		fmt.Println("\033[1m*ping_retries\033[0m")
		fmt.Println(" -> " + db.GetValue("ping_retries"))
		fmt.Println("    \033[3mNumber of packets sent during a single ping scan.\033[0m")
	}

	if key == "all" || key == "ignore_ssl_errors" {
		found = true
		fmt.Println("\033[1m*ignore_ssl_errors\033[0m")
		fmt.Println(" -> " + db.GetValue("ignore_ssl_errors"))
		fmt.Println("    \033[3mWhether invalid or self-signed TLS certificates are accepted during HTTP scans.\033[0m")
	}

	if key == "all" || key == "max_http_body_size" {
		found = true
		fmt.Println("\033[1m*max_http_body_size\033[0m")
		fmt.Println(" -> " + db.GetValue("max_http_body_size"))
		fmt.Println("    \033[3mMaximum response body size (in MB) read during HTTP scans.\033[0m")
	}

	if key == "all" || key == "output_http_info" {
		found = true
		fmt.Println("\033[1m*output_http_info\033[0m")
		fmt.Println(" -> " + db.GetValue("output_http_info"))
		fmt.Println("    \033[3mWhether the status code and response body are printed by <kprobe test http>.\033[0m")
	}

	if found {
		fmt.Println("\n\033[3m(keys marked with\033[0m \033[1m*\033[0m \033[3mcan be changed with <kprobe keys set <key> <value>>)\033[0m")
	} else {
		helpers.PrintError(true, "Key "+key+" not found")
	}
}

func SetKeys(key string, value string) {
	switch key {
	case "probe_name":
		if len(value) < 3 || len(value) > 32 {
			helpers.PrintError(true, "probe_name must be between 3 and 32 characters long")
		}

		db.InsertValue("probe_name", value)
		helpers.PrintInfo("Run <sudo kprobe api restart> to apply the change")

	case "delete_after":
		valInt, correct := helpers.StrToInt(value)
		if !correct || valInt < 1 || valInt > 36500 {
			helpers.PrintError(true, "delete_after must be between 1 and 36500 days")
		}

		value = helpers.IntToStr(valInt)
		db.InsertValue("delete_after", value)

	case "api_port":
		valInt, correct := helpers.StrToInt(value)
		if !correct || valInt < 1 || valInt > 65535 {
			helpers.PrintError(true, "api_port must be between 1 and 65535")
		}

		value = helpers.IntToStr(valInt)
		db.InsertValue("api_port", value)
		helpers.PrintInfo("Run <sudo kprobe api restart> to apply the change")

	case "editor_endpoint":
		if value != "true" && value != "false" {
			helpers.PrintError(true, "editor_endpoint must be true or false")
		}

		db.InsertValue("editor_endpoint", value)
		helpers.PrintInfo("Run <sudo kprobe api restart> to apply the change")

	case "ping_retries":
		valInt, correct := helpers.StrToInt(value)
		if !correct || valInt < 1 || valInt > 100 {
			helpers.PrintError(true, "ping_retries must be between 1 and 100")
		}

		value = helpers.IntToStr(valInt)
		db.InsertValue("ping_retries", value)

	case "ignore_ssl_errors":
		if value != "true" && value != "false" {
			helpers.PrintError(true, "ignore_ssl_errors must be true or false")
		}

		db.InsertValue("ignore_ssl_errors", value)

	case "max_http_body_size":
		valInt, correct := helpers.StrToInt(value)
		if !correct || valInt < 1 || valInt > 1024 {
			helpers.PrintError(true, "max_http_body_size must be between 1 and 1024 (MB)")
		}

		db.InsertValue("max_http_body_size", value)

	case "output_http_info":
		if value != "true" && value != "false" {
			helpers.PrintError(true, "output_http_info must be true or false")
		}

		db.InsertValue("output_http_info", value)

	default:
		helpers.PrintError(true, "Key "+key+" not found or cannot be changed")
	}

	helpers.PrintSuccess("Key " + key + " set to " + value)
}
