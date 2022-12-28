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

func main() {
	fmt.Println("hello world client")
	rand.Seed(time.Now().UnixNano())
	var pid = strconv.Itoa(os.Getpid())
	var salt = strconv.Itoa(rand.Int())
	var hash = md5.Sum([]byte(pid + salt))
	var peerID string = hex.EncodeToString(hash[:])[:10]

	client, err := rpc.DialHTTP("tcp", REG_ADDR+":"+REG_PORT)
	if err != nil {
		log.Fatal("dialing:", err)
	}
	// Synchronous call
	var reply string
	err = client.Call("ServiceRegistry.Join", peerID, &reply)
	if err != nil {
		log.Fatal("arith error:", err)
	}
	fmt.Printf("PeerID: %s=%s", peerID, reply)

}
