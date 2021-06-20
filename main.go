package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"

	userpb "github.com/cs3org/go-cs3apis/cs3/identity/user/v1beta1"
	"github.com/hashicorp/go-hclog"
	hashPlugin "github.com/hashicorp/go-plugin"
	"github.com/jimil749/reva-plugin-benchmark/pkg/shared"
)

// this is just for the purpose of testing
func main() {
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
	manager := raw.(shared.UserManager)

	user, _ := manager.GetUser(&userpb.UserId{OpaqueId: "4c510ada-c86b-4815-8820-42cdf82c3d51", Idp: "cernbox.cern.ch"})
	fmt.Printf(user.DisplayName)
}
