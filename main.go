package main

import (
	"fmt"
	"github.com/arberiii/bittorrent-client/peers"
	"github.com/arberiii/bittorrent-client/torrent"
	"net/http"
)

func main() {
	fmt.Println("Hello, world!")
	client := &http.Client{}
	torrent, err := torrent.GetTorrentInfoFromFile("torrent/sample.torrent")
	peers, err := peers.GetTrackingPeers(torrent, client)
	fmt.Println(peers, err)
}
