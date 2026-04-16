package eightid_test

import (
	"slices"
	"sync"
	"testing"
	"time"

	eightid "github.com/chi07/eight-id-go"
)

func TestGenerateID_UniqueAndValidFormat(t *testing.T) {
	seen := make(map[string]struct{}, 100000)

	for i := 0; i < 100000; i++ {
		id := eightid.New()

		if len(id) != 8 {
			t.Errorf("ID length invalid: got %d, want 8", len(id))
		}

		if !eightid.IsValid(id) {
			t.Errorf("Invalid ID format: %s", id)
		}

		if _, ok := seen[id]; ok {
			t.Errorf("Duplicate ID generated: %s", id)
		}
		seen[id] = struct{}{}
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
	if eightid.IsValid("abcdEF12") != true {
		t.Fatal("expected mixed-case alphanumeric ID to be valid")
	}
	if eightid.IsValid("abcdEF1-") != false {
		t.Fatal("expected non-alphanumeric ID to be invalid")
	}
}

func TestID_IsSortable(t *testing.T) {
	ids := make([]string, 4096)
	for i := range ids {
		ids[i] = eightid.New()
	}

	if !slices.IsSorted(ids) {
		t.Fatal("generated IDs are not lexicographically sorted")
	}
}

func TestNewWithTime_IsDeterministicAndSortable(t *testing.T) {
	t1 := time.Date(2026, time.January, 2, 3, 4, 5, 120*int(time.Millisecond), time.UTC)
	t2 := t1.Add(20 * time.Millisecond)

	id1 := eightid.NewWithTime(t1)
	id2 := eightid.NewWithTime(t1)
	id3 := eightid.NewWithTime(t2)

	if id1 != id2 {
		t.Fatalf("expected deterministic IDs for same timestamp: %s != %s", id1, id2)
	}
	if id1 >= id3 {
		t.Fatalf("expected later timestamp to sort after earlier timestamp: %s >= %s", id1, id3)
	}
	if !eightid.IsValid(id1) || !eightid.IsValid(id3) {
		t.Fatal("NewWithTime should generate valid IDs")
	}
}

func TestID_IsUniqueInParallel(t *testing.T) {
	const workers = 8
	const perWorker = 10000

	ids := make(chan string, workers*perWorker)
	var wg sync.WaitGroup

	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < perWorker; j++ {
				ids <- eightid.New()
			}
		}()
	}

	wg.Wait()
	close(ids)

	seen := make(map[string]struct{}, workers*perWorker)
	for id := range ids {
		if _, ok := seen[id]; ok {
			t.Fatalf("duplicate ID generated in parallel: %s", id)
		}
		seen[id] = struct{}{}
	}
}

func BenchmarkNew(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = eightid.New()
	}
}

func BenchmarkIsValid(b *testing.B) {
	id := eightid.New()
	for i := 0; i < b.N; i++ {
		_ = eightid.IsValid(id)
	}
}

func BenchmarkNewWithTime(b *testing.B) {
	ts := time.Date(2026, time.January, 2, 3, 4, 5, 0, time.UTC)
	for i := 0; i < b.N; i++ {
		_ = eightid.NewWithTime(ts)
	}
}
