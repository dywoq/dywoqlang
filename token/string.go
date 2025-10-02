package token

import "encoding/json"

// ToString converts t into the Json string presentation.
// Returns <nil> instead of the string if t is nil.
func ToString(t *Token) string {
	if t == nil {
		return "<nil>"
	}
	bytes, err := json.MarshalIndent(t, "", "  ")
	if err != nil {
		return err.Error()
	}
	return string(bytes)
}
