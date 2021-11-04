package IPFilter

import "C"
import (
	"math"
	"strconv"
	"strings"
)

type treeNode struct {
	children []*treeNode
	isBanned bool
}

func TreeNode() *treeNode {
	return &treeNode{
		children: make([]*treeNode, 2),
		isBanned: false,
	}
}

func (this *treeNode) New(addresses []string) {
	for _, item := range addresses {
		this.Insert(item)
	}
}

func (this *treeNode) Insert(ipAddress string) error {

	var fullValue uint32

	position := strings.Index(ipAddress, "/")
	subnetBits, _ := strconv.Atoi(ipAddress[position+1:])
	ipAddress = ipAddress[:position]

	for i := 0; i < 4; i++ {
		var fragment uint64
		var index int

		for j := range ipAddress {
			if ipAddress[j] == '.' {
				index = j
				break
			}
			index = 0
		}

		if index == 0 {
			fragment, _ = strconv.ParseUint(ipAddress, 10, 32)
		} else {
			fragment, _ = strconv.ParseUint(ipAddress[:index], 10, 32)
		}
		fullValue = fullValue << 8
		fullValue += uint32(fragment)
	}

	current := this

	for i := 0; i < subnetBits; i++ {
		next := uint32(fullValue << (32 - i) >> (31))
		child := current.children[next]

		if child == nil {
			child = TreeNode()

			current.children[next] = child
		}
		current = child
	}

	current.isBanned = true
	return nil
}

//func (this *treeNode) addNetworkChild(numOfNodes, subnetBits int, ipAddress string, last, network bool) error {
//	if numOfNodes == 1 && last == true && network == true {
//		childFragment, _ := strconv.ParseUint(ipAddress, 10, 8)
//		for _, children := range this.children {
//			if children.value == childFragment {
//				return nil
//			}
//		}
//
//		child := &treeNode{
//			value:    childFragment,
//			children: nil,
//			end:      true,
//			min:      childFragment,
//			max:      childFragment,
//		}
//
//		this.children = append(this.children, child)
//		return nil
//	}
//
//	if numOfNodes == 1 && last == true && network == false {
//		childFragment, _ := strconv.ParseUint(ipAddress, 10, 64)
//		this.addMinMaxChild(childFragment, subnetBits)
//		return nil
//	}
//
//	for i := range ipAddress {
//		if ipAddress[i] == '.' {
//			childFragment, _ := strconv.ParseUint(ipAddress[:i], 10, 64)
//
//			ipAddress = ipAddress[i+1:]
//
//			for _, child := range this.children {
//				if child.value == childFragment {
//					numOfNodes--
//					if numOfNodes == 1 {
//						last = true
//					}
//					child.addNetworkChild(numOfNodes, subnetBits, ipAddress, last, network)
//					return nil
//				}
//			}
//
//			child := &treeNode{
//				value: childFragment,
//			}
//			this.children = append(this.children, child)
//
//			numOfNodes--
//			if numOfNodes == 1 {
//				last = true
//			}
//			child.addNetworkChild(numOfNodes, subnetBits, ipAddress, last, network)
//			return nil
//		}
//	}
//	return nil
//}

//func (this *treeNode) addMinMaxChild(minValue uint64, subnetBits int) error {
//	var maxValue uint64
//
//	if subnetBits > 8 && subnetBits < 16 {
//		changeableBits := 16 - subnetBits
//
//		valueToAdd := powInt(2, changeableBits) - 1
//
//		maxValue = minValue + valueToAdd
//	}
//	if subnetBits > 16 && subnetBits < 24 {
//		changeableBits := 24 - subnetBits
//
//		valueToAdd := powInt(2, changeableBits) - 1
//
//		maxValue = minValue + valueToAdd
//	}
//	if subnetBits > 24 && subnetBits < 32 {
//		changeableBits := 32 - subnetBits
//
//		valueToAdd := powInt(2, changeableBits) - 1
//
//		maxValue = minValue + valueToAdd
//	}
//
//	for _, child := range this.children {
//		if child.min <= minValue && child.max >= maxValue {
//			return nil
//		}
//	}
//
//	child := &treeNode{
//		min:      minValue,
//		max:      maxValue,
//		children: nil,
//		end:      true,
//	}
//
//	this.children = append(this.children, child)
//
//	return nil
//}

func powInt(x, y int) uint64 {
	return uint64(math.Pow(float64(x), float64(y)))
}

func (this *treeNode) Search(ipAddress string) bool {
	var fullValue uint32

	for i := 0; i < 4; i++ {
		var fragment uint64
		var index int

		for j := range ipAddress {
			if ipAddress[j] == '.' {
				index = j
				break
			}
			index = 0
		}

		if index == 0 {
			fragment, _ = strconv.ParseUint(ipAddress, 10, 32)
		} else {
			fragment, _ = strconv.ParseUint(ipAddress[:index], 10, 32)
		}
		fullValue = fullValue << 8
		fullValue += uint32(fragment)
	}

	current := this

	for i := 0; i < 32; i++ {
		next := uint32(fullValue << (32 - i) >> (31))
		child := current.children[next]

		if child == nil {
			return false
		}

		current = child
		if current.isBanned {
			return true
		}
	}
	return false
}
