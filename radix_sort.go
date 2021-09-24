package IPFilter

import (
	"strconv"
	"strings"
)

type treeNode struct {
	value    string     // Partial portion of IPAddress, the value of the node
	network  bool       // This is what tells us if it is a network or not
	children []treeNode // The potential children of a node
}

/*
Becomes:        Network:
210 <- head     210
111 <- child    111
12  <- child    No children
12  <- child
*/

//Insert
/*
Will take in an IP Address as a string and will insert it into the tree accordingly
*/
func (this *treeNode) Insert(IPAddress string) error {
	//we may not need this. I need to walk through to confirm
	if len(IPAddress) == 0 {
		return nil
	}

	//Check to see if there is a slash
	if position := strings.Index(IPAddress, "/"); position != -1 {
		//We may want to do this a little differently, but for now this is what it looks like.
		subnetBits, _ := strconv.Atoi(IPAddress[position:])
		if subnetBits == 8 {
			//Then the head of the node is 8 bits long (the first few numbers)
			this.network = true
			this.value = IPAddress[:strings.Index(IPAddress, ".")]
			this.children = nil
			//TODO: We need to return the node here - no children after this
			return nil
		}
		if position == 16 {
			//everything up to the second dot can be a node
			index := findIndex(IPAddress, 16) //5
			this.network = true
			this.value = IPAddress[:index] //this current returns 210.111  --> without the second dot
			this.children = nil
			return nil
		}
		if position == 24 {
			//everything up to the third dot can be a node
			index := findIndex(IPAddress, 24)
			this.network = true
			this.value = IPAddress[:index] //this current returns 210.111  --> without the second dot
			this.children = nil
			return nil
		}
	}

	//We want to set the values in the case that it does NOT have "/"
	IPFragment := IPAddress[:strings.Index(IPAddress, ".")]
	//this.value = IPFragment
	//this.network = false

	//reduce the IPAddress to grab the next section to add?
	IPAddress = IPAddress[(strings.Index(IPAddress, ".") + 1):] //+1 would get 111.12.12 instead of .111.12.12 I think

	IPFragment = IPAddress[:strings.Index(IPAddress, ".")] //grab the next fragment up until the next "." --> 111

	if len(IPAddress) > 0 { //if we have nothing else to add we need to return nil
		this.AddChild(IPFragment, IPAddress)
	}

	return nil
}

func (this *treeNode) AddChild(IPFragment string, IPAddress string) error {

	// IPFragment --> 111
	for _, children := range this.children { //loop through all the children already attached to 210
		if children.value == IPFragment {
			return children.Insert(IPAddress) // if 210 has a child that is equal to 111 then skip this one and call insert to add the rest of the sections to that child
		}
	}

	Child := &treeNode{ //if the child did not exist already, create another treeNode
		value:    IPFragment,
		network:  false,
		children: nil,
	}

	this.children = append(this.children, *Child) //append the node to the parent

	return nil
}

func findIndex(IPAddress string, position int) int {
	var count int
	for i := range IPAddress {
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
//We will not be given a "/"
func (this *treeNode) Search(IPAddress string) bool {
	// 210.111.12.12

	if len(IPAddress) == 0 {
		return false
	}

	indexes := findIndexes(IPAddress) //now I know where all the "." are...

	if this.networkCheck(IPAddress, indexes) {
		return true
	}

	//so if we get here we did not find a matching network node and will have to just check through the rest of the nodes
	for i, child := range this.children {

		//check to make sure the child isn't a network we are comparing? Reduce the amount of nodes to check?
		if child.network == false{
			fragment := IPAddress[:indexes[i]]
			if child.value != fragment { //It will keep looping as the value does not match the fragment
				continue
			}
			//The fragment does match -- let's reduce the exiting Address to check the next fragment
			remainingAddress := IPAddress[len(fragment)+1:] // reduce the IPAddress 111.12.12

			//Recursively call the search function for the next child?
			if true == this.Search(remainingAddress) {
				return true
			}
		}
	}

	return false
}

//helper function to find all the indexes
func findIndexes(IPAddress string) []int {
	var indexes []int
	for i, _ := range IPAddress {
		if IPAddress[i] == '.' { // once we find a "." we will add that index to the slice of indexes
			indexes = append(indexes, i)
		}
	}
	return indexes
}

func (this *treeNode) networkCheck(IPAddress string, Indexes []int) bool {
	var fragment string

	//maybe we should hard code this?
	for i := range Indexes { //i --> 0, 1, 2 : indexes --> 3, 6, 9
		fragment = IPAddress[:Indexes[i]] // this should grab each section - but then get the next section and the next

		//this goes through all the children to find a match ... //how do we make sure we are starting with the very very top of our tree? If we are on the top then we look for the matching
		for _, child := range this.children {
			if child.value != fragment { //It will keep looping as the value does not match the fragment
				continue
			}

			if child.network == true {
				return true // if the node is a /8 network then we return true e basta così
			}
		}
	}
	return false
}
