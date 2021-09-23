package IPFilter

import "strings"


type treeNode struct {
	//Do we need to keep the value of the potential head? < why would we need this? what benefit would we have with this?
	value    string     // Partial portion of IPAddress, the value of the node
	network  bool       // This is what tells us if it is a network or not
	networkValue string
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
func (tree *treeNode) Insert(IPAddress string) {
	//Check to see if there is a slash
	//if there is a slash then we will call the validate/add until we are done.
	if position := strings.Count("/", IPAddress); position != 0 {
		//We may want to do this a little differently, but for now this is what it looks like.
		if position == 8 {
			//Then the head of the node is 8 bits long (the first few numbers)
			tree.network = true
		}
		if position == 16 {
			tree.network = true
		}
		if position == 24 {
			tree.network = true
		}
	}
	//Start by seeing if the head is already in the tree, if it is not, then we keep moving.
	//Finally if there is no slash continue to add the node through the tree.
}

//validateNode
/*
Will validate the node. Is this a valid node? Does it already exist in our tree? Etc.
*/
func (tree *treeNode) validateNode(IPAddress string) {

}

//Delete
/*
Will take in an IPAddress as a string and will remove it from the tree/list of nodes
*/
func (tree *treeNode) Delete(IPAddress string) {}

//Search
/*
Will search the tree for a specific IPAddress within our nodes
*/
func (tree *treeNode) Search(IPAddress string) {}
