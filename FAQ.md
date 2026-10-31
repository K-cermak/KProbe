# FAQ Guide
- [Installation](#installation)
- [Editor](#editor)
- [Creating a Scan, Setting Up Cron](#creating-a-scan-setting-up-cron)
- [API, Reverse Proxy](#api-reverse-proxy)
- [Integration with Uptime Kuma](#integration-with-uptime-kuma)
- [List of All Commands](#list-of-all-commands)
- [Uninstallation](#uninstallation)
- [Migration](#migration)

<br>

## Installation
- **System Requirements**: We recommend a dedicated virtual machine (VM) or physical server. The requirements are minimal: **1 GB RAM**, **1 CPU core** and **1 GB of free disk space** for the SQLite database (more may be needed if you run many scans or keep a long history).
- **Supported Operating Systems**: KProbe is compatible with any modern Linux distribution using **systemd** (`systemctl`). Native packages are available for:
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

- <b>2. Verify That the Service Is Running</b>
  - The database is initialized automatically during installation. To check that the API service is running, use:
  ```bash
  kprobe api test service
  ```

<br>
<br>

## Editor
- Use the editor to create a scan configuration. The editor is available [here](https://github.com/K-cermak/KProbe/tree/main/web-editor) (download the folder and open `index.html` in your browser).

> [!NOTE]
> You can also access the editor at `http://YOUR_SERVER_IP/editor`.

- The editor is straightforward to use – all the necessary information is provided within it. To download the configuration file, click **"Download Config"** (the values are verified first). To edit an existing configuration file, click **"Load Config"**.

<img src="https://cdn.karlosoft.com/cdn-data/ks/img/kprobe/editor2.png" width="700" alt="KProbe Editor">



<br>
<br>

## Creating a Scan, Setting Up Cron
- After exporting the configuration file from the editor, load it into KProbe with this command:
```
kprobe config replace <path_to_config_file>
```

- To verify the configuration file before loading it, run:
```
kprobe config verify <path_to_config_file>
```

> [!NOTE]
> The configuration is copied into the database, so KProbe does not need the original file. You can safely delete it after loading.


<br>

- You can run scans manually with:
```
kprobe scan <type>
```

Available types:
- `all` – run all scans
- `all_except:<names>` – run all scans except the ones specified (separate names with commas, no spaces)
- `only:<names>` – run only the specified scans (separate names with commas, no spaces)

To set up automated scanning, open the cron editor:
```
crontab -e
```

Then add a line such as:
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

- Your monitoring tool (such as Uptime Kuma) must be able to reach the API server. Depending on your network, you can achieve this with a reverse proxy, port forwarding or a VPN tunnel.
- To change the API port, run:
```
kprobe keys set api_port <port>
sudo kprobe api restart
```

<br>
<br>

## Integration with Uptime Kuma
- KProbe provides a standard JSON endpoint for each scan, so it is easy to connect to external monitoring dashboards and notification tools.
- To integrate KProbe with **Uptime Kuma**, create a new monitor with these settings:
    - <b>Type:</b> HTTP(s) - Keyword
    - <b>URL:</b> `http://YOUR_SERVER_IP/status/<scan_name>`
    - <b>Keyword:</b> `"status":"true"`

- Uptime Kuma will then query the KProbe API at regular intervals and send notifications whenever a scan fails.

<br>
<br>

## List of All Commands
- To display the full list of commands, run `kprobe help`.

<br>

```
kprobe scan <type>
```
- Run the selected scans once.
    - Use `all` to run all scans.
    - Use `all_except:<names>` to run all scans except the specified ones (separate names with commas, no spaces).
    - Use `only:<names>` to run only the specified scans (separate names with commas, no spaces).

<br>

```
kprobe state
```
- Show the latest result of each scan.

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
- This runs automatically during package installation.
- If you need to run it manually, use `sudo` (the database is stored in `/opt/kprobe/`).

```
sudo kprobe db reset
```
- Reset the database. **This will delete all data!**
- Requires `sudo` privileges.

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
- Show all keys and their values.
- Keys with the `*` prefix can be modified.

```
kprobe keys view <key>
```
- Show the value of the specified key.
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
- Ping the specified address.
- The timeout is in milliseconds and applies to each packet.

```
kprobe test http <address> <timeout_ms>
```
- Send an HTTP request to the specified address.
- The timeout is in milliseconds.

<br>
<br>

```
kprobe api test [service|http]
```
- Check the API service.
- Use `service` to check the service status via systemctl.
- Use `http` to send a test HTTP request to the API.

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
- Print the help message.

<br>
<br>

## Uninstallation
- To remove KProbe, run the command for your package manager:

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

- This stops the API service and removes the binaries and system files.

> [!NOTE]
> The SQLite database file at `/opt/kprobe/db.sqlite` is kept during uninstallation to prevent accidental data loss. For a completely clean removal, delete the directory manually: `sudo rm -rf /opt/kprobe`.

<br>
<br>

## Migration

### From Uptime Kuma Probe Extension (1.0.X) to KProbe (2.0.0)
KProbe does not include an automatic migration tool, so the configuration has to be migrated manually:

1. **Uninstall the legacy application**: Run the uninstall script from the legacy repository to remove the old version:
   ```bash
   curl -sSL https://raw.githubusercontent.com/K-cermak/Uptime-Kuma-Probe/61ff4e8d88f423294359dbd472846eb81cde5d1b/scripts/uninstall.sh -o uninstall.sh
   sudo bash uninstall.sh
   rm uninstall.sh
   sudo rm -rf /opt/kprobe
   ```
2. **Install KProbe**: Download and install the new package (`.deb` or `.rpm`) following the [Installation](#installation) guide.
3. **Upgrade your configuration**: Load your old configuration file into the new configuration editor, or recreate your scans there. The configuration format has changed – for example, the keyword system has been replaced by the new Variables + Expressions engine. The editor converts legacy `keyword` parameters automatically.
4. **Download the new configuration**: Verify your values in the editor and click **Download Config**.
5. **Apply the configuration**: Load your new configuration file into KProbe:
   ```bash
   kprobe config replace <path_to_new_config_file>
   ```
6. **Update your cron jobs and scripts**: The legacy command `kprobe cron` has been renamed to `kprobe scan`. If you use cron or custom scripts, update them to call the new command:
   ```bash
   # Example cron line:
   */5 * * * * /usr/local/bin/kprobe scan all
   ```