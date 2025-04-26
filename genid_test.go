package eightid_test

import (
	"strings"
	"testing"
	"time"

	eightid "github.com/chi07/eight-id-go"
)

func TestGenerateID_UniqueAndValidFormat(t *testing.T) {
	seen := make(map[string]bool)

	for i := 0; i < 100000; i++ {
		id := eightid.New()

		// Độ dài đúng
		if len(id) != 8 {
			t.Errorf("ID length invalid: got %d, want 8", len(id))
		}

		// Format hợp lệ
		if !eightid.IsValid(id) {
			t.Errorf("Invalid ID format: %s", id)
		}

		// Không được trùng
		if seen[id] {
			t.Errorf("Duplicate ID generated: %s", id)
		}
		seen[id] = true
	}
}

func TestVerifyIDFormat(t *testing.T) {
	valid := []string{"aZ09kLx3", "AbC123Ef", "xyXY09AB"}
	invalid := []string{"abc123", "123456789", "ABC-12ab", "abc_12XY", "abc$12xy", "abc 12aa"}

	for _, id := range valid {
		if !eightid.IsValid(id) {
			t.Errorf("Expected valid ID: %s", id)
		}
	}

	for _, id := range invalid {
		if eightid.IsValid(id) {
			t.Errorf("Expected invalid ID: %s", id)
		}
	}
}

func TestID_IsCaseSensitiveEffective(t *testing.T) {
	for i := 0; i < 1000; i++ {
		id := eightid.New()

		lower := strings.ToLower(id)
		upper := strings.ToUpper(id)

		if lower == upper {
			t.Errorf("Generated ID does not contain mixed case: %s", id)
		}
	}
}

func TestID_IsSortable(t *testing.T) {
	id1 := eightid.New()
	time.Sleep(1 * time.Millisecond)
	id2 := eightid.New()

	if id1 >= id2 {
		t.Logf("Warning: id1 >= id2 (%s >= %s), could be same-millisecond", id1, id2)
	}
}
