package peers

import (
	"bytes"
	"github.com/arberiii/bittorrent-client/torrent"
	"io/ioutil"
	"net/http"
	"testing"
)

type MockHTTPClient struct {
	Response *http.Response
	Err      error
}

func (m *MockHTTPClient) Get(url string) (*http.Response, error) {
	return m.Response, m.Err
}

func TestGetTorrentInfoFromFile(t *testing.T) {
	t.Run("Peers from torrent", func(t *testing.T) {
		mockResponse := &http.Response{
			StatusCode: http.StatusOK, // Adjust the status code as needed.
			Body:       ioutil.NopCloser(bytes.NewReader([]byte{100, 56, 58, 99, 111, 109, 112, 108, 101, 116, 101, 105, 51, 101, 49, 48, 58, 105, 110, 99, 111, 109, 112, 108, 101, 116, 101, 105, 49, 101, 56, 58, 105, 110, 116, 101, 114, 118, 97, 108, 105, 54, 48, 101, 49, 50, 58, 109, 105, 110, 32, 105, 110, 116, 101, 114, 118, 97, 108, 105, 54, 48, 101, 53, 58, 112, 101, 101, 114, 115, 49, 56, 58, 178, 62, 85, 20, 201, 33, 178, 62, 82, 89, 201, 14, 165, 232, 33, 77, 201, 11, 101})),
		}
		mockClient := &MockHTTPClient{
			Response: mockResponse,
			Err:      nil,
		}

		simpleTorrent := torrent.Torrent{
			Announce: "http://bittorrent-test-tracker.codecrafters.io/announce",
			Info: torrent.Info{
				Length:      92063,
				PieceLength: 16384,
				PieceHashes: []string{
					"0x123",
				},
			},
		}
		actual, err := GetTrackingPeers(simpleTorrent, mockClient)
		if err != nil {
			t.Errorf("Error: %s", err)
		}

		expected := []string{"178.62.85.20:51489", "178.62.82.89:51470", "165.232.33.77:51467"}
		for i, v := range actual {
			if v != expected[i] {
				t.Errorf("Expected: %s, got: %s", expected[i], v)
			}
		}
	})
}
