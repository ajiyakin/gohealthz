# GoHealthz

Web application that provide you a simple interface to show healthiness of the
website.

## Requirements

- [GNU Make](https://www.gnu.org/software/make/manual/make.html)
- [Golang](https://golang.org/)

## How to Run

In order to run the apps, please build it first so it will produce binary
executable file.

### Build

`make build`

This will produce executable binary file under `cmd/gohealthz` folder, called
`gohealthz`.

### Run

After build, the executable can be executed directly by executing the binary
file itself. However, there's a command to start the server by simply run
this command:

`make run`

### Clean

If there's a need to clean all the resources created during build and run,
use this command:

`make clean`