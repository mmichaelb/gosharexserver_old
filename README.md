# gosharexserver [![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT) [![stability-beta](https://img.shields.io/badge/stability-beta-33bbff.svg)](https://github.com/mkenney/software-guides/blob/master/STABILITY-BADGES.md#beta) [![GoDoc](https://godoc.org/github.com/mmichaelb/gosharexserver?status.svg)](https://godoc.org/github.com/mmichaelb/gosharexserver) [![Build Status](https://travis-ci.org/mmichaelb/gosharexserver.svg?branch=master)](https://travis-ci.org/mmichaelb/gosharexserver) [![Go Report Card](https://goreportcard.com/badge/github.com/mmichaelb/gosharexserver)](https://goreportcard.com/report/github.com/mmichaelb/gosharexserver)
Lightweight upload server for the ShareX client (https://getsharex.com/).

# Description
This application can be used as a standalone server side endpoint for your ShareX client. It is written in Go and designed to be lightweight and easy to understand. If you are a Golang developer, you can also use this project as your dependency and use the code in your own project.

# Features
- [x] upload and share images or in general files with the ShareX client
- [x] MongoDB GridFS file storage
- [ ] MySQL-driven file storage
- [x] mime type whitelisting
- [x] run behind a reverse proxy
- [x] delete entries
- [x] limit access by offering authorization 
- [ ] user system
- [x] Docker image/compose 

# Installation
## Getting the binaries
In order to install the ShareX server you have to get the binaries. There two possible methods of getting them:
- download a release file from the [GitHub releases page](https://github.com/mmichaelb/gosharexserver/releases)
- compile the source manually on your own (see [Compilation](https://github.com/mmichaelb/gosharexserver#compilation))
## Download default configuration files
In order to adjust values of the application's runtime, you should download the default configurations to get an orientation. The downloads can be found in the [config directory](https://github.com/mmichaelb/gosharexserver/tree/master/configs). After downloading the configuration you should rename it and adjust the values according to the [TOML conventions](https://github.com/toml-lang/toml).
## Running the application
At the moment the only parameter which the application accepts on startup is -config - you can specify the path to your configuration file. If you do not specify one, the default path (`./config.toml`) is used. An example of running the application would be:
```bash
./gosharexserver-executable -config=./my-custom-config.toml
```
Have fun and feel free to open up an issue if you have a problem with running your application.

# Installation with docker compose
When using Docker compose, the installer script can be used:
```bash
mkdir /opt/gosharexserver
cd /opt/gosharexserver
curl -s https://raw.githubusercontent.com/mmichaelb/gosharexserver/master/scripts/docker-compose-installer.sh | bash
```

# Compilation
The compilation of this code was successful with Go `1.8`-`1.10.1`.

In general, there are two ways of building the application:
## Makefile
When using `make`, life is easy and you can just run and a executable should be dropped in your working directory:
```bash
make build
```
## go build command
When compiling with the standard `go build` command, you can use the extracted command from the Makefile. Because with `make` the ld flags are parsed automatically, you have to replace them on your own when running `go build` manually.
```bash
go build -ldflags '-X "main.applicationName=gosharexserver" -X "main.version=<version>" -X "main.branch=<branch>" -X "main.commit=<commit>"' ./cmd/gosharexserver
```

# Using ShareX server as a dependency
To use this project as a dependency for your own project, you can just `go get` the `cmd/gosharexserver` package:
```bash
go get -u github.com/mmichaelb/gosharexserver/cmd/gosharexserver
```
Make sure to check out the [examples package](https://github.com/mmichaelb/gosharexserver/tree/master/examples/) for implemented examples and use cases.

# Example configuration for ShareX client
```
{
  "DestinationType": "ImageUploader, TextUploader, FileUploader",
  "RequestURL": "http://example.com/upload",
  "FileFormName": "file",
  "Headers": {
    "Authorization": "1337#Secure_Token"
  },
  "URL": "http://example.com/$json:call_reference$",
  "DeletionURL": "http://example.com/delete/$json:delete_reference$"
}
```

# Contribution
Feel free to contribute and help this project to grow. You can also just suggest features/enhancements - for more details check the [contributing file](https://github.com/mmichaelb/gosharexserver/tree/master/.github/CONTRIBUTING.md).
