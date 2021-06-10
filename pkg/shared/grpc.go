package shared

import (
	"context"

	userpb "github.com/cs3org/go-cs3apis/cs3/identity/user/v1beta1"
	proto "github.com/jimil749/reva-plugin-benchmark/pkg/proto"
)

type GRPCClient struct{ client proto.UserAPIClient }

func (m *GRPCClient) OnLoad(userFile string) error {
	_, err := m.client.OnLoad(context.Background(), &proto.OnLoadRequest{
		UserFile: userFile,
	})
	return err
}

func (m *GRPCClient) GetUser(uid *userpb.UserId) (*userpb.User, error) {
	resp, err := m.client.GetUser(context.Background(), &proto.GetUserRequest{
		UserId: &proto.UserId{
			Idp:      uid.Idp,
			OpaqueId: uid.OpaqueId,
			Type:     proto.UserType(uid.Type),
		},
	})
	if err != nil {
		return nil, err
	}
	// unmarshal the manager user to userpb.User
	user := UnmarshalUser(resp.User)
	return user, nil
}

func (m *GRPCClient) GetUserByClaim(claim, value string) (*userpb.User, error) {
	resp, err := m.client.GetUserByClaim(context.Background(), &proto.GetUserByClaimRequest{
		Value: value,
		Claim: claim,
	})
	if err != nil {
		return nil, err
	}

	user := UnmarshalUser(resp.User)
	return user, nil
}

func (m *GRPCClient) GetUserGroups(uid *userpb.UserId) ([]string, error) {
	resp, err := m.client.GetUserGroups(context.Background(), &proto.GetUserGroupsRequest{
		UserId: &proto.UserId{
			Idp:      uid.Idp,
			OpaqueId: uid.OpaqueId,
			Type:     proto.UserType(uid.Type),
		},
	})
	if err != nil {
		return nil, err
	}

	return resp.Groups, nil
}

func (m *GRPCClient) FindUsers(filter string) ([]*userpb.User, error) {
	resp, err := m.client.FindUsers(context.Background(), &proto.FindUsersRequest{
		Filter: filter,
	})
	if err != nil {
		return nil, err
	}

	var users []*userpb.User
	for _, user := range resp.Users {
		users = append(users, UnmarshalUser(user))
	}
	return users, nil
}

// gRPC server that gRPC client talks to.
type GRPCServer struct {
	Impl Manager
	proto.UnimplementedUserAPIServer
}

func (m *GRPCServer) OnLoad(ctx context.Context, req *proto.OnLoadRequest) (*proto.OnLoadResponse, error) {
	return &proto.OnLoadResponse{}, m.Impl.OnLoad(req.UserFile)
}

func (m *GRPCServer) GetUser(ctx context.Context, req *proto.GetUserRequest) (*proto.GetUserResponse, error) {
	userpb, err := m.Impl.GetUser(UnmarshalUserId(req.UserId))
	user := MarshalUser(userpb)
	return &proto.GetUserResponse{User: user}, err
}

func (m *GRPCServer) GetUserByClaim(ctx context.Context, req *proto.GetUserByClaimRequest) (*proto.GetUserByClaimResponse, error) {
	userpb, err := m.Impl.GetUserByClaim(req.Claim, req.Value)
	user := MarshalUser(userpb)
	return &proto.GetUserByClaimResponse{User: user}, err
}

func (m *GRPCServer) GetUserGroups(ctx context.Context, req *proto.GetUserGroupsRequest) (*proto.GetUserGroupsResponse, error) {
	groups, err := m.Impl.GetUserGroups(UnmarshalUserId(req.UserId))
	return &proto.GetUserGroupsResponse{Groups: groups}, err
}

func (m *GRPCServer) FindUsers(ctx context.Context, req *proto.FindUsersRequest) (*proto.FindUsersResponse, error) {
	userspb, err := m.Impl.FindUsers(req.Filter)
	var users []*proto.User
	for _, user := range userspb {
		users = append(users, MarshalUser(user))
	}
	return &proto.FindUsersResponse{Users: users}, err
}

// UnmarshalUserId "unmarshals" the proto.UserId to *userpb.UserId
func UnmarshalUserId(uid *proto.UserId) *userpb.UserId {
	return &userpb.UserId{
		Idp:      uid.Idp,
		OpaqueId: uid.OpaqueId,
		Type:     userpb.UserType(uid.Type),
	}
}

// UnmarshalUser "unmarshals" the proto.User type to *userpb.User
func UnmarshalUser(user *proto.User) *userpb.User {
	return &userpb.User{
		Id: &userpb.UserId{
			Idp:      user.Id.Idp,
			OpaqueId: user.Id.OpaqueId,
			Type:     userpb.UserType(user.Id.Type),
		},
		Username:     user.Username,
		Mail:         user.Mail,
		MailVerified: user.MailVerified,
		DisplayName:  user.DisplayName,
		Groups:       user.Groups,
		UidNumber:    user.UidNumber,
		GidNumber:    user.GidNumber,
	}
}

// MarshalUser marshals userpb.User to proto.User
func MarshalUser(user *userpb.User) *proto.User {
	return &proto.User{
		Id: &proto.UserId{
			Idp:      user.Id.Idp,
			OpaqueId: user.Id.OpaqueId,
			Type:     proto.UserType(user.Id.Type),
		},
		Username:     user.Username,
		Mail:         user.Mail,
		MailVerified: user.MailVerified,
		DisplayName:  user.DisplayName,
		Groups:       user.Groups,
		UidNumber:    user.UidNumber,
		GidNumber:    user.GidNumber,
	}
}
