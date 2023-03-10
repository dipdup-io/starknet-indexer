package storage

// ClassType -
type ClassType uint64

// class types
const (
	ClassTypeERC20 ClassType = 1 << iota
	ClassTypeERC721
	ClassTypeERC721Metadata
	ClassTypeERC721Receiver
	ClassTypeERC1155
	ClassTypeProxy
	ClassTypeArgentX0
	ClassTypeArgentX
)

// Set -
func (ct *ClassType) Set(types ...ClassType) {
	for i := range types {
		*ct |= types[i]
	}
}

// Is -
func (ct ClassType) Is(typ ClassType) bool {
	return ct&typ > 0
}

// OneOf -
func (ct ClassType) OneOf(typ ...ClassType) bool {
	for i := range typ {
		if ct&typ[i] > 0 {
			return true
		}
	}
	return false
}

// NewClassType -
func NewClassType(interfaces ...string) ClassType {
	var ct ClassType

	for i := range interfaces {
		switch interfaces[i] {
		case "erc20":
			ct.Set(ClassTypeERC20)
		case "erc721":
			ct.Set(ClassTypeERC721)
		case "erc721_metadata":
			ct.Set(ClassTypeERC721Metadata)
		case "erc721_receiver":
			ct.Set(ClassTypeERC721Receiver)
		case "proxy":
			ct.Set(ClassTypeProxy)
		case "erc1155":
			ct.Set(ClassTypeERC1155)
		case "argentx":
			ct.Set(ClassTypeArgentX)
		case "argentx_0":
			ct.Set(ClassTypeArgentX0)
		}
	}

	return ct
}
