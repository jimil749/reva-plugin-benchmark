package main

import (
	"encoding/json"
	"io/ioutil"

	userpb "github.com/cs3org/go-cs3apis/cs3/identity/user/v1beta1"
	"github.com/hashicorp/go-plugin"
	"github.com/jimil749/reva-plugin-benchmark/pkg/shared"
)

// Here is a real implementation of KV that writes to a local file with
// the key name and the contents are the value of the key.
type Manager struct {
	users []*userpb.User
}

func (m *Manager) OnLoad(userFile string) error {
	f, err := ioutil.ReadFile(userFile)
	if err != nil {
		return err
	}

	users := []*userpb.User{}

	err = json.Unmarshal(f, &users)
	if err != nil {
		return err
	}

	m.users = users

	return nil
}

func (m *Manager) GetUser(uid *userpb.UserId) (*userpb.User, error) {
	for _, u := range m.users {
		if (u.Id.GetOpaqueId() == uid.OpaqueId || u.Username == uid.OpaqueId) && (uid.Idp == "" || uid.Idp == u.Id.GetIdp()) {
			return u, nil
		}
	}
	return nil, nil
}

func main() {
	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: shared.Handshake,
		Plugins: map[string]plugin.Plugin{
			"json": &shared.JSONPlugin{Impl: &Manager{}},
		},
	})
}
