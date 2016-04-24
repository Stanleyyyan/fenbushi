package libkademlia

// Contains the core kademlia type. In addition to core state, this type serves
// as a receiver for the RPC methods, which is required by that package.

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"strconv"
	"errors"
	// "os"
)

const (
	alpha = 3
	b     = 8 * IDBytes
	// k     = 20
)

type Pair struct {
	Key ID
	Value []byte
}

// Kademlia type. You can put whatever state you need in this.
type Kademlia struct {
	NodeID      ID
	SelfContact Contact
	K_buckets	RoutingTable
	PingChan	chan *Contact
	HashChan	chan Pair
	H_Table		map[ID][]byte
}

func NewKademliaWithId(laddr string, nodeID ID) *Kademlia {
	k := new(Kademlia)
	k.NodeID = nodeID
	k.K_buckets = RoutingTable{buckets: make([][]Contact, 160)}
	k.PingChan = make(chan *Contact)
	k.HashChan = make(chan Pair)
	k.H_Table = make(map [ID][]byte)
	// TODO: Initialize other state here as you add functionality.

	// Set up RPC server
	// NOTE: KademliaRPC is just a wrapper around Kademlia. This type includes
	// the RPC functions.

	s := rpc.NewServer()
	s.Register(&KademliaRPC{k})
	hostname, port, err := net.SplitHostPort(laddr)
	if err != nil {
		return nil
	}
	s.HandleHTTP(rpc.DefaultRPCPath+port,
		rpc.DefaultDebugPath+port)
	l, err := net.Listen("tcp", laddr)
	if err != nil {
		log.Fatal("Listen: ", err)
	}

	// Run RPC server forever.
	go http.Serve(l, nil)

	// Add self contact
	hostname, port, _ = net.SplitHostPort(l.Addr().String())
	port_int, _ := strconv.Atoi(port)
	ipAddrStrings, err := net.LookupHost(hostname)
	var host net.IP
	for i := 0; i < len(ipAddrStrings); i++ {
		host = net.ParseIP(ipAddrStrings[i])
		if host.To4() != nil {
			break
		}
	}
	k.SelfContact = Contact{k.NodeID, host, uint16(port_int)}
	return k
}

func NewKademlia(laddr string) *Kademlia {
	return NewKademliaWithId(laddr, NewRandomID())
}

type ContactNotFoundError struct {
	id  ID
	msg string
}

func (e *ContactNotFoundError) Error() string {
	return fmt.Sprintf("%x %s", e.id, e.msg)
}

func (k *Kademlia) FindContact(nodeId ID) (*Contact, error) {
	// TODO: Search through contacts, find specified ID
	// Find contact with provided ID
	dis := nodeId.Xor(k.NodeID)
	numOfBucket := 159 - dis.PrefixLen()
	fmt.Println("num of Bucket is :" ,numOfBucket)
	fmt.Println(len(k.K_buckets.buckets))
	for _, c1 := range k.K_buckets.buckets[numOfBucket] {
		if c1.NodeID.Equals(nodeId) {
			return &c1, nil
		}
	}
	return nil, &ContactNotFoundError{nodeId, "Not found"}
}

type CommandFailed struct {
	msg string
}

func (e *CommandFailed) Error() string {
	return fmt.Sprintf("%s", e.msg)
}

func (k *Kademlia) DoPing(host net.IP, port uint16) (*Contact, error) {
	// TODO: Implement
	portnum := strconv.Itoa(int(port))
	temp := host.String() + ":" + portnum
	fmt.Println("DoPing:", temp)
	//

	client, err := rpc.DialHTTPPath("tcp", host.String() + ":" + portnum,
		rpc.DefaultRPCPath + portnum)

	if err != nil {
		fmt.Println("error!")
		log.Fatal("dialing:", err)
	}
	pim := PingMessage{k.SelfContact, NewRandomID()}
	pom := new(PongMessage)
  	err = client.Call("KademliaRPC.Ping", pim, pom)
	if err != nil {
		return nil, &CommandFailed{
			"Unable to ping " + fmt.Sprintf("%s:%v", host.String(), port)}
	} else {
		k.PingChan <- &(pom.Sender)
		fmt.Println("updating Done")
		return &(pom.Sender), nil
	}
}

