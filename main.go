package main

import (
	"fmt"
	"github.com/arberiii/bittorrent-client/torrent"
)

func main() {
	fmt.Println("Hello, world!")
	//val, err := bencode.DecodeBencode("i42e")

	torrent, err := torrent.GetTorrentInfoFromFile("torrent/sample.torrent")
	fmt.Println(torrent, err)
}
