package snowflake

// NodeIdStorager store the latest worker node id
type NodeIdStorager interface {
	// NextNodeId return next worker node id
	NextNodeId() (int64, error)
}
