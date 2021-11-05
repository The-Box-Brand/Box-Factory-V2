package boxes

import wr "github.com/mroth/weightedrand"

// Takes an array of choices and creates the chooser with them
func createChooser(choices []wr.Choice) (interface{}, error) {
	chooser, err := wr.NewChooser(choices...)
	if err != nil {
		return Attribute{}, err
	}
	return chooser.Pick(), nil
}

// Takes an array of attributes and turns it into an array of choices with the item as each attribute
// and the weight as the attributes rarity
func attributesToChoices(attributes []Attribute) []wr.Choice {
	var choices []wr.Choice
	for _, attribute := range attributes {
		choices = append(choices, wr.Choice{
			Weight: uint(attribute.Rarity),
			Item:   attribute,
		})
	}
	return choices
}

// Takes an array of attributes and chooses a random one based off their rarities weighted
func generateRandomAttribute(attributes []Attribute) (Attribute, error) {
	choices := attributesToChoices(attributes)
	chooserInterface, err := createChooser(choices)
	if err != nil {
		return Attribute{}, err
	}
	return chooserInterface.(Attribute), err
}
