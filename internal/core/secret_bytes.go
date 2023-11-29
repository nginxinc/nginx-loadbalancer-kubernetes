package core

import (
	"encoding/json"
)

type SecretBytes []byte

func (sb SecretBytes) String() string {
	return "[REDACTED]"
}

func (sb SecretBytes) MarshalJSON() ([]byte, error) {
	return json.Marshal("[REDACTED]")
}
