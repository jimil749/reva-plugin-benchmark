package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"

	userpb "github.com/cs3org/go-cs3apis/cs3/identity/user/v1beta1"
	"github.com/hashicorp/go-plugin"
	"github.com/jimil749/reva-plugin-benchmark/pkg/shared"
)

// this is just for the purpose of testing
func main() {
	// We don't want to see the plugin logs.
	log.SetOutput(ioutil.Discard)

	// We're a host. Start by launching the plugin process.
	client := plugin.NewClient(&plugin.ClientConfig{
		HandshakeConfig: shared.Handshake,
		Plugins:         shared.PluginMap,
		Cmd:             exec.Command("./hashicorp-plugin"),
		AllowedProtocols: []plugin.Protocol{
			plugin.ProtocolNetRPC, plugin.ProtocolGRPC},
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

	// We should have a KV store now! This feels like a normal interface
	// implementation but is in fact over an RPC connection.
	manager := raw.(shared.Manager)

	err = manager.OnLoad("/home/jimil/Desktop/reva-plugin-benchmark/file/user.demo.json")
	if err != nil {
		fmt.Println(err)
	}

	user, err := manager.GetUser(&userpb.UserId{OpaqueId: "4c510ada-c86b-4815-8820-42cdf82c3d51", Idp: "cernbox.cern.ch"})
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(user.DisplayName)

}
