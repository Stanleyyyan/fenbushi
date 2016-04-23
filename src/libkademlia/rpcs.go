package libkademlia

// Contains definitions mirroring the Kademlia spec. You will need to stick
// strictly to these to be compatible with the reference implementation and
// other groups' code.

import (
	"net"
)

type KademliaRPC struct {
	kademlia *Kademlia
}

// Host identification.
type Contact struct {
	NodeID ID
	Host   net.IP
	Port   uint16
}

///////////////////////////////////////////////////////////////////////////////
// PING
///////////////////////////////////////////////////////////////////////////////
type PingMessage struct {
	Sender Contact
	MsgID  ID
}

type PongMessage struct {
	MsgID  ID
	Sender Contact
}

func (k *KademliaRPC) Ping(ping PingMessage, pong *PongMessage) error {
	// TODO: Finish implementation
	pong.MsgID = CopyID(ping.MsgID)
	// Specify the sender
	pong.Sender = k.Kademlia.SelfContact
	// Update contact, etc
	c := PingMessage.SelfContact
	dis := ping.Sender.NodeID.Xor(k.Kademlia.NodeID)
	numOfBucket := 159-dis.PrefixLen()
	containsender := false
	idx := 0
	for index, c1 := range k.Kademlia.K_buckets.buckets[numOfBucket] {
				if c1.NodeID == c.NodeID {
						containsender = true
						idx = index
				}
	}
	if containsender {
		k.Kademlia.K_buckets.buckets[numOfBucket] = k.Kademlia.K_buckets.buckets[numOfBucket][:,idx - 1] +
									k.Kademlia.K_buckets.buckets[numOfBucket][idx + 1,:]
									+ k.Kademlia.K_buckets.buckets[numOfBucket][idx]
	} else {
		if len(k.Kademlia.K_buckets.buckets[numOfBucket]) < 20 {
			k.Kademlia.K_buckets.buckets[numOfBucket].append(c)
		} else {
			k.Kademlia.DoPing(k.Kademlia.K_buckets.buckets[0].Host, k.Kademlia.K_buckets.buckets[0].Port)

		}
	}
	return nil
}


///////////////////////////////////////////////////////////////////////////////
// STORE
///////////////////////////////////////////////////////////////////////////////
type StoreRequest struct {
	Sender Contact
	MsgID  ID
	Key    ID
	Value  []byte
}

type StoreResult struct {
	MsgID ID
	Err   error
}

func (k *KademliaRPC) Store(req StoreRequest, res *StoreResult) error {
	// TODO: Implement.
	return nil
}

///////////////////////////////////////////////////////////////////////////////
// FIND_NODE
///////////////////////////////////////////////////////////////////////////////
type FindNodeRequest struct {
	Sender Contact
	MsgID  ID
	NodeID ID
}

type FindNodeResult struct {
	MsgID ID
	Nodes []Contact
	Err   error
}

func (k *KademliaRPC) FindNode(req FindNodeRequest, res *FindNodeResult) error {
	// TODO: Implement.
	return nil
}

///////////////////////////////////////////////////////////////////////////////
// FIND_VALUE
///////////////////////////////////////////////////////////////////////////////
type FindValueRequest struct {
	Sender Contact
	MsgID  ID
	Key    ID
}

// If Value is nil, it should be ignored, and Nodes means the same as in a
// FindNodeResult.
type FindValueResult struct {
	MsgID ID
	Value []byte
	Nodes []Contact
	Err   error
}

func (k *KademliaRPC) FindValue(req FindValueRequest, res *FindValueResult) error {
	// TODO: Implement.
	return nil
}

// For Project 3

type GetVDORequest struct {
	Sender Contact
	VdoID  ID
	MsgID  ID
}

type GetVDOResult struct {
	MsgID ID
	VDO   VanashingDataObject
}

func (k *KademliaRPC) GetVDO(req GetVDORequest, res *GetVDOResult) error {
	// TODO: Implement.
	return nil
}
