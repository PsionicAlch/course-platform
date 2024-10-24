package gatekeeper

type GatekeeperPassword struct {
	ArgonVersion int
	Hash         []byte
	Salt         []byte
	Iterations   uint8
	Memory       uint32
	Threads      uint8
	KeyLength    uint8
}
