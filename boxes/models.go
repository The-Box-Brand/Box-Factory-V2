package boxes

type Box struct {
	Background Attribute

	Color    Attribute
	Cutouts  []Attribute
	Bindings []Attribute
}

type Attribute struct {
	Name      string
	Rarity    int
	ImagePath string
}
