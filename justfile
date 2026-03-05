# Just commands for sflib

# Variables
gocmd := "go"
gotest := gocmd + " test"
test_opts := "-count=1 -p 1 -shuffle=on -coverprofile=coverage.txt -covermode=atomic"

# Default recipe (runs when you just type 'just')
default: test

# Run the tests of the project
test:
    {{gotest}} {{test_opts}} ./...
