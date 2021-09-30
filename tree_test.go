package IPFilter

import (
	"reflect"
	"testing"
)

func TestNetworks(t *testing.T) {
	tree := NewTreeNode() //Creates the blank tree
	Assert(t).That(len(tree.children)).Equals(0)

	tree.Insert(IPNetwork8)
	Assert(t).That(tree.children[0]).Equals(&treeNode{
		value:    NetworkValue8,
		network:  true,
		children: nil,
	})
	Assert(t).That(len(tree.children)).Equals(1)

	tree.Insert(IPNetwork16)
	Assert(t).That(tree.children[1]).Equals(&treeNode{
		value:    NetworkValue16,
		network:  true,
		children: nil,
	})
	Assert(t).That(len(tree.children)).Equals(2)

	tree.Insert(IPNetwork24)
	Assert(t).That(tree.children[2]).Equals(&treeNode{
		value:    NetworkValue24,
		network:  true,
		children: nil,
	})
	Assert(t).That(len(tree.children)).Equals(3)

	tree.Insert(IPNetwork32)
	Assert(t).That(tree.children[3]).Equals(&treeNode{
		value:    NetworkValue32,
		network:  true,
		children: nil,
	})
	Assert(t).That(len(tree.children)).Equals(4)
}

func TestNonNetworks(t *testing.T) {
	tree := NewTreeNode() //Creates the blank tree
	Assert(t).That(len(tree.children)).Equals(0)

	tree.Insert("3.144.0.0/13")
	Assert(t).That(tree.children[0]).Equals(&treeNode{
		value:   "3",
		network: false,
		children: []*treeNode{
			{
				value:   "144",
				network: false,
				children: []*treeNode{
					{
						value:   "0",
						network: false,
						children: []*treeNode{
							{
								value:    "0",
								network:  false,
								children: nil,
							},
						},
					},
				},
			},
		},
	})
	Assert(t).That(len(tree.children)).Equals(1)
}

func TestChildAlreadyAdded(t *testing.T) {
	tree := NewTreeNode() //Creates the blank tree
	Assert(t).That(len(tree.children)).Equals(0)

	tree.Insert("3.144.0.0/13")
	Assert(t).That(tree.children[0]).Equals(&treeNode{
		value:   "3",
		network: false,
		children: []*treeNode{
			{
				value:   "144",
				network: false,
				children: []*treeNode{
					{
						value:   "0",
						network: false,
						children: []*treeNode{
							{
								value:    "0",
								network:  false,
								children: nil,
							},
						},
					},
				},
			},
		},
	})
}

func TestFindNonNetworkIPAddress(t *testing.T) {
	tree := NewTreeNode() //Creates the blank tree
	Assert(t).That(len(tree.children)).Equals(0)

	tree.Insert("3.144.0.0/13")
	Assert(t).That(tree.children[0]).Equals(treeNode{
		value:   "3",
		network: false,
		children: []*treeNode{
			{
				value:   "144",
				network: false,
				children: []*treeNode{
					{
						value:   "0",
						network: false,
						children: []*treeNode{
							{
								value:    "0",
								network:  false,
								children: nil,
							},
						},
					},
				},
			},
		},
	})
	Assert(t).That(len(tree.children)).Equals(1)

	exists := tree.Search("3.144.0.0")

	Assert(t).That(exists).Equals(true)
}

const (
	IPNetwork8  = "10.0.0.0/8"
	IPNetwork16 = "54.168.0.0/16"
	IPNetwork24 = "150.222.10.0/24"
	IPNetwork32 = "52.93.126.244/32"

	NetworkValue8  = "10"
	NetworkValue16 = "54.168"
	NetworkValue24 = "150.222.10"
	NetworkValue32 = "52.93.126.244"
)

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
type That struct{ t *testing.T }
type Assertion struct {
	*testing.T
	actual interface{}
}

func Assert(t *testing.T) *That                       { return &That{t: t} }
func (this *That) That(actual interface{}) *Assertion { return &Assertion{T: this.t, actual: actual} }

func (this *Assertion) IsNil() {
	this.Helper()
	if this.actual != nil && !reflect.ValueOf(this.actual).IsNil() {
		this.Equals(nil)
	}
}
func (this *Assertion) Equals(expected interface{}) {
	this.Helper()
	if !reflect.DeepEqual(this.actual, expected) {
		this.Errorf("\nExpected: %#v\nActual:   %#v", expected, this.actual)
	}
}
