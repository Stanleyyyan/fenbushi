package libkademlia

// Contains the core kademlia type. In addition to core state, this type serves
// as a receiver for the RPC methods, which is required by that package.

// Ounan Ma, omg049
// Chong Yan, cyu422
// Wenjie Zhang, wzm416

import (
	"bytes"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"strconv"
	"errors"
	"sort"
	"time"
	// "os"
)

const (
	alpha = 3
	b     = 8 * IDBytes
	// k     = 20
)

type Pair struct {
	Key 	ID
	Value 	[]byte
	MsgID	ID
}

type FindValueMsg struct {
	Value 		[]byte
	QueryNode	Contact
	Contacts 	[]Contact
	Err 		error
}

type FindNodeMsg struct {
	QueryNode	Contact
	Contacts 	[]Contact
	Err 		error
}

type CandidateCon struct {
	Con 		Contact
	Distance	ID
}

type ConAry []CandidateCon
// Kademlia type. You can put whatever state you need in this.
type Kademlia struct {
	NodeID      		ID
	SelfContact 		Contact
	K_buckets			RoutingTable
	PingChan			chan *UpdateMessage
	HashChan			chan Pair
	FindReqChan			chan FindNodeRequest
	FindResChan 		chan *FindNodeResult
	FindValueReqChan	chan FindValueRequest
	FindValueResChan 	chan *FindValueResult
	AckChan				chan AckMessage
	HTAckChan			chan AckMessage
	ConChan				chan FindNodeMsg
	ValChan 			chan *FindValueMsg
	H_Table				map[ID][]byte
	CandiateList		[]CandidateCon//20 - len(ShortList)
	VisitedCon			map[ID]bool
	ShortList			[]Contact
	VDOChan				chan VanashingDataObject
	VDOs 				[]VanashingDataObject
	// StoreChan			chan StoreMessage
}

//Sort Help Function
func (s ConAry) Len() int {
    return len(s)
}

func (s ConAry) Less(i, j int) bool {
	ret := s[i].Distance.Compare(s[j].Distance)
    return ret == -1
}

func (s ConAry) Swap(i, j int) {
    s[i], s[j] = s[j], s[i]
}


