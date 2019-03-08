# GoHealthz

Web application that provide you a simple interface to show healthiness of the
website.

## Requirements

- [GNU Make](https://www.gnu.org/software/make/manual/make.html)
- [Golang](https://golang.org/)
- [Google's UUID](github.com/google/uuid): Install it using `go get -u github.com/google/uuid`

## How to Run, Build, and Clean Up

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

## Playground

When the application is running, there's an endpoint called `/swagger` that
will provide playground for REST API for this application. It is provided so
that developer can easily play with REST API of this apps without having
to test it manually through postman, curl, etc. One thing to note is that
the API specs (open API) is not automatically generated from the routes
of the application, it is manually written _as for now_.

## Library Dependencies

Currently, the only external package that being imported by this application
is only [Google's uuid](github.com/google/uuid). It is used to generate UUID
for the database.

## Database

This apps, for now is only support in-memory database, which will works out
of the box withouth having to install other third party database storage.
However, the pakcage for the storage is modular and layered with an `interface`
so that developer can easily writes and switch the database driver without
having to worry changing so many lines of code.
