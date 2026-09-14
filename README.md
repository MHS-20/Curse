# curse

Swear at your terminal and it swears back.

The **Trust-Me-Bro-Studies** shown that imprecating in your terminal reduce stress by 65% 
and lowers burnout rate by 50%.

`curse` draws a curse word as big ASCII art that dances in place, changing color on every frame.

## Usage

```sh
curse               # pick a random curse from the built-in list
curse "dio paguro"  # draw your own
```

Press `Ctrl+C` to stop.

## Install globally

Requires Go 1.21 or newer.

```sh
make build
sudo make install
```

Build as your normal user first: `sudo` usually can't find `go`, and `make install` only copies the already-built binary to `/usr/local/bin/curse`.

To install somewhere else, for example without `sudo`:

```sh
make install PREFIX=~/.local/bin
```

To remove it: `sudo make uninstall`.
