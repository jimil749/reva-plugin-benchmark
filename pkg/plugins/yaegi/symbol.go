package yaegi

import (
	"reflect"

	"github.com/jimil749/reva-plugin-benchmark/pkg/manager"
)

// _github_com_jimil749_reva_plugin_benchmark_pkg_manager_UserManager is an interface wrapper for UserManager type
type _github_com_jimil749_reva_plugin_benchmark_pkg_manager_UserManager struct {
	WFindUsers      func(query string) ([]*manager.User, error)
	WGetUser        func(a0 *manager.UserId) (*manager.User, error)
	WGetUserByClaim func(claim string, value string) (*manager.User, error)
	WGetUserGroups  func(a0 *manager.UserId) ([]string, error)
}

func (W _github_com_jimil749_reva_plugin_benchmark_pkg_manager_UserManager) FindUsers(query string) ([]*manager.User, error) {
	return W.WFindUsers(query)
}
func (W _github_com_jimil749_reva_plugin_benchmark_pkg_manager_UserManager) GetUser(a0 *manager.UserId) (*manager.User, error) {
	return W.WGetUser(a0)
}
func (W _github_com_jimil749_reva_plugin_benchmark_pkg_manager_UserManager) GetUserByClaim(claim string, value string) (*manager.User, error) {
	return W.WGetUserByClaim(claim, value)
}
func (W _github_com_jimil749_reva_plugin_benchmark_pkg_manager_UserManager) GetUserGroups(a0 *manager.UserId) ([]string, error) {
	return W.WGetUserGroups(a0)
}

func ManagerSymbols() map[string]map[string]reflect.Value {
	return map[string]map[string]reflect.Value{
		"github.com/jimil749/reva-plugin-benchmark/pkg/manager/manager": {
			// type definitions
			"Opaque":      reflect.ValueOf((*manager.Opaque)(nil)),
			"OpaqueEntry": reflect.ValueOf((*manager.OpaqueEntry)(nil)),
			"User":        reflect.ValueOf((*manager.User)(nil)),
			"UserId":      reflect.ValueOf((*manager.UserId)(nil)),
			"UserManager": reflect.ValueOf((*manager.UserManager)(nil)),

			// interface wrapper definitions
			"_UserManager": reflect.ValueOf((*_github_com_jimil749_reva_plugin_benchmark_pkg_manager_UserManager)(nil)),
		},
	}
}
