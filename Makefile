GOCMD = go
GOTEST = $(GOCMD) test 
TEST_OPTS = -count=1 -p 1 -shuffle=on -coverprofile=coverage.txt -covermode=atomic

all: test

## Test
test: ## Run the tests of the project
	$(GOTEST) $(TEST_OPTS) ./...

