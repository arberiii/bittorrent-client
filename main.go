package main

import (
	"fmt"
	"github.com/arberiii/bittorrent-client/client"
	"github.com/arberiii/bittorrent-client/peers"
	"github.com/arberiii/bittorrent-client/torrent"
	"net/http"
)

func main() {
	httpClient := &http.Client{}
	torrent_, err := torrent.GetTorrentInfoFromFile("torrent/sample.torrent")
	if err != nil {
		fmt.Println("Failed to get torrent")
		return
	}
	peers_, err := peers.GetTrackingPeers(torrent_, httpClient)
	if err != nil {
		fmt.Println("Failed to get peers")
		return
	}
	torrentClient := client.TorrentClient{Torrent: torrent_, Peers: peers_}
	torrentConnection := torrentClient.DoHandshake(0)
	if torrentConnection == nil {
		fmt.Println("Failed to establish connection")
		return
	}

	torrentConnection = torrentClient.WriteInterested(torrentConnection)
	if torrentConnection == nil {
		fmt.Println("Failed to write interested")
		return
	}

	torrentConnection = torrentClient.ReadUnchoke(torrentConnection)
	if torrentConnection == nil {
		fmt.Println("Failed to read unchoke")
		return
	}

	torrentClient.Download(torrentConnection)
}
