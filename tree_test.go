package IPFilter

import (
	"reflect"
	"testing"
)

func TestNetworks(t *testing.T) {
	tree := NewTreeNode()
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
	tree := NewTreeNode()
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

	tree.Insert("54.233.0.0/18")
	Assert(t).That(tree.children[1]).Equals(&treeNode{
		value:   "54",
		network: false,
		children: []*treeNode{
			{
				value:   "233",
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
	},
	)
	Assert(t).That(len(tree.children)).Equals(2)
}

func TestMultipleChildren(t *testing.T) {
	tree := NewTreeNode()
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
	Assert(t).That(len(tree.children[0].children)).Equals(1)

	tree.Insert("3.145.0.0/13")
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
			{
				value:   "145",
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
	Assert(t).That(len(tree.children[0].children)).Equals(2)
}

func TestMultiMultiLevelChildren(t *testing.T) {
	tree := NewTreeNode()
	Assert(t).That(len(tree.children)).Equals(0)

	tree.Insert("52.144.218.60/26")
	tree.Insert("52.144.218.61/26")
	tree.Insert("52.144.218.62/26")
	tree.Insert("52.144.218.63/26")
	tree.Insert("52.144.218.64/26")
	Assert(t).That(tree.children[0]).Equals(&treeNode{
		value:   "52",
		network: false,
		children: []*treeNode{
			{
				value:   "144",
				network: false,
				children: []*treeNode{
					{
						value:   "218",
						network: false,
						children: []*treeNode{
							{
								value:    "60",
								network:  false,
								children: nil,
							},
							{
								value:    "61",
								network:  false,
								children: nil,
							},
							{
								value:    "62",
								network:  false,
								children: nil,
							},
							{
								value:    "63",
								network:  false,
								children: nil,
							},
							{
								value:    "64",
								network:  false,
								children: nil,
							},
						},
					},
				},
			},
		},
	})
	Assert(t).That(len(tree.children[0].children[0].children[0].children)).Equals(5)
	tree.Insert("52.144.227.192/26")
	tree.Insert("52.144.229.64/26")
	tree.Insert("52.144.225.128/26")
	tree.Insert("52.144.197.192/26")
	tree.Insert("52.144.199.128/26")
	Assert(t).That(tree.children[0]).Equals(&treeNode{
		value:   "52",
		network: false,
		children: []*treeNode{
			{
				value:   "144",
				network: false,
				children: []*treeNode{
					{
						value:   "218",
						network: false,
						children: []*treeNode{
							{
								value:    "60",
								network:  false,
								children: nil,
							},
							{
								value:    "61",
								network:  false,
								children: nil,
							},
							{
								value:    "62",
								network:  false,
								children: nil,
							},
							{
								value:    "63",
								network:  false,
								children: nil,
							},
							{
								value:    "64",
								network:  false,
								children: nil,
							},
						},
					},
					{
						value:   "227",
						network: false,
						children: []*treeNode{
							{
								value:    "192",
								network:  false,
								children: nil,
							},
						},
					},
					{
						value:   "229",
						network: false,
						children: []*treeNode{
							{
								value:    "64",
								network:  false,
								children: nil,
							},
						},
					},
					{
						value:   "225",
						network: false,
						children: []*treeNode{
							{
								value:    "128",
								network:  false,
								children: nil,
							},
						},
					},
					{
						value:   "197",
						network: false,
						children: []*treeNode{
							{
								value:    "192",
								network:  false,
								children: nil,
							},
						},
					},
					{
						value:   "199",
						network: false,
						children: []*treeNode{
							{
								value:    "128",
								network:  false,
								children: nil,
							},
						},
					},
				},
			},
		},
	})
	Assert(t).That(len(tree.children[0].children[0].children)).Equals(6)
}

func TestChildAlreadyAdded(t *testing.T) {
	tree := NewTreeNode()
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

	tree.Insert("3.144.0.0/13")
	Assert(t).That(len(tree.children)).Equals(1)
}

func TestErrors(t *testing.T) {
	tree := NewTreeNode()
	Assert(t).That(len(tree.children)).Equals(0)
	err := tree.Insert("")
	Assert(t).That(err).Equals(ErrInvalidIPAddress)
}

func TestNetworkAndNonNetworkChildrenAdded(t *testing.T) {
	tree := NewTreeNode()
	Assert(t).That(len(tree.children)).Equals(0)
	tree.Insert("52.144.218.60/26")
	Assert(t).That(tree.children[0]).Equals(&treeNode{
		value:   "52",
		network: false,
		children: []*treeNode{
			{
				value:   "144",
				network: false,
				children: []*treeNode{
					{
						value:   "218",
						network: false,
						children: []*treeNode{
							{
								value:    "60",
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
	tree.Insert("54.168.0.0/16")
	Assert(t).That(tree.children[1]).Equals(&treeNode{
		value:    "54.168",
		network:  true,
		children: nil,
	})
	Assert(t).That(len(tree.children)).Equals(2)
}

func TestFindNonNetworkIPAddress(t *testing.T) {
	tree := NewTreeNode()
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

	exists := tree.Search("3.144.0.0")
	Assert(t).That(exists).Equals(true)
}

func TestFindNetworkIPAddress(t *testing.T) {
	tree := NewTreeNode()
	Assert(t).That(len(tree.children)).Equals(0)

	tree.Insert(IPNetwork8) //10.0.0.0/8
	Assert(t).That(tree.children[0]).Equals(&treeNode{
		value:    "10",
		network:  true,
		children: nil,
	})
	Assert(t).That(len(tree.children)).Equals(1)
	exists := tree.Search("10.0.0.0")
	Assert(t).That(exists).Equals(true)

	tree.Insert(IPNetwork16) //54.168.0.0/16
	Assert(t).That(tree.children[1]).Equals(&treeNode{
		value:    "54.168",
		network:  true,
		children: nil,
	})
	Assert(t).That(len(tree.children)).Equals(2)
	exists = tree.Search("54.168.0.0")
	Assert(t).That(exists).Equals(true)

	tree.Insert(IPNetwork24) //150.222.10.0/24
	Assert(t).That(tree.children[2]).Equals(&treeNode{
		value:    "150.222.10",
		network:  true,
		children: nil,
	})
	Assert(t).That(len(tree.children)).Equals(3)
	exists = tree.Search("150.222.10.0")
	Assert(t).That(exists).Equals(true)

	tree.Insert(IPNetwork32) //52.93.126.244
	Assert(t).That(tree.children[3]).Equals(&treeNode{
		value:    "52.93.126.244",
		network:  true,
		children: nil,
	})
	Assert(t).That(len(tree.children)).Equals(4)
	exists = tree.Search("52.93.126.244")
	Assert(t).That(exists).Equals(true)
}

func TestFindNetworkAndNonNetwork(t *testing.T) {
	tree := NewTreeNode()
	Assert(t).That(len(tree.children)).Equals(0)

	tree.Insert(IPNetwork16) //54.168.0.0/16
	Assert(t).That(tree.children[0]).Equals(&treeNode{
		value:    "54.168",
		network:  true,
		children: nil,
	})
	exists := tree.Search("54.168.0.0")
	Assert(t).That(exists).Equals(true)

	tree.Insert("3.144.0.0/13")
	Assert(t).That(tree.children[1]).Equals(&treeNode{
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
	exists = tree.Search("3.144.0.0")
	Assert(t).That(exists).Equals(true)
	Assert(t).That(len(tree.children)).Equals(2)
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
