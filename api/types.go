package api

import (
	"encoding/json"
	"fmt"
	"strconv"
)

type BoolInt bool

// UnmarshalJSON implements the json.Unmarshaler interface.
func (b *BoolInt) UnmarshalJSON(data []byte) error {
	// First, try to unmarshal as a standard boolean.
	var boolVal bool
	if err := json.Unmarshal(data, &boolVal); err == nil {
		*b = BoolInt(boolVal)
		return nil
	}

	// If that fails, try to unmarshal as a number.
	var intVal int
	if err := json.Unmarshal(data, &intVal); err == nil {
		if intVal == 1 {
			*b = true
		} else {
			*b = false
		}
		return nil
	}

	// As a fallback, handle quoted numbers like "1" or "0"
	var stringVal string
	if err := json.Unmarshal(data, &stringVal); err == nil {
		val, err := strconv.Atoi(stringVal)
		if err == nil {
			*b = val == 1
			return nil
		}
	}

	return fmt.Errorf("cannot unmarshal %s into a boolean", data)
}
