package main

import (
	"fmt"
	"net/http"
	"net/rpc"
	"sort"
)

const REG_PORT = "1234"

type PeerInfo struct {
	Uid  string
	Ip   string
	Port string
}

var peerList = make([]PeerInfo, 0)

func findSuccessor(peer PeerInfo, peerList []PeerInfo) PeerInfo {
	if len(peerList) == 0 {
		return peer
	}

	for _, v := range peerList {
		if v.Uid > peer.Uid {
			return v
		}
	}
	return peerList[0]
}

// RPC, return to the joining node his successor
type ServiceRegistry PeerInfo

func (t *ServiceRegistry) Join(newPeerInfo *PeerInfo, succ *PeerInfo) error {

	//Ricerca e ritorno del successpre
	*succ = findSuccessor(*newPeerInfo, peerList)

	//Aggiungo il nuovo peer e riordino la lista per uid
	peerList = append(peerList, *newPeerInfo)

	sort.Slice(peerList, func(i, j int) bool {
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
