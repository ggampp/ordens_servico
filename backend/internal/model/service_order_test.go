package model

import "testing"

func TestCanTransition(t *testing.T) {
	cases := []struct {
		from, to string
		want     bool
	}{
		{StatusOpen, StatusAssigned, true},
		{StatusOpen, StatusInProgress, true},
		{StatusOpen, StatusCompleted, false},
		{StatusAssigned, StatusInProgress, true},
		{StatusInProgress, StatusCompleted, true},
		{StatusInProgress, StatusOpen, false},
		{StatusCompleted, StatusOpen, false},
		{StatusCancelled, StatusInProgress, false},
		{StatusInProgress, StatusInProgress, true}, // no-op allowed
	}
	for _, c := range cases {
		if got := CanTransition(c.from, c.to); got != c.want {
			t.Errorf("CanTransition(%q,%q) = %v, want %v", c.from, c.to, got, c.want)
		}
	}
}

func TestPagination(t *testing.T) {
	p := Pagination{Page: 3, PageSize: 25}
	if p.Offset() != 50 {
		t.Errorf("Offset() = %d, want 50", p.Offset())
	}
	p2 := Pagination{Page: 0, PageSize: 0}
	p2.Normalize()
	if p2.Page != 1 || p2.PageSize != 20 {
		t.Errorf("Normalize() = %+v, want page 1 size 20", p2)
	}
	p3 := Pagination{Page: 1, PageSize: 500}
	p3.Normalize()
	if p3.PageSize != 100 {
		t.Errorf("Normalize() PageSize = %d, want 100", p3.PageSize)
	}
}
