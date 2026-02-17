package internal

import (
	"testing"
	"time"
)

func TestWPMHistoryTracking(t *testing.T) {
	stats := NewStats()

	// Start the test
	stats.Start()

	// Simulate typing with delays to trigger snapshots
	for i := 0; i < 30; i++ {
		stats.RecordKeystroke(true)

		// Every 10 keystrokes, advance time by more than snapshot interval
		if i%10 == 0 && i > 0 {
			time.Sleep(1100 * time.Millisecond) // Longer than 1 second interval
		}
	}

	// Get history
	history := stats.GetWPMHistory()

	// Should have some snapshots (we waited twice for 1.1 seconds)
	if len(history) < 1 {
		t.Errorf("Expected at least 1 WPM snapshot, got %d", len(history))
	}

	// All snapshots should have timestamps after start time
	startTime := stats.GetStartTime()
	for i, snapshot := range history {
		if snapshot.Timestamp.Before(startTime) {
			t.Errorf("Snapshot %d has timestamp before start time", i)
		}

		// WPM should be non-negative
		if snapshot.WPM < 0 {
			t.Errorf("Snapshot %d has negative WPM: %.2f", i, snapshot.WPM)
		}
	}
}

func TestWPMHistoryEmpty(t *testing.T) {
	stats := NewStats()

	// Get history before any keystrokes
	history := stats.GetWPMHistory()

	// Should be empty
	if len(history) != 0 {
		t.Errorf("Expected empty history, got %d snapshots", len(history))
	}
}

func TestWPMHistoryImmutable(t *testing.T) {
	stats := NewStats()
	stats.Start()

	// Record some keystrokes to generate history
	for i := 0; i < 10; i++ {
		stats.RecordKeystroke(true)
	}
	time.Sleep(1100 * time.Millisecond) // Force a snapshot
	stats.RecordKeystroke(true)

	// Get history
	history1 := stats.GetWPMHistory()
	originalLen := len(history1)

	// Modify the returned slice
	if len(history1) > 0 {
		history1[0].WPM = 9999.0
	}

	// Get history again
	history2 := stats.GetWPMHistory()

	// Should have same length
	if len(history2) != originalLen {
		t.Errorf("History length changed: %d -> %d", originalLen, len(history2))
	}

	// Modification should not affect internal state
	if len(history2) > 0 && history2[0].WPM == 9999.0 {
		t.Error("Returned history is not a copy, internal state was modified")
	}
}
