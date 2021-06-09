# Pie Plugin


# Benchmarks

This section contains the benchmark of each of the methods of the `Manager` interface that are called using the pie go-plugin framework.

| Method Name                       | Operations  | ns/op       |
| -------------------------- |:-----------:| -----------:|
| OnLoad()                   | 10000        | 109514 ns/op  |
| GetUser()                  | 16585       | 72331 ns/op |
| GetUserByClaim()           | 16987       | 70628 ns/op |
| GetUserGroups()            | 22274       | 51122 ns/op |
| FindUser()                 | 16581      | 75336 ns/op |


```
goos: linux
goarch: amd64
pkg: github.com/jimil749/reva-plugin-benchmark
cpu: Intel(R) Core(TM) i5-7200U CPU @ 2.50GHz
BenchmarkPiePlugin/OnLoad-4         	   10000	    109514 ns/op
BenchmarkPiePlugin/GetUser-4        	   16585	     72331 ns/op
BenchmarkPiePlugin/GetUserByClaim-4 	   16987	     70628 ns/op
BenchmarkPiePlugin/GetUserGroups-4  	   22274	     51122 ns/op
BenchmarkPiePlugin/FindUser-4       	   16581	     75336 ns/op
PASS
ok  	github.com/jimil749/reva-plugin-benchmark	8.634s

```