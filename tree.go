package IPFilter

import "C"
import (
	"strconv"
	"strings"
)

type treeNode struct {
	value    string
	minValue string
	maxValue string
	children []*treeNode
}

func NewTreeNode() *treeNode {
	return &treeNode{
		value:    "",
		minValue: "",
		maxValue: "",
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
			this.addChildNetwork("", "", IPFragment)
			return nil
		}
		if subnetBits == 16 {
			index := findIndex(IPAddress, 16) //5
			IPFragment = IPAddress[:index]    //this current returns 210.111  --> without the second dot
			this.addChildNetwork("", "", IPFragment)
			return nil
		}
		if subnetBits == 24 {
			//everything up to the third dot can be a node
			index := findIndex(IPAddress, 24)
			IPFragment = IPAddress[:index] //this current returns 210.111  --> without the second dot
			this.addChildNetwork("", "", IPFragment)
			return nil
		}
		if subnetBits == 32 {
			IPFragment = IPAddress[:position]
			this.addChildNetwork("", "", IPFragment)
			return nil
		}
		IPAddress = IPAddress[:position]

		if subnetBits > 8 && subnetBits < 16 {
			minValue, _ := strconv.Atoi(IPAddress[:strings.Index(IPAddress, ".")])
			changeableBits := 16 - subnetBits

			valueToAdd := (2 ^ changeableBits) - 1

			maxValue := strconv.Itoa(minValue + valueToAdd)

			IPFragment = IPAddress[:strings.Index(IPAddress, ".")]
			this.addChildNetwork(maxValue, strconv.Itoa(minValue), IPFragment)
			this.addChildNoNetwork(maxValue, strconv.Itoa(minValue))
		}
		if subnetBits > 16 && subnetBits < 24 {
			minValue, _ := strconv.Atoi(IPAddress[:findIndex(IPAddress, 16)])
			changeableBits := 24 - subnetBits

			valueToAdd := (2 ^ changeableBits) - 1

			maxValue := strconv.Itoa(minValue + valueToAdd)

			this.addChildNoNetwork(maxValue, strconv.Itoa(minValue))
		}
		if subnetBits > 24 && subnetBits < 32 {
			minValue, _ := strconv.Atoi(IPAddress[:findIndex(IPAddress, 24)])
			changeableBits := 24 - subnetBits

			valueToAdd := (2 ^ changeableBits) - 1

			maxValue := strconv.Itoa(minValue + valueToAdd)
			this.addChildNoNetwork(maxValue, strconv.Itoa(minValue))
		}

	}
	return nil
}

func (this *treeNode) addChildNetwork(maxValue, minValue, IPFragment string) error {

	for _, children := range this.children {
		if children.value == IPFragment {
			return nil
		}
	}

	child := &treeNode{
		value:    IPFragment,
		children: nil,
	}
	this.children = append(this.children, child)

	if maxValue == "" && minValue == "" {
		return nil
	}

	child.addChildNoNetwork(maxValue, minValue)
	return nil
}
func (this *treeNode) addChildNoNetwork(maxValue, minValue string) error {

	for _, children := range this.children {
		if children.minValue == minValue && children.maxValue == maxValue {
			return nil
		}
	}

	child := &treeNode{
		minValue: minValue,
		maxValue: maxValue,
		children: nil,
	}

	this.children = append(this.children, child)

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
}

//Search
/*
Will search the tree for a specific IPAddress within our nodes
*/
func (this *treeNode) Search(IPAddress string) bool {
	var fragment string
	if len(IPAddress) == 0 {
		return false //TODO need to return error here?
	}

	indexes := findIndexes(IPAddress) //TODO: add a check here?

	for i, child := range this.children {
		//if child.network == true {
		//	if this.searchNetworkChild(IPAddress, indexes) {
		//		return true
		//	}
		//}

		if indexes == nil { //TODO: forse just add a check if the indexes are empty
			if child.value != IPAddress {
				continue
			}
			return true
		}

		fragment = IPAddress[:indexes[i]] //TODO: add a check here for out of range //ALSO can I just hard card a 0?

		if child.value != fragment {
			continue
		}

		remainingAddress := IPAddress[len(fragment)+1:]

		if len(remainingAddress) > 0 { //TODO: Do I still need this one?
			if true == child.Search(remainingAddress) {
				return true
			}
		}
	}

	return false
}

func findIndexes(IPAddress string) []int {
	var indexes []int
	for i, _ := range IPAddress {
		if IPAddress[i] == '.' {
			indexes = append(indexes, i)
		}
	}
	return indexes
}

func (this *treeNode) searchNetworkChild(IPAddress string, Indexes []int) bool { // indexes will Always be 3
	var fragment string
	i := 0
	for i < 4 {
		if i == 3 {
			fragment = IPAddress
			i = 5
		} else {
			fragment = IPAddress[:Indexes[i]]
			i++
		}

		for _, child := range this.children {
			if child.value != fragment {
				continue
			}
			//if child.network == true {
			//	return true
			//}
		}
	}
	return false
}
