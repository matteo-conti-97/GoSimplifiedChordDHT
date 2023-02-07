package main

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"log"
	"net/http"
	"net/rpc"
	"os"
	"sort"
	"strconv"
	"sync"
)

const REG_ADDR = "127.0.0.1"
const REG_PORT = "1234"
const ADDR = "127.0.0.1"
const PUT = true
const GET = false

type PeerInfo struct {
	Uid  string
	Ip   string
	Port string
}

type ManagedKeys struct {
	Mux  sync.Mutex
	Keys []Resource
}

type Resource struct {
	Key   string
	Value string
}

type Request struct {
	ResId   string
	ReqPeer PeerInfo
	Type    bool
}

type RpcRegistryOutPut struct {
	Out []PeerInfo
}

type GetSuccResOutPut struct {
	Out []Resource
}

// Peer RPC server thread function
func peerServer(peerInfo *PeerInfo) {
	fmt.Println("Starting Peer RPC server")
	peerServer := new(Peer)
	rpc.Register(peerServer)
	rpc.HandleHTTP()
	err := http.ListenAndServe(ADDR+":"+peerInfo.Port, nil)
	if err != nil {
		fmt.Println(err.Error())
	}
}

// Uid generation
func getUid(str string) string {
	var hash = md5.Sum([]byte(str))
	return hex.EncodeToString(hash[:])
}

func getPeerInfo(addr string, port string) *PeerInfo {
	p := PeerInfo{Uid: getUid(addr + port), Ip: addr, Port: port}
	return &p
}

type Peer []string

// Insert resourse to remote peer managed keys
func (t *Peer) RemoteGet(req Request, retRes *Resource) error {
	for _, v := range managedKeys.Keys {
		if v.Key == req.ResId {
			*retRes = v
		}
	}
	return nil
}

// Insert resourse to remote peer managed keys
func (t *Peer) RemotePut(res Resource, err *PeerInfo) error {
	managedKeys.Mux.Lock()
	managedKeys.Keys = append(managedKeys.Keys, res)

	//Sort  managed Keys
	sort.Slice(managedKeys.Keys, func(i, j int) bool {
		return managedKeys.Keys[i].Key < managedKeys.Keys[j].Key
	})
	managedKeys.Mux.Unlock()
	return nil
}

// Update predecessor successor
func (t *Peer) UpdatePred(newPeer *PeerInfo, res *string) error {
	succInfo = newPeer
	fmt.Print("\nNew successor: ")
	fmt.Println(newPeer)
	return nil
}

// Get keys to manage from the successor and update of successor predecessor
func (t *Peer) GetSuccRes(newPeer *PeerInfo, res *GetSuccResOutPut) error {
	var ret []Resource
	var temp []Resource

	managedKeys.Mux.Lock()

	//Sort  managed Keys
	sort.Slice(managedKeys.Keys, func(i, j int) bool {
		return managedKeys.Keys[i].Key < managedKeys.Keys[j].Key
	})

	//Select of resources with Uid<PeerUid
	for _, v := range managedKeys.Keys {
		if newPeer.Uid > v.Key {
			ret = append(ret, v) //res to return to new peerInfo
		} else {
			temp = append(temp, v) //res remaining to successor
		}
	}

	managedKeys.Keys = temp
	managedKeys.Mux.Unlock()

	//Predecessor update
	predInfo = newPeer
	fmt.Print("\nNew predecessor: ")
	fmt.Println(newPeer)

	//Insert Keys to send in return buffer
	res.Out = ret

	return nil
}

