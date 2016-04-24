package libkademlia

// Contains the core kademlia type. In addition to core state, this type serves
// as a receiver for the RPC methods, which is required by that package.

import (
	 "fmt"
	// "log"
	// "net"
	// "net/http"
	// "net/rpc"
	// "strconv"
	// "container/list"
	// "math"
)


// Kademlia type. You can put whatever state you need in this.
type RoutingTable struct {
	buckets	[][]Contact
	//Self_contact	Contact
	addChan 	chan Contact
	delChan 	chan Contact

}

func NewRoutingTable(node Contact) *RoutingTable {
	tb := new(RoutingTable)
	tb.buckets = make([][]Contact, 160)
	fmt.Println("the bucket lenght is ", len(tb.buckets))
	// tb.addChan = make(chan *Contact)
	// tb.delChan = make(chan *Contact)
	return tb
}

// func (tb *Routing_Table) HandleChange() {
// 	for{
// 		select {
// 			//Add Contact
// 		case c := <-tb.addChan:
// 			tb.AddContact(c)
// 		case c := <-tb.delChan:
// 			tb.DelContact(c)
// 		}
// 	}
// }


// func (tb *Routing_Table) AddContact(c Contact) {
// 	dis := tb.NodeID.Xor(c.NodeID)
// 	numOfBucket := 159 - dis.PrefixLen()
// 	tb.buckets[numberOfBucket]
// 	addHelper(numberOfBucket, tb, c)
// }


// func addHelper(numOfBucket int, tb RoutingTable, c Contact) {
// 	containsC := false
// 	idx := 0
// 	for index, c1 := range tb.buckets[numOfBucket] {
//         if c1.NodeID == c.NodeID {
//             containsC = true
// 						idx = index
//         }
//   }
// 	if containsC {
// 		tb.buckets[numOfBucket] = tb.buckets[numOfBucket][:idx - 1] +
// 									tb.buckets[numOfBucket][idx + 1:] + tb.buckets[numOfBucket][idx]
// 	} else {
// 		if len(tb.buckets[numOfBucket]) < 20 {
// 			tb.buckets[numOfBucket].append(c)
// 		} else {
// 			// TODO:ping

// 		}
// 	}
// }

// func (tb Routing_Table)DelContact(c *Contact) {
// 	dis := tb.NodeID.Xor(c.NodeID)
// 	numOfBucket := 159 - dis.PrefixLen()
// 	tb.buckets[numberOfBucket]
// 	delHelper(numberOfBucket, tb, c)
// }


// func delHelper(numOfBucket int, tb RoutingTable, c Contact) {
// 	containsC := false
// 	idx := 0
// 	for index, c1 := range tb.buckets[numOfBucket] {
//         if c1.NodeID == c.NodeID {
//             containsC = true
// 			idx = index
//         }
//     }
// 	if containsC {
// 		tb.buckets[numOfBucket] = tb.buckets[numOfBucket][:,idx - 1] + tb.buckets[numOfBucket][idx + 1,:]
// 	}
// }
