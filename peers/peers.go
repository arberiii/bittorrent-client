package peers

import (
	"encoding/binary"
	"fmt"
	"github.com/arberiii/bittorrent-client/bencode"
	"github.com/arberiii/bittorrent-client/torrent"
	"io"
	"net"
	"net/http"
	"net/url"
	"strconv"
)

const MyPeerId = "19981996199920002023"

type Peer struct {
	IP   string
	Port int
}

type HTTPClient interface {
	Get(url string) (*http.Response, error)
}

func GetTrackingPeers(torrent torrent.Torrent, client HTTPClient) ([]string, error) {
	params := url.Values{}
	infoHash := fmt.Sprintf("%s", torrent.Info.Hash)
	params.Add("info_hash", infoHash)
	params.Add("peer_id", MyPeerId)
	params.Add("port", "6881")
	params.Add("uploaded", "0")
	params.Add("downloaded", "0")
	params.Add("left", strconv.Itoa(torrent.Info.Length))
	params.Add("compact", "1")
	baseUrl, err := url.Parse(torrent.Announce)
	if err != nil {
		return nil, err
	}
	baseUrl.RawQuery = params.Encode()

	response, err := client.Get(baseUrl.String())
	if err != nil {
		return nil, err
	}
	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	decodedResponse, err := bencode.DecodeBencode(string(responseBody))
	if err != nil {
		return nil, err
	}
	decodedResponseMap, _ := decodedResponse.(map[string]interface{})

	var ret []string
	peers := getAllPeers(decodedResponseMap["peers"].(string))
	for _, peer := range peers {
		ret = append(ret, fmt.Sprintf("%s:%d", peer.IP, peer.Port))
	}

	return ret, nil
}

func getAllPeers(peers string) []Peer {
	index := 0
	var ret []Peer
	for index+6 <= len(peers) {
		peer := Peer{}
		peer.IP = net.IPv4(peers[index], peers[index+1], peers[index+2], peers[index+3]).String()
		peer.Port = int(binary.BigEndian.Uint16([]byte(peers[index+4 : index+6])))
		ret = append(ret, peer)
		index += 6
	}
	return ret
}
