package main

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"log"
	"math/rand"
	"net/rpc"
	"os"
	"strconv"
	"sync"
	"time"
)

const REG_ADDR = "127.0.0.1"
const REG_PORT = "1234"
const ADDR = "127.0.0.1"

type PeerInfo struct {
	Uid  string
	Ip   string
	Port string
}

// Uid generation
func getUid() string {
	rand.Seed(time.Now().UnixNano())
	var pid = strconv.Itoa(os.Getpid())
	var salt = strconv.Itoa(rand.Int())
	var hash = md5.Sum([]byte(pid + salt))
	return hex.EncodeToString(hash[:])[:10]
}

func getPeerInfo(addr string, port string) *PeerInfo {
	p := PeerInfo{Uid: getUid(), Ip: addr, Port: port}
	return &p
}

func main() {

	// Keys managed from the node with associated mutex for thread shared access
	//https://notes.shichao.io/gopl/ch9/ vedi qui per continuare
	var (
		mux  sync.Mutex
		keys string
	)

	fmt.Println("hello world client")
	if len(os.Args) < 2 {
		log.Fatal("Missing argument, usage-> go run node.go portNum")
	}
	var port, _ = strconv.Atoi(os.Args[1])
	if port <= 1024 {
		log.Fatal("Invalid port error")
	}

	var peer *PeerInfo = getPeerInfo(ADDR, os.Args[1])

	client, err := rpc.DialHTTP("tcp", REG_ADDR+":"+REG_PORT)
	if err != nil {
		log.Fatal("dialing:", err)
	}

	// Synchronous call to service registry
	var succ PeerInfo
	err = client.Call("ServiceRegistry.Join", peer, &succ)
	if err != nil {
		log.Fatal("Join error:", err)
	}
	fmt.Printf("PeerID: Sended %s Received %s\n", peer.Uid, succ.Uid)
	fmt.Printf("PeerID: Sended %s Received %s\n", peer.Ip, succ.Ip)
	fmt.Printf("PeerID: Sended %s Received %s\n", peer.Port, succ.Port)
	//1TODO CONTATTARE IL SUCCESSORE E PRENDERE LE CHIAVI DI CUI CI SI DEVE OCCUPARE (CLIENT DI 2)
	//2TODO LANCIARE LA GOROUTINE PER ASCOLTARE LE RICHIESTE DI NUOVI PREDECESSORI E PASSARGLI LE CHIAVI (SERVER DI 1)
	//3TODO LANCIARE LA GOROUTINE PER ASCOLTARE RICHIESTE DI RISORSE ED EVENTUALMENTE FARE ROUTING
	for {
		//TODO MENU CLIENT PER FARE PUT O GET DI RISORSE ED EVENTUALMENTE INIZIARE IL ROUTING
	}

}
