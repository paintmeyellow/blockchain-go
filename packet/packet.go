package packet

type Transport struct {
	KeeperID string `json:"keeper_id"`
	Handler  string `json:"handler"`
}

type Connection struct {
	KeeperID string `json:"keeper_id"`
}