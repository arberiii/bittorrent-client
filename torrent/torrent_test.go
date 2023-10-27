package torrent

import (
	"testing"
)

func TestGetTorrentInfoFromFile(t *testing.T) {
	t.Run("Simple torrent", func(t *testing.T) {
		result, err := GetTorrentInfoFromFile("sample.torrent")
		if err != nil {
			t.Errorf("Error: %s", err)
		}

		if result.Announce != "http://bittorrent-test-tracker.codecrafters.io/announce" {
			t.Errorf("Expected: %s, got: %s", "http://bittorrent-test-tracker.codecrafters.io/announce", result.Announce)
		}

		if result.Info.Length != 92063 {
			t.Errorf("Expected: %d, got: %d", 92063, result.Info.Length)
		}
	})
}
