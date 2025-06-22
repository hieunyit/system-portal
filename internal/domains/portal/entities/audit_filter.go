package entities

import "time"

// AuditFilter defines filtering and pagination for audit logs
//
// Limit and Page determine pagination; Offset is calculated from them.
// Default values are Page=1 and Limit=20 if not provided.
// FromTime and ToTime can be nil to omit the bound.
//
// Username filters by username, UserGroup by group name, IPAddress by IP.
type AuditFilter struct {
	Username  string
	UserGroup string
	IPAddress string
	FromTime  *time.Time
	ToTime    *time.Time
	Page      int
	Limit     int
	Offset    int
}

// SetDefaults ensures pagination defaults and calculates offset.
func (f *AuditFilter) SetDefaults() {
	if f.Page <= 0 {
		f.Page = 1
	}
	if f.Limit <= 0 {
		f.Limit = 20
	}
	f.Offset = (f.Page - 1) * f.Limit
}
