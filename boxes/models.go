package boxes

import "sync"

type Box struct {
	Background Attribute
	Color      Attribute
	Cutouts    []Attribute
	Adhesives  []Attribute
	Label      Attribute
	Straps     []Attribute
}

type Attribute struct {
	Name      string
	Rarity    int
	ImagePath string
}

type Factory struct {
	sync.RWMutex
}
