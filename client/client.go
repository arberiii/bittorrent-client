package client

import (
	"crypto/sha1"
	"fmt"
	"github.com/arberiii/bittorrent-client/peers"
	"github.com/arberiii/bittorrent-client/torrent"
	"net"
	"os"
	"strconv"
)

type TorrentClient struct {
	Torrent torrent.Torrent
	Peers   []string
}

func (tc *TorrentClient) DoHandshake(peerIndex int) *net.TCPConn {
	if peerIndex >= len(tc.Peers) {
		return nil
	}
	peerAddr := tc.Peers[peerIndex]
	tcpAddr, err := net.ResolveTCPAddr("tcp", peerAddr)
	if err != nil {
		fmt.Println("Resolving the address failed:", err.Error())
		return nil
	}

	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		fmt.Println("Establishing the connection failed:", err.Error())
		return nil
	}

	echo := []byte{}
	echo = append(echo, 19)
	echo = append(echo, []byte("BitTorrent protocol")...)
	echo = append(echo, make([]byte, 8)...)
	echo = append(echo, tc.Torrent.Info.Hash[:]...)
	echo = append(echo, []byte(peers.MyPeerId)...)

	_, err = conn.Write(echo)
	if err != nil {
		fmt.Println("Write to server failed:", err.Error())
		return nil
	}

	reply := make([]byte, 1024)

	_, err = conn.Read(reply)
	if err != nil {
		fmt.Println("Failed to read the reply:", err.Error())
		return nil
	}
	fmt.Println("Connection established to peer: ", peerAddr)
	return conn
}

func (tc *TorrentClient) WriteInterested(conn *net.TCPConn) *net.TCPConn {
	interested := []byte{0, 0, 0, 1, 2}
	_, err := conn.Write(interested)
	if err != nil {
		fmt.Println("Write to server failed:", err.Error())
		return nil
	}

	fmt.Println("Sent interested to peer")
	return conn
}

func (tc *TorrentClient) ReadUnchoke(conn *net.TCPConn) *net.TCPConn {
	reply := make([]byte, 1024)
	_, err := conn.Read(reply)
	if err != nil {
		fmt.Println("Failed to read the reply:", err.Error())
		return nil
	}
	fmt.Println("Peer unchoked us")
	return conn
}

func (tc *TorrentClient) Download(conn *net.TCPConn) {
	totalLength := tc.Torrent.Info.Length
	pieceLength := tc.Torrent.Info.PieceLength
	downloadingPiece := 0
	totalDownloaded := 0
	for totalDownloaded < totalLength {
		if totalDownloaded+pieceLength > totalLength {
			pieceLength = totalLength - totalDownloaded
		}
		result, err := downloadPiece(conn, downloadingPiece, pieceLength)
		if err != nil {
			fmt.Println("Failed to download piece:", err.Error())
			return
		}
		if !checkPieceHash(result, tc.Torrent.Info.PieceHashes[downloadingPiece]) {
			fmt.Println("Piece hash doesn't match")
			return
		}
		err = storePiece(result, downloadingPiece)
		if err != nil {
			fmt.Println("Failed to store piece:", err.Error())
			return
		}
		totalDownloaded += pieceLength
		downloadingPiece += 1
		fmt.Printf("Downloaded piece: %d/%d\n", downloadingPiece, len(tc.Torrent.Info.PieceHashes))
	}
}

func downloadPiece(conn *net.TCPConn, pieceIndex int, pieceLength int) ([]byte, error) {
	var downloadedPiece []byte
	downloaded := 0
	for downloaded < pieceLength {
		request := []byte{0, 0, 0, 13, 6}

		request = append(request, convertNumberToByte(pieceIndex)...)
		request = append(request, convertNumberToByte(downloaded)...)

		length := 16 * 1024
		if downloaded+length >= pieceLength {
			length = pieceLength - downloaded
		}
		request = append(request, convertNumberToByte(length)...)

		_, err := conn.Write(request)
		if err != nil {
			return nil, err
		}

		pieceReply := make([]byte, length+13)
		totalRead := 0
		for totalRead < length+13 {
			n, err := conn.Read(pieceReply[totalRead:])
			if err != nil {
				return nil, err
			}
			totalRead += n
		}
		downloaded += length
		downloadedPiece = append(downloadedPiece, pieceReply[13:]...)
	}

	return downloadedPiece, nil
}

func checkPieceHash(piece []byte, hash string) bool {
	hashDownloadedPiece := sha1.Sum(piece)
	hashDownloadedString := fmt.Sprintf("%x", hashDownloadedPiece)
	return hashDownloadedString == hash
}

func storePiece(piece []byte, pieceIndex int) error {
	err := os.WriteFile("/tmp/dat"+strconv.Itoa(pieceIndex), piece, 0644)
	if err != nil {
		return err
	}
	return nil
}

func convertNumberToByte(n int) []byte {
	var b [4]byte
	b[0] = byte(n >> 24)
	b[1] = byte(n >> 16)
	b[2] = byte(n >> 8)
	b[3] = byte(n)
	return b[:]
}
