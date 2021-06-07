# Hashicorp Go-Plugin

# Benchmarks

This section contains the benchmark of each of the methods of the `Manager` interface that are called using the hashicorp go-plugin framework.

| Method Name                       | Operations  | ns/op       |
| -------------------------- |:-----------:| -----------:|
| OnLoad()                   | 7863        | 145501 ns/op  |
| GetUser()                  | 14107       | 85400 ns/op |
| GetUserByClaim()           | 14268       | 84915 ns/op |
| GetUserGroups()            | 15367       | 79499 ns/op |
| FindUser()                 | 264476      | 6429 ns/op |


```
$ ./run.sh
goos: linux
goarch: amd64
pkg: github.com/jimil749/reva-plugin-benchmark
cpu: Intel(R) Core(TM) i5-7200U CPU @ 2.50GHz
BenchmarkHashicorpPluginRPC/OnLoad-4         	    7863	    145501 ns/op
BenchmarkHashicorpPluginRPC/GetUser-4        	   14107	     85400 ns/op
BenchmarkHashicorpPluginRPC/GetUserByClaim-4 	   14268	     84915 ns/op
BenchmarkHashicorpPluginRPC/GetUserGroups-4  	   15367	     79499 ns/op
BenchmarkHashicorpPluginRPC/FindUser-4       	  264476	      6429 ns/op
PASS
ok  	github.com/jimil749/reva-plugin-benchmark	9.066s

```