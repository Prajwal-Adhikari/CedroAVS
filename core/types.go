package core

import (
	"encoding/hex"
	"log/slog"
)

type Bytes32 [32]byte

func (m Bytes32) LogValue() slog.Value {
	return slog.StringValue(m.String())
}

func (m Bytes32) String() string {
	return hex.EncodeToString(m[:])
}

func (m *Bytes32) UnderlyingType() [32]byte {
	return *m
}