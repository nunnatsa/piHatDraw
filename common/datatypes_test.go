package common

import (
	"bytes"
	"encoding/json"
	"testing"
)

func TestColor_MarshalJSON(t *testing.T) {
	c := Color(0)
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(c); err != nil {
		t.Fatal(err)
	}

	if res := buf.String(); "\"#000000\"\n" != res {
		t.Fatalf(`should be "#000000" but it's %#v'`, res)
	}

	buf.Reset()
	c = Color(0xFFFFFF)
	if err := json.NewEncoder(&buf).Encode(c); err != nil {
		t.Fatal(err)
	}

	if res := buf.String(); "\"#ffffff\"\n" != res {
		t.Fatalf(`should be "#ffffff" but it's %#v'`, res)
	}

	buf.Reset()
	c = Color(0xFFFFFFFF)
	if err := json.NewEncoder(&buf).Encode(c); err != nil {
		t.Fatal(err)
	}

	if res := buf.String(); "\"#ffffff\"\n" != res {
		t.Fatalf(`should be "#ffffff" but it's %#v'`, res)
	}

	buf.Reset()
	cs := []Color{0xFF0000, 0x00FF00, 0x0000FF}

	if err := json.NewEncoder(&buf).Encode(cs); err != nil {
		t.Fatal(err)
	}

	if res := buf.String(); "[\"#ff0000\",\"#00ff00\",\"#0000ff\"]\n" != res {
		t.Fatalf(`should be "["#ff0000","#00ff00","#0000ff"]" but it's %#v'`, res)
	}
}

func TestColor_UnmarshalJSON(t *testing.T) {
	val := `"#000000"`
	var c Color
	err := json.Unmarshal([]byte(val), &c)
	if err != nil {
		t.Fatal(err)
	}

	if c != 0 {
		t.Fatalf("should be 0, but its %x", c)
	}

	val = `"#ffffff"`
	err = json.Unmarshal([]byte(val), &c)
	if err != nil {
		t.Fatal(err)
	}

	if c != 0xffffff {
		t.Fatalf("should be 0xffffff, but its %x", c)
	}

	vals := `["#ff0000", "#00ff00", "#0000ff"]`
	cs := make([]Color, 0, 10)
	err = json.Unmarshal([]byte(vals), &cs)
	if err != nil {
		t.Fatal(err)
	}

	if len(cs) != 3 {
		t.Fatalf("should be length of 3, but it's %d", len(cs))
	}

	if cs[0] != 0xFF0000 {
		t.Fatalf("should be 0xFF0000, but it's %X", cs[0])
	}

	if cs[1] != 0x00FF00 {
		t.Fatalf("should be 0x00FF00, but it's %X", cs[1])
	}

	if cs[2] != 0x0000FF {
		t.Fatalf("should be 0x0000FF, but it's %X", cs[2])
	}
}
