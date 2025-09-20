# clima
Weather forecast TUI
- Written in Go.
- Built with the Bubble Tea framework.
- Integrated with the Open-Meteo forecast and geocoding HTTP APIs.
> The Open-Meteo APIs do not require a key, but are subject to usage limits.

## Develop
Run the program from the main file with `go run main.go`.

>You will not see logs in stdout due to the nature of TUI apps occupying that stream. Pass the `--debug` flag to make the program write the messages received by the `Update` function to a log file at `dev/debug.log`.

**Restart automatically on changes:**

Make sure `watchexec` is installed and available. Run the watch script.
```bash
./watch.sh
```
