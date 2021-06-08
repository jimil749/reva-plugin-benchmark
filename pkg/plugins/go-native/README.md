# Native Go Plugins

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