package IPFilter

import "C"
import (
	"strconv"
	"strings"
)

type treeNode struct {
	value    string
	network  bool
	children []*treeNode
}

func NewTreeNode() *treeNode {
	return &treeNode{
		value:    "",
		network:  false,
		children: nil,
	}
}

//Insert
/*
Will take in an IP Address as a string and will insert it into the tree accordingly
*/
func (this *treeNode) Insert(IPAddress string) error {
	if len(IPAddress) == 0 {
		return ErrInvalidIPAddress
	}

	var IPFragment string
	if position := strings.Index(IPAddress, "/"); position != -1 {
		subnetBits, _ := strconv.Atoi(IPAddress[position+1:])
		if subnetBits == 8 {
			IPFragment = IPAddress[:strings.Index(IPAddress, ".")]
			this.addChildNetwork(IPFragment)
			return nil
		}
		if subnetBits == 16 {
			index := findIndex(IPAddress, 16) //5
			if index == -1 {
				return ErrInvalidIndex
			}
			IPFragment = IPAddress[:index] //this current returns 210.111  --> without the second dot
			this.addChildNetwork(IPFragment)
			return nil
		}
		if subnetBits == 24 {
			//everything up to the third dot can be a node
			index := findIndex(IPAddress, 24)
			if index == -1 {
				return ErrInvalidIndex
			}
			IPFragment = IPAddress[:index] //this current returns 210.111  --> without the second dot
			this.addChildNetwork(IPFragment)
			return nil
		}
		if subnetBits == 32 {
			IPFragment = IPAddress[:position]
			this.addChildNetwork(IPFragment)
			return nil
		}
		IPAddress = IPAddress[:position]
	}

	this.addChildNoNetwork(IPAddress)

	return nil
}

func (this *treeNode) addChildNetwork(IPFragment string) error {
	for _, children := range this.children {
		if children.value == IPFragment {
			return nil //Network has already been found we don't care about the children.
		}
	}

	child := &treeNode{ //if the child did not exist already, create another treeNode
		value:    IPFragment,
		network:  true,
		children: nil,
	}
	this.children = append(this.children, child)

	return nil
}
func (this *treeNode) addChildNoNetwork(IPAddress string) error {
	if len(IPAddress) == 0 {
		return nil
	}

	var IPFragment string
	var index int

	if index = strings.Index(IPAddress, "."); index != -1 {
		IPFragment = IPAddress[:index]
	} else {
		IPFragment = IPAddress
		IPAddress = ""
	}
	for _, child := range this.children { //loop through all the children already attached to 210
		if child.value == IPFragment {
			if index != -1 {
				IPAddress = IPAddress[(index + 1):]
			}
			return child.addChildNoNetwork(IPAddress) // if 210 has a child that is equal to 111 then skip this one and call insert to add the rest of the sections to that child
		}
	}

	child := &treeNode{ //if the child did not exist already, create another treeNode
		value:    IPFragment,
		network:  false,
		children: nil,
	}

	this.children = append(this.children, child)

	if index != -1 {
		IPAddress = IPAddress[(index + 1):]
	}

	return child.addChildNoNetwork(IPAddress)
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
func (this *treeNode) Search(IPAddress string) bool {
	if len(IPAddress) == 0 {
		return false //TODO need to return error here?
	}

	indexes := findIndexes(IPAddress)

	for i, child := range this.children {
		if child.network == true {
			if this.searchNetworkChild(IPAddress, indexes) {
				return true
			}
		}
		//check to make sure the child isn't a network we are comparing? Reduce the amount of nodes to check?
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

	return false
}

func findIndexes(IPAddress string) []int {
	var indexes []int
	for i, _ := range IPAddress {
		if IPAddress[i] == '.' { // once we find a "." we will add that index to the slice of indexes
			indexes = append(indexes, i)
		}
	}
	return indexes
}

func (this *treeNode) searchNetworkChild(IPAddress string, Indexes []int) bool {
	var fragment string
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
