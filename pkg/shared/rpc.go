package shared

import (
	"net/rpc"

	userpb "github.com/cs3org/go-cs3apis/cs3/identity/user/v1beta1"
)

// RPCClient is an implementation of Manager that talks over RPC.
type RPCClient struct{ Client *rpc.Client }

// OnLoadArg for RPC
type OnLoadArg struct {
	UserFile string
}

// OnLoadReply for RPC
type OnLoadReply struct {
	Err error
}

func (m *RPCClient) OnLoad(userFile string) error {
	args := OnLoadArg{UserFile: userFile}
	resp := OnLoadReply{}
	err := m.Client.Call("Plugin.OnLoad", args, &resp)
	if err != nil {
		return err
	}
	return resp.Err
}

type GetUserArg struct {
	Uid *userpb.UserId
}

type GetUserReply struct {
	User *userpb.User
	Err  error
}

func (m *RPCClient) GetUser(uid *userpb.UserId) (*userpb.User, error) {
	args := GetUserArg{Uid: uid}
	resp := GetUserReply{}
	err := m.Client.Call("Plugin.GetUser", args, &resp)
	if err != nil {
		return nil, err
	}
	return resp.User, resp.Err
}

type GetUserByClaimArg struct {
	Claim string
	Value string
}

type GetUserByClaimReply struct {
	User *userpb.User
	Err  error
}

func (m *RPCClient) GetUserByClaim(claim, value string) (*userpb.User, error) {
	args := GetUserByClaimArg{Claim: claim, Value: value}
	resp := GetUserByClaimReply{}
	err := m.Client.Call("Plugin.GetUserByClaim", args, &resp)
	if err != nil {
		return nil, err
	}
	return resp.User, resp.Err
}

type GetUserGroupsArg struct {
	User *userpb.UserId
}

type GetUserGroupsReply struct {
	Group []string
	Err   error
}

func (m *RPCClient) GetUserGroups(user *userpb.UserId) ([]string, error) {
	args := GetUserGroupsArg{User: user}
	resp := GetUserGroupsReply{}
	err := m.Client.Call("Plugin.GetUserGroups", args, &resp)
	if err != nil {
		return nil, err
	}
	return resp.Group, resp.Err
}

type FindUsersArg struct {
	Query string
}

type FindUsersReply struct {
	User []*userpb.User
	Err  error
}

func (m *RPCClient) FindUsers(query string) ([]*userpb.User, error) {
	args := FindUsersArg{Query: query}
	resp := FindUsersReply{}
	err := m.Client.Call("Plugin.FindUsers", args, &resp)
	if err != nil {
		return nil, err
	}
	return resp.User, resp.Err
}

// Here is the RPC server that RPCClient talks to, conforming to
// the requirements of net/rpc
type RPCServer struct {
	// This is the real implementation
	Impl Manager
}

func (m *RPCServer) OnLoad(args OnLoadArg, resp *OnLoadReply) error {
	resp.Err = m.Impl.OnLoad(args.UserFile)
	return nil
}

func (m *RPCServer) GetUser(args GetUserArg, resp *GetUserReply) error {
	resp.User, resp.Err = m.Impl.GetUser(args.Uid)
	return nil
}

func (m *RPCServer) GetUserByClaim(args GetUserByClaimArg, resp *GetUserByClaimReply) error {
	resp.User, resp.Err = m.Impl.GetUserByClaim(args.Claim, args.Value)
	return nil
}

func (m *RPCServer) GetUserGroups(args GetUserGroupsArg, resp *GetUserGroupsReply) error {
	resp.Group, resp.Err = m.Impl.GetUserGroups(args.User)
	return nil
}

func (m *RPCServer) FindUsers(args FindUsersArg, resp *FindUsersReply) error {
	resp.User, resp.Err = m.Impl.FindUsers(args.Query)
	return nil
}
