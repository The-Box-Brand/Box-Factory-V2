package boxes

import (
	"image/png"
	"os"

	wr "github.com/mroth/weightedrand"
	gim "github.com/ozankasikci/go-image-merge"
)

var Traits = make(map[string][]Attribute)

func CreateBox() (Box, error) {
	var box Box
	var err error

	// Required traits
	box.Background, err = generateRandomAttribute(Traits["background"])
	if err != nil {
		return box, err
	}

	box.Color, err = generateRandomAttribute(Traits["color"])
	if err != nil {
		return box, err
	}

	// Extra traits
	extras, err := generateExtras()
	if err != nil {
		return box, err
	}

	// For each extra trait get a random attribute for it
	for _, extra := range extras {
		attribute, err := generateRandomAttribute(Traits[extra])
		if err != nil {
			return box, err
		}

		switch extra {
		case "binding":
			box.Binding = attribute
		}
	}

	return box, err
}

// Saves the box as a png, will need to do lots more in future
func (box Box) SaveBox() error {

	rgba, err := gim.New([]*gim.Grid{
		{
			ImageFilePath: box.Background.ImagePath,
			Grids:         box.createGrids(),
		},
	}, 1, 1).Merge()

	if err != nil {
		return err
	}
	f, err := os.Create("image.png")
	if err != nil {
		return err
	}

	err = png.Encode(f, rgba)
	if err != nil {
		return err
	}

	return err
}

// Create grid helper
func (box Box) createGrids() (grid []*gim.Grid) {
	grid = append(grid, &gim.Grid{
		ImageFilePath: box.Color.ImagePath,
	})
	if box.Binding.ImagePath != "" {
		grid = append(grid, &gim.Grid{
			ImageFilePath: box.Binding.ImagePath,
		})
	}

	return
}

// Returns an array of extra traits that will then be used to add to the box
func generateExtras() (extras []string, err error) {
	chooserInterface, err := createChooser(NumberOfTraitsConfig)
	if err != nil {
		return extras, err
	}
	numberOfExtras := chooserInterface.(int)

	extraChoicesMap := ExtrasConfig

	for i := 0; i < numberOfExtras; i++ {

		var extraChoices []wr.Choice
		for _, val := range extraChoicesMap {
			extraChoices = append(extraChoices, val)
		}

		chooserInterface, err := createChooser(extraChoices)
		if err != nil {
			return extras, err
		}

		// Don't allow duplicate extras
		for i := range extras {
			delete(extraChoicesMap, extras[i])
		}

		extras = append(extras, chooserInterface.(string))
	}
	return extras, err
}
