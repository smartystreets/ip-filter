package IPFilter

import "C"
import (
	"math"
	"strconv"
	"strings"
)

type treeNode struct {
	value    uint64
	children []*treeNode
}

func TreeNode() *treeNode {
	return &treeNode{
		value:    0,
		children: nil,
	}
}

func (this *treeNode) New(addresses []string) {
	for _, item := range addresses {
		this.Insert(item)
	}
}

func (this *treeNode) Insert(ipAddress string) error {
	if len(ipAddress) == 0 {
		return ErrInvalidIPAddress
	}
	var last bool
	var network bool
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
			last = true
			network = true
			numOfNodes = 1
			ipAddress = ipAddress[:indexes[0]]
			this.addNetworkChild(numOfNodes, subnetBits, ipAddress, last, network)
			return nil
		}
		numOfNodes = 2
		ipAddress = ipAddress[:indexes[1]]
		this.addNetworkChild(numOfNodes, subnetBits, ipAddress, last, network)
		return nil
	}
	if subnetBits >= 16 && subnetBits <= 23 {
		if subnetBits == 16 {
			network = true
			numOfNodes = 2
			ipAddress = ipAddress[:indexes[1]]
			this.addNetworkChild(numOfNodes, subnetBits, ipAddress, last, network)
			return nil
		}
		numOfNodes = 3
		ipAddress = ipAddress[:indexes[2]]
		this.addNetworkChild(numOfNodes, subnetBits, ipAddress, last, network)
		return nil
	}
	if subnetBits >= 24 && subnetBits <= 31 {
		if subnetBits == 24 {
			network = true
			numOfNodes = 3
			ipAddress = ipAddress[:indexes[2]]
			this.addNetworkChild(numOfNodes, subnetBits, ipAddress, last, network)
			return nil
		}
		numOfNodes = 4
		this.addNetworkChild(numOfNodes, subnetBits, ipAddress, last, network)
		return nil
	}
	if subnetBits == 32 {
		numOfNodes = 4
		network = true
		this.addNetworkChild(numOfNodes, subnetBits, ipAddress, last, network)
		return nil
	}
	return nil
}

func (this *treeNode) addNetworkChild(numOfNodes, subnetBits int, ipAddress string, last, network bool) error {
	if numOfNodes == 1 && last == true && network == true {
		childFragment, _ := strconv.ParseUint(ipAddress, 10, 64)
		for _, children := range this.children {
			if children.value == childFragment {
				return nil
			}
		}

		child := &treeNode{
			value:    childFragment,
			children: nil,
		}

		this.children = append(this.children, child)
		return nil
	}

	if numOfNodes == 1 && last == true && network == false {
		childFragment, _ := strconv.ParseUint(ipAddress, 10, 64)
		this.addMinMaxChild(childFragment, subnetBits)
		return nil
	}

	for i := range ipAddress {
		if ipAddress[i] == '.' {
			childFragment, _ := strconv.ParseUint(ipAddress[:i], 10, 64)

			ipAddress = ipAddress[i+1:]

			for _, child := range this.children {
				if child.value == childFragment {
					return nil
				}
			}

			child := &treeNode{
				value: childFragment,
			}
			this.children = append(this.children, child)

			numOfNodes--
			if numOfNodes == 1 {
				last = true
			}
			child.addNetworkChild(numOfNodes, subnetBits, ipAddress, last, network)
			return nil
		}
	}
	return nil
}

func (this *treeNode) addMinMaxChild(minValue uint64, subnetBits int) error {
	var maxValue uint64

	if subnetBits > 8 && subnetBits < 16 {
		changeableBits := 16 - subnetBits

		valueToAdd := powInt(2, changeableBits) - 1

		maxValue = minValue + valueToAdd
	}
	if subnetBits > 16 && subnetBits < 24 {
		changeableBits := 24 - subnetBits

		valueToAdd := powInt(2, changeableBits) - 1

		maxValue = minValue + valueToAdd
	}
	if subnetBits > 24 && subnetBits < 32 {
		changeableBits := 32 - subnetBits

		valueToAdd := powInt(2, changeableBits) - 1

		maxValue = minValue + valueToAdd
	}

	for i := minValue; i <= maxValue; i++ {

		for _, child := range this.children {
			if child.value == i {
				return nil
			}
		}

		child := &treeNode{
			value:    i,
			children: nil,
		}

		this.children = append(this.children, child)
	}
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
