package libkademlia

import (
	"bytes"
	"net"
	"strconv"
	"testing"
	"time"
	"fmt"
)

func StringToIpPort(laddr string) (ip net.IP, port uint16, err error) {
	hostString, portString, err := net.SplitHostPort(laddr)
	if err != nil {
		return
	}
	ipStr, err := net.LookupHost(hostString)
	if err != nil {
		return
	}
	for i := 0; i < len(ipStr); i++ {
		ip = net.ParseIP(ipStr[i])
		if ip.To4() != nil {
			break
		}
	}
	portInt, err := strconv.Atoi(portString)
	port = uint16(portInt)
	return
}

// func TestPing(t *testing.T) {
// 	instance1 := NewKademlia("localhost:7890")
// 	instance2 := NewKademlia("localhost:7891")
// 	host2, port2, _ := StringToIpPort("localhost:7891")
// 	contact2, err := instance2.FindContact(instance2.NodeID)
// 	if err != nil {
// 		t.Error("A node cannot find itself's contact info")
// 	}
// 	contact2, err = instance2.FindContact(instance1.NodeID)
// 	if err == nil {
// 		t.Error("Instance 2 should not be able to find instance " +
// 			"1 in its buckets before ping instance 1")
// 	}
// 	instance1.DoPing(host2, port2)
//  //  	duration := time.Duration(1)*time.Second
//  //  	time.Sleep(duration)
// 	contact2, err = instance1.FindContact(instance2.NodeID)
// 	if err != nil {
// 		t.Error("Instance 2's contact not found in Instance 1's contact list")
// 		return
// 	}
// 	wrong_ID := NewRandomID()
// 	_, err = instance2.FindContact(wrong_ID)
// 	if err == nil {
// 		t.Error("Instance 2 should not be able to find a node with the wrong ID")
// 	}

// 	contact1, err := instance2.FindContact(instance1.NodeID)
// 	if err != nil {
// 		t.Error("Instance 1's contact not found in Instance 2's contact list")
// 		return
// 	}
// 	if contact1.NodeID != instance1.NodeID {
// 		t.Error("Instance 1 ID incorrectly stored in Instance 2's contact list")
// 	}
// 	if contact2.NodeID != instance2.NodeID {
// 		t.Error("Instance 2 ID incorrectly stored in Instance 1's contact list")
// 	}
// 	return
// }

// func TestStore(t *testing.T) {
// 	// test Dostore() function and LocalFindValue() function
// 	instance1 := NewKademlia("localhost:7892")
// 	instance2 := NewKademlia("localhost:7893")
// 	host2, port2, _ := StringToIpPort("localhost:7893")
// 	instance1.DoPing(host2, port2)
//   	duration := time.Duration(1)*time.Second
//   	time.Sleep(duration)
// 	contact2, err := instance1.FindContact(instance2.NodeID)
// 	if err != nil {
// 		t.Error("Instance 2's contact not found in Instance 1's contact list")
// 		return
// 	}
// 	key := NewRandomID()
// 	value := []byte("Hello World")
// 	err = instance1.DoStore(contact2, key, value)
// 	if err != nil {
// 		t.Error("Can not store this value")
// 	}
// 	storedValue, err := instance2.LocalFindValue(key)
// 	if err != nil {
// 		t.Error("Stored value not found!")
// 	}
// 	if !bytes.Equal(storedValue, value) {
// 		t.Error("Stored value did not match found value")
// 	}
// 	return
// }

// func TestFindNode(t *testing.T) {
// 	// tree structure;
// 	// A->B->tree
// 	/*
// 	         C
// 	      /
// 	  A-B -- D
// 	      \
// 	         E
// 	*/

