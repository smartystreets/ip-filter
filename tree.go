package IPFilter

import "C"
import (
	"math"
	"strconv"
	"strings"
)

type treeNode struct {
	value    string
	minValue uint64
	maxValue uint64
	children []*treeNode
}

func NewTreeNode() *treeNode {
	return &treeNode{
		value:    "",
		minValue: 0,
		maxValue: 0,
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
			this.addChildNetwork(0, 0, IPFragment)
			return nil
		}
		if subnetBits == 16 {
			index := findIndex(IPAddress, 16) //5
			IPFragment = IPAddress[:index]    //this current returns 210.111  --> without the second dot
			this.addChildNetwork(0, 0, IPFragment)
			return nil
		}
		if subnetBits == 24 {
			//everything up to the third dot can be a node
			index := findIndex(IPAddress, 24)
			IPFragment = IPAddress[:index] //this current returns 210.111  --> without the second dot
			this.addChildNetwork(0, 0, IPFragment)
			return nil
		}
		if subnetBits == 32 {
			IPFragment = IPAddress[:position]
			this.addChildNetwork(0, 0, IPFragment)
			return nil
		}
		IPAddress = IPAddress[:position]

		if subnetBits > 8 && subnetBits < 16 {
			firstSection := strings.Index(IPAddress, ".")
			secondSection := findIndex(IPAddress, 16)

			minValue, _ := strconv.ParseUint(string(IPAddress[firstSection+1:secondSection]), 10, 64)

			changeableBits := 16 - subnetBits

			valueToAdd := powInt(2, changeableBits) - 1

			maxValue := minValue + valueToAdd

			IPFragment = IPAddress[:strings.Index(IPAddress, ".")]

			this.addChildNetwork(maxValue, minValue, IPFragment)
			return nil
		}
		if subnetBits > 16 && subnetBits < 24 {
			firstSection := findIndex(IPAddress, 16)
			secondSection := findIndex(IPAddress, 24)
			minValue, _ := strconv.ParseUint(IPAddress[(firstSection+1):secondSection], 10, 64)

			changeableBits := 24 - subnetBits

			valueToAdd := powInt(2, changeableBits) - 1

			maxValue := minValue + valueToAdd

			IPFragment = IPAddress[:findIndex(IPAddress, 16)]

			this.addChildNetwork(maxValue, minValue, IPFragment)
			return nil
		}
		if subnetBits > 24 && subnetBits < 32 {
			minValue, _ := strconv.ParseUint(IPAddress[findIndex(IPAddress, 24)+1:], 10, 64)

			changeableBits := 32 - subnetBits

			valueToAdd := powInt(2, changeableBits) - 1

			maxValue := minValue + valueToAdd

			IPFragment = IPAddress[:findIndex(IPAddress, 24)]

			this.addChildNetwork(maxValue, minValue, IPFragment)
			return nil
		}

	}
	return nil
}

func (this *treeNode) addChildNetwork(maxValue, minValue uint64, IPFragment string) error {

	for _, children := range this.children {
		if children.value == IPFragment {
			for _, childrenRanges := range children.children {
				if childrenRanges.minValue == minValue && childrenRanges.maxValue == maxValue {
					return nil
				}
				children.addChildNoNetwork(maxValue, minValue)
				return nil
			}
		}
	}

	child := &treeNode{
		value:    IPFragment,
		children: nil,
	}

	this.children = append(this.children, child)

	if maxValue == 0 {
		return nil
	}

	child.addChildNoNetwork(maxValue, minValue)
	return nil
}
func (this *treeNode) addChildNoNetwork(maxValue, minValue uint64) error {

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

func powInt(x, y int) uint64 {
	return uint64(math.Pow(float64(x), float64(y)))
}

//Search
/*
Will search the tree for a specific IPAddress within our nodes
*/
func (this *treeNode) Search(IPAddress string) bool {
	var childRange string
	var fragment string
	if len(IPAddress) == 0 {
		return false
	}

	indexes := findIndexes(IPAddress)
	for i := 0; i < 4; i++ {

		if i == 3 {
			fragment = IPAddress
		} else {
			fragment = IPAddress[:indexes[i]]
		}

		for j, child := range this.children {

			if child.value != fragment {
				continue
			}

			// if the IPAddress fragment is found then
			//check to see if the network is set to true
			if child.children == nil {
				//if set to true then return true
				return true
			}

			//Get the section
			//TODO: There is an error
			//TODO: This is where we are getting thrown off for the tests
			if j == 0 {
				childRange = IPAddress[(indexes[j] + 1):indexes[(j+1)]]
			} else {
				childRange = IPAddress[(indexes[(j-1)] + 1):indexes[j]]
			}

			//turn that magic section to a number
			parsed, err := strconv.ParseUint(childRange, 10, 32)

			if err != nil {
				return false //TODO: return error?
			}

			//loop through it's children
			for _, c := range child.children {
				if parsed >= c.minValue && parsed <= c.maxValue {
					return true
				}
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