func (t *Peer) RouteRes(req Request, resPeerInfo *PeerInfo) error {

	//Request went accross all DHT, resource not in the DHT for GET or resource has to be managed from requesting peer for PUT
	if req.Type == GET {
		check := false
		for _, v := range managedKeys.Keys {
			if v.Key == req.ResId {
				check = true
				break
			}
		}

		if check {
			//Current peer responsible for resource
			*resPeerInfo = *peerInfo
			return nil
		} else if !check && req.ReqPeer.Uid == succInfo.Uid {
			*resPeerInfo = req.ReqPeer //Risorsa non presente nella dht il prossimo peer é il richiedente
			return nil
		} else {
			//Call RoutRes to successor
			succ, err := rpc.DialHTTP("tcp", succInfo.Ip+":"+succInfo.Port)
			if err != nil {
				log.Fatal("dialing Get RouteRes:", err)
			}
			var nestedResPeerInfo PeerInfo
			err = succ.Call("Peer.RouteRes", req, &nestedResPeerInfo)
			if err != nil {
				log.Fatal("Join DHT error:", err)
			}
			*resPeerInfo = nestedResPeerInfo
			succ.Close()
		}

	} else if req.Type == PUT {
		if req.ResId <= peerInfo.Uid && req.ResId > predInfo.Uid {
			//Current peer responsible for resource
			*resPeerInfo = *peerInfo
			return nil
		} else if succInfo.Uid < peerInfo.Uid {
			*resPeerInfo = *succInfo //Risorsa con id piu grande dell'id del nodo con id piu grande, la metto sul primo nodo
			return nil
		} else if req.ResId > peerInfo.Uid && req.ReqPeer.Uid != succInfo.Uid {
			//Call RouteRes to successor
			succ, err := rpc.DialHTTP("tcp", succInfo.Ip+":"+succInfo.Port)
			if err != nil {
				log.Fatal("dialing Put RouteRes:", err)
			}
			var nestedResPeerInfo PeerInfo
			err = succ.Call("Peer.RouteRes", req, &nestedResPeerInfo)
			if err != nil {
				log.Fatal("Join DHT error:", err)
			}
			*resPeerInfo = nestedResPeerInfo
			succ.Close()
		}
	}

	return nil
}

var managedKeys ManagedKeys
var peerInfo *PeerInfo
var succInfo *PeerInfo
var predInfo *PeerInfo

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Missing argument, usage-> go run node.go portNum, port 1234 reserved for service registry")
	}
	var port, _ = strconv.Atoi(os.Args[1])
	if port <= 1024 {
		log.Fatal("Invalid port error")
	}

	peerInfo = getPeerInfo(ADDR, os.Args[1])

	//Dial to service registry
	serviceRegistry, err := rpc.DialHTTP("tcp", REG_ADDR+":"+REG_PORT)
	if err != nil {
		log.Fatal("dialing JoinDHT:", err)
	}

	// Synchronous call to service registry
	var posInfo RpcRegistryOutPut
	err = serviceRegistry.Call("ServiceRegistry.JoinDHT", peerInfo, &posInfo)
	if err != nil {
		log.Fatal("Join DHT error:", err)
	}

	serviceRegistry.Close()
	predInfo = &(posInfo.Out[0])
	succInfo = &(posInfo.Out[1])
	fmt.Print("\nCurrent node info: ")
	fmt.Println(peerInfo)
	fmt.Print("Successor node info: ")
	fmt.Println(succInfo)
	fmt.Print("Predecessor node info: ")
	fmt.Println(predInfo)

	//If there is a successor and a predecessor retrieve resource to manage and update their pointers
	if succInfo.Uid != peerInfo.Uid {
		//Dial to successor
		succ, err := rpc.DialHTTP("tcp", succInfo.Ip+":"+succInfo.Port)
		if err != nil {
			log.Fatal("dialing GetSuccRes:", err)
		}

		managedKeys.Mux.Lock()
		// Synchronous call to service successor
		var keyToManage GetSuccResOutPut
		err = succ.Call("Peer.GetSuccRes", peerInfo, &keyToManage)
		if err != nil {
			log.Fatal("GetSuccRes error:", err)
		}
		succ.Close()
		if len(keyToManage.Out) != 0 {
			managedKeys.Keys = keyToManage.Out
		}
		managedKeys.Mux.Unlock()

	}
	//Update predecessor successor
	if predInfo.Uid != peerInfo.Uid {
		var ret string

		//Dial to predecessor
		pred, err := rpc.DialHTTP("tcp", predInfo.Ip+":"+predInfo.Port)
		if err != nil {
			log.Fatal("dialing UpdatePred:", err)
		}

		err = pred.Call("Peer.UpdatePred", peerInfo, &ret)
		if err != nil {
			log.Fatal("UpdatePred error:", err)
		}
		pred.Close()
	}

	//Peer RPC server thread
	go peerServer(peerInfo)

	//Wait for user commands
	var cmd string
	for {
		fmt.Println("What can i do for you?")
		fmt.Println("1) GET")
		fmt.Println("2) PUT")
		fmt.Println("\nInsert the number of the operation to be perfomed.")
		fmt.Scanln(&cmd)
		switch cmd {
		case "1":
			fmt.Println("Retrieved resource: " + get().Value + "\n\n")
		case "2":
			fmt.Println("Resource inserted on node: " + put().Uid + "\n\n")
		default:
			fmt.Println("Error, insert a valid command")
		}
	}

}

