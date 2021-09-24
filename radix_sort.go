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
func (this *treeNode) Insert(IPAddress string) error{
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
			this.value = IPAddress[:strings.Index(IPAddress,".")]
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
	this.value = IPFragment
	this.network = false

	//reduce the IPAddress to grab the next section to add?
	IPAddress = IPAddress[(strings.Index(IPAddress, ".") + 1):] //+1 would get 111.12.12 instead of .111.12.12 I think

	IPFragment = IPAddress[:strings.Index(IPAddress, ".")] //grab the next fragment up until the next "." --> 111

	if len(IPAddress) > 0{ //if we have nothing else to add we need to return nil
		this.AddChild(IPFragment, IPAddress)
	}

	return nil
}

func(this *treeNode) AddChild(IPFragment string, IPAddress string) error{

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
//We will not be given a "/"
func (this *treeNode) Search(IPAddress string) bool {
	// 210.111.12.12

	if len(IPAddress) == 0 {
		return false
	}

	//var dotIndexes []int // to hold all the nodes we are looking for?
	//TODO: Get rid of the nested for loop
	var fragment string
	for i, _ := range IPAddress { // range through every single letter... Yep this is just an idea
		if IPAddress[i] == '.' { // once we find a "." we will do some stuff
			fragment = IPAddress[:i] // grab the fragment leading up to the "." --> this will be our node!

			//now we need to go look to see if that node exists?
			//except how do we make sure we are starting with the very very top of our tree? If we are on the top then we look for the matching

			for _, child := range this.children{
				if child.value != fragment { //It will keep looping as the value does not match the fragment
					continue
				}
				//So we found a match at this.child ... is this child a network? //Houston we have a problem, some of our networks are longer than one dot... this only works for the networks that are /8
				if child.network == true {
					return true // if the node is a /8 network then we return true e basta così
				}

				//add the fragment we just looked at and the next to check the next node
				
				//one option is to use the following function, but we would have to take it out of the current for loop

				//once it finds a match!
				remainingAddress := IPAddress[len(fragment)+1:] // reduce the IPAddress 111.12.12

				if true == this.Search(remainingAddress) {
					return true
				}
				break
			}
		}

	}

	//if there is no network provided then search node by node...
	//need to find a way to get the first node (get everything up until the "."

	//maybe loop through the string IP? yikes and find the "."

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