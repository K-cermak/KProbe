# FAQ Guide
- [Installation](#installation)
- [Editor](#editor)
- [Creating a Scan, Setting Up Cron](#creating-a-scan-setting-up-cron)
- [API, Reverse Proxy](#api-reverse-proxy)
- [Connect to Uptime Kuma](#connect-to-uptime-kuma)
- [List of All Commands](#list-of-all-commands)
- [Uninstallation](#uninstallation)
- [Migration](#migration)

<br>

## Installation
- **System Requirements**: We recommend setting up a virtual server (VM) or dedicating physical hardware. The baseline requirements are minimal: **1 GB RAM**, **1 CPU core**, and **1 GB of free disk space** for the SQLite database (more space might be needed if you plan to run many scans or retain long history records).
- **Supported Operating Systems**: KProbe is compatible with any modern Linux distribution using **systemd** (`systemctl`). We offer native packages for:
  - **Debian / Ubuntu** (and derivatives) via `.deb` packages.
  - **CentOS / RHEL / Rocky Linux / Fedora** (and derivatives) via `.rpm` packages.

<br>

- <b>1. Download and Install the Package</b>
  - Get the latest package for your architecture from the GitHub [Releases](https://github.com/K-cermak/KProbe/releases) page.
  
  **On Debian/Ubuntu-based systems (.deb):**
  ```bash
  sudo dpkg -i kprobe_*.deb
  ```
  
  **On RHEL/CentOS/Rocky/Fedora-based systems (.rpm):**
  ```bash
  sudo rpm -i kprobe_*.rpm
  # Or using dnf:
  sudo dnf install ./kprobe_*.rpm
  ```

> [!NOTE]
> KProbe is installed to `/opt/kprobe`. The dedicated API server runs as a systemd service (`kprobe.service`) and is configured to start automatically on system boot.

- <b>2. Initialize the Database and Start the Service</b>
  - Once installed, initialize KProbe's local database and restart the background API service:
  ```bash
  kprobe db init
  sudo kprobe api restart
  ```

<br>
<br>

## Editor
- You can open the editor to create a scan configuration. The editor is available [here](https://github.com/K-cermak/KProbe/blob/main/web-editor/editor.html) (download the folder and open it in your browser).

> [!NOTE]
> You can also access the editor at `http://YOUR_SERVER_IP/editor`.

- The editor is straightforward to use – all the necessary information is provided within it. To download the configuration file, click on the **"Verify Values"** button and then on **"Download Config"**. To reload an existing configuration file, click on **"Load Config"**.

<img src="https://cdn.karlosoft.com/cdn-data/ks/img/kprobe/editor2.png" width="700" alt="KProbe Editor">



<br>
<br>

## Creating a Scan, Setting Up Cron
- After exporting the configuration file from the editor, load it into KProbe with this command:
```
kprobe config replace <path_to_config_file>
```

- If you want to verify the configuration file before loading it, run:
```
kprobe config verify <path_to_config_file>
```

> [!NOTE]
> The configuration file is stored in the database and the original file is not used by KProbe itself. You can safely delete it after loading.


<br>

- You can run scans manually with:
```
kprobe scan <type>
```

As a type, use:
- `all` – run all scans
- `all_except:<names>` – run all scans except the ones specified (separate names with commas, no spaces)
- `only:<names>` – run only the specified scans (separate names with commas, no spaces)

To set up automated scanning, open the cron editor:
```
crontab -e
```

And add, for example, the following line:
```
*/5 * * * * /usr/local/bin/kprobe scan all
```

This will run all scans every 5 minutes.



<br>
<br>


## API, Reverse Proxy
- The API runs on port **80** by default. You can access scan results at `http://YOUR_SERVER_IP/status/<scan_name>`.
- The response is in JSON format and looks like this:
```json
{
    "probe_name": "Probe Name",
    "time": "2025-01-01 01:23:59",
    "scan_name": "scan_name",
    "check": "2025-01-01T00:20:02Z",
    "status": "true"
}
```

| Field | Description |
|---|---|
| `probe_name` | Name of this probe instance |
| `time` | Current server time (YYYY-MM-DD HH:MM:SS) |
| `scan_name` | Name of the requested scan |
| `check` | Last scan execution time (ISO 8601) |
| `status` | `"true"` if the scan passed, `"false"` if it failed |

- You need to set up network access to the API server from your external monitoring dashboard (such as Uptime Kuma). Depending on your local network architecture, this can be achieved using a reverse proxy, port forwarding, or a VPN tunnel.
- If you want to change the API port, run:
```
kprobe keys set api_port <port>
sudo kprobe api restart
```

<br>
<br>

## Integration with External Monitoring (e.g., Uptime Kuma)
- KProbe provides a standard JSON endpoint for each scan, making it straightforward to connect with external monitoring dashboards and notification tools.
- To set up integration with **Uptime Kuma**, create a new monitor with these parameters:
    - <b>Type:</b> HTTP(s) - Keyword
    - <b>URL:</b> `http://YOUR_SERVER_IP/status/<scan_name>`
    - <b>Keyword:</b> `"status":"true"`

- Once set up, Uptime Kuma will query KProbe's API at regular intervals and handle alerts or notifications if any scan fails.

<br>
<br>

## List of All Commands
- You can display the full list of commands by running: `kprobe help`.

<br>

```
kprobe scan <type>
```
- Start a single pass of scans with the specified type.
    - Use `all` to run all scans.
    - Use `all_except:<names>` to run all scans except the specified ones (separate names with commas, no spaces).
    - Use `only:<names>` to run only the specified scans (separate names with commas, no spaces).

<br>

```
kprobe state
```
- View the current state of all scans.

<br>
<br>

```
kprobe history <scan_name> <from> <to>
```
- View the history of the specified scan.
- For `<from>` and `<to>`, use the format `YYYY-MM-DD HH:MM:SS`.

<br>
<br>

```
kprobe db init
```
- Initialize the database.

```
kprobe db reset
```
- Reset the database. **This will delete all data!**

<br>
<br>

```
kprobe config verify <path>
```
- Verify the configuration file at the specified path.

```
kprobe config replace <path>
```
- Replace the current configuration with the one at the specified path.
- The file is copied to the database, so you can delete the original file afterwards.

```
kprobe config view
```
- View the current configuration.

<br>
<br>

```
kprobe keys view all
```
- View all keys with their values in the database.
- Keys with the `*` prefix can be modified.

```
kprobe keys view <key>
```
- View the value of the specified key.
- If the key has the `*` prefix, it can be modified.

```
kprobe keys set <key> <value>
```
- Set the value of the specified key.

<br>
<br>

```
kprobe test ping <address> <timeout_ms>
```
- Test a ping to the specified address with the specified timeout.
- Timeout is in milliseconds.

```
kprobe test http <address> <timeout_ms>
```
- Test an HTTP request to the specified address with the specified timeout.
- Timeout is in milliseconds.

<br>
<br>

```
kprobe api test [service|http]
```
- Test the API service.
- Use `service` to test the API service via systemctl.
- Use `http` to test the API service via an HTTP request.

```
kprobe api restart
```
- Restart the API service.
- This command requires sudo privileges.

<br>
<br>


```
kprobe help
```
- Print this help message.

<br>
<br>

## Uninstallation
- To completely remove KProbe, run the package manager command for your system:

  **On Debian/Ubuntu-based systems:**
  ```bash
  sudo dpkg -r kprobe
  ```

  **On RHEL/CentOS/Rocky/Fedora-based systems:**
  ```bash
  sudo rpm -e kprobe
  # Or using dnf:
  sudo dnf remove kprobe
  ```

- This will stop the API service and clean up all binary and configuration files.

> [!NOTE]
> The SQLite database file at `/opt/kprobe/db.sqlite` is kept during uninstallation to prevent accidental data loss. If you want a completely clean removal, delete it manually: `sudo rm -rf /opt/kprobe`.

<br>
<br>

## Migration

### From Uptime Kuma Probe Extension (1.0.X) to KProbe (2.0.0)
Because KProbe does not include an automated migration tool, you must migrate your configuration manually. Please follow these steps:

1. **Uninstall the legacy application**: Run the uninstall script from the legacy repository to remove the old version:
   ```bash
   curl -sSL https://raw.githubusercontent.com/K-cermak/Uptime-Kuma-Probe/61ff4e8d88f423294359dbd472846eb81cde5d1b/scripts/uninstall.sh -o uninstall.sh
   sudo bash uninstall.sh
   rm uninstall.sh
   ```
2. **Install KProbe**: Download and install the new package (`.deb` or `.rpm`) following the [Installation](#installation) guide.
3. **Reinitialize the database**: Recreate the database structure. Note that database schema changes prevent a direct data migration, meaning previous scan history will be lost:
   ```bash
   kprobe db reset
   ```
4. **Upgrade your configuration**: Open your old configuration file, import or recreate your settings in the new KProbe configuration editor (since the configuration schema has been upgraded, e.g., the keyword system has been replaced by the new Variables + Expressions engine).
5. **Download the new configuration**: Verify your values in the editor and click **Download Config**.
6. **Apply the configuration**: Load your new configuration file into KProbe:
   ```bash
   kprobe config replace <path_to_new_config_file>
   ```
7. **Update your cron jobs / automated scripts**: The legacy command `kprobe cron` has been renamed to `kprobe scan`. If you use cron or custom scripts, edit them to call the new command:
   ```bash
   # Example cron line:
   */5 * * * * /usr/local/bin/kprobe scan all
   ```