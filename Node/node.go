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
	"time"
)

const REG_ADDR = "127.0.0.1"
const REG_PORT = "1234"

type PeerInfo struct {
	Uid  string
	Ip   string
	Port string
}

func getUid() string {
	//Id generation
	rand.Seed(time.Now().UnixNano())
	var pid = strconv.Itoa(os.Getpid())
	var salt = strconv.Itoa(rand.Int())
	var hash = md5.Sum([]byte(pid + salt))
	return hex.EncodeToString(hash[:])[:10]
}

func getPeerInfo() *PeerInfo {
	p := PeerInfo{Uid: getUid(), Ip: REG_ADDR, Port: "Prova"}
	return &p
}

func main() {
	fmt.Println("hello world client")

	var peer *PeerInfo = getPeerInfo()

	client, err := rpc.DialHTTP("tcp", REG_ADDR+":"+REG_PORT)
	if err != nil {
		log.Fatal("dialing:", err)
	}

	// Synchronous call to service registry
	var reply PeerInfo
	err = client.Call("ServiceRegistry.Join", peer, &reply)
	if err != nil {
		log.Fatal("arith error:", err)
	}
	fmt.Printf("PeerID: Sended %s Received %s", peer.Uid, reply.Uid)
	fmt.Printf("PeerID: Sended %s Received %s", peer.Ip, reply.Ip)
	fmt.Printf("PeerID: Sended %s Received %s", peer.Port, reply.Port)

}
