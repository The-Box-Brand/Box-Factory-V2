package boxes

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
