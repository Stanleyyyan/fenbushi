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

// Kademlia type. You can put whatever state you need in this.
type Kademlia struct {
	NodeID      ID
	SelfContact Contact
	K_buckets	RoutingTable
	PingChan	chan *Contact
	H_Table		map[ID][]byte
}

func NewKademliaWithId(laddr string, nodeID ID) *Kademlia {
	k := new(Kademlia)
	k.NodeID = nodeID
	k.K_buckets = RoutingTable{buckets: make([][]Contact, 160)}
	k.PingChan = make(chan *Contact)
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
	if nodeId == k.SelfContact.NodeID {
		return &k.SelfContact, nil
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
		// k.Update(&(pom.Sender))
		fmt.Println("updating Done")
		return &(pom.Sender), nil
	}
}

func (k *Kademlia) Update(c *Contact) error {
	fmt.Println("NodeID: ", c.NodeID)
	dis := c.NodeID.Xor(k.NodeID)
	numOfBucket := 159 - dis.PrefixLen()
	containSender := false
	idx := 0
	fmt.Println("num of Bucket is :" ,numOfBucket)
	fmt.Println(len(k.K_buckets.buckets))
	for index, c1 := range k.K_buckets.buckets[numOfBucket] {
		if c1.NodeID == c.NodeID {
			containSender = true
			idx = index
		}
	}
	if containSender {
		if len(k.K_buckets.buckets[numOfBucket]) == 1 {
			fmt.Println("size is 1, add done")
			return nil
		}
		fmt.Println("index is :", idx)
		temp := k.K_buckets.buckets[numOfBucket][idx]
		k.K_buckets.buckets[numOfBucket] = 
				append(k.K_buckets.buckets[numOfBucket][:idx - 1],
						k.K_buckets.buckets[numOfBucket][idx + 1:]...)
		k.K_buckets.buckets[numOfBucket] = append(k.K_buckets.buckets[numOfBucket], temp)
				fmt.Println("Moved to Tail!")
				return errors.New("Move to tail")
	} else {
		if len(k.K_buckets.buckets[numOfBucket]) < 20 {
			k.K_buckets.buckets[numOfBucket] = append(k.K_buckets.buckets[numOfBucket], *c)
			fmt.Println("k buckets not full, add to tail")
			return errors.New("k buckets not full, add to tail")
		} else {
			_, err := 
				k.DoPing(k.K_buckets.buckets[numOfBucket][0].Host, k.K_buckets.buckets[numOfBucket][0].Port)
			if err != nil {
				k.K_buckets.buckets[numOfBucket] = 
					k.K_buckets.buckets[numOfBucket][1:]
				k.K_buckets.buckets[numOfBucket] = append(k.K_buckets.buckets[numOfBucket], *c)
				fmt.Println("head dead, replace head")
				return errors.New("head dead, replace head")
			}else{
				fmt.Println("Discard!")
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
			k.Update(newContact)
		}
	}
}

func (k *Kademlia) DoStore(contact *Contact, key ID, value []byte) error {
	// // TODO: Implement
	// portnum := strconv.Itoa(int(contact.Port))
	// temp := contact.Host.String() + ":" + portnum
	// fmt.Println("DoStore:", temp)
	// //
	
	// conn, err := rpc.DialHTTPPath("tcp", host.String() + ":" + portnum,
	// 	rpc.DefaultRPCPath + portnum)

	// if err != nil {
	// 	fmt.Println("error!")
	// 	log.Fatal("dialing:", err)
	// }

	return &CommandFailed{"Not implemented"}
}

func (k *Kademlia) DoFindNode(contact *Contact, searchKey ID) ([]Contact, error) {
	// TODO: Implement
	return nil, &CommandFailed{"Not implemented"}
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
