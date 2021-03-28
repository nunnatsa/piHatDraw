package common

import (
	"encoding/json"
	"fmt"
)

const (
	// the Sense HAT display is 8X8 matrix
	WindowSize = 8
)

// Color is the Color of one pixel in the Canvas
type Color uint32

func (c Color) MarshalJSON() ([]byte, error) {
	r := (c >> 16) & 0xFF
	g := (c >> 8) & 0xFF
	b := c & 0xFF

	s := fmt.Sprintf(`"#%02x%02x%02x"`, r, g, b)
	return []byte(s), nil
}

func (c *Color) UnmarshalJSON(bt []byte) error {
	var r, g, b uint32
	var s string
	err := json.Unmarshal(bt, &s)
	if err != nil {
		return err
	}

	_, err = fmt.Sscanf(s, `#%02x%02x%02x`, &r, &g, &b)
	if err != nil {
		return err
	}

	*c = Color(r<<16 | g<<8 | b)

	return nil
}
