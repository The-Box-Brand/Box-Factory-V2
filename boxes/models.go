package boxes

import (
	"image"
	"sync"
)

type Box struct {
	Background Attribute
	Color      Attribute
	State      Attribute
	Cutouts    []Attribute
	Adhesives  []Attribute
	Label      Attribute
	Straps     []Attribute

	Secret Attribute

	IMG image.Image
}

type Attribute struct {
	Name      string
	Rarity    int
	ImagePath string
}

type Factory struct {
	sync.RWMutex
}