func get() Resource {
	var ret Resource
	var resName string
	var scanErr = "err"
	fmt.Println("\nInsert the name of the resource to get from the DHT, replace spaces with underscores")
	fmt.Scanln(&resName, &scanErr)
	if scanErr != "err" {
		fmt.Println("Error input contains blank spaces, replace them with underscores")
	}
	req := Request{ResId: getUid(resName), ReqPeer: *peerInfo, Type: GET}
	fmt.Println("Searching resource: " + resName + " with Uid: " + req.ResId)

	for _, v := range managedKeys.Keys { //If i have the resource no need to search in the DHT
		if v.Key == req.ResId {
			return v
		}
	}

	var resPeerInfo PeerInfo
	//Call RoutRes to successor
	succ, err := rpc.DialHTTP("tcp", succInfo.Ip+":"+succInfo.Port)
	if err != nil {
		log.Fatal("dialing Get starting RouteRes:", err)
	}
	err = succ.Call("Peer.RouteRes", req, &resPeerInfo)
	if err != nil {
		log.Fatal("RouteRes error:", err)
	}
	succ.Close()

	if resPeerInfo.Uid == peerInfo.Uid {
		for _, v := range managedKeys.Keys {
			if v.Key == req.ResId {
				ret = v
				return ret
			}
		}
		ret.Value = "Resource not in the DHT"

	} else {
		//Call RoutRes to successor
		resPeer, err := rpc.DialHTTP("tcp", resPeerInfo.Ip+":"+resPeerInfo.Port)
		if err != nil {
			log.Fatal("dialing RemoteGet with node: ", resPeerInfo, "Error: ", err)
		}
		err = resPeer.Call("Peer.RemoteGet", req, &ret)
		if err != nil {
			log.Fatal("Remote Get error:", err)
		}
		succ.Close()
	}

	return ret
}

func put() PeerInfo {
	var ret PeerInfo
	var resName string
	var err = "err"
	fmt.Println("\nInsert the name of the resource to put in the DHT, replace spaces with underscores")
	fmt.Scanln(&resName, &err)
	if err != "err" {
		fmt.Println("Error input contains blank spaces, replace them with underscores")
	}
	req := Request{ResId: getUid(resName), ReqPeer: *peerInfo, Type: PUT}
	res := Resource{Key: req.ResId, Value: resName}
	fmt.Println("Putting resource: " + resName + " with Uid: " + req.ResId)

	if (req.ResId <= peerInfo.Uid && req.ResId > predInfo.Uid) || (succInfo.Uid == peerInfo.Uid) { //Redundant control which eventually can save an rpc call throug all the DHT
		managedKeys.Mux.Lock()
		managedKeys.Keys = append(managedKeys.Keys, res)
		ret = *peerInfo
		//Sort  managed Keys
		sort.Slice(managedKeys.Keys, func(i, j int) bool {
			return managedKeys.Keys[i].Key < managedKeys.Keys[j].Key
		})
		managedKeys.Mux.Unlock()

	} else {
		fmt.Println("Routing")
		var resPeerInfo PeerInfo
		//Call RoutRes to successor
		succ, err := rpc.DialHTTP("tcp", succInfo.Ip+":"+succInfo.Port)
		if err != nil {
			log.Fatal("dialing Put starting RouteRes:", err)
		}
		err = succ.Call("Peer.RouteRes", req, &resPeerInfo)
		if err != nil {
			log.Fatal("RouteRes error:", err)
		}
		fmt.Println(resPeerInfo)
		succ.Close()

		if resPeerInfo.Uid == peerInfo.Uid {
			managedKeys.Mux.Lock()
			managedKeys.Keys = append(managedKeys.Keys, res)
			ret = *peerInfo
			//Sort  managed Keys
			sort.Slice(managedKeys.Keys, func(i, j int) bool {
				return managedKeys.Keys[i].Key < managedKeys.Keys[j].Key
			})
			managedKeys.Mux.Unlock()

		} else {
			//Call RoutRes to successor
			resPeer, err := rpc.DialHTTP("tcp", resPeerInfo.Ip+":"+resPeerInfo.Port)
			if err != nil {
				log.Fatal("dialing RemotePut with node: ", resPeerInfo, " Error: ", err)
			}
			err = resPeer.Call("Peer.RemotePut", res, &ret)
			if err != nil {
				log.Fatal("Remote Put error:", err)
			}
			succ.Close()
			ret = resPeerInfo
		}
	}
	return ret
}
