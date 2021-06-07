package shared

import (
	"net/rpc"

	userpb "github.com/cs3org/go-cs3apis/cs3/identity/user/v1beta1"
)

// RPCClient is an implementation of Manager that talks over RPC.
type RPCClient struct{ client *rpc.Client }

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
	err := m.client.Call("Plugin.OnLoad", args, &resp)
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
	err := m.client.Call("Plugin.GetUser", args, &resp)
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
