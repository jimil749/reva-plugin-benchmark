package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/rpc/jsonrpc"
	"os"
	"os/exec"
	"plugin"
	"testing"

	userpb "github.com/cs3org/go-cs3apis/cs3/identity/user/v1beta1"
	"github.com/hashicorp/go-hclog"
	hashPlugin "github.com/hashicorp/go-plugin"
	"github.com/jimil749/reva-plugin-benchmark/pkg/shared"
	"github.com/natefinch/pie"
)

// BenchmarkGoPlugin benchmarks the native go-plugin
func BenchmarkGoPlugin(b *testing.B) {
	// Open the plugin shared object library
	sym, err := plugin.Open("./go-plugin.so")
	if err != nil {
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
