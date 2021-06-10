# Benchmarks

```
$ ./run.sh
goos: linux
goarch: amd64
pkg: github.com/jimil749/reva-plugin-benchmark
cpu: Intel(R) Core(TM) i5-7200U CPU @ 2.50GHz
BenchmarkHashicorpPlugingRPC/OnLoad-4         	    5725	    185584 ns/op
BenchmarkHashicorpPlugingRPC/GetUser-4        	   10000	    105572 ns/op
BenchmarkHashicorpPlugingRPC/GetUserByClaim-4 	   12775	     93748 ns/op
BenchmarkHashicorpPlugingRPC/GetUserGroups-4  	   12098	     99894 ns/op
BenchmarkHashicorpPlugingRPC/FindUser-4       	   12608	     96465 ns/op
PASS
ok  	github.com/jimil749/reva-plugin-benchmark	8.712s

```