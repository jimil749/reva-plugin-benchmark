package yaegi

import (
	"fmt"
	"reflect"

	"github.com/jimil749/reva-plugin-benchmark/pkg/manager"
	"github.com/traefik/yaegi/interp"
	"github.com/traefik/yaegi/stdlib"
)

func Load() {
	i := interp.New(interp.Options{
		GoPath: "/home/jimil/go",
	})

	i.Use(stdlib.Symbols)
	i.Use(managerSymbols())

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

	userFile := "/home/jimil/Desktop/reva-plugin-benchmark/file/user.demo.json"
	args := []reflect.Value{reflect.ValueOf(userFile)}
	results := fnNew.Call(args)

	if len(results) > 1 && results[1].Interface() != nil {
		panic(results[1].Interface().(error))
	}

	mgr, ok := results[0].Interface().(manager.UserManager)
	if !ok {
		panic(fmt.Errorf("invalid interface type"))
	}

	user, err := mgr.GetUser(&manager.UserId{OpaqueId: "4c510ada-c86b-4815-8820-42cdf82c3d51", Idp: "cernbox.cern.ch"})
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	fmt.Println(user.DisplayName)

}
