# Examples

Each subdirectory holds a small runnable `main` package that exercises one
part of the API. Every example accepts the same connection flags:

`-endpoint <url>`  connect to a running Firefox Remote Agent
`-exec <path>`     path to the Firefox binary to launch
`-timeout <d>`     overall command timeout (default 2m)

A BiDi-capable Firefox (Firefox Nightly with the Remote Agent enabled) is
required. To point an example at an installed Nightly binary:

    go run ./examples/basic -exec /path/to/firefox

To connect to a Firefox you started yourself, launch it first:

    firefox --remote-debugging-port 9222 --headless

then run, for example:

    go run ./examples/connect -endpoint http://127.0.0.1:9222

| Directory   | Demonstrates                                     |
| ----------- | ------------------------------------------------ |
| connect/    | WebDriver classic handshake, `session.status`    |
| basic/      | open a page, navigate, screenshot                |
| evaluate/   | run JS, decode `RemoteValue`                     |
| input/      | click, type text, key presses                    |
| events/     | subscribe to browsingContext/script events       |
| network/    | intercept and modify requests                    |
| screenshot/ | viewport and full-page capture                   |