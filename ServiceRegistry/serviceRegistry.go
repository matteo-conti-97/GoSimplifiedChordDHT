package main

import (
	"fmt"
	"net/http"
	"net/rpc"
	"sort"
)

const REG_PORT = "1234"
const REG_ADDR = "127.0.0.1"

type PeerInfo struct {
	Uid  string
	Ip   string
	Port string
}

type RpcRegistryOutPut struct {
	Out []PeerInfo
}

var peerList = make([]PeerInfo, 0)

// Search joining node successor and predecessor
func findSuccAndPred(peer PeerInfo, peerList []PeerInfo) (PeerInfo, PeerInfo) {

	succ, pred := peer, peer
	if len(peerList) > 1 {
		for i := range peerList {
			if peerList[i].Uid > peer.Uid {
				succ = peerList[i]
				pred = peerList[(i-1)%len(peerList)]
				break
			}
		}
	} else if len(peerList) == 1 {
		succ, pred = peerList[0], peerList[0]
	}

	return succ, pred
}

// RPC, return to the joining node his successor
type ServiceRegistry PeerInfo

func (t *ServiceRegistry) JoinDHT(newPeerInfo *PeerInfo, ret *RpcRegistryOutPut) error {
	var succ, pred PeerInfo
	//Ricerca e ritorno del successpre
	succ, pred = findSuccAndPred(*newPeerInfo, peerList)
	ret.Out = append(ret.Out, succ)
	ret.Out = append(ret.Out, pred)

	//Aggiungo il nuovo peer e riordino la lista per uid
	peerList = append(peerList, *newPeerInfo)
	fmt.Print("\nNode ")
	fmt.Print(newPeerInfo)
	fmt.Println(" joined the DHT")
	fmt.Println("\nDHT updated node list:")
	fmt.Println(peerList)

	sort.Slice(peerList, func(i, j int) bool {
		return peerList[i].Uid < peerList[j].Uid
	})

	return nil
}

func main() {
	fmt.Println("Starting ServiceRegistry RPC server")
	//Listen di nuovi nodi
	serviceRegistryServer := new(ServiceRegistry)
	rpc.Register(serviceRegistryServer)
	rpc.HandleHTTP()
	err := http.ListenAndServe(REG_ADDR+":"+REG_PORT, nil)
	if err != nil {
		fmt.Println(err.Error())
	}
}
