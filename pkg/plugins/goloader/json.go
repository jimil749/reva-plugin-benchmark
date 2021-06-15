package json

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/jimil749/reva-plugin-benchmark/pkg/manager"
)

type Manager struct {
	users []*manager.User
}

// New returns a user manager implementation that reads a json file to provide user metadata.
func New(userFile string) (manager.UserManager, error) {

	f, err := ioutil.ReadFile(userFile)
	if err != nil {
		return nil, err
	}

	users := []*manager.User{}

	err = json.Unmarshal(f, &users)
	if err != nil {
		return nil, err
	}

	return &Manager{
		users: users,
	}, nil
}

func (m *Manager) GetUser(uid *manager.UserId) (*manager.User, error) {
	for _, u := range m.users {
		if (u.Id.GetOpaqueId() == uid.OpaqueId || u.Username == uid.OpaqueId) && (uid.Idp == "" || uid.Idp == u.Id.GetIdp()) {
			return u, nil
		}
	}
	return nil, fmt.Errorf("not found")
}

func (m *Manager) GetUserByClaim(claim, value string) (*manager.User, error) {
	for _, u := range m.users {
		if userClaim, err := extractClaim(u, claim); err == nil && value == userClaim {
			return u, nil
		}
	}
	return nil, fmt.Errorf("not found")
}

func extractClaim(u *manager.User, claim string) (string, error) {
	switch claim {
	case "mail":
		return u.Mail, nil
	case "username":
		return u.Username, nil
	case "uid":
		if u.Opaque != nil && u.Opaque.Map != nil {
			if uidObj, ok := u.Opaque.Map["uid"]; ok {
				if uidObj.Decoder == "plain" {
					return string(uidObj.Value), nil
				}
			}
		}
	}
	return "", fmt.Errorf("json: invalid field")
}

// TODO(jfd) search Opaque? compare sub?
func userContains(u *manager.User, query string) bool {
	query = strings.ToLower(query)
	return strings.Contains(strings.ToLower(u.Username), query) || strings.Contains(strings.ToLower(u.DisplayName), query) ||
		strings.Contains(strings.ToLower(u.Mail), query) || strings.Contains(strings.ToLower(u.Id.OpaqueId), query)
}

func (m *Manager) FindUsers(query string) ([]*manager.User, error) {
	users := []*manager.User{}
	for _, u := range m.users {
		if userContains(u, query) {
			users = append(users, u)
		}
	}
	return users, nil
}

func (m *Manager) GetUserGroups(uid *manager.UserId) ([]string, error) {
	user, err := m.GetUser(uid)
	if err != nil {
		return nil, err
	}
	return user.Groups, nil
}
