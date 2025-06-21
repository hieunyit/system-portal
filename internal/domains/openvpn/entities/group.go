package entities

type VpnGroup struct {
	GroupName     string   `json:"groupName"`
	AuthMethod    string   `json:"authMethod"`
	MFA           string   `json:"mfa"`
	Role          string   `json:"role"`
	DenyAccess    string   `json:"denyAccess"`
	AccessControl []string `json:"accessControl"`
	GroupSubnet   []string `json:"groupSubnet"`
	GroupRange    []string `json:"groupRange"`
}

type VpnGroupFilter struct {
	GroupName  string
	AuthMethod string
	IsEnabled  *bool
	Role       string
	Limit      int
	Offset     int
	Page       int
}

// Methods
func (g *VpnGroup) IsAccessDenied() bool {
	return g.DenyAccess == "true"
}

func (g *VpnGroup) SetDenyAccess(deny bool) {
	if deny {
		g.DenyAccess = "true"
	} else {
		g.DenyAccess = "false"
	}
}

func (g *VpnGroup) SetMFA(enabled bool) {
	if enabled {
		g.MFA = "true"
	} else {
		g.MFA = "false"
	}
}

func (g *VpnGroup) HasAccessControl() bool {
	return len(g.AccessControl) > 0
}

func (g *VpnGroup) HasGroupSubnet() bool {
	return len(g.GroupSubnet) > 0
}

func (g *VpnGroup) HasGroupRange() bool {
	return len(g.GroupRange) > 0
}

func NewVpnGroup(groupName, authMethod string) *VpnGroup {
	return &VpnGroup{
		GroupName:   groupName,
		AuthMethod:  authMethod,
		Role:        UserRoleUser,
		DenyAccess:  "false",
		MFA:         "true",
		GroupSubnet: []string{},
		GroupRange:  []string{},
	}
}

// Backward compatibility aliases
type Group = VpnGroup
type GroupFilter = VpnGroupFilter

func NewGroup(groupName, authMethod string) *Group { return NewVpnGroup(groupName, authMethod) }
