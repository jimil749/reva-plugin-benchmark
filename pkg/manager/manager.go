package manager

// UserManager is interface we are exposing to Yaegi interpreter. (This interface uses custom types rather than userpb types)
type UserManager interface {
	GetUser(*UserId) (*User, error)
	GetUserByClaim(claim, value string) (*User, error)
	GetUserGroups(*UserId) ([]string, error)
	FindUsers(query string) ([]*User, error)
}

// These are "custom" types that I am using for the sake of benchmarking. They should be replaced
// by protobuf implementation but Yaegi panics on using protobuf, hence this workaround. (At least for the time being.)

type UserId struct {
	Idp      string
	OpaqueId string
}

func (m *UserId) GetIdp() string {
	if m != nil {
		return m.Idp
	}
	return ""
}

func (m *UserId) GetOpaqueId() string {
	if m != nil {
		return m.OpaqueId
	}
	return ""
}

type User struct {
	Id           *UserId
	Username     string
	Mail         string
	MainVerified string
	DisplayName  string
	Groups       []string
	Opaque       *Opaque
	UidNumber    int64
	GidNumber    int64
}

type Opaque struct {
	Map map[string]*OpaqueEntry
}

type OpaqueEntry struct {
	Decoder string
	Value   []byte
}