// 	instance1 := NewKademlia("localhost:7894")
// 	instance2 := NewKademlia("localhost:7895")
// 	host2, port2, _ := StringToIpPort("localhost:7895")
// 	instance1.DoPing(host2, port2)
//   	duration := time.Duration(1)*time.Second
//   	time.Sleep(duration)
// 	contact2, err := instance1.FindContact(instance2.NodeID)
// 	if err != nil {
// 		t.Error("Instance 2's contact not found in Instance 1's contact list")
// 		return
// 	}
// 	tree_node := make([]*Kademlia, 10)
// 	for i := 0; i < 10; i++ {
// 		address := "localhost:" + strconv.Itoa(7896+i)
// 		tree_node[i] = NewKademlia(address)
// 		host_number, port_number, _ := StringToIpPort(address)
// 		instance2.DoPing(host_number, port_number)
// 	}
// 	key := NewRandomID()
// 	contacts, err := instance1.DoFindNode(contact2, key)
// 	if err != nil {
// 		t.Error("Error doing FindNode")
// 	}

// 	if contacts == nil || len(contacts) == 0 {
// 		t.Error("No contacts were found")
// 	}
// 	// TODO: Check that the correct contacts were stored
// 	//       (and no other contacts)
// 		//			Self Contact included
// 	if  len(contacts) != 10 {
// 		t.Error("Contacts don't match ")
// 	}
// 	Nodes := make(map[ID]int)
// 	for _, c := range tree_node {
// 		Nodes[c.NodeID]++
// 	}
// 	for _, c1 := range contacts {
// 		for  m := range Nodes {
// 			if c1.NodeID == m {
// 				Nodes[m]--
// 			}
// 		}
// 	}

// 	for m := range Nodes {
// 		if Nodes[m] != 0 {
// 			t.Error("Contacts don't match ")
// 		}
// 	}
// 	return
// }

// func TestFindValue(t *testing.T) {
// 	// tree structure;
// 	// A->B->tree
// 	/*
// 	         C
// 	      /
// 	  A-B -- D
// 	      \
// 	         E
// 	*/
// 	instance1 := NewKademlia("localhost:7926")
// 	instance2 := NewKademlia("localhost:7927")
// 	host2, port2, _ := StringToIpPort("localhost:7927")
// 	instance1.DoPing(host2, port2)
//   	duration := time.Duration(1)*time.Second
//   	time.Sleep(duration)
// 	contact2, err := instance1.FindContact(instance2.NodeID)
// 	if err != nil {
// 		t.Error("Instance 2's contact not found in Instance 1's contact list")
// 		return
// 	}

// 	tree_node := make([]*Kademlia, 10)
// 	for i := 0; i < 10; i++ {
// 		address := "localhost:" + strconv.Itoa(7928+i)
// 		tree_node[i] = NewKademlia(address)
// 		host_number, port_number, _ := StringToIpPort(address)
// 		instance2.DoPing(host_number, port_number)
// 	}

// 	key := NewRandomID()
// 	value := []byte("Hello world")
// 	err = instance2.DoStore(contact2, key, value)
// 	if err != nil {
// 		t.Error("Could not store value")
// 	}

// 	// Given the right keyID, it should return the value
// 	foundValue, contacts, err := instance1.DoFindValue(contact2, key)
// 	if !bytes.Equal(foundValue, value) {
// 		t.Error("Stored value did not match found value")
// 	}

// 	//Given the wrong keyID, it should return k nodes.
// 	wrongKey := NewRandomID()
// 	foundValue, contacts, err = instance1.DoFindValue(contact2, wrongKey)
// 	if contacts == nil || len(contacts) < 10 {
// 		t.Error("Searching for a wrong ID did not return contacts")
// 	}
// 	// for _, c := range contacts {
// 	// 	t.Error(c)
// 	// }

// 	// TODO: Check that the correct contacts were stored
// 	//       (and no other contacts)
// 	//			Self Contact included
// 	if  len(contacts) != 10 {
// 		t.Error("Contacts don't match ")
// 	}
// 	Nodes := make(map[ID]int)
// 	for _, c := range tree_node {
// 		Nodes[c.NodeID]++
// 	}
// 	for _, c1 := range contacts {
// 		for  m := range Nodes {
// 			if c1.NodeID == m {
// 				Nodes[m]--
// 			}
// 		}
// 	}

