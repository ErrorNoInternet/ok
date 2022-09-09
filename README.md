# ok
This is a command that simply prints "ok" onto your screen whenever you run the `ok` command

`ok` 2.0 is in development! The new version will be written in Rust, however the database (and server) will not be backwards-compatible.

![Screenshot](https://raw.githubusercontent.com/ErrorNoInternet/ok/main/ok.png)

--------------------

## Installation (Linux)
Download the latest release into your Downloads folder and open a shell.
```sh
chmod +x ~/Downloads/ok
sudo cp ~/Downloads/ok /usr/bin/
```

## Installation (Windows)
Download `ok.exe` from the releases tab and copy it to any folder in `%PATH%` (for example `system32`).

## Installation (Android)
Download `ok-aarch64` from v1.4.2-termux into your downloads folder and open Termux.
```sh
cd
cp storage/downloads/ok-aarch64 ./ok
chmod 777 ok
alias ok=./ok
```

--------------------

Use `ok help` to get a list of commands

<sub>If you would like to modify or use this repository (including its code) in your own project, please be sure to credit!</sub>
