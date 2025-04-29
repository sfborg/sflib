package coldp

// Data can be used when we conversion from one format to another
// does correspond to several structs in coldp. SQLite is not very
// good and concurrent writes, and this structure allows us to keep data
// writes sequential.
type Data struct {
	NameUsages    []NameUsage
	Distributions []Distribution
	References    []Reference
}
