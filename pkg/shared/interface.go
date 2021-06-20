package shared

import (
	"context"
	"net/rpc"

	userpb "github.com/cs3org/go-cs3apis/cs3/identity/user/v1beta1"
	"github.com/hashicorp/go-plugin"
	proto "github.com/jimil749/reva-plugin-benchmark/pkg/proto"
	"google.golang.org/grpc"
)

// Handshake is a common handshake that is shared by plugin and host.
var Handshake = plugin.HandshakeConfig{
	// This isn't required when using VersionedPlugins
	ProtocolVersion:  1,
	MagicCookieKey:   "BASIC_PLUGIN",
	MagicCookieValue: "hello",
}

// PluginMap is the map of plugins we can dispense.
var PluginMap = map[string]plugin.Plugin{
	"json":      &JSONPlugin{},
	"json_grpc": &JSONPluginGRPC{},
}

// Manager is the interface that we're exposing as a plugin. This interface is only used for plugin
// systems using RPC for communication.
type Manager interface {
	// instead of this, use some interface/method that returns an initialized rpc server
	OnLoad(userFile string) error
	GetUser(*userpb.UserId) (*userpb.User, error)
	GetUserByClaim(claim, value string) (*userpb.User, error)
	GetUserGroups(*userpb.UserId) ([]string, error)
	FindUsers(query string) ([]*userpb.User, error)
}

// This is the implementation of plugin.Plugin so we can serve/consume this.
type JSONPlugin struct {
	// Concrete implementation, written in Go. This is only used for plugins
	// that are written in Go.
	Impl UserManager
}

func (p *JSONPlugin) Server(*plugin.MuxBroker) (interface{}, error) {
	return &RPCServer{Impl: p.Impl}, nil
}

func (*JSONPlugin) Client(b *plugin.MuxBroker, c *rpc.Client) (interface{}, error) {
	return &RPCClient{Client: c}, nil
}

// ------------------------------ For Non-RPC Plugins--------------------------------------
// UserManager is the interface we're exposing as a plugin for plugin systems NOT using RPC.
type UserManager interface {
	GetUser(*userpb.UserId) (*userpb.User, error)
	GetUserByClaim(claim, value string) (*userpb.User, error)
	GetUserGroups(*userpb.UserId) ([]string, error)
	FindUsers(query string) ([]*userpb.User, error)
}

type JSONPluginGRPC struct {
	plugin.Plugin
	Impl Manager
}

func (p *JSONPluginGRPC) GRPCServer(broker *plugin.GRPCBroker, s *grpc.Server) error {
	proto.RegisterUserAPIServer(s, &GRPCServer{Impl: p.Impl})
	return nil
}

func (p *JSONPluginGRPC) GRPCClient(ctx context.Context, broker *plugin.GRPCBroker, c *grpc.ClientConn) (interface{}, error) {
	return &GRPCClient{client: proto.NewUserAPIClient(c)}, nil
}
