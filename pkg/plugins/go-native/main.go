package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/jimil749/reva-plugin-benchmark/pkg/shared"

	userpb "github.com/cs3org/go-cs3apis/cs3/identity/user/v1beta1"
	"github.com/cs3org/reva/pkg/errtypes"
)

type Manager struct {
	users []*userpb.User
}

// New returns a user manager implementation that reads a json file to provide user metadata.
func New(userFile string) (shared.UserManager, error) {

	f, err := ioutil.ReadFile(userFile)
	if err != nil {
		return nil, err
	}

	users := []*userpb.User{}

	err = json.Unmarshal(f, &users)
	if err != nil {
		return nil, err
	}

	return &Manager{
		users: users,
	}, nil
}

func (m *Manager) GetUser(uid *userpb.UserId) (*userpb.User, error) {
	for _, u := range m.users {
		if (u.Id.GetOpaqueId() == uid.OpaqueId || u.Username == uid.OpaqueId) && (uid.Idp == "" || uid.Idp == u.Id.GetIdp()) {
			return u, nil
		}
	}
	return nil, errtypes.NotFound(uid.OpaqueId)
}

func (m *Manager) GetUserByClaim(claim, value string) (*userpb.User, error) {
	for _, u := range m.users {
		if userClaim, err := extractClaim(u, claim); err == nil && value == userClaim {
			return u, nil
		}
	}
	return nil, errtypes.NotFound(value)
}

func extractClaim(u *userpb.User, claim string) (string, error) {
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
func userContains(u *userpb.User, query string) bool {
	query = strings.ToLower(query)
	return strings.Contains(strings.ToLower(u.Username), query) || strings.Contains(strings.ToLower(u.DisplayName), query) ||
		strings.Contains(strings.ToLower(u.Mail), query) || strings.Contains(strings.ToLower(u.Id.OpaqueId), query)
}

func (m *Manager) FindUsers(query string) ([]*userpb.User, error) {
	users := []*userpb.User{}
	for _, u := range m.users {
		if userContains(u, query) {
			users = append(users, u)
		}
	}
	return users, nil
}

func (m *Manager) GetUserGroups(uid *userpb.UserId) ([]string, error) {
	user, err := m.GetUser(uid)
	if err != nil {
		return nil, err
	}
	return user.Groups, nil
}

func main() {}
