# simplewall utilities

> Simple utilities for [simplewall](https://github.com/henrypp/simplewall)

## Introduction

[Windows Filtering Platform](https://docs.microsoft.com/en-us/windows/win32/fwp/windows-filtering-platform-start-page), which simplewall is based on, does not allow rules based on process names without a full path. This becomes an issue when you have a program that drops a randomly named updater executable, which exits before you can allow it in the firewall.

To work around this, you can look at the logs, and allow the destination IPs globally. This can be done by hand, but it's a **very** tedious process, especially if you care about de-duplicating your IPs.

Discussion: [simplewall#136](https://github.com/henrypp/simplewall/issues/136)

## Usage

This program allows you to parse the log file and generate unique rules for a process name that you specify, with or without full path.

1. Enable logging from simplewall: \
   `Settings > Packets log > Enable packets logging to a file`

2. Let the problematic program run and fail. Don't allow or deny the prompt, simply close or ignore it.

3. Exit simplewall from the tray to make sure it doesn't overwrite our changes

4. Use this tool:

   ```powershell
   .\simplewall-utils.exe allow --help

   Usage:
   simplewall-utils allow [flags]

   Flags:
   -a, --append                Append to existing rules instead of overwriting
   -h, --help                  help for allow
   -l, --log-path string       Path to log file (default "%USERPROFILE%\\simplewall.log")
   -n, --process-name string   Process name to allow
   -p, --profile-path string   Path to profile file (default "%APPDATA%\\Henry++\\simplewall\\profile.xml")
   ```

5. Start simplewall once again. Observe the new rules appear under the `User rules` tab.

## Example

```powershell
  simplewall-utils allow -n "docker desktop installer.exe"
  simplewall-utils allow -n "backgrounddownload.exe"
```
