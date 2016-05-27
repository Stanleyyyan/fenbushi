// Ounan Ma, omg049
// Chong Yan, cyu422
// Wenjie Zhang, wzm416

package libkademlia

import (
	"bytes"
	"net"
	"strconv"
	"testing"
	// "time"
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
/*
func TestIterativeFindNode(t *testing.T) {
	kNum := 40
	targetIdx := kNum - 5
	// instance1 := NewKademlia("localhost:7305")
	host1, port1, _ := StringToIpPort("localhost:7305")
	tree_node := make([]*Kademlia, kNum)
	for i := 0; i < kNum; i++ {
		address := "localhost:" + strconv.Itoa(7306+i)
		tree_node[i] = NewKademlia(address)
		tree_node[i].DoPing(host1, port1)
	}
	for i := 0; i < kNum; i++ {
		if i != targetIdx {
			tree_node[targetIdx].DoPing(tree_node[i].SelfContact.Host, tree_node[i].SelfContact.Port)
		}
	}
	SearchKey := tree_node[targetIdx].SelfContact.NodeID
	res, err := tree_node[5].DoIterativeFindNode(SearchKey)
	if err != nil {
		t.Error(err.Error())
	}

	if res == nil || len(res) == 0 {
		t.Error("No contacts were found")
	}
	find := false
	for _, value := range res {
		if value.NodeID.Equals(SearchKey) {
			find = true
		}
	}
	if len(res) != 20 {
		t.Log("K list is not full")
		t.Error("error")
	}
	if !find {
		t.Error("Find wrong id")
	}
	return
}

func TestIterativeFindNodeFail(t *testing.T) {
	Trie1Num := 20
	ContactIdx := Trie1Num - 10
	instance1 := NewKademlia("localhost:7926")
	instance2 := NewKademlia("localhost:7927")
	host2, port2, _ := StringToIpPort("localhost:7927")
	instance1.DoPing(host2, port2)
	duration := time.Duration(1)*time.Second
  	time.Sleep(duration)
	tree_node_trie1 := make([]*Kademlia, Trie1Num)
	for i := 0; i < Trie1Num; i++ {
		fmt.Println("First Ping Loop")
		address := "localhost:" + strconv.Itoa(7906+i)
		tree_node_trie1[i] = NewKademlia(address)
		tree_node_trie1[i].DoPing(host2, port2)
	}
	Trie2Num := 10
	tree_node_trie2 := make([]*Kademlia, Trie2Num)
	host22, port22, _ := StringToIpPort(tree_node_trie1[ContactIdx].SelfContact.Host.String())
	for i := 0; i < Trie2Num; i++ {
		fmt.Println("First Ping Loop")
		address := "localhost:" + strconv.Itoa(6906+i)
		tree_node_trie2[i] = NewKademlia(address)
		tree_node_trie2[i].DoPing(host22, port22)
	}

	Trie3Num := 10
	tree_node_trie3 := make([]*Kademlia, Trie2Num)
	host33, port33, _ := StringToIpPort(tree_node_trie2[3].SelfContact.Host.String())
	for i := 0; i < Trie3Num; i++ {
		fmt.Println("First Ping Loop")
		address := "localhost:" + strconv.Itoa(5106+i)
		tree_node_trie3[i] = NewKademlia(address)
		tree_node_trie3[i].DoPing(host33, port33)
	}

	SearchKey := tree_node_trie3[3].SelfContact.NodeID
	res, err := instance2.DoIterativeFindNode(SearchKey)

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
	}

	if len(res) != 20 {
		t.Log("K list has no 20 ")
		t.Error("error")
	}
	if find {
		t.Error("Find wrong id")
	}
	return
}

func TestIterativeStore(t *testing.T) {
	kNum := 40
	targetIdx := kNum - 10
	instance1 := NewKademlia("localhost:4926")
	// instance2 := NewKademlia("localhost:7927")
	host2, port2, _ := StringToIpPort("localhost:4927")
	instance1.DoPing(host2, port2)
	duration := time.Duration(1)*time.Second
  	time.Sleep(duration)
	tree_node := make([]*Kademlia, kNum)
	for i := 0; i < kNum; i++ {
		fmt.Println("First Ping Loop")
		address := "localhost:" + strconv.Itoa(4306+i)
		tree_node[i] = NewKademlia(address)
		tree_node[i].DoPing(host2, port2)
	}
	for i := 0; i < kNum; i++ {
		if i != targetIdx {
			fmt.Println("Second Ping Loop")
			tree_node[targetIdx].DoPing(tree_node[i].SelfContact.Host, tree_node[i].SelfContact.Port)
		}
	}
	TargetKey := tree_node[targetIdx].SelfContact.NodeID
	value := []byte("Hello world")
	res, _ := tree_node[2].DoIterativeStore(TargetKey, value)
	for _, c := range res {
		value, _, _ := tree_node[2].DoFindValue(&c, TargetKey)
		if !bytes.Equal(value, []byte("Hello world")) {
			t.Error("Stored value did not match found value")
		}
	}
	return
}


func TestIterativeFindValue(t *testing.T) {
	kNum := 30
	targetIdx := kNum - 10
	instance1 := NewKademlia("localhost:3926")
	instance2 := NewKademlia("localhost:3927")
	host2, port2, _ := StringToIpPort("localhost:3927")
	instance1.DoPing(host2, port2)
	duration := time.Duration(1)*time.Second
  	time.Sleep(duration)
	tree_node := make([]*Kademlia, kNum)
	for i := 0; i < kNum; i++ {
		fmt.Println("First Ping Loop")
		address := "localhost:" + strconv.Itoa(1306+i)
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

func TestIterativeFindValueFail(t *testing.T) {
	kNum := 30
	targetIdx := kNum - 10
	instance1 := NewKademlia("localhost:2926")
	instance2 := NewKademlia("localhost:2927")
	host2, port2, _ := StringToIpPort("localhost:2927")
	instance1.DoPing(host2, port2)
	duration := time.Duration(1)*time.Second
  	time.Sleep(duration)
	tree_node := make([]*Kademlia, kNum)
	for i := 0; i < kNum; i++ {
		fmt.Println("First Ping Loop")
		address := "localhost:" + strconv.Itoa(2306+i)
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
	value, e := tree_node[2].DoIterativeFindValue(SearchKey)

	closest := tree_node[0].SelfContact.NodeID
	closestDis := tree_node[0].SelfContact.NodeID.Xor(SearchKey)
	for _, c := range tree_node  {
		dis := c.SelfContact.NodeID.Xor(SearchKey)
		if c.SelfContact.NodeID != SearchKey && dis.Compare(closestDis) < 0{
			closest = c.SelfContact.NodeID
			closestDis = dis
		}
	}
	dis := instance1.SelfContact.NodeID.Xor(SearchKey)
	if dis.Compare(closestDis) < 0{
		closest = instance1.SelfContact.NodeID
		closestDis = dis
	}
	dis = instance2.SelfContact.NodeID.Xor(SearchKey)
	if dis.Compare(closestDis) < 0{
		closest = instance2.SelfContact.NodeID
		closestDis = dis
	}
	if !bytes.Equal(value, []byte("")) || e.Error() != closest.AsString() {
		t.Error(e.Error())
	}
	return
}

*/


func TestVanish_Unvanish(t *testing.T) {
	instance1 := NewKademlia("localhost:7927")
	//instance2 := NewKademlia("localhost:7927")
	host1, port1, _ := StringToIpPort("localhost:7927")
	tree_node_trie1 := make([]*Kademlia, 30)
	for i := 0; i < 30; i++ {
		address := "localhost:" + strconv.Itoa(7930+i)
		tree_node_trie1[i] = NewKademlia(address)
		tree_node_trie1[i].DoPing(host1, port1)
	}
	SearchKey := instance1.SelfContact.NodeID
	VdoID := NewRandomID()
	instance1.Vanish(VdoID, []byte("AAAAAA"), 20, 15, 1)
	ciphertext := tree_node_trie1[10].Unvanish(SearchKey, VdoID)
	// data := []byte("Hello world")
	// numberKeys := byte(10)
	// threshold := byte(9)
	// timeoutSeconds := 100
	// tree_node_trie1[3].Vanish(VdoID, data , numberKeys, threshold, timeoutSeconds)
	// ciphertext := tree_node_trie1[3].Unvanish(SearchKey, VdoID)
	fmt.Println(string(ciphertext) + "is result")
	if (!bytes.Equal(ciphertext, []byte("AAAAAA"))) {
		t.Error("Unvanish error")
	}

	return
}
