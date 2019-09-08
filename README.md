# Update Go

Program to update your installation of Golang. It will also install Go from scratch if not installed.

This program automates the [recommended steps](https://golang.org/doc/install#install) by the Go Team to install / update Golang. It checks if there's a new version of Go for your machine, and if there is, downloads the installer / archive with the new version, removes the old one and installs the new one.

## Installation and Usage

The following methods are available to install this program on your computer. Choose the one you are most comfortable with.

### Binary releases

This is probably the easiest way to install the program since it has no pre-requisites.

* Go to the [releases section](https://github.com/andrerfcsantos/update-go/releases) of this repository
* Download the binary for your OS / Architecture
* Run the binary
* (Optional) Change the name of the binary and add it to your PATH for easier use in the future

### Installation via `go get`

Alternatively, if you have Go installed, you can install the program with `go get`:

  ```
  $ go get github.com/andrerfcsantos/update-go/cmd/updatego
  $ updatego
 ```

**Note:** This method assumes the place where Go places binaries is already on your PATH. Usually this will be your `$GOPATH/bin` directory, which is typically `~/go/bin` or `/usr/local/go/bin` on Unix systems and `C:\Users\<your_username>\go\bin` on Windows systems.

### Building from source

You can also compile the binary from source. This method also requires Go previously installed.

  ```
  $ git clone git@github.com:andrerfcsantos/update-go.git
  $ cd update-go
  $ go build ./...
  $ ./updatego
 ```
 