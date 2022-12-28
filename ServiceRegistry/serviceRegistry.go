package main

import (
	"fmt"
	"net/http"
	"net/rpc"
	"sort"
)

const PORT = "1234"

const NUM_PEERS = 1024

var peerList = make([]string, 0, NUM_PEERS)

type ServiceRegistry string

func (t *ServiceRegistry) Join(newPeerID string, succ *string) error {
	*succ = newPeerID                      //Ritorno il valore del successore
	peerList = append(peerList, newPeerID) //Aggiungo il nuovo peer e riordino la lista
	sort.Strings(peerList)
	return nil
}

func main() {
	fmt.Println("hello world server")
	serviceRegistry := new(ServiceRegistry)
	rpc.Register(serviceRegistry)
	rpc.HandleHTTP()
	err := http.ListenAndServe(":"+PORT, nil)
	if err != nil {
		fmt.Println(err.Error())
	}
}
