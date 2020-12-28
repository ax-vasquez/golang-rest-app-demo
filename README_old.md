## Original README

> This is only included for reference - the "working" documentation for this project is all in [`README.md`](README.md)

### Golang Starter Framework

#### Getting Started With Go
We currently use the latest version of Go 1.13. You can [download it from the go website](https://golang.org/dl/). 
Click on "Archived versions" to find 1.13.

Why do we use it? Go 1.13 is the latest version supported by Google Cloud Functions. 
Since we use their cloud functions we pin our development version to theirs.

Feel free to use the latest stable version. It'll probably work just fine, but please update the version in go.mod
if you do so, so we can easily test your code.
 
#### Makefile targets
We've created a basic Makefile that is already setup with standard go tools. The following targets are available:

**clean**: cleans the project

**build**: builds cmd/main.go

**run**: starts the project webserver (defaults to port 8080)

**lint**: lints your project for formatting issues

**vet**: looks for suspicious constructs

**secure**: looks for security problems

**test**: runs tests in the project

**show-coverage**: opens a web browser with a report of testing coverage

**race**: tests for race conditions