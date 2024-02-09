package subnet

import (
	"encoding/json"
	"os"
	"fmt"
)
// SubnetNode represents a node in the subnet division tree.
type Subnet struct {
	Address uint32
	MaskLen uint32
	Parent  *Subnet `json:"-"`
	Left    *Subnet
	Right   *Subnet
	Labels  []string
}

// divide splits a subnet node into two subnets.
func (n *Subnet) Divide() {
    if n.MaskLen >= 32 {
        return // Cannot divide further
    }
    // No change for the left child; it starts at the same address as the parent subnet.
    n.Left = &Subnet{
        Address: n.Address,
        MaskLen: n.MaskLen + 1,
		Parent:  n,
    }

    // Calculate the starting address for the right child.
    // The offset is determined by 2^(32 - (n.MaskLen + 1)), 
    // which is the size of the new (smaller) subnets after division.
    offset := 1 << (32 - (n.MaskLen + 1))

    n.Right = &Subnet{
        Address: n.Address + uint32(offset), // Correctly calculate the offset for the right child.
        MaskLen: n.MaskLen + 1,
		Parent:  n,
    }
}

// merge combines two child subnets into their parent subnet.
func (n *Subnet) Join() {
    if n.Left == nil || n.Right == nil {
        return
    }
    // Assuming the caller ensures that n is the correct parent of Left and Right,
    // and they are adjacent, thus can be merged.
    n.Left = nil
    n.Right = nil
}
// findNode searches for a node with the specified address and mask length.
func (n *Subnet) Find(address uint32, maskLen uint32) *Subnet {
    
	if n.Address == address && n.MaskLen == maskLen {
		return n
	}
	if n.Left != nil {
		if found := n.Left.Find(address, maskLen); found != nil {
			return found
		}
	}
	if n.Right != nil {
		if found := n.Right.Find(address, maskLen); found != nil {
			return found
		}
	}
	return nil // Node not found
}

// iterate applies a function to each node in the tree.
func (n *Subnet) Iterate(f func(*Subnet)) {
	if n.Left == nil && n.Right == nil {
		f(n) // Apply function only if it's a leaf node
	} else {
		// If not a leaf node, recursively iterate over children
		if n.Left != nil {
			n.Left.Iterate(f)
		}
		if n.Right != nil {
			n.Right.Iterate(f)
		}
	}
}

//print
func (n *Subnet) Print() {
	n.Iterate(func(n *Subnet) {
		fmt.Printf("%s/%d\n", InetNtoa(n.Address), n.MaskLen)
	})
}
// SaveTree saves the subnet tree to a file in JSON format.
func SaveTree(root *Subnet, filename string) error {
    data, err := json.MarshalIndent(root, "", "  ")
    if err != nil {
        return err
    }
    return os.WriteFile(filename, data, 0644)
}

// LoadTree loads the subnet tree from a file in JSON format.
func LoadTree(filename string) (*Subnet, error) {
    data, err := os.ReadFile(filename)
    if err != nil {
        return nil, err
    }

    var root Subnet
    if err := json.Unmarshal(data, &root); err != nil {
        return nil, err
    }

    reconstructParent(&root, nil) // Reconstruct parent pointers

    return &root, nil
}

// reconstructParent helps to set the Parent field after loading from JSON.
func reconstructParent(node *Subnet, parent *Subnet) {
    if node != nil {
        node.Parent = parent
        reconstructParent(node.Left, node)
        reconstructParent(node.Right, node)
    }
}


