package ipfilter

import (
	"strconv"
	"strings"
)

type treeNode2 struct {
	children  []*treeNode2
	isAllowed bool
}

func New2(addresses ...string) WhiteList {
	this := newNode2()

	for _, item := range addresses {
		this.add(item)
	}

	return this
}

func newNode2() *treeNode2 {
	return &treeNode2{children: make([]*treeNode2, 2)}
}

func (this *treeNode2) add(subnetMask string) {
	if len(subnetMask) == 0 {
		return
	}

	index := strings.Index(subnetMask, subnetMaskSeparator)
	if index == -1 {
		return
	}

	subnetBits, _ := strconv.Atoi(subnetMask[index+1:])

	if subnetBits < 24 {
		return
	}

	baseIPAddress := subnetMask[:index]
	if !isNumeric(baseIPAddress) {
		return
	}

	numericIP := parseIPAddress(baseIPAddress)
	if numericIP == 0 {
		return
	}

	current := this
	for i := 0; i < subnetBits; i++ {
		nextBit := uint32(numericIP << i >> ipv4BitMask)
		child := current.children[nextBit]

		if child == nil {
			child = newNode2()
			current.children[nextBit] = child
		}
		current = child
	}

	current.isAllowed = true
}

func (this *treeNode2) Contains(ipAddress string) bool {
	if len(ipAddress) == 0 {
		return false
	}

	var numericIP uint32
	if numericIP = parseIPAddress(ipAddress); numericIP == 0 {
		return false
	}

	current := this
	for i := 0; i < ipv4BitCount; i++ {
		nextBit := uint32(numericIP << i >> ipv4BitMask)
		child := current.children[nextBit]

		if child == nil {
			return false
		}

		current = child
		if current.isAllowed {
			return true
		}
	}

	return false
}

func (this *treeNode2) Remove2(numericIP, i uint32, subnetBits int, current *treeNode2) bool {
	if numericIP == 0 {
		return false
	}
	if i == uint32(subnetBits) {
		if current.isAllowed == true {
			if current.children[0] == nil && current.children[1] == nil {
				return true
			}
			current.isAllowed = false
			return false
		}
	}

	nextBit := numericIP << i >> ipv4BitMask
	child := current.children[nextBit]

	if d := this.Remove2(numericIP, i+1, subnetBits, child); d == true {

		current.children[nextBit] = nil //
		if current.isAllowed == true || current.children[0] != nil || current.children[1] != nil {
			return false
		}
		return true
	}
	return false
}

func (this *treeNode2) Remove(ipAddress string) bool {
	var numericIP uint32

	index := strings.Index(ipAddress, subnetMaskSeparator)
	if index == -1 {
		return false
	}

	subnetBits, _ := strconv.Atoi(ipAddress[index+1:])
	baseIPAddress := ipAddress[:index]
	if !isNumeric(baseIPAddress) {
		return false
	}
	numericIP = parseIPAddress(baseIPAddress)
	return this.Remove2(numericIP, 0, subnetBits, this)
}
