package boxes

type Box struct {
	Background Attribute
	Color      Attribute
	Binding    Attribute
}

type Attribute struct {
	Name      string
	Rarity    int
	ImagePath string
}
