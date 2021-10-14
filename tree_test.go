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
		children: nil,
	})
	Assert(t).That(len(tree.children)).Equals(1)

	tree.Insert(IPNetwork16)
	Assert(t).That(tree.children[1]).Equals(&treeNode{
		value:    NetworkValue16,
		children: nil,
	})
	Assert(t).That(len(tree.children)).Equals(2)

	tree.Insert(IPNetwork24)
	Assert(t).That(tree.children[2]).Equals(&treeNode{
		value:    NetworkValue24,
		children: nil,
	})
	Assert(t).That(len(tree.children)).Equals(3)

	tree.Insert(IPNetwork32)
	Assert(t).That(tree.children[3]).Equals(&treeNode{
		value:    NetworkValue32,
		children: nil,
	})
	Assert(t).That(len(tree.children)).Equals(4)
}

func TestNonNetworks(t *testing.T) {
	tree := NewTreeNode()
	Assert(t).That(len(tree.children)).Equals(0)

	tree.Insert("3.144.0.0/13")
	Assert(t).That(tree.children[0]).Equals(&treeNode{
		value: "3",
		children: []*treeNode{
			{
				minValue: 144,
				maxValue: 151,
			},
		},
	})
	Assert(t).That(len(tree.children)).Equals(1)

	tree.Insert("54.233.0.0/18")
	Assert(t).That(tree.children[1]).Equals(&treeNode{
		value: "54.233",
		children: []*treeNode{
			{
				minValue: 0,
				maxValue: 63,
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
		value: "3",
		children: []*treeNode{
			{
				minValue: 144,
				maxValue: 151,
			},
		},
	})
	Assert(t).That(len(tree.children[0].children)).Equals(1)

	tree.Insert("3.152.0.0/13")
	Assert(t).That(tree.children[0]).Equals(&treeNode{
		value: "3",
		children: []*treeNode{
			{
				minValue: 144,
				maxValue: 151,
			},
			{
				minValue: 152,
				maxValue: 159,
			},
		},
	})
	Assert(t).That(len(tree.children[0].children)).Equals(2)

	tree.Insert("3.108.0.0/14")
	Assert(t).That(tree.children[0]).Equals(&treeNode{
		value: "3",
		children: []*treeNode{
			{
				minValue: 144,
				maxValue: 151,
			},
			{
				minValue: 152,
				maxValue: 159,
			},
			{
				minValue: 108,
				maxValue: 111,
			},
		},
	})
	Assert(t).That(len(tree.children[0].children)).Equals(3)
}

func TestChildAlreadyAdded(t *testing.T) {
	tree := NewTreeNode()
	Assert(t).That(len(tree.children)).Equals(0)

	tree.Insert("3.144.0.0/13")
	Assert(t).That(tree.children[0]).Equals(&treeNode{
		value: "3",
		children: []*treeNode{
			{
				minValue: 144,
				maxValue: 151,
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
	tree.Insert("52.144.229.64/26")
	Assert(t).That(tree.children[0]).Equals(&treeNode{
		value: "52.144.229",
		children: []*treeNode{
			{
				minValue: 64,
				maxValue: 127,
			},
		},
	})
	Assert(t).That(len(tree.children)).Equals(1)
	tree.Insert("54.168.0.0/16")
	Assert(t).That(tree.children[1]).Equals(&treeNode{
		value:    "54.168",
		children: nil,
	})
	Assert(t).That(len(tree.children)).Equals(2)
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func TestFindIPAddressWithMinMax(t *testing.T) {
	tree := NewTreeNode()
	Assert(t).That(len(tree.children)).Equals(0)

	tree.Insert("3.144.0.0/13")
	Assert(t).That(tree.children[0]).Equals(&treeNode{
		value: "3",
		children: []*treeNode{
			{
				minValue: 144,
				maxValue: 151,
			},
		},
	})
	Assert(t).That(len(tree.children)).Equals(1)

	tree.Insert("3.5.140.0/22")
	Assert(t).That(tree.children[1]).Equals(&treeNode{
		value: "3.5",
		children: []*treeNode{
			{
				minValue: 140,
				maxValue: 143,
			},
		},
	})
	Assert(t).That(len(tree.children)).Equals(2)

	tree.Insert("13.34.37.64/27")
	Assert(t).That(tree.children[2]).Equals(&treeNode{
		value: "13.34.37",
		children: []*treeNode{
			{
				minValue: 64,
				maxValue: 95,
			},
		},
	})
	Assert(t).That(len(tree.children)).Equals(3)

	tree.Insert("52.219.170.0/23")
	Assert(t).That(tree.children[3]).Equals(&treeNode{
		value: "52.219",
		children: []*treeNode{
			{
				minValue: 170,
				maxValue: 171,
			},
		},
	})
	Assert(t).That(len(tree.children)).Equals(4)

	tree.Insert("52.94.76.0/22")
	Assert(t).That(tree.children[4]).Equals(&treeNode{
		value: "52.94",
		children: []*treeNode{
			{
				minValue: 76,
				maxValue: 79,
			},
		},
	})
	Assert(t).That(len(tree.children)).Equals(5)

	tree.Insert("52.95.36.0/22")
	Assert(t).That(tree.children[5]).Equals(&treeNode{
		value: "52.95",
		children: []*treeNode{
			{
				minValue: 36,
				maxValue: 39,
			},
		},
	})
	Assert(t).That(len(tree.children)).Equals(6)

	tree.Insert("120.52.22.96/27")
	Assert(t).That(tree.children[6]).Equals(&treeNode{
		value: "120.52.22",
		children: []*treeNode{
			{
				minValue: 96,
				maxValue: 127,
			},
		},
	})
	Assert(t).That(len(tree.children)).Equals(7)

	tree.Insert("150.222.11.86/31")
	Assert(t).That(tree.children[7]).Equals(&treeNode{
		value: "150.222.11",
		children: []*treeNode{
			{
				minValue: 86,
				maxValue: 87,
			},
		},
	})
	Assert(t).That(len(tree.children)).Equals(8)

	tree.Insert("13.34.11.32/27")
	Assert(t).That(tree.children[8]).Equals(&treeNode{
		value: "13.34.11",
		children: []*treeNode{
			{
				minValue: 32,
				maxValue: 63,
			},
		},
	})
	Assert(t).That(len(tree.children)).Equals(9)

	tree.Insert("15.230.39.60/31")
	Assert(t).That(tree.children[9]).Equals(&treeNode{
		value: "15.230.39",
		children: []*treeNode{
			{
				minValue: 60,
				maxValue: 61,
			},
		},
	})
	Assert(t).That(len(tree.children)).Equals(10)

	exists := tree.Search("3.144.0.0")
	Assert(t).That(exists).Equals(true)
	exists = tree.Search("3.5.140.0")
	Assert(t).That(exists).Equals(true)
	exists = tree.Search("13.34.37.64")
	Assert(t).That(exists).Equals(true)
	exists = tree.Search("52.219.170.0")
	Assert(t).That(exists).Equals(true)
	exists = tree.Search("52.94.76.0")
	Assert(t).That(exists).Equals(true)
	exists = tree.Search("52.95.36.0")
	Assert(t).That(exists).Equals(true)
	exists = tree.Search("120.52.22.96")
	Assert(t).That(exists).Equals(true)
	exists = tree.Search("150.222.11.86")
	Assert(t).That(exists).Equals(true)
	exists = tree.Search("13.34.11.32")
	Assert(t).That(exists).Equals(true)
	exists = tree.Search("15.230.39.60")
	Assert(t).That(exists).Equals(true)
}

func TestFindNonExistentNetwork(t *testing.T) {
	tree := NewTreeNode()
	Assert(t).That(len(tree.children)).Equals(0)

	tree.Insert("3.144.0.0/13")
	Assert(t).That(tree.children[0]).Equals(&treeNode{
		value: "3",
		children: []*treeNode{
			{
				minValue: 144,
				maxValue: 151,
			},
		},
	})
	Assert(t).That(len(tree.children)).Equals(1)

	tree.Insert("3.5.140.0/22")
	Assert(t).That(tree.children[1]).Equals(&treeNode{
		value: "3.5",
		children: []*treeNode{
			{
				minValue: 140,
				maxValue: 143,
			},
		},
	})
	Assert(t).That(len(tree.children)).Equals(2)

	exists := tree.Search("3.152.0.0")
	Assert(t).That(exists).Equals(false)
	exists = tree.Search("3.5.144.0")
	Assert(t).That(exists).Equals(false)

}

func TestFindIPAddressWithoutMinMax(t *testing.T) {
	tree := NewTreeNode()
	Assert(t).That(len(tree.children)).Equals(0)

	tree.Insert(IPNetwork8)
	Assert(t).That(tree.children[0]).Equals(&treeNode{
		value:    "10",
		children: nil,
	})
	Assert(t).That(len(tree.children)).Equals(1)
	exists := tree.Search("10.0.0.0")
	Assert(t).That(exists).Equals(true)

	tree.Insert(IPNetwork16)
	Assert(t).That(tree.children[1]).Equals(&treeNode{
		value:    "54.168",
		children: nil,
	})
	Assert(t).That(len(tree.children)).Equals(2)
	exists = tree.Search("54.168.0.0")
	Assert(t).That(exists).Equals(true)

	tree.Insert(IPNetwork24)
	Assert(t).That(tree.children[2]).Equals(&treeNode{
		value:    "150.222.10",
		children: nil,
	})
	Assert(t).That(tree.children[2]).Equals(&treeNode{
		value:    "150.222.10",
		children: nil,
	})
	Assert(t).That(len(tree.children)).Equals(3)
	exists = tree.Search("150.222.10.0")
	Assert(t).That(exists).Equals(true)

	tree.Insert(IPNetwork32)
	Assert(t).That(tree.children[3]).Equals(&treeNode{
		value:    "52.93.126.244",
		children: nil,
	})
	Assert(t).That(len(tree.children)).Equals(4)
	exists = tree.Search("52.93.126.244")
	Assert(t).That(exists).Equals(true)
}

func TestFindWithAndWithOutMinMax(t *testing.T) {
	tree := NewTreeNode()
	Assert(t).That(len(tree.children)).Equals(0)

	tree.Insert(IPNetwork16) //54.168.0.0/16
	Assert(t).That(tree.children[0]).Equals(&treeNode{
		value:    "54.168",
		children: nil,
	})
	exists := tree.Search("54.168.0.0")
	Assert(t).That(exists).Equals(true)

	tree.Insert("3.144.0.0/13")
	Assert(t).That(tree.children[1]).Equals(&treeNode{
		value: "3",
		children: []*treeNode{
			{
				minValue: 144,
				maxValue: 151,
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

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func BenchmarkTreeTest(b *testing.B) {
	IpToSearch := "13.108.0.0/14"

	Tree := NewTreeNode()

	Tree.Insert("3.5.140.0/22")
	Tree.Insert("13.34.37.64/27")
	Tree.Insert("2.93.178.234/32")
	Tree.Insert("52.95.36.0/22")
	Tree.Insert("52.94.76.0/22")
	Tree.Insert("150.222.81.0/24")
	Tree.Insert("3.108.0.0/14")

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		_ = Tree.Search(IpToSearch)
	}
	b.ReportAllocs()
}
