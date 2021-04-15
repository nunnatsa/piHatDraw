package state

import "testing"

func (s *changeStack) len() int {
	if s == nil {
		return 0
	}

	res := 0
	for p := s.head; p != nil; p = p.next {
		res++
	}

	return res
}

func TestChangeStack(t *testing.T) {
	s := changeStack{}
	if s.len() != 0 {
		t.Errorf("should be empty")
	}

	s.push(&Change{ToolName: "first"})
	if s.len() != 1 {
		t.Errorf("should be with len of 1")
	}
	s.push(&Change{ToolName: "second"})
	if s.len() != 2 {
		t.Errorf("should be with len of 2")
	}

	// check LIFO:
	chng := s.pop()
	if chng == nil {
		t.Fatal("should no be bil")
	}
	if chng.ToolName != "second" {
		t.Errorf("chng.ToolName should be 'second'")
	}
	chng = s.pop()
	if chng.ToolName != "first" {
		t.Errorf("chng.ToolName should be 'first'")
	}
	chng = s.pop()
	if chng != nil {
		t.Error("chng should be nil")
	}
}
