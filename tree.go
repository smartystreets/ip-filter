package IPFilter

import "C"
import (
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
	if len(ipAddress) == 0 {
		return ErrInvalidIPAddress
	}
	position := strings.Index(ipAddress, "/")
	subnetBits, _ := strconv.Atoi(ipAddress[position+1:])
	ipAddress = ipAddress[:position]

	numericIP := parseIP(ipAddress)

	current := this

	for i := 0; i < subnetBits; i++ {
		next := uint32(numericIP << (i) >> (31))
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

func (this *treeNode) Search(ipAddress string) bool {
	if len(ipAddress) == 0{
		return false
	}

	numericIP := parseIP(ipAddress)

	current := this

	for i := 0; i < 32; i++ {
		nextBit := uint32(numericIP << (i) >> (31))
		child := current.children[nextBit]

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

func parseIP(ipAddress string) uint32 {
	var numericIP uint32
	for i := 0; i < 4; i++ {
		var fragment uint64
		var index int

		for j := range ipAddress {
			if ipAddress[j] == '.' {
				index = j
				break
			}
			continue
		}

		if index == 0 {
			fragment, _ = strconv.ParseUint(ipAddress, 10, 32)
		} else {
			fragment, _ = strconv.ParseUint(ipAddress[:index], 10, 32)
		}

		ipAddress = ipAddress[index+1:]

		numericIP = numericIP << 8
		numericIP += uint32(fragment)
	}
	return numericIP
}

//TODO: try and figure out how to kill it
//is it long enough is it too long?
//
