package boxes

func (attribute Attribute) checkIfInside(attributes []Attribute) bool {
	for i := range attributes {
		if attributes[i].ImagePath == attribute.ImagePath {
			return true
		}
	}
	return false
}
