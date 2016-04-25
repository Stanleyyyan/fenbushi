package libkademlia

// Contains definitions mirroring the Kademlia spec. You will need to stick
// strictly to these to be compatible with the reference implementation and
// other groups' code.

import (
	"fmt"
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

type UpdateMessage struct {
	MsgID		ID
	NewContact	Contact
}

type AckMessage struct {
	MsgID		ID
}

func (k *KademliaRPC) Ping(ping PingMessage, pong *PongMessage) error {
	// TODO: Finish implementation
	fmt.Println("Ping !!")
	pong.MsgID = CopyID(ping.MsgID)
	// Specify the sender
	pong.Sender = k.kademlia.SelfContact
	//chan

	updateMessage := new(UpdateMessage)
	updateMessage.MsgID = ping.MsgID
	updateMessage.NewContact = ping.Sender

	k.kademlia.PingChan <- updateMessage
	flag := true
	for flag {
		select {
		case ack := <- k.kademlia.AckChan:
			if ack.MsgID.Equals(ping.MsgID) {
				flag = false
			}else {
				k.kademlia.AckChan <- ack
			}
		}
	}
	fmt.Println("Pong")
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
	fmt.Println("Store!!")
	res.MsgID = CopyID(req.MsgID)
	res.Err = nil
	k.kademlia.HashChan <- Pair{Key: req.Key, Value: req.Value, MsgID: req.MsgID}
	flag := true
	for flag {
		select {
		case ack := <- k.kademlia.HTAckChan:
			if ack.MsgID.Equals(req.MsgID) {
				flag = false
			}else {
				k.kademlia.HTAckChan <- ack
			}
		}
	}

	updateMessage := new(UpdateMessage)
	updateMessage.MsgID = req.MsgID
	updateMessage.NewContact = req.Sender
	k.kademlia.PingChan <- updateMessage
	flag = true
	for flag {
		select {
		case ack := <- k.kademlia.AckChan:
			if ack.MsgID.Equals(req.MsgID) {
				flag = false
			}else {
				k.kademlia.AckChan <- ack
			}
		}
	}
	fmt.Println("Store Done")
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
	fmt.Println("FindNode")
	res.MsgID = CopyID(req.MsgID)
	dis := req.NodeID.Xor(k.kademlia.NodeID)
	fmt.Println("res MsgID is ",res.MsgID.AsString())
	bucketIdx := 159 - dis.PrefixLen()
	// containSender := false
	// counter := 0
	// idx := 0
	fmt.Println("num of Bucket is :" ,bucketIdx)
	fmt.Println(len(k.kademlia.K_buckets.buckets))
	k.kademlia.FindReqChan <- req
	flag := true
	for flag {
		select {
		case ret := <- k.kademlia.FindResChan:
			if ret.MsgID.Equals(res.MsgID){
				res.Nodes = ret.Nodes
				flag = false
			}else {
				k.kademlia.FindResChan <- ret
			}
		}
	}

	updateMessage := new(UpdateMessage)
	updateMessage.MsgID = req.MsgID
	updateMessage.NewContact = req.Sender

	k.kademlia.PingChan <- updateMessage
	flag = true
	for flag {
		select {
		case ack := <- k.kademlia.AckChan:
			if ack.MsgID.Equals(req.MsgID) {
				flag = false
			}else {
				k.kademlia.AckChan <- ack
			}
		}
	}
	fmt.Println("FindNode Done")
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
	// TODO: Implement
	fmt.Println("Enter FINDVALUE")
	k.kademlia.FindValueReqChan <- req
	fmt.Println("pass req to chan")
	res.MsgID = CopyID(req.MsgID)
	// fmt.Println("ret MsgID is ", ret.MsgID.AsString())
	flag := true
	for flag {
		fmt.Println("wait for res")
		select {
		case ret := <- k.kademlia.FindValueResChan:
			fmt.Println("ret MsgID is ", ret.MsgID.AsString())
			if ret.MsgID.Equals(res.MsgID){
				res.Nodes = ret.Nodes
				res.Value = ret.Value
				flag = false
			}else {
				k.kademlia.FindValueResChan <- ret
			}
		}
	}

	fmt.Println("K is ", len(res.Nodes))
	updateMessage := new(UpdateMessage)
	updateMessage.MsgID = req.MsgID
	updateMessage.NewContact = req.Sender

	k.kademlia.PingChan <- updateMessage
	flag = true
	for flag {
		select {
		case ack := <- k.kademlia.AckChan:
			if ack.MsgID.Equals(req.MsgID) {
				flag = false
			}else {
				k.kademlia.AckChan <- ack
			}
		}
	}
	fmt.Println("FindNode Done")
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
