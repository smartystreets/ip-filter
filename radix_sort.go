package IPFilter

import (
	"strconv"
	"strings"
)


type treeNode struct {
	//Do we need to keep the value of the potential head? < why would we need this? what benefit would we have with this?
	value    string     // Partial portion of IPAddress, the value of the node
	network  bool       // This is what tells us if it is a network or not
	children []treeNode // The potential children of a node
}

//Check to see if there is a /8 (whatever)
/*
If the slash exists, get the network house and then collapse it.
If the slash does not exist, go by .
*/
//210.111.12.12/8
/*
Becomes:
210 <- head
111 <- child
12  <- child
12  <- child

Collapse on the net host
*/

//Insert
/*
Will take in an IP Address as a string and will insert it into the tree accordingly
*/
func (this *treeNode) Insert(IPAddress string) error{
	//Check to see if there is a slash
	//if there is a slash then we will call the validate/add until we are done.
	if position := strings.Index(IPAddress, "/"); position != -1 {
		//We may want to do this a little differently, but for now this is what it looks like.
		subBits, _ := strconv.Atoi(IPAddress[position:])
		if subBits == 8 {
			//Then the head of the node is 8 bits long (the first few numbers)
			//everything up to the first dot can be a node

			this.network = true
			this.value = IPAddress[:strings.Index(IPAddress,".")]
			this.children = nil
				//TODO: We need to return the node here - no children after this
			return nil
		}
		if position == 16 {
			//everything up to the second dot can be a node
			index := findIndex(IPAddress, 16)
			this.network = true
			this.value = IPAddress[:index]
			this.children = nil
			return nil
		}
		if position == 24 {
			//everything up to the third dot can be a node
			index := findIndex(IPAddress, 24)
			this.network = true
			this.value = IPAddress[:index]
			this.children = nil
			return nil
		}
	}
	return nil
	//Start by seeing if the head is already in the this, if it is not, then we keep moving.
	//Finally if there is no slash continue to add the node through the this.
}

func findIndex(IPAddress string, position int) int{
	var count int
	for i := range IPAddress{
		if IPAddress[i] == '.' {
			count++
		}
		if count == 2 && position == 16 {
			return i
		}
		if count == 3 && position == 24 {
			return i
		}
	}
	return -1
	//TODO: Return error
}

//validateNode
/*
Will validate the node. Is this a valid node? Does it already exist in our tree? Etc.
*/
func (this *treeNode) validateNode(IPAddress string) {

}

//Delete
/*
Will take in an IPAddress as a string and will remove it from the tree/list of nodes
*/
func (this *treeNode) Delete(IPAddress string) {}

//Search
/*
Will search the tree for a specific IPAddress within our nodes
*/
func (this *treeNode) Search(IPAddress string) {}
