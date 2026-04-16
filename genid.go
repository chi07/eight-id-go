package eightid

import (
	"runtime"
	"sync/atomic"
	"time"
)

const (
	idLength        = 8
	timeWidth       = 6
	sequenceWidth   = 2
	sequenceBase    = 62 * 62
	sequenceMask    = 1<<12 - 1
	stateShift      = 12
	customEpochUnix = 1735689600 // 2025-01-01T00:00:00Z
)

const charset = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

var state atomic.Uint64

func init() {
	now := currentTick()
	state.Store(now << stateShift)
}

func currentTick() uint64 {
	return tickFromUnixMilli(time.Now().UnixMilli())
}

func tickFromUnixMilli(unixMilli int64) uint64 {
	now := unixMilli - customEpochUnix*1000
	if now < 0 {
		return 0
	}
	return uint64(now / 10)
}

func appendBase62Fixed(dst []byte, value uint64, width int) {
	for i := width - 1; i >= 0; i-- {
		dst[i] = charset[value%62]
		value /= 62
	}
}

// New generates a unique, 8-character, lexicographically sortable, case-sensitive ID.
func New() string {
	tick, seq := nextState()
	return encodeID(tick, seq)
}

// NewWithTime generates an ID using the provided time for the sortable prefix.
// It is useful for tests, deterministic fixtures, or backfilling data.
func NewWithTime(t time.Time) string {
	return encodeID(tickFromUnixMilli(t.UnixMilli()), 0)
}

func encodeID(tick uint64, seq uint32) string {
	var buf [idLength]byte
	appendBase62Fixed(buf[:timeWidth], tick, timeWidth)
	appendBase62Fixed(buf[timeWidth:], uint64(seq), sequenceWidth)
	return string(buf[:])
}

func nextState() (uint64, uint32) {
	for {
		now := currentTick()
		current := state.Load()
		lastTick := current >> stateShift
		lastSeq := uint32(current & sequenceMask)

		nextTick := now
		nextSeq := uint32(0)

		switch {
		case now < lastTick:
			nextTick = lastTick
			if lastSeq >= sequenceBase-1 {
				runtime.Gosched()
				continue
			}
			nextSeq = lastSeq + 1
		case now == lastTick:
			if lastSeq >= sequenceBase-1 {
				runtime.Gosched()
				continue
			}
			nextTick = lastTick
			nextSeq = lastSeq + 1
		}

		next := (nextTick << stateShift) | uint64(nextSeq)
		if state.CompareAndSwap(current, next) {
			return nextTick, nextSeq
		}
	}
}

// IsValid checks whether id is 8 ASCII alphanumeric characters.
func IsValid(id string) bool {
	if len(id) != idLength {
		return false
	}
	for i := 0; i < len(id); i++ {
		ch := id[i]
		if !((ch >= 'a' && ch <= 'z') ||
			(ch >= 'A' && ch <= 'Z') ||
			(ch >= '0' && ch <= '9')) {
			return false
		}
	}
	return true
}
