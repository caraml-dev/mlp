package types

// Policy authorization policy
type Policy struct {
	ID        string
	Subjects  []string
	Resources []string
	Actions   []string
}
