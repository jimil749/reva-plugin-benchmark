# Native Go Plugins

A Go plugin is a Go main package with exported functions and variables that has been compiled using the `-buildmode=plugin` build flag to produce a shared object (`.so`) library file. The exported functions and variables, in the Go package can be looked up and be bound to at runtime using the `plugin` package.

# Usage

A Go plugin is just a Go package. The only difference is how it is built. The plugin package is supposed to export functions/variables, which can be used by the main program. A go plugin package can be compiled using the following command:
```
go build -buildmode=plugin -o plug.so ./test.go
```
This command will generate an shared obejct library from `test.go` by the name of `plug.so`. All the exported function/variable can be accessed from the main program using the plugin package: 

1. Open the plugin file with `plugin.Open()`
2. Lookup the exported symbol with `plugin.Lookup(<Symbol name>)`
3. Assert that loaded symbol is of a desired type.
4. Use the module

# Pros

1. Super Fast! Runs as fast as the native go code
2. Easy to use.
3. Uses interface type as boundary

# Cons

1. Only supported on Linux
2. Plugin once loaded cannot be "unloaded".
3. Version dependency: All the packages used by the plugin and the host should have the same version. (Making the host and the plugin tightly coupled.)

# Benchmarks

This section contains the benchmark of each of the methods of the `UserManager` interface that are called using the native go plugin.

| Method Name                       | Operations  | ns/op       |
| -------------------------- |:-----------:| -----------:|
| GetUser()                  | 16151228        | 69.20 ns/op  |
| GetUserByClaim()           | 100000000       | 10.73 ns/op |
| GetUserGroups()            | 16016910        | 73.48 ns/op |
| FindUser()                 | 1769547         | 677.9 ns/op |

```
$ ./run.sh
goos: linux
goarch: amd64
pkg: github.com/jimil749/reva-plugin-benchmark
cpu: Intel(R) Core(TM) i5-7200U CPU @ 2.50GHz
BenchmarkGoPlugin/GetUser-4     	        16151228	        69.20 ns/op
BenchmarkGoPlugin/GetUserByClaim-4         	100000000	        10.73 ns/op
BenchmarkGoPlugin/GetUserGroups-4          	16016910	        73.48 ns/op
BenchmarkGoPlugin/FindUser-4               	 1769547	       677.9 ns/op

```