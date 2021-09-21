package IPFilter

// Sort and store the data being passed in

type RadixTree struct {
	tree []*node //Made this a slice of pointers to make it easier to read through?
}
//TODO: Fix this struct, the node needs to have a different setup
type node struct {
	isLeaf bool // Checks to see if it is already a leaf
	//Example uses a hashmap, but I don't think we will want to use a hashmap.
	//HashMap<Character, Edge> edges (This will store the character as well as its associated edges)
	//Needs a label for the address
	label string
	next *node
}

//Insert
/*
Will take in an IP Address as a string and will insert it into the tree accordingly
*/
func (tree *RadixTree) Insert(IPAddress string) {}

//Delete
/*
Will take in an IPAddress as a string and will remove it from the tree/list of nodes
*/
func(tree *RadixTree) Delete(IPAddress string){}

//Search
/*
Will search the tree for a specific IPAddress within our nodes
*/
func(tree *RadixTree) Search(IPAddress string){}