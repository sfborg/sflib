# `sflib`

`sflib` is a Go library containing shared functionality for Species File Group (SFG) projects. It primarily focuses on handling Species File Group Archives (SFGAs).

## Overview

This library provides tools for:

* Creating SFG Archives.
* Importing data into SFG Archives.
* Connecting to and managing SFG Archive databases.
* Checking compatibility between an application and the SFGA schema version.

## Installation

To use `sflib` in your Go project, run:

```bash
go get github.com/sfborg/sflib
```

## Usage Examples

### Working with Species File Group Archives (SFGA)

This section demonstrates how to interact with SFG Archives using `sflib`.

#### Creating a new SFGA

 ```go
import (
  "fmt"
  "github.com/sfborg/sflib/pkg/sfgaio"
)

func main() {
  sfga := sfgaio.New()
  err := sfga.Create(dif)
  ...
  _, err = sfga.Connect()
  ...
  defer sfga.Close()
  fmt.Println("SFGA created successfully.")
}
```

#### Import

This example demonstrates how to import data into an existing SFGA.

```go
import (
  "fmt"
  "github.com/sfborg/sflib/pkg/sfgaio"
)

func main() {
  sfga := sfgaio.New()
  err := sfga.SetDb(dbPath) // Set the path to the existing SFGA database
  ...
  err = sfga.Import(filePath) // Import data from filePath
  ...
  _, err = sfga.Connect() // Connect to the database
  ...
  defer sfga.Close() // Close the connection when done
  fmt.Println("Data imported successfully.")
}
```

#### Ad-hoc connection

This example demonstrates how to connect to an existing SFGA database without
creating a new one.

```go
import (
  "fmt"
  "github.com/sfborg/sflib/pkg/sfgaio"
)

func main() {
  sfga := sfgaio.New()
  err := sfga.SetDb(filePath) // Set the path to the existing SFGA database
  ...
  _, err = sfga.Connect() // Connect to the database
  ...
  defer sfga.Close() // Close the connection when done
  fmt.Println("Connected to SFGA successfully.")
}
```

## Testing

As the library modifies file system, running tests in parallel might
create running conditions, that will break some tests. To make sure
running only one thread with tests either use

```sh
make test
```

or run tests with `-p 1` option:

```sh
go test ./... -p 1
```
