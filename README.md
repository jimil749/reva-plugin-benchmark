# Reva Plugin Benchmarking

A comparison of plugin implementations for Golang. The reason for this is to evaluate plugin options for [Reva](https://github.com/cs3org/reva), which is an interoperability platform, connecting storage, sync and share platforms and application providers. This benchmarking will give a good idea about viability of each plugin framework and its performance.

The following packages have been (ongoing) benchmarked:
- [Hashicorp go-plugin](https://github.com/hashicorp/go-plugin)
- [Natefinch pie-plugin](https://github.com/natefinch/pie)
- [Native golang Plugin](https://golang.org/pkg/plugin/)
- [Traefik Yaegi](https://github.com/traefik/yaegi)
- [Goloader](https://github.com/pkujhd/goloader)

# Plugin

The "core" program exposes the following interface (defined in `pkg/shared/interface.go`) for the plugin to implement:

```go
type Manager interface {
	OnLoad(userFile string) error
	GetUser(*userpb.UserId) (*userpb.User, error)
	GetUserByClaim(claim, value string) (*userpb.User, error)
	GetUserGroups(*userpb.UserId) ([]string, error)
	FindUsers(query string) ([]*userpb.User, error)
}
```

The responsibility of the plugin is to implement the above interface and make the implementation available to the main program. The framework used for communication b/w the host ("core") and the plugin is being benchmarked.


 
# Benchmarks

For the purpose of benchmarking, this project uses the existing [reva plugin driver](https://github.com/cs3org/reva/tree/master/pkg/user/manager/json). Details about each plugin package and the corresponding benchmarks can be found in the respective directories, which can be accessed from the directory structure below.

# File Directory Structure
- __reva\-plugin\-benchmark__
   - [README.md](README.md)
   - __file__
     - [user.demo.json](file/user.demo.json)
   - [go.mod](go.mod)
   - [go.sum](go.sum)
   - [hashicorp\-plugin](hashicorp-plugin)
   - [main.go](main.go)
   - [main\_test.go](main_test.go)
   - __pkg__
     - __plugins__
       - __hashicorp__
         - __netrpc__
           - [README.md](pkg/plugins/hashicorp/netrpc/README.md)
           - [main.go](pkg/plugins/hashicorp/netrpc/main.go)
     - __shared__
       - [interface.go](pkg/shared/interface.go)
       - [rpc.go](pkg/shared/rpc.go)
   - [run.sh](run.sh)

# Running locally

In order to run the benchmarks locally, follow:

1. Clone this repository
```
git clone https://github.com/jimil749/reva-plugin-benchmark.git
cd reva-plugin-benchmark
```
2. Run `run.sh`
```
./run.sh
```

That's it!