package main

import (
	"encoding/json"
	"fmt"
	"os"
	"unsafe"

	"github.com/jimil749/reva-plugin-benchmark/pkg/plugins/goloader/manager"
	"github.com/pkujhd/goloader"
)

// this is just for the purpose of testing
func main() {

	files := []string{
		"/home/jimil/go/pkg/linux_amd64/github.com/jimil749/reva-plugin-benchmark/pkg/plugins/goloader/manager.a",
		"json.o",
	}
	pkgPath := []string{
		"github.com/jimil749/reva-plugin-benchmark/pkg/plugins/goloader/manager",
		"",
	}

	symPtr := make(map[string]uintptr)
	err := goloader.RegSymbol(symPtr)
	if err != nil {
		fmt.Println("fails in regsymbol")
		panic(err)
	}

	goloader.RegTypes(symPtr, os.ReadFile)
	goloader.RegTypes(symPtr, json.Unmarshal)

	linker, err := goloader.ReadObjs(files, pkgPath)
	if err != nil {
		fmt.Println("fails in readobjs")
		panic(err)
	}

	var mmapByte []byte
	codeModule, err := goloader.Load(linker, symPtr)
	fmt.Println(err)
	if err != nil {
		fmt.Println("fails in loading")
		panic(err)
	}

	runFuncPtr := codeModule.Syms["json.New"]
	if runFuncPtr == 0 {
		panic("Load error! Function not found: json.New")
	}
	funcPtrContainer := (uintptr)(unsafe.Pointer(&runFuncPtr))
	runFunc := *(*func(string) (manager.UserManager, error))(unsafe.Pointer(&funcPtrContainer))

	manager, err := runFunc("./file/user.demo.json")
	if err != nil {
		panic(err)
	}

	fmt.Printf("%+v", manager)
	codeModule.Unload()

	if mmapByte == nil {
		mmapByte, err = goloader.Mmap(1024)
		if err != nil {
			panic(err)
		}
		b := make([]byte, 1024)
		copy(mmapByte, b)
	} else {
		goloader.Munmap(mmapByte)
		mmapByte = nil
	}

	// We don't want to see the plugin logs.
	// log.SetOutput(ioutil.Discard)

	// // We're a host. Start by launching the plugin process.
	// client := plugin.NewClient(&plugin.ClientConfig{
	// 	HandshakeConfig: shared.Handshake,
	// 	Plugins:         shared.PluginMap,
	// 	Cmd:             exec.Command("./hashicorp-plugin-grpc"),
	// 	AllowedProtocols: []plugin.Protocol{
	// 		plugin.ProtocolNetRPC, plugin.ProtocolGRPC},
	// })
	// defer client.Kill()

	// // Connect via RPC
	// rpcClient, err := client.Client()
	// if err != nil {
	// 	fmt.Println("Error:", err.Error())
	// 	os.Exit(1)
	// }

	// // Request the plugin
	// raw, err := rpcClient.Dispense("json_grpc")
	// if err != nil {
	// 	fmt.Println("Error:", err.Error())
	// 	os.Exit(1)
	// }

	// // We should have the Manager now! This feels like a normal interface
	// // implementation but is in fact over an RPC connection.
	// manager := raw.(shared.Manager)

	// err = manager.OnLoad("./file/user.demo.json")
	// if err != nil {
	// 	fmt.Println(err)
	// }

	// user, err := manager.GetUser(&userpb.UserId{OpaqueId: "4c510ada-c86b-4815-8820-42cdf82c3d51", Idp: "cernbox.cern.ch"})
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// user, _ = manager.GetUserByClaim("mail", "einstein@cern.ch")
	// fmt.Println(user.DisplayName)

}