// 	for m := range Nodes {
// 		if Nodes[m] != 0 {
// 			t.Error("Contacts don't match ")
// 		}
// 	}
// }

/*

func TestIterativeFindNode(t *testing.T) {

	kNum := 40
	targetIdx := kNum - 10
	instance2 := NewKademlia("localhost:7305")
	host2, port2, _ := StringToIpPort("localhost:7305")
	//	instance2.DoPing(host2, port2)
	tree_node := make([]*Kademlia, kNum)
	fmt.Print("Marker111")
	//t.Log("Before loop")
	for i := 0; i < kNum; i++ {
		address := "localhost:" + strconv.Itoa(7306+i)
		tree_node[i] = NewKademlia(address)
		tree_node[i].DoPing(host2, port2)
		t.Log("ID:" + tree_node[i].SelfContact.NodeID.AsString())
	}
	for i := 0; i < kNum; i++ {
		if i != targetIdx {
			tree_node[targetIdx].DoPing(tree_node[i].SelfContact.Host, tree_node[i].SelfContact.Port)
		}
	}
	SearchKey := tree_node[targetIdx].SelfContact.NodeID
	//t.Log("Wait for connect")
	//Connect(t, tree_node, kNum)
	//t.Log("Connect!")
	time.Sleep(100 * time.Millisecond)
	//cHeap := PriorityQueue{instance2.SelfContact, []Contact{}, SearchKey}
	//t.Log("Wait for iterative")
	res, err := tree_node[2].DoIterativeFindNode(SearchKey)


	fmt.Print("Marker")

	if err != nil {
		t.Error(err.Error())
	}
	t.Log("SearchKey:" + SearchKey.AsString())
	if res == nil || len(res) == 0 {
		t.Error("No contacts were found")
	}
	find := false
	fmt.Print("# of results:  ")
	fmt.Println(len(res))
	for _, value := range res {
		t.Log(value.NodeID.AsString())
		if value.NodeID.Equals(SearchKey) {
			find = true
		}
//		heap.Push(&cHeap, value)
	}
//	c := cHeap.Pop().(Contact)
//	t.Log("Closet Node:" + c.NodeID.AsString())

//	t.Log(strconv.Itoa(cHeap.Len()))
	if len(res) != 20 {
		t.Log("K list has no 20 ")
		t.Error("error")
	}
	if !find {
		t.Log("2:" + instance2.NodeID.AsString())
		t.Error("Find wrong id")
	}
	return
}
*/
func TestIterativeFindValue(t *testing.T) {
	kNum := 30
	targetIdx := kNum - 10
	instance1 := NewKademlia("localhost:7926")
	instance2 := NewKademlia("localhost:7927")
	host2, port2, _ := StringToIpPort("localhost:7927")

	instance1.DoPing(host2, port2)
	duration := time.Duration(1)*time.Second
  	time.Sleep(duration)
	tree_node := make([]*Kademlia, kNum)
	for i := 0; i < kNum; i++ {
		fmt.Println("First Ping Loop")
		address := "localhost:" + strconv.Itoa(7306+i)
		tree_node[i] = NewKademlia(address)
		tree_node[i].DoPing(host2, port2)
	}
	for i := 0; i < kNum; i++ {
		if i != targetIdx {
			fmt.Println("Second Ping Loop")
			tree_node[targetIdx].DoPing(tree_node[i].SelfContact.Host, tree_node[i].SelfContact.Port)
		}
	}

	SearchKey := tree_node[targetIdx].SelfContact.NodeID
	value := []byte("Hello world")
	err := instance2.DoStore(&tree_node[targetIdx].SelfContact, SearchKey, value)
	if err != nil {
		t.Error("Could not store value")
	}
	time.Sleep(100 * time.Millisecond)
	value, e := tree_node[2].DoIterativeFindValue(SearchKey)

	if e != nil || !bytes.Equal(value, []byte("Hello world")) {
		t.Error(e.Error())
	}
	return
}