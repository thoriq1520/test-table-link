package domain

type RoleRights struct {
	ID      int64
	RoleID  int64
	Route   string
	Section string
	Path    string
	RCreate bool
	RRead   bool
	RUpdate bool
	RDelete bool
}
