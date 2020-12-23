# README

This repository is modeled loosely after a takehome task I had for an interview I did with Unity Technologies. As a Node.js developer, I did the task in Node, but I would like to learn Golang. This seemed like a perfect opportunity to leverage my REST experience so that I could learn Golang.

The task is to do the following:
1. [ ] An endpoint where users can leave feedback for a specific game session
   1. [ ] A user can only submit one review per game session
   2. [ ] User MUST leave a rating of 1-5 if providing feedback
   3. [ ] User MAY add a comment when providing feedback
   4. [ ] Multiple players can rate the same session
2. [ ] Folloing RESTful principles, create HTTP endpoints to allow:
   1. [ ] Players to add feedback for a session
   2. [ ] Ops team members to see recent feedback left by players
   3. [ ] Allow filtering by rating
3. [ ] Include a README that includes at least the following:
   1. [ ] API Documentation
   2. [ ] Instructions for launching and testing your API locally (if not the built-in scripts)
4. [ ] Bonus items
   1.  [ ] A simple front-end
   2.  [ ] Tests
   3.  [ ] Deployment scripts/tools
   4.  [ ] Authentication
   5.  [ ] User permissions

## Original README

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