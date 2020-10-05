package identity

//Domain model represents domains which users can be part of
type Domain struct {
	ID   *int64  `db:"id" dbignoreinsert:"" json:"id"`
	Name *string `db:"name" json:"name"`
}

// Domains is a map where the key represents the domains's id
type Domains map[int64]struct {
	Role *string `json:"role,omitempty"`
}
