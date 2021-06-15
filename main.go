package main

import (
	"encoding/json"
	"fmt"
	"os"
	"unsafe"

	"github.com/jimil749/reva-plugin-benchmark/pkg/manager"
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

	runFuncPtr := codeModule.Syms["main.New"]
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
}
