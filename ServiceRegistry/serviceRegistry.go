package main

import (
	"fmt"
	"net/http"
	"net/rpc"
	"sort"
)

const REG_PORT = "1234"

const NUM_PEERS = 1024

type PeerInfo struct {
	Uid  string
	Ip   string
	Port string
}

var peerList = make([]PeerInfo, 0, NUM_PEERS)

type ServiceRegistry PeerInfo

func (t *ServiceRegistry) Join(newPeerInfo *PeerInfo, succ *PeerInfo) error {

	*succ = *newPeerInfo                      //Ritorno il valore del successore
	peerList = append(peerList, *newPeerInfo) //Aggiungo il nuovo peer e riordino la lista

	sort.Slice(peerList, func(i, j int) bool { //Ordino la lista di peerInfo per uid
		return peerList[i].Uid < peerList[j].Uid
	})

	return nil
}

func main() {
	fmt.Println("hello world server")
	serviceRegistry := new(ServiceRegistry)
	rpc.Register(serviceRegistry)
	rpc.HandleHTTP()
	err := http.ListenAndServe(":"+REG_PORT, nil)
	if err != nil {
		fmt.Println(err.Error())
	}
}
