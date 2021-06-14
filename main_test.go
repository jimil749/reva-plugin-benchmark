package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/rpc/jsonrpc"
	"os"
	"os/exec"
	"plugin"
	"reflect"
	"testing"

	userpb "github.com/cs3org/go-cs3apis/cs3/identity/user/v1beta1"
	"github.com/hashicorp/go-hclog"
	hashPlugin "github.com/hashicorp/go-plugin"
	"github.com/jimil749/reva-plugin-benchmark/pkg/manager"
	"github.com/jimil749/reva-plugin-benchmark/pkg/plugins/yaegi"
	"github.com/jimil749/reva-plugin-benchmark/pkg/shared"
	"github.com/natefinch/pie"
	"github.com/traefik/yaegi/interp"
	"github.com/traefik/yaegi/stdlib"
)

func BenchmarkYaegi(b *testing.B) {
	i := interp.New(interp.Options{
		GoPath: "/home/jimil/go",
	})

	i.Use(stdlib.Symbols)
	i.Use(yaegi.ManagerSymbols())

	fmt.Println(i.Symbols("github.com/jimil749/reva-plugin-benchmark/pkg/manager/manager"))

	_, err := i.Eval(`import "github.com/jimil749/reva-yaegi-benchmark"`)
	if err != nil {
		fmt.Println("failed to import plugin: ", err)
		panic(err)
	}

	_, err = i.Eval(`package wrapper
import (
	json "github.com/jimil749/reva-yaegi-benchmark"
	"github.com/jimil749/reva-plugin-benchmark/pkg/manager"
)

func NewWrapper(userFile string) (manager.UserManager, error) {
	mgr, err := json.New(userFile)
	var m manager.UserManager = mgr
	return m, err
}
`)

	if err != nil {
		fmt.Println("error: ", err)
		panic(err)
	}

	fnNew, err := i.Eval("wrapper.NewWrapper")
	if err != nil {
		panic(err)
	}

	userFile := "./file/user.demo.json"
	args := []reflect.Value{reflect.ValueOf(userFile)}
	results := fnNew.Call(args)

	if len(results) > 1 && results[1].Interface() != nil {
		panic(results[1].Interface().(error))
	}

	mgr, ok := results[0].Interface().(manager.UserManager)
	if !ok {
		panic(fmt.Errorf("invalid interface type"))
	}

	b.Run("GetUser", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, _ = mgr.GetUser(&manager.UserId{OpaqueId: "4c510ada-c86b-4815-8820-42cdf82c3d51", Idp: "cernbox.cern.ch"})
		}
	})

	b.Run("GetUserByClaim", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, _ = mgr.GetUserByClaim("mail", "einstein@cern.ch")
		}
	})
	b.Run("GetUserGroups", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, _ = mgr.GetUserGroups(&manager.UserId{OpaqueId: "4c510ada-c86b-4815-8820-42cdf82c3d51", Idp: "cernbox.cern.ch"})
		}
	})
	b.Run("FindUser", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, _ = mgr.FindUsers("einstein")
		}
	})
}

// BenchmarkGoPlugin benchmarks the native go-plugin
func BenchmarkGoPlugin(b *testing.B) {
	// Open the plugin shared object library
	sym, err := plugin.Open("go-plugin.so")
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	// Lookup the "new" function that returns the
	fn, err := sym.Lookup("New")
	if err != nil {
		panic(err)
	}

	newFn, ok := fn.(func(string) (shared.UserManager, error))
	if !ok {
		panic("unexpected type from module symbol")
	}

	manager, err := newFn("./file/user.demo.json")
	if err != nil {
		panic(err)
	}

	b.Run("GetUser", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, _ = manager.GetUser(&userpb.UserId{OpaqueId: "4c510ada-c86b-4815-8820-42cdf82c3d51", Idp: "cernbox.cern.ch"})
		}
	})

	b.Run("GetUserByClaim", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, _ = manager.GetUserByClaim("mail", "einstein@cern.ch")
		}
	})
	b.Run("GetUserGroups", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, _ = manager.GetUserGroups(&userpb.UserId{OpaqueId: "4c510ada-c86b-4815-8820-42cdf82c3d51", Idp: "cernbox.cern.ch"})
		}
	})
	b.Run("FindUser", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, _ = manager.FindUsers("einstein")
		}
	})

}

