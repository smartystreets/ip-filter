package IPFilter

import "C"
import (
	"strconv"
	"strings"
	"unicode"
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
	var index int
	var numericIP uint32

	if len(ipAddress) == 0 {
		return ErrInvalidIPAddress
	}

	if !NumbersOnly(ipAddress) {
		return ErrInvalidIPAddress
	}

	if index = strings.Index(ipAddress, "/"); index == -1 {
		return ErrInvalidIPAddress
	}

	subnetBits, _ := strconv.Atoi(ipAddress[index+1:])
	ipAddress = ipAddress[:index]

	if numericIP = parseIP(ipAddress); numericIP == 0 {
		return ErrInvalidIPAddress
	}

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
	var numericIP uint32

	if len(ipAddress) == 0 {
		return false
	}

	if numericIP = parseIP(ipAddress); numericIP == 0 {
		return false
	}

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
	var count int

	for i := 0; i < 4; i++ {
		var fragment uint64
		var index int

		for j := range ipAddress {
			if ipAddress[j] == '.' {
				index = j
				count++
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

		if len(ipAddress) == 0 && count < 3 {
			return 0
		}
		numericIP = numericIP << 8
		numericIP += uint32(fragment)
	}

	if count > 3 || count < 3 {
		return 0
	}

	return numericIP
}

func NumbersOnly(IPAddress string) bool {
	nonLetter := func(c rune) bool { return unicode.IsLetter(c) }
	words := strings.FieldsFunc(IPAddress, nonLetter)
	return IPAddress == strings.Join(words, "")
}
