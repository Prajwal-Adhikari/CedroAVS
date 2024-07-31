package protocols

type IProtocol interface {
	PrepareBytes() []byte

	GetBytes() []byte
}


