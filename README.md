# `sflib`

This library contains functionality shared by Species File Group Go projects.

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
