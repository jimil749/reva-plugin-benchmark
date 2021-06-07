# Hashicorp Go-Plugin

Hashicorp `go-plugin` is a Go Plugin system over RPC. The plugin system is only designed to work over a local reliable network. For more information about the go-plugin architecture and features visit https://github.com/hashicorp/go-plugin#readme

# Usage

To use the plugin system, following steps should be followed (These are high level steps): 

1. Choose the interface you want to expose for the plugins. (`Manager` in our case, defined in `pkg/shared/interface.go`)
2. For each interface, implement an implementation of that interface that communicates over a `net/rpc` connection. Both, the client and server implementation are required for the communication to take place. (defined in `pkg/shared/rpc.go`)
3. Create a `Plugin` implementation that knows how to create the RPC client/server for a given plugin type. `Plugin` is an inbuilt type defined by hashicorp/go-plugins. 
4. Plugin authors i.e the plugin code calls `plugin.Serve` to serve the plugin.
5. Plugin users i.e the core program use `plugin.Client` to launch a subprocess and request an interface implementation over RPC

# Pros

1. Plugins use RPC/gRPC to communicate with the host process, hence plugins can't crash the host process. A panic in a plugin doesn't panic the plugin user.
2. Plugins are a Go interface implementation: hence writing and consuming plugins feel very "natural"!
3. Used in production by many hashicorp projects: Vault, Terraform etc

More features [here](https://github.com/hashicorp/go-plugin#features)

# Cons
1. Performance overhead is HUGE! (Refer the benchmarks below)
2. Lack of documentations

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