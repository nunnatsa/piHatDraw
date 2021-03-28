package hat

import "testing"

func TestToHatColorRed(t *testing.T) {
	if res := toHatColor(0xFFFFFF); res != 0xFFFF {
		t.Errorf("should be 0xFFFF, but it 0x%X", res)
	}

	if res := toHatColor(0b000001110000000000000000); res != 0 {
		t.Errorf("should be 0, but it 0x%X", res)
	}
}

func TestToHatColorGreen(t *testing.T) {
	if res := toHatColor(0b111110000000000000000000); res != 0xF800 {
		t.Errorf("should be 0xF800, but it 0x%X", res)
	}

	if res := toHatColor(0b000000000000001100000000); res != 0 {
		t.Errorf("should be 0, but it 0x%X", res)
	}

	if res := toHatColor(0b000000001111110000000000); res != 0x07E0 {
		t.Errorf("should be 0x07E0, but it 0x%X", res)
	}
}

func TestToHatColorBlue(t *testing.T) {
	if res := toHatColor(0b000000000000000000000111); res != 0 {
		t.Errorf("should be 0, but it 0x%X", res)
	}

	if res := toHatColor(0b000000000000000011111000); res != 0x01F {
		t.Errorf("should be 0x1F, but it 0x%X", res)
	}
}