func (k *Kademlia) UpdateRT(c *Contact) error {
	fmt.Println("")
	fmt.Println("NodeID: ", c.NodeID)
	dis := c.NodeID.Xor(k.NodeID)
	numOfBucket := 159 - dis.PrefixLen()
	containSender := false
	idx := 0
	fmt.Println("num of Bucket is :" ,numOfBucket)
	fmt.Println(len(k.K_buckets.buckets))
	for index, c1 := range k.K_buckets.buckets[numOfBucket] {
		if c1.NodeID.Equals(c.NodeID) {
			containSender = true
			idx = index
		}
	}
	if containSender {
		if len(k.K_buckets.buckets[numOfBucket]) == 1 {
			fmt.Println("size is 1, add done")
			fmt.Printf("kademlia> ")
			return nil
		}
		fmt.Println("index is :", idx)
		temp := k.K_buckets.buckets[numOfBucket][idx]
		k.K_buckets.buckets[numOfBucket] =
				append(k.K_buckets.buckets[numOfBucket][:idx - 1],
						k.K_buckets.buckets[numOfBucket][idx + 1:]...)
		k.K_buckets.buckets[numOfBucket] = append(k.K_buckets.buckets[numOfBucket], temp)
				fmt.Println("Moved to Tail!")
				fmt.Printf("kademlia> ")
				return errors.New("Move to tail")
	} else {
		if len(k.K_buckets.buckets[numOfBucket]) < 20 {
			k.K_buckets.buckets[numOfBucket] = append(k.K_buckets.buckets[numOfBucket], *c)
			fmt.Println("k buckets not full, add to tail")
			fmt.Printf("kademlia> ")
			return errors.New("k buckets not full, add to tail")
		} else {
			_, err :=
				k.DoPing(k.K_buckets.buckets[numOfBucket][0].Host, k.K_buckets.buckets[numOfBucket][0].Port)
			if err != nil {
				k.K_buckets.buckets[numOfBucket] =
					k.K_buckets.buckets[numOfBucket][1:]
				k.K_buckets.buckets[numOfBucket] = append(k.K_buckets.buckets[numOfBucket], *c)
				fmt.Println("head dead, replace head")
				fmt.Printf("kademlia> ")
				return errors.New("head dead, replace head")
			}else{
				fmt.Println("Discard!")
				fmt.Printf("kademlia> ")
				return errors.New("Discard")
			}
		}
	}
	return nil
}

func (k *Kademlia) HandlePing(){
	for {
		select {
		case newContact := <- k.PingChan:
			k.UpdateRT(newContact)
		default:
		}
	}
}

func (k *Kademlia) DoStore(contact *Contact, key ID, value []byte) error {
	//TODO: Implement
	portnum := strconv.Itoa(int(contact.Port))
	temp := contact.Host.String() + ":" + portnum
	fmt.Println("DoStore:", temp)
	//

	conn, err := rpc.DialHTTPPath("tcp", contact.Host.String() + ":" + portnum,
		rpc.DefaultRPCPath + portnum)

	if err != nil {
		fmt.Println("error!")
		log.Fatal("dialing:", err)
	}
	req := StoreRequest{Sender: k.SelfContact, MsgID: NewRandomID(), Key: key, Value: value}
	res := new(StoreResult)
  	err = conn.Call("KademliaRPC.Store", req, res)
  	if err != nil {
  		fmt.Println(err)
		log.Fatal("RPC:", err)
		return &CommandFailed{"DoStore Failed"}
  	} else {
  		k.PingChan <- contact
  		fmt.Println("DoStore complete!")
  	}
  	return nil
}

func (k *Kademlia) UpdateHT(key ID, value []byte) error {
	k.H_Table[key] = value
	return nil
}

func (k *Kademlia) HandleStore(){
	for {
		select {
		case newPair := <- k.HashChan:
			k.UpdateHT(newPair.Key, newPair.Value)
		}
	}
}


func (k *Kademlia) DoFindNode(contact *Contact, searchKey ID) ([]Contact, error) {
	// TODO: Implement

	portnum := strconv.Itoa(int(contact.Port))
	temp := contact.Host.String() + ":" + portnum
	fmt.Println("DoStore:", temp)
	//
	conn, err := rpc.DialHTTPPath("tcp", contact.Host.String() + ":" + portnum,
		rpc.DefaultRPCPath + portnum)
	if err != nil {
		fmt.Println("error!")
		log.Fatal("dialing:", err)
	}
	req := FindNodeRequest{Sender: k.SelfContact, MsgID: NewRandomID(), NodeID: searchKey}
	res := new(FindNodeResult)
  	err = conn.Call("KademliaRPC.FindNode", req, res)
  	if err != nil {
		return nil, &CommandFailed{"Not implemented"}
  	}else {
  		fmt.Println("Find Node Completed")
  		k.PingChan <- contact
  		return res.Nodes, nil
  	}
}

func (k *Kademlia) DoFindValue(contact *Contact,
	searchKey ID) (value []byte, contacts []Contact, err error) {
	// TODO: Implement
	return nil, nil, &CommandFailed{"Not implemented"}
}

func (k *Kademlia) LocalFindValue(searchKey ID) ([]byte, error) {
	// TODO: Implement
	return []byte(""), &CommandFailed{"Not implemented"}
}

// For project 2!
func (k *Kademlia) DoIterativeFindNode(id ID) ([]Contact, error) {
	return nil, &CommandFailed{"Not implemented"}
}
func (k *Kademlia) DoIterativeStore(key ID, value []byte) ([]Contact, error) {
	return nil, &CommandFailed{"Not implemented"}
}
func (k *Kademlia) DoIterativeFindValue(key ID) (value []byte, err error) {
	return nil, &CommandFailed{"Not implemented"}
}

// For project 3!
func (k *Kademlia) Vanish(data []byte, numberKeys byte,
	threshold byte, timeoutSeconds int) (vdo VanashingDataObject) {
	return
}

func (k *Kademlia) Unvanish(searchKey ID) (data []byte) {
	return nil
}
