package models

// StorageDriver ...
type StorageDriver struct {
	ID     string   `json:"id"`
	Name   string   `json:"name"`
	Type   string   `json:"type"`
	Mounts []*Mount `json:"mounts"`
}

// Mount ...
type Mount struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Type     string `json:"type"`
	Location string `json:"mountLocation"`
}
