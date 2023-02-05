package main

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"log"
	"net/rpc"
	"os"
	"sort"
	"strconv"
	"sync"
)

const REG_ADDR = "127.0.0.1"
const REG_PORT = "1234"
const ADDR = "127.0.0.1"
const UID_DIM = 10

type PeerInfo struct {
	Uid  string
	Ip   string
	Port string
}

type ManagedKeys struct {
	mux  sync.Mutex
	keys []string
}

// Uid generation
func getUid(port string) string {
	var hash = md5.Sum([]byte(ADDR + port))
	return hex.EncodeToString(hash[:])[:10]
}

func getPeerInfo(addr string, port string) *PeerInfo {
	p := PeerInfo{Uid: getUid(port), Ip: addr, Port: port}
	return &p
}

type Peer []string

func (t *Peer) GetSuccRes(newPeerUid string, keys *[]string) error {
	var ret []string
	var temp []string
	managedKeys.mux.Lock()
	defer managedKeys.mux.Unlock()

	//Sort delle managed keys, FORSE INUTILE
	sort.Slice(managedKeys.keys, func(i, j int) bool {
		return managedKeys.keys[i] < managedKeys.keys[j]
	})

	//Select of resources with Uid<PeerUid
	for _, v := range managedKeys.keys {
		if newPeerUid > v {
			ret = append(ret, v) //res to return to new peer
		} else {
			temp = append(temp, v) //res remaining to successor
		}
	}

	managedKeys.keys = temp
	managedKeys.mux.Unlock()

	//Sort del ret
	sort.Slice(ret, func(i, j int) bool {
		return ret[i] < ret[j]
	})

	//Insert keys to send in return buffer
	*keys = ret

	return nil
}

var managedKeys ManagedKeys

func main() {

	//List of node managed keys with associated mutex for thread shared access
	//https://notes.shichao.io/gopl/ch9/ vedi qui per continuare
	//https://stackoverflow.com/questions/73859360/lock-slice-before-reading-and-modifying-it

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
	err = client.Call("ServiceRegistry.GetSuccRes", peer, &succ)
	if err != nil {
		log.Fatal("GetSuccRes error:", err)
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