// BenchmarkHashicorpPlugingRPC benchmarks hashicorp gRPC plugin
func BenchmarkHashicorpPlugingRPC(b *testing.B) {
	log.SetOutput(ioutil.Discard)

	logger := hclog.New(&hclog.LoggerOptions{
		Name:   "plugin",
		Output: os.Stdout,
		Level:  hclog.NoLevel,
	})

	// We're a host. Start by launching the plugin process.
	client := hashPlugin.NewClient(&hashPlugin.ClientConfig{
		HandshakeConfig: shared.Handshake,
		Plugins:         shared.PluginMap,
		Cmd:             exec.Command("./hashicorp-plugin-grpc"),
		AllowedProtocols: []hashPlugin.Protocol{
			hashPlugin.ProtocolNetRPC, hashPlugin.ProtocolGRPC},
		Logger: logger,
	})
	defer client.Kill()

	// Connect via RPC
	rpcClient, err := client.Client()
	if err != nil {
		fmt.Println("Error:", err.Error())
		os.Exit(1)
	}

	// Request the plugin
	raw, err := rpcClient.Dispense("json_grpc")
	if err != nil {
		fmt.Println("Error:", err.Error())
		os.Exit(1)
	}

	// We should have the manager now! This feels like a normal interface
	// implementation but is in fact over an RPC connection.
	manager := raw.(shared.Manager)

	b.Run("OnLoad", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = manager.OnLoad("./file/user.demo.json")
		}
	})
	b.Run("GetUser", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, _ = manager.GetUser(&userpb.UserId{OpaqueId: "4c510ada-c86b-4815-8820-42cdf82c3d51", Idp: "cernbox.cern.ch"})
		}
	})
	b.Run("GetUserByClaim", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, _ = manager.GetUserByClaim("mail", "einstein@cern.ch")
		}
	})
	b.Run("GetUserGroups", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, _ = manager.GetUserGroups(&userpb.UserId{OpaqueId: "4c510ada-c86b-4815-8820-42cdf82c3d51", Idp: "cernbox.cern.ch"})
		}
	})
	b.Run("FindUser", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, _ = manager.FindUsers("einstein")
		}
	})
}

// BenchmarkHashicorpPluginRPC benchmarks hashicorp rpc plugin.
func BenchmarkHashicorpPluginRPC(b *testing.B) {
	log.SetOutput(ioutil.Discard)

	logger := hclog.New(&hclog.LoggerOptions{
		Name:   "plugin",
		Output: os.Stdout,
		Level:  hclog.NoLevel,
	})

	// We're a host. Start by launching the plugin process.
	client := hashPlugin.NewClient(&hashPlugin.ClientConfig{
		HandshakeConfig: shared.Handshake,
		Plugins:         shared.PluginMap,
		Cmd:             exec.Command("./hashicorp-plugin"),
		AllowedProtocols: []hashPlugin.Protocol{
			hashPlugin.ProtocolNetRPC},
		Logger: logger,
	})
	defer client.Kill()

	// Connect via RPC
	rpcClient, err := client.Client()
	if err != nil {
		fmt.Println("Error:", err.Error())
		os.Exit(1)
	}

	// Request the plugin
	raw, err := rpcClient.Dispense("json")
	if err != nil {
		fmt.Println("Error:", err.Error())
		os.Exit(1)
	}

	// We should have the manager now! This feels like a normal interface
	// implementation but is in fact over an RPC connection.
	manager := raw.(shared.Manager)

	b.Run("OnLoad", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = manager.OnLoad("./file/user.demo.json")
		}
	})
	b.Run("GetUser", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, _ = manager.GetUser(&userpb.UserId{OpaqueId: "4c510ada-c86b-4815-8820-42cdf82c3d51", Idp: "cernbox.cern.ch"})
		}
	})
	b.Run("GetUserByClaim", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, _ = manager.GetUserByClaim("mail", "einstein@cern.ch")
		}
	})
	b.Run("GetUserGroups", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, _ = manager.GetUserGroups(&userpb.UserId{OpaqueId: "4c510ada-c86b-4815-8820-42cdf82c3d51", Idp: "cernbox.cern.ch"})
		}
	})
	b.Run("FindUser", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, _ = manager.FindUsers("einstein")
		}
	})
}

// BenchmarkPiePlugin benchmarks pie plugin
func BenchmarkPiePlugin(b *testing.B) {
	// plugin provider path (bin exe)
	path := "./pieplugin"

	// we are client and communicate with the plugin using JSON-RPC.
	client, err := pie.StartProviderCodec(jsonrpc.NewClientCodec, os.Stderr, path)
	if err != nil {
		fmt.Println("Error running plugin")
		panic(err)
	}
	defer client.Close()

	manager := shared.RPCClient{client}

	b.Run("OnLoad", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = manager.OnLoad("./file/user.demo.json")
		}
	})
	b.Run("GetUser", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, _ = manager.GetUser(&userpb.UserId{OpaqueId: "4c510ada-c86b-4815-8820-42cdf82c3d51", Idp: "cernbox.cern.ch"})
		}
	})
	b.Run("GetUserByClaim", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, _ = manager.GetUserByClaim("mail", "einstein@cern.ch")
		}
	})
	b.Run("GetUserGroups", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, _ = manager.GetUserGroups(&userpb.UserId{OpaqueId: "4c510ada-c86b-4815-8820-42cdf82c3d51", Idp: "cernbox.cern.ch"})
		}
	})
	b.Run("FindUser", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, _ = manager.FindUsers("einstein")
		}
	})
}
