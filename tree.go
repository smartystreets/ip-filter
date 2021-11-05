package IPFilter

import "C"
import (
	"strconv"
	"strings"
	"unicode"
)

type IPFilter interface {
	Contains(string) bool
}

type treeNode struct {
	children []*treeNode
	isBanned bool
}

func New(addresses ...string) IPFilter {
	this := newNode()

	for _, item := range addresses {
		this.insertSubnetMask(item)
	}

	return this
}
func newNode() *treeNode {
	return &treeNode{children: make([]*treeNode, 2)}
}

func (this *treeNode) insertSubnetMask(ipAddress string) {
	if len(ipAddress) == 0 {
		return
	}

	if !isNumeric(ipAddress) {
		return
	}

	index := strings.Index(ipAddress, "/")
	if index == -1 {
		return
	}

	subnetBits, _ := strconv.Atoi(ipAddress[index+1:])
	ipAddress = ipAddress[:index]

	numericIP := parseIPAddress(ipAddress)
	if numericIP == 0 {
		return
	}

	current := this

	for i := 0; i < subnetBits; i++ {
		next := uint32(numericIP << (i) >> (31))
		child := current.children[next]

		if child == nil {
			child = newNode()
			current.children[next] = child
		}
		current = child
	}

	current.isBanned = true
}
func isNumeric(value string) bool {
	// TODO: check explicitly for '0' through '9' and also '.'
	nonLetter := func(c rune) bool { return unicode.IsLetter(c) }
	words := strings.FieldsFunc(value, nonLetter)
	return value == strings.Join(words, "")
}

func (this *treeNode) Contains(ipAddress string) bool {
	var numericIP uint32

	if len(ipAddress) == 0 {
		return false
	}

	if numericIP = parseIPAddress(ipAddress); numericIP == 0 {
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
func parseIPAddress(value string) uint32 {
	var numericIP uint32
	var count int

	for i := 0; i < 4; i++ {
		var fragment uint64
		var index int

		for x := range value {
			if value[x] == '.' {
				index = x
				count++
				break
			}
			continue
		}

		if index == 0 {
			fragment, _ = strconv.ParseUint(value, 10, 32)
		} else {
			fragment, _ = strconv.ParseUint(value[:index], 10, 32)
		}

		value = value[index+1:]

		if len(value) == 0 && count < 3 {
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