func NewKademliaWithId(laddr string, nodeID ID) *Kademlia {
	k := new(Kademlia)
	k.NodeID 			= nodeID
	k.K_buckets 		= RoutingTable{buckets: make([][]Contact, 160)}
	k.PingChan 			= make(chan *UpdateMessage)
	k.HashChan 			= make(chan Pair)
	k.FindReqChan 		= make(chan FindNodeRequest)
	k.FindResChan		= make(chan *FindNodeResult)
	k.FindValueReqChan 	= make(chan FindValueRequest)
	k.FindValueResChan	= make(chan *FindValueResult)
	k.AckChan			= make(chan AckMessage)
	k.HTAckChan			= make(chan AckMessage)
	k.ConChan			= make(chan FindNodeMsg)
	k.ValChan			= make(chan *FindValueMsg)
	k.H_Table 			= make(map 	[ID][]byte)
	k.ShortList			= []Contact{}
	k.CandiateList		= []CandidateCon{}
	k.VisitedCon		= make(map 	[ID]bool)
	k.VDOChan 			= make(chan VanashingDataObject)
	// k.StoreChan			= make(chan StoreMessage)
	go k.Handler()
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
	// k.DoPing(host, uint16(port_int))
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
	if k.SelfContact.NodeID.Equals(nodeId) {
		return &k.SelfContact, nil
	}
	dis := nodeId.Xor(k.NodeID)
	numOfBucket := 159 - dis.PrefixLen()
	fmt.Println("num of Bucket is :" ,numOfBucket)
	fmt.Println(len(k.K_buckets.buckets))
	for _, c1 := range k.K_buckets.buckets[numOfBucket] {
		if c1.NodeID.Equals(nodeId) {
			fmt.Println("target node Found!")
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

func (k *Kademlia) PingAliveTest(host net.IP, port uint16) (*Contact, error) {
	portnum := strconv.Itoa(int(port))
	_, err := rpc.DialHTTPPath("tcp", host.String() + ":" + portnum,
		rpc.DefaultRPCPath + portnum)
	if err != nil {
		fmt.Println("error!")
		// log.Fatal("dialing:", err)
		return nil, err
	} else {
		return nil, nil
	}
}

func (k *Kademlia) DoPing(host net.IP, port uint16) (*Contact, error) {
	// TODO: Implement
	portnum := strconv.Itoa(int(port))
	client, err := rpc.DialHTTPPath("tcp", host.String() + ":" + portnum,
		rpc.DefaultRPCPath + portnum)

	if err != nil {
		fmt.Println("error!")
		// log.Fatal("dialing:", err)
		return nil, err
	}
	pim := PingMessage{k.SelfContact, NewRandomID()}
	fmt.Println("Ping Self port is: ", k.SelfContact.Port)
	fmt.Println("Ping Target port is: ", port)
	pom := new(PongMessage)
  	err = client.Call("KademliaRPC.Ping", pim, pom)

	updateMessage := new(UpdateMessage)
	updateMessage.MsgID = pim.MsgID
	updateMessage.NewContact = pom.Sender
	if err != nil {
		return nil, &CommandFailed{
			"Unable to ping " + fmt.Sprintf("%s:%v", host.String(), port)}
	} else {
		k.PingChan <- updateMessage
		flag := true
		for flag {
			select {
			case ack := <- k.AckChan:
				if ack.MsgID.Equals(pim.MsgID) {
					flag = false
				}else {
					k.AckChan <- ack
				}
			}
		}
		return &(pom.Sender), nil
	}
}

func (k *Kademlia) UpdateRT(c *UpdateMessage) error {
	ack := AckMessage{MsgID: c.MsgID}
  	if c.NewContact.NodeID.Equals(k.NodeID) {
		k.AckChan <- ack
		return nil
	}
	dis := c.NewContact.NodeID.Xor(k.NodeID)
	numOfBucket := 159 - dis.PrefixLen()
	containSender := false
	idx := 0
	for index, c1 := range k.K_buckets.buckets[numOfBucket] {
		if c1.NodeID.Equals(c.NewContact.NodeID) {
			containSender = true
			idx = index
		}
	}
	if containSender {
		if len(k.K_buckets.buckets[numOfBucket]) == 1 {
			k.AckChan <- ack
			return nil
		}
		if idx != 0 {
			temp := k.K_buckets.buckets[numOfBucket][idx]
			k.K_buckets.buckets[numOfBucket] =
					append(k.K_buckets.buckets[numOfBucket][:idx - 1],
							k.K_buckets.buckets[numOfBucket][idx + 1:]...)
			k.K_buckets.buckets[numOfBucket] = append(k.K_buckets.buckets[numOfBucket], temp)
		}else {
			k.K_buckets.buckets[numOfBucket] = append(k.K_buckets.buckets[numOfBucket][idx + 1:], k.K_buckets.buckets[numOfBucket][idx])
		}
		k.AckChan <- ack
		return errors.New("Move to tail")
	} else {
		if len(k.K_buckets.buckets[numOfBucket]) < 20 {
			k.K_buckets.buckets[numOfBucket] = append(k.K_buckets.buckets[numOfBucket], (c.NewContact))
			k.AckChan <- ack
			return errors.New("k buckets not full, add to tail")
		} else {
			fmt.Println("WTF!!!")
			_, err :=
				k.PingAliveTest(k.K_buckets.buckets[numOfBucket][0].Host, k.K_buckets.buckets[numOfBucket][0].Port)
			if err != nil {
				k.K_buckets.buckets[numOfBucket] =
					k.K_buckets.buckets[numOfBucket][1:]
				k.K_buckets.buckets[numOfBucket] = append(k.K_buckets.buckets[numOfBucket], (c.NewContact))
				k.AckChan <- ack
				fmt.Println("Head Dead")
				return errors.New("head dead, replace head")
			}else{
				k.AckChan <- ack
				fmt.Println("Head Live")
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
		// log.Fatal("dialing:", err)
		return err
	}
	req := StoreRequest{Sender: k.SelfContact, MsgID: NewRandomID(), Key: key, Value: value}
	res := new(StoreResult)
  	err = conn.Call("KademliaRPC.Store", req, res)

	updateMessage := new(UpdateMessage)
	updateMessage.MsgID = req.MsgID
	updateMessage.NewContact = *contact

	// storeMessage := new(storeMessage)
  	if err != nil {
  		fmt.Println(err)
		log.Fatal("RPC:", err)
		// storeMessage.Err = res.Err
		// storeMessage.NodeID = req.Key
		// k.StoreChan <- storeMessage
		return &CommandFailed{"DoStore Failed"}
  	} else {
  		k.PingChan <- updateMessage
		flag := true
		for flag {
			select {
			case ack := <- k.AckChan:
				if ack.MsgID.Equals(req.MsgID){
					flag = false
				}else {
					k.AckChan <- ack
				}
			}
		}
		// storeMessage.Err = res.Err
		// storeMessage.NodeID = req.Key
		// k.StoreChan <- storeMessage
  		fmt.Println("DoStore complete!")
  	}
  	return nil
}

func (k *Kademlia) UpdateHT(key ID, value []byte, MsgID ID) AckMessage {
	k.H_Table[key] = value
	ack := AckMessage{MsgID: MsgID}
	return ack
}

func (k *Kademlia) Handler(){
	for {
		select {
		case newPair := <- k.HashChan:
			k.HTAckChan <- k.UpdateHT(newPair.Key, newPair.Value, newPair.MsgID)
		case newContact := <- k.PingChan:
			k.UpdateRT(newContact)
		case newFindNodeReq := <- k.FindReqChan:
			k.FindResChan <- k.GetNode(newFindNodeReq)
		case newFindValueReq := <- k.FindValueReqChan:
			k.FindValueResChan <- k.GetValue(newFindValueReq)
		}
	}
}

func (k *Kademlia) DoFindNode(contact *Contact, searchKey ID) ([]Contact, error) {
	// TODO: Implement
	portnum := strconv.Itoa(int(contact.Port))
	temp := contact.Host.String() + ":" + portnum
	fmt.Println("DoFind:", temp)
	//
	conn, err := rpc.DialHTTPPath("tcp", contact.Host.String() + ":" + portnum,
		rpc.DefaultRPCPath + portnum)
	if err != nil {
		fmt.Println("error!")
		// log.Fatal("dialing:", err)
		return nil, err
	}
	req := FindNodeRequest{Sender: k.SelfContact, MsgID: NewRandomID(), NodeID: searchKey}
	res := new(FindNodeResult)
  	err = conn.Call("KademliaRPC.FindNode", req, res)

	updateMessage := new(UpdateMessage)
	updateMessage.MsgID = req.MsgID
	updateMessage.NewContact = *contact
  	if err != nil {
		return nil, &CommandFailed{"Not implemented"}
  	}else {
  		k.PingChan <- updateMessage
			flag := true
			for flag {
				select {
				case ack := <- k.AckChan:
					if ack.MsgID.Equals(req.MsgID){
						flag = false
					}else {
						k.AckChan <- ack
					}
				}
		}
		fmt.Println("Find Node Completed")
		for i := 0; i < len(res.Nodes); i++ {
			fmt.Println(res.Nodes[i].Port)
			updateMessage := new(UpdateMessage)
			updateMessage.MsgID = NewRandomID()
			updateMessage.NewContact = res.Nodes[i]
			k.PingChan <- updateMessage
			flag := true
			for flag {
				select {
				case ack := <- k.AckChan:
					if ack.MsgID.Equals(updateMessage.MsgID){
						flag = false
					}else {
						k.AckChan <- ack
					}
				}
			fmt.Println(res.Nodes[i].Port)
			}
		}
  		return res.Nodes, nil
  	}
}

func (k *Kademlia) GetNode(req FindNodeRequest) *FindNodeResult {
	// TODO: Implement.
	res := new(FindNodeResult)
	res.MsgID = CopyID(req.MsgID)

	dis := req.NodeID.Xor(k.NodeID)
	bucketIdx := 159 - dis.PrefixLen()
	fmt.Println("distance is ", bucketIdx)
	for _, c1 := range k.K_buckets.buckets[bucketIdx] {
		fmt.Println("current nodeID ", c1.NodeID.AsString())
		if !c1.NodeID.Equals(req.Sender.NodeID){
			fmt.Println("append nodeID ", c1.NodeID.AsString())
			res.Nodes = append(res.Nodes, c1)
		}
	}
	for i := bucketIdx - 1; len(res.Nodes) < 20 && i >= 0; i-- {
		for _, c1 := range k.K_buckets.buckets[i] {
			if !c1.NodeID.Equals(req.Sender.NodeID){
				fmt.Println("append nodeID ", c1.NodeID.AsString())
				res.Nodes = append(res.Nodes, c1)
			}
			if len(res.Nodes) == 20 {
				fmt.Println("K is 20")
				return res
			}
		}
	}
	for i := bucketIdx + 1; len(res.Nodes) < 20 && i < 160; i++ {
		for _, c1 := range k.K_buckets.buckets[i] {
			if !c1.NodeID.Equals(req.Sender.NodeID){
				fmt.Println("append nodeID ", c1.NodeID.AsString())
				res.Nodes = append(res.Nodes, c1)
			}
			if len(res.Nodes) == 20 {
				fmt.Println("K is 20")
				return res;
			}
		}
	}
	fmt.Println("not matched")
	return res
}


func (k *Kademlia) DoFindValue(contact *Contact,
	searchKey ID) (value []byte, contacts []Contact, err error) {
	// TODO: Implement
	portnum := strconv.Itoa(int(contact.Port))
	temp := contact.Host.String() + ":" + portnum
	fmt.Println("DoFindValue:", temp)
	//
	conn, err := rpc.DialHTTPPath("tcp", contact.Host.String() + ":" + portnum,
		rpc.DefaultRPCPath + portnum)

	if err != nil {
		fmt.Println("error!")
		// log.Fatal("dialing:", err)
		return nil, nil, err
	}

	req := FindValueRequest{Sender: k.SelfContact, MsgID: NewRandomID(), Key: searchKey}
	res := new(FindValueResult)
  	err = conn.Call("KademliaRPC.FindValue", req, res)

	updateMessage := new(UpdateMessage)
	updateMessage.MsgID = req.MsgID
	updateMessage.NewContact = *contact
  	if err != nil {
		return nil, nil, &CommandFailed{"Not implemented"}
  	}else {
  		k.PingChan <- updateMessage
		flag := true
		for flag {
			select {
			case ack := <- k.AckChan:
				if ack.MsgID.Equals(req.MsgID){
					flag = false
				}else {
					k.AckChan <- ack
				}
			}
		}
		for i := 0; i < len(res.Nodes); i++ {
			fmt.Println(res.Nodes[i].Port)
			updateMessage := new(UpdateMessage)
			updateMessage.MsgID = NewRandomID()
			updateMessage.NewContact = res.Nodes[i]
			k.PingChan <- updateMessage
			flag := true
			for flag {
				select {
				case ack := <- k.AckChan:
					if ack.MsgID.Equals(updateMessage.MsgID){
						flag = false
					}else {
						k.AckChan <- ack
					}
				}
			fmt.Println(res.Nodes[i].Port)
			}
		}
  		return res.Value, res.Nodes, nil
  	}
	return nil, nil, &CommandFailed{"Not Found"}
}

func (k *Kademlia) GetValue(req FindValueRequest) *FindValueResult{
	res := new(FindValueResult)
	res.MsgID = CopyID(req.MsgID)
	value, err := k.LocalFindValue(req.Key)
	if err != nil {
		fmt.Println("can not find value")
		tempReq := FindNodeRequest{Sender: req.Sender, MsgID: req.MsgID, NodeID: req.Key}
		FindNodeResult := k.GetNode(tempReq)
		res.Nodes = FindNodeResult.Nodes
	} else {
		fmt.Println("find value!!")
		res.Value = value
	}
	return res
}

func (k *Kademlia) LocalFindValue(searchKey ID) ([]byte, error) {
	// TODO: Implement
	for key, value := range k.H_Table{
		if key.Equals(searchKey) {
			return value, nil
		}
	}
	return nil, &CommandFailed{"Not Find"}
}

// For project 2!
func (k *Kademlia) DoIterativeFindNode(id ID) ([]Contact, error) {
	k.ShortList = []Contact{}
	k.CandiateList = []CandidateCon{}
	k.VisitedCon = map[ID]bool{}
	dis := k.NodeID.Xor(id)
	bucketIdx := 159 - dis.PrefixLen()
	fmt.Println("distance is ", bucketIdx)
	for i := 0; len(k.CandiateList) < 20 && i < len(k.K_buckets.buckets[bucketIdx]); i++ {
		var element = CandidateCon{Con: k.K_buckets.buckets[bucketIdx][i], Distance: id.Xor(k.K_buckets.buckets[bucketIdx][i].NodeID)}
		k.CandiateList = append(k.CandiateList, element)
	}

	for i := bucketIdx - 1; len(k.CandiateList) < 20 && i >= 0; i-- {
		for j := 0; len(k.CandiateList) < 20 && j < len(k.K_buckets.buckets[i]); j++ {
			var element = CandidateCon{Con: k.K_buckets.buckets[i][j], Distance: id.Xor(k.K_buckets.buckets[i][j].NodeID)}
			k.CandiateList = append(k.CandiateList, element)
		}
	}

	for i := bucketIdx + 1; len(k.CandiateList) < 20 && i < 160; i++ {
		for j := 0; len(k.CandiateList) < 20 && j < len(k.K_buckets.buckets[i]); j++ {
			var element = CandidateCon{Con: k.K_buckets.buckets[i][j], Distance: id.Xor(k.K_buckets.buckets[i][j].NodeID)}
			k.CandiateList = append(k.CandiateList, element)
		}
	}

	// start cycles of three parallel DoFindNode()
	cycle := map[ID]bool {
		NewRandomID(): true,
		NewRandomID(): true,
		NewRandomID(): true,
	}
	start := time.Now()
	start = start.Add(300000000)
	timer := time.Now()
	timer = timer.Add(3000000000)
	terminator := false		// the flag of ending condition
	for len(k.ShortList) < 20 {
		s := true
		for _, i := range cycle {
			s = s && i
		}
		curTime := time.Now()
		timeOut := curTime.After(start)
		STOP := curTime.After(timer)
		if timeOut || s && len(k.CandiateList) == 0  {
			break
		}
		if STOP {
			break
		}
		if s || timeOut{
			conList := []Contact{}
			cycle = map[ID]bool{}
			for i := 0; i < 3 && i < len(k.CandiateList); i++ {
				conList = append(conList, k.CandiateList[i].Con)
				k.VisitedCon[conList[i].NodeID] = true
				cycle[conList[i].NodeID] = false
			}
			k.CandiateList = k.CandiateList[len(conList):]
			for i := 0; i < len(conList); i++ {
				go k.FindNodeHandler(&conList[i], id)
			}
			start = time.Now()
			start = start.Add(300000000)
		}
		// receive the results
		ret := <- k.ConChan
		_, ok := cycle[ret.QueryNode.NodeID]
		if !ok {
			continue
		} else {
			terminator = !k.RcvNodeHandler(ret, id, terminator, /*ok*/)
			cycle[ret.QueryNode.NodeID] = true
		}

	}

	return k.ShortList, nil
}

// handle the parallel DoFindNode() and receive results
func (k *Kademlia) FindNodeHandler(contact *Contact, searchKey ID) {
	contacts, err := k.DoFindNode(contact, searchKey)
	k.ConChan <- FindNodeMsg{QueryNode: *contact, Contacts: contacts, Err: err}
}

// update the ShortList and CandiateList
func (k *Kademlia) RcvNodeHandler(ret FindNodeMsg, id ID, terminator bool) (res bool) {
	if ret.Err == nil {
		k.ShortList = append(k.ShortList, ret.QueryNode)
		if terminator == true {
			return false
		}
		for _, c := range ret.Contacts {
			if k.VisitedCon[c.NodeID] == true {
				continue
			}
			contained := false
			for i := 0; i < len(k.CandiateList); i++ {
				if k.CandiateList[i].Con.NodeID.Equals(c.NodeID) {
					contained = true;
					break;
				}
			}
			if contained == false {
				k.CandiateList = append(k.CandiateList, CandidateCon{Con:c, Distance: id.Xor(c.NodeID)})
			}
		}
		sort.Sort(ConAry(k.CandiateList))
		idx := len(k.CandiateList)
		if idx > 20 - len(k.ShortList){
			idx = 20 - len(k.ShortList)
		}

		fmt.Println("k ShortList Len is :", len(k.ShortList))
		fmt.Println("k CandiateList Len is :", len(k.CandiateList))
		fmt.Println("idx is :", idx)
		k.CandiateList = k.CandiateList[:idx]
		if len(k.CandiateList) == 0{
			return true
		}
		if len(ret.Contacts) == 0 {
			return true
		}
		return id.Xor(ret.Contacts[0].NodeID).Compare(k.CandiateList[len(k.CandiateList) - 1].Distance) < 1
	}
	return !(false || terminator)
}


func (k *Kademlia) DoIterativeStore(key ID, value []byte) ([]Contact, error) {
	contacts, _:= k.DoIterativeFindNode(key)
	res := []Contact{}
	for _, c := range contacts {
		err := k.DoStore(&c, key, value)
		if err == nil {
			res = append(res, c)
		}
	}
	return res, nil
}

// similar to DoIterativeFindNode
func (k *Kademlia) DoIterativeFindValue(key ID) (value []byte, err error) {
	k.ShortList = []Contact{}
	k.CandiateList = []CandidateCon{}
	k.VisitedCon = map[ID]bool{}
	dis := k.NodeID.Xor(key)
	bucketIdx := 159 - dis.PrefixLen()
	fmt.Println("distance is ", bucketIdx)
	for i := 0; len(k.CandiateList) < 20 && i < len(k.K_buckets.buckets[bucketIdx]); i++ {
		var element = CandidateCon{Con: k.K_buckets.buckets[bucketIdx][i], Distance: key.Xor(k.K_buckets.buckets[bucketIdx][i].NodeID)}
		k.CandiateList = append(k.CandiateList, element)
	}

	for i := bucketIdx - 1; len(k.CandiateList) < 20 && i >= 0; i-- {
		for j := 0; len(k.CandiateList) < 20 && j < len(k.K_buckets.buckets[i]); j++ {
			var element = CandidateCon{Con: k.K_buckets.buckets[i][j], Distance: key.Xor(k.K_buckets.buckets[i][j].NodeID)}
			k.CandiateList = append(k.CandiateList, element)
		}
	}

	for i := bucketIdx + 1; len(k.CandiateList) < 20 && i < 160; i++ {
		for j := 0; len(k.CandiateList) < 20 && j < len(k.K_buckets.buckets[i]); j++ {
			var element = CandidateCon{Con: k.K_buckets.buckets[i][j], Distance: key.Xor(k.K_buckets.buckets[i][j].NodeID)}
			k.CandiateList = append(k.CandiateList, element)
		}
	}

	cycle := map[ID]bool {
		NewRandomID(): true,
		NewRandomID(): true,
		NewRandomID(): true,
	}
	start := time.Now()
	start = start.Add(300000000)
	terminator := false
	ret := new(FindValueMsg)
	foundValue := []byte("")
	for len(k.ShortList) < 20 {
		s := true
		// restCon := false
		for _, i := range cycle {
			s = s && i
		}
		curTime := time.Now()
		timeOut := curTime.After(start)
		if terminator && (s || timeOut) {
			break
		}
		if s || timeOut{
			conList := []Contact{}
			cycle = map[ID]bool{}
			for i := 0; i < 3 && i < len(k.CandiateList); i++ {
				conList = append(conList, k.CandiateList[i].Con)
				k.VisitedCon[conList[i].NodeID] = true
				cycle[conList[i].NodeID] = false
			}
			k.CandiateList = k.CandiateList[len(conList):]
			for i := 0; i < len(conList); i++ {
				go k.FindValueHandler(&conList[i], key)
			}
			start = time.Now()
			start = start.Add(300000000)
		}
		ret = <- k.ValChan
		if ret.Value != nil {
			foundValue = ret.Value
			terminator = true
		}
		_, ok := cycle[ret.QueryNode.NodeID]
		if !ok {
			continue
		} else {
			terminator =  !k.RcvValueHandler(*ret, key, terminator, /*ok*/) || terminator
			cycle[ret.QueryNode.NodeID] = true
		}
	}

	closest := -1
	closestDis := NewRandomID()
	for idx, c := range k.ShortList {
		if k.ShortList[idx].NodeID != key {
			closest = idx
			closestDis = c.NodeID.Xor(key)
		}
	}
	for idx, c := range k.ShortList {
		dis := c.NodeID.Xor(key)
		if c.NodeID != key && dis.Compare(closestDis) < 0{
			closest = idx
			closestDis = dis
		}
	}

	// not find value
	if bytes.Equal(foundValue, []byte("")){
		foundValue = nil
	}
	if foundValue == nil {
		return nil, &CommandFailed{k.ShortList[closest].NodeID.AsString()}
	} else {
		k.DoStore(&k.ShortList[closest], key, foundValue)
		return foundValue, nil
	}

}

// handle the parallel DoFindValue() and receive results
func (k *Kademlia) FindValueHandler(contact *Contact, searchKey ID) {
	value, contacts, err := k.DoFindValue(contact, searchKey)
	k.ValChan <- &FindValueMsg{Value: value, QueryNode: *contact, Contacts: contacts, Err: err}// pointer? argument?
}

func (k *Kademlia) RcvValueHandler(ret FindValueMsg, id ID, terminator bool, /*ok bool*/) (res bool) {
	if ret.Err == nil {
		k.ShortList = append(k.ShortList, ret.QueryNode)
		if terminator == true {
			return false
		}
		for _, c := range ret.Contacts {
			if k.VisitedCon[c.NodeID] == true {
				continue
			}
			contained := false
			for i := 0; i < len(k.CandiateList); i++ {
				if k.CandiateList[i].Con.NodeID.Equals(c.NodeID) {
					contained = true;
					break;
				}
			}
			if contained == false {
				k.CandiateList = append(k.CandiateList, CandidateCon{Con:c, Distance: id.Xor(c.NodeID)})
			}
		}
		sort.Sort(ConAry(k.CandiateList))
		idx := len(k.CandiateList)
		if idx > 20 - len(k.ShortList){
			idx = 20 - len(k.ShortList)
		}

		fmt.Println("k ShortList Len is :", len(k.ShortList))
		fmt.Println("k CandiateList Len is :", len(k.CandiateList))
		fmt.Println("idx is :", idx)
		k.CandiateList = k.CandiateList[:idx]
		// if !ok {
		// 	return !(false || terminator)
		// }
		if len(k.CandiateList) == 0{
			return false
		}
		if len(ret.Contacts) == 0 {
			return false
		}
		return id.Xor(ret.Contacts[0].NodeID).Compare(k.CandiateList[len(k.CandiateList) - 1].Distance) < 1
	}
	return !(false || terminator)
	return !(false || terminator)
}

// For project 3!
func (k *Kademlia) Vanish(data []byte, numberKeys byte,
	threshold byte, timeoutSeconds int) (vdo VanashingDataObject) {
	return
}

func (k *Kademlia) Unvanish(searchKey ID) (data []byte) {
	return nil
}
