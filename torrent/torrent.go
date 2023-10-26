package torrent

import (
	"crypto/sha1"
	"fmt"
	"github.com/arberiii/bittorrent-client/bencode"
	"os"
)

type Info struct {
	Name        string
	PieceLength int
	PieceHashes []string
	Length      int
	Hash        [20]byte
}

type Torrent struct {
	Announce string
	Info     Info
}

func GetTorrentInfoFromFile(filePath string) (Torrent, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return Torrent{}, err
	}

	decoded, err := bencode.DecodeBencode(string(data))
	if err != nil {
		return Torrent{}, err
	}

	decodedMap, _ := decoded.(map[string]interface{})
	infoMap, _ := decodedMap["info"].(map[string]interface{})

	info := Info{}
	info.Length = infoMap["length"].(int)
	info.PieceLength = infoMap["piece length"].(int)
	info.PieceHashes = getAllHashesFromPieces(infoMap["pieces"].(string))

	d := bencode.EncodeBencode(infoMap)
	info.Hash = sha1.Sum([]byte(d))

	torrent := Torrent{
		Announce: decodedMap["announce"].(string),
		Info:     info,
	}

	return torrent, nil
}

func getAllHashesFromPieces(pieces string) []string {
	index := 0
	var ret []string
	for index+20 <= len(pieces) {
		hashBytes := []byte(pieces[index : index+20])

		hashString := fmt.Sprintf("%x", hashBytes)
		ret = append(ret, hashString)
		index += 20
	}
	return ret
}
