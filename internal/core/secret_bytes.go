package core

import (
	"encoding/json"
)

// Wraps byte slices which potentially could contain
// sensitive data that should not be output to the logs.
// This will output [REDACTED] if attempts are made
// to print this type in logs, serialize to JSON, or
// otherwise convert it to a string.
// Usage: core.SecretBytes(myByteSlice)
type SecretBytes []byte

func (sb SecretBytes) String() string {
	return "[REDACTED]"
}

func (sb SecretBytes) MarshalJSON() ([]byte, error) {
	return json.Marshal("[REDACTED]")
}
