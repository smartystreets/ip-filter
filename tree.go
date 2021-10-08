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

func (this *treeNode) Insert(IPAddress string) error {
	if len(IPAddress) == 0 {
		return ErrInvalidIPAddress
	}

	var IPFragment string
	if position := strings.Index(IPAddress, "/"); position != -1 {

		subnetBits, _ := strconv.Atoi(IPAddress[position+1:])

		if subnetBits == 8 {
			IPFragment = IPAddress[:strings.Index(IPAddress, ".")]
			this.addNetworkChild(0, 0, IPFragment)
			return nil
		}
		if subnetBits == 16 {
			index := findIndex(IPAddress, 16)
			IPFragment = IPAddress[:index]
			this.addNetworkChild(0, 0, IPFragment)
			return nil
		}
		if subnetBits == 24 {
			index := findIndex(IPAddress, 24)
			IPFragment = IPAddress[:index]
			this.addNetworkChild(0, 0, IPFragment)
			return nil
		}
		if subnetBits == 32 {
			IPFragment = IPAddress[:position]
			this.addNetworkChild(0, 0, IPFragment)
			return nil
		}
		IPAddress = IPAddress[:position]

		if subnetBits > 8 && subnetBits < 16 {
			firstSection := strings.Index(IPAddress, ".")
			secondSection := findIndex(IPAddress, 16)

			minValue, _ := strconv.ParseUint(IPAddress[firstSection+1:secondSection], 10, 64)

			changeableBits := 16 - subnetBits

			valueToAdd := powInt(2, changeableBits) - 1

			maxValue := minValue + valueToAdd

			IPFragment = IPAddress[:strings.Index(IPAddress, ".")]

			this.addNetworkChild(maxValue, minValue, IPFragment)
			return nil
		}
		if subnetBits > 16 && subnetBits < 24 {
			firstSection := findIndex(IPAddress, 16)
			secondSection := findIndex(IPAddress, 24)
			minValue, _ := strconv.ParseUint(IPAddress[firstSection+1:secondSection], 10, 64)

			changeableBits := 24 - subnetBits

			valueToAdd := powInt(2, changeableBits) - 1

			maxValue := minValue + valueToAdd

			IPFragment = IPAddress[:findIndex(IPAddress, 16)]

			this.addNetworkChild(maxValue, minValue, IPFragment)
			return nil
		}
		if subnetBits > 24 && subnetBits < 32 {
			minValue, _ := strconv.ParseUint(IPAddress[findIndex(IPAddress, 24)+1:], 10, 64)

			changeableBits := 32 - subnetBits

			valueToAdd := powInt(2, changeableBits) - 1

			maxValue := minValue + valueToAdd

			IPFragment = IPAddress[:findIndex(IPAddress, 24)]

			this.addNetworkChild(maxValue, minValue, IPFragment)
			return nil
		}

	}
	return nil
}

func (this *treeNode) addNetworkChild(maxValue, minValue uint64, IPFragment string) error {

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
func (this *treeNode) addMinMaxChild(maxValue, minValue uint64) error {

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
	var childRange string
	var fragment string
	if len(IPAddress) == 0 {
		return false
	}

	indexes := findIndexes(IPAddress)

	for i := 3; i >= 0; i-- {

		if i == 3 {
			fragment = IPAddress
		} else {
			fragment = IPAddress[:indexes[i]]
		}

		for _, child := range this.children {

			if child.value != fragment {
				continue
			}

			if child.children == nil {
				return true
			}

			firstSection := len(fragment) + 1
			var secondSection int

			for t, x := range indexes {
				if t == 2 {
					secondSection = len(IPAddress)
					break
				}
				if firstSection > x && firstSection < indexes[t+1] {
					secondSection = indexes[t+1]
					break
				}
			}

			childRange = IPAddress[firstSection:secondSection]

			parsed, err := strconv.ParseUint(childRange, 10, 32)

			if err != nil {
				return false //TODO: return error?
			}

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
