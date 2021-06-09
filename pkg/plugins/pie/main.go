package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/rpc/jsonrpc"
	"strings"

	userpb "github.com/cs3org/go-cs3apis/cs3/identity/user/v1beta1"
	"github.com/cs3org/reva/pkg/errtypes"
	"github.com/jimil749/reva-plugin-benchmark/pkg/shared"
	"github.com/natefinch/pie"
)

// Here is a real implementation of Manager
type Manager struct {
	users []*userpb.User
}

func (m *Manager) OnLoad(args shared.OnLoadArg, reply *shared.OnLoadReply) error {
	f, err := ioutil.ReadFile(args.UserFile)
	if err != nil {
		return err
	}

	users := []*userpb.User{}

	err = json.Unmarshal(f, &users)
	if err != nil {
		reply.Err = err
		return err
	}

	m.users = users

	return nil
}

func (m *Manager) GetUser(args shared.GetUserArg, reply *shared.GetUserReply) error {
	for _, u := range m.users {
		if (u.Id.GetOpaqueId() == args.Uid.OpaqueId || u.Username == args.Uid.OpaqueId) && (args.Uid.Idp == "" || args.Uid.Idp == u.Id.GetIdp()) {
			reply.User = u
			reply.Err = nil
			return nil
		}
	}
	return nil
}

func (m *Manager) GetUserByClaim(args shared.GetUserByClaimArg, reply *shared.GetUserByClaimReply) error {
	for _, u := range m.users {
		if userClaim, err := extractClaim(u, args.Claim); err == nil && args.Value == userClaim {
			reply.User = u
			reply.Err = nil
			return nil
		}
	}
	reply.Err = errtypes.NotFound(args.Value)
	reply.User = nil
	return nil
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
	return "", errors.New("json: invalid field")
}

// TODO(jfd) search Opaque? compare sub?
func userContains(u *userpb.User, query string) bool {
	query = strings.ToLower(query)
	return strings.Contains(strings.ToLower(u.Username), query) || strings.Contains(strings.ToLower(u.DisplayName), query) ||
		strings.Contains(strings.ToLower(u.Mail), query) || strings.Contains(strings.ToLower(u.Id.OpaqueId), query)
}

func (m *Manager) FindUsers(args shared.FindUsersArg, reply *shared.FindUsersReply) error {
	users := []*userpb.User{}
	for _, u := range m.users {
		if userContains(u, args.Query) {
			users = append(users, u)
		}
	}
	reply.User = users
	return nil
}

// func (m Manager) GetUserGroups(args shared.GetUserGroupsArg, reply *shared.GetUserGroupsReply) (error) {
// 	user, err := m.GetUser(args.User)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return user.Groups, nil
// }

func main() {
	// We are plugin provider!
	p := pie.NewProvider()
	if err := p.RegisterName("Plugin", &Manager{}); err != nil {
		panic(err)
	}
	p.ServeCodec(jsonrpc.NewServerCodec)
}
