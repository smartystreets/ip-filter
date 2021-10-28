package IPFilter

import "C"
import (
	"math"
	"strconv"
	"strings"
)

type treeNode struct {
	value    uint64
	minValue uint64
	maxValue uint64
	children []*treeNode
}

func TreeNode() *treeNode {
	return &treeNode{
		value:    0,
		minValue: 0,
		maxValue: 0,
		children: nil,
	}
}

func (this *treeNode) New(addresses []string) {
	for _, item := range addresses {
		this.Insert(item)
	}
}

func (this *treeNode) Insert(ipAddress string) error {
	var last bool
	var numOfNodes int
	position := strings.Index(ipAddress, "/")
	subnetBits, _ := strconv.Atoi(ipAddress[position+1:])
	ipAddress = ipAddress[:position]

	if subnetBits > 32 {
		return ErrInvalidIPAddress
	}

	indexes := findIndexes(ipAddress)

	if subnetBits >= 8 && subnetBits <= 15 {
		if subnetBits == 8 {
			last = false
			numOfNodes = 1
			ipAddress = ipAddress[:indexes[0]]
			this.addNetworkChild(numOfNodes, subnetBits, ipAddress, last)
		}
		numOfNodes = 2
		ipAddress = ipAddress[:indexes[1]]
		this.addNetworkChild(numOfNodes, subnetBits, ipAddress, last)
	}
	if subnetBits >= 16 && subnetBits <= 23 {
		if subnetBits == 16 {
			numOfNodes = 2
			ipAddress = ipAddress[:indexes[1]]
		}
		numOfNodes = 3
		ipAddress = ipAddress[:indexes[2]]
		this.addNetworkChild(numOfNodes, subnetBits, ipAddress, last)
	}
	if subnetBits >= 24 && subnetBits <= 31 {
		if subnetBits == 24 {
			numOfNodes = 3
			ipAddress = ipAddress[:indexes[2]]
		}
		numOfNodes = 4
		this.addNetworkChild(numOfNodes, subnetBits, ipAddress, last)
	}
	if subnetBits == 32 {
		numOfNodes = 4
		this.addNetworkChild(numOfNodes, subnetBits, ipAddress, last)
	}
	return nil
}

func (this *treeNode) addNetworkChild(numOfNodes, subnetBits int, ipAddress string, last bool) error {
	if numOfNodes == 1 && last == false{
		childFragment, _ := strconv.ParseUint(ipAddress, 10, 64)
		for _, children := range this.children {
			if children.value == childFragment {
				return nil
			}
		}

		child := &treeNode{
			children: nil,
		}

		this.children = append(this.children, child)
		return nil
	}
	//I may have to do this for the min max children
	if numOfNodes == 1 && last == true{
		childFragment, _ := strconv.ParseUint(ipAddress, 10, 64)
		this.addMinMaxChild(childFragment, subnetBits)
	}

	//TODO: I also don't believe that I should have this numOfNodes... maybe try and do range of ipAddress
	for i := range ipAddress {
		if ipAddress[i] == '.' {
			childFragment, _ := strconv.ParseUint(ipAddress[:i], 10, 64)
			ipAddress = ipAddress[i+1:]
			for _, children := range this.children {
				if children.value == childFragment {
					return nil
				}
			}

			child := &treeNode{
				children: nil,
			}

			this.children = append(this.children, child)
			//We need a way to check if this is the last one to add... then we can change the bool "last"
			numOfNodes--
			if numOfNodes == 1 {
				last = true
			}
			child.addNetworkChild(numOfNodes, subnetBits, ipAddress, last)
		}
	}

	for _, child := range this.children {
		if child.value != ipFragment {
			continue
		}
	}

	for _, children := range this.children {
		if children.value == IPFragment {
			for _, childrenRanges := range children.children {
				if childrenRanges.minValue == minValue && childrenRanges.maxValue == maxValue {
					return nil
				}
				children.addMinMaxChild(maxValue, minValue)
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

	child.addMinMaxChild(maxValue, minValue)
	return nil
}
func (this *treeNode) addMinMaxChild(minValue uint64, subnetBits int) error {

	//if subnetBits > 8 && subnetBits < 16 {
	//firstSection := strings.Index(IPAddress, ".")
	//secondSection := findIndex(IPAddress, 16)
	//
	//minValue, _ := strconv.ParseUint(IPAddress[firstSection+1:secondSection], 10, 64)
	//
	//changeableBits := 16 - subnetBits
	//
	//valueToAdd := powInt(2, changeableBits) - 1
	//
	//maxValue := minValue + valueToAdd
	//
	//IPFragment = IPAddress[:strings.Index(IPAddress, ".")]
	//
	//this.addNetworkChild(maxValue, minValue, IPFragment)


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

func (this *treeNode) Search(IPAddress string) bool {
	//var childRange string
	//var countOfPeriods int
	//var fragment string
	//first := true
	//for i, _ := range IPAddress {
	//	if countOfPeriods != 3 {
	//		if IPAddress[i] != '.' {
	//			continue
	//		} else {
	//			countOfPeriods++
	//		}
	//	}
	//
	//	if countOfPeriods == 3 && first == true{
	//		fragment = IPAddress[:i]
	//		first = false
	//	} else if countOfPeriods == 3 && first == false{
	//		fragment = IPAddress
	//	} else {
	//		fragment = IPAddress[:i]
	//	}
	//
	//	for _, child := range this.children {
	//		if child.value != fragment {
	//			continue
	//		}
	//
	//		if child.children == nil {
	//			return true
	//		}
	//
	//		var indexOfSecond int
	//
	//
	//		if countOfPeriods == 3 && first == true{
	//			childRange = IPAddress[(len(fragment)+1):]
	//		} else if countOfPeriods == 3 && first != true{
	//			childRange = IPAddress[(len(fragment)+1):]
	//		} else {
	//			count := i
	//			for _, _ = range IPAddress {
	//				count++
	//				if IPAddress[count] == '.' {
	//					indexOfSecond = count
	//					break
	//				}
	//			}
	//			childRange = IPAddress[(len(fragment) + 1):indexOfSecond]
	//		}
	//
	//		changeableRange, _ := strconv.ParseUint(childRange, 10, 64)
	//
	//		for _, child2 := range child.children {
	//			if changeableRange >= child2.minValue && changeableRange <= child2.maxValue {
	//				return true
	//			}
	//		}
	//	}
	//}
	return false
}

func findIndexes(IPAddress string) [3]int {
	var indexes [3]int
	count := 0
	i := 0
	for count < 3 {
		if IPAddress[i] == '.' {
			indexes[count] = i
			count++
		}
		i++
	}
	return indexes
}
