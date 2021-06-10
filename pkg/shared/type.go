package shared

type User struct {
	Id           *UserId
	Username     string
	Mail         string
	MailVerified bool
	DisplayName  string
	Groups       []string
	Opaque       *Opaque
	UidNumber    int64
	GidNumber    int64
}

// A UserId represents a user.
type UserId struct {
	Idp      string
	OpaqueId string
	Type     UserType
}

// The type of user.
type UserType int32

const (
	// The user is invalid, for example, is missing primary attributes.
	UserType_USER_TYPE_INVALID UserType = 0
	// A primary user.
	UserType_USER_TYPE_PRIMARY UserType = 1
	// A secondary user for cases with multiple identities.
	UserType_USER_TYPE_SECONDARY UserType = 2
	// A user catering to specific services.
	UserType_USER_TYPE_SERVICE UserType = 3
	// A user to be used by specific applications.
	UserType_USER_TYPE_APPLICATION UserType = 4
	// A guest user not affiliated to the IDP.
	UserType_USER_TYPE_GUEST UserType = 5
	// A federated user provided by external IDPs.
	UserType_USER_TYPE_FEDERATED UserType = 6
	// A lightweight user account without access to various major functionalities.
	UserType_USER_TYPE_LIGHTWEIGHT UserType = 7
)

type Opaque struct {
	// REQUIRED.
	Map map[string]*OpaqueEntry
}

type OpaqueEntry struct {
	Decoder string `protobuf:"bytes,1,opt,name=decoder,proto3" json:"decoder,omitempty"`
	Value   []byte `protobuf:"bytes,2,opt,name=value,proto3" json:"value,omitempty"`
}
