# Pie Plugin

Package pie provides a toolkit for creating plugins for Go applications. Pie plugins communicate with the original process via RPC over the plugin application's Stdin and Stdout.

# Usage

Pie plugin package provides 2 kinds of functions:
1. `NewProvider`: Executed by the plugin to set-up it's end of the communications. `NewProvider` starts an JSON RPC Server, which is responsible for serving plugin "methods" over RPC
2. `StartProvider`: Executed by the main application to set up it's end of the communications and start a plugin executable. The main application is the RPC client, consuming the methods served by the plugin. (which is the RPC Client).

Method Definition [here](https://github.com/natefinch/pie#func-startprovider)

# Pros

1. Uses RPC for communication, hence the plugin code can easily be "unloaded" by killing the client. (Hot Reloading can easily be applied here.)
2. Simple API driven plugins, which are intuitive and easy to write.
3. Flexible, in the sense that the plugin shouldn't necesserily be in Go, can be in any language as long as the plugin can provide an RPC API.

# Cons

1. Performance Overhead. 
2. No longer maintained.
3. Lack of Documentation.

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
$ ./run.sh
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