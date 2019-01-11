package models

// Group ...
type Group struct {
	ID       string `json:"id,omitempty"`
	ParentID string `json:"parent_id,omitempty"`
	Name     string `json:"name,omitempty"`
	Type     string `json:"type,omitempty"`
	Logo     string `json:"logo,omitempty"`
}

// GroupLineage represents a set of groups connected by consecutive parent-child relationships, where idx 0 is the org (i.e. has no parent)
type GroupLineage []*Group

// Leaf returns the last group in the lineage
func (gl GroupLineage) Leaf() *Group {
	if len(gl) == 0 {
		return nil
	}

	return gl[len(gl)-1]
}
