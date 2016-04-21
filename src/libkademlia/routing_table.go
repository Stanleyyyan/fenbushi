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
	"container/list"
	"math"
)

const (
	alpha = 3
	b     = 8 * IDBytes
	k     = 20
)

// Kademlia type. You can put whatever state you need in this.
type RoutingTable struct {
	var buckets	[][]Contact
	//Self_contact	Contact
	addChan 	chan Contact
	delChan 	chan Contact

}

func NewRoutingTable(node Contact) *Routing_Table {
	tb := new(Routing_Table)
	tb.buckets = make([][]Contact, 160)
	tb.addChan = make(chan *Contact)
	tb.delChan = make(chan *Contact)
}

func (tb Routing_Table)AddContact() {
	for{
		select{
		case c := <-tb.addChan:
			dis := tb.NodeID.Xor(c.NodeID)
			numOfBucket := 159 - dis.PrefixLen()
			tb.buckets[numberOfBucket]
			addHelper(numberOfBucket, tb, c)
		}
	}
}

func addHelper(numOfBucket int, tb RoutingTable, c Contact) {
	containsC := false
	idx := 0
	for index, c1 := range tb.buckets[numOfBucket] {
        if c1.NodeID == c.NodeID {
            containsC = true
			idx = index
        }
    }
	if containsC {
		tb.buckets[numOfBucket] = tb.buckets[numOfBucket][:,idx - 1] + tb.buckets[numOfBucket][idx + 1,:] + tb.buckets[numOfBucket][idx]
	} else {
		if len(tb.buckets[numOfBucket]) < 20 {
			tb.buckets[numOfBucket].append(c)
		} else {
			// TODO:ping
		}
	}
}

func (tb Routing_Table)DelContact() {
	for{
		select{
		case c := <-tb.delChan:
			dis := tb.NodeID.Xor(c.NodeID)
			numOfBucket := 159 - dis.PrefixLen()
			tb.buckets[numberOfBucket]
			delHelper(numberOfBucket, tb, c)
		}
	}
}

func delHelper(numOfBucket int, tb RoutingTable, c Contact) {
	containsC := false
	idx := 0
	for index, c1 := range tb.buckets[numOfBucket] {
        if c1.NodeID == c.NodeID {
            containsC = true
			idx = index
        }
    }
	if containsC {
		tb.buckets[numOfBucket] = tb.buckets[numOfBucket][:,idx - 1] + tb.buckets[numOfBucket][idx + 1,:]
	} 
}
