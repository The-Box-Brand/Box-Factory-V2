package boxes

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"image/png"
	"os"
	"strings"

	"github.com/disintegration/imaging"
	wr "github.com/mroth/weightedrand"
	gim "github.com/ozankasikci/go-image-merge"
)

var Traits = make(map[string][]Attribute)

func CreateBox() (Box, error) {
	var box Box
	var err error

	// Required traits
	err = box.getBackground()
	if err != nil {
		return box, err
	}

	err = box.getColor()
	if err != nil {
		return box, err
	}

	// Extras
	err = box.getExtras()
	if err != nil {
		return box, err
	}

	return box, err
}

func (box *Box) getBackground() (err error) {
	box.Background, err = generateRandomAttribute(Traits["background"])
	if err != nil {
		return err
	}
	return nil
}

func (box *Box) getColor() (err error) {
	box.Color, err = generateRandomAttribute(Traits["color"])
	if err != nil {
		return err
	}
	return nil
}

func (box *Box) getExtras() (err error) {
	// Extra traits
	extras, err := generateExtras()
	if err != nil {
		return err
	}

	fmt.Println(extras)
	// For each extra trait get a random attribute for it
	for _, extra := range extras {
		attribute, err := generateRandomAttribute(Traits[extra])
		if err != nil {
			return err
		}

		switch extra {
		case "binding":
			box.Bindings = append(box.Bindings, attribute)
		case "cutout":
			box.Cutouts = append(box.Cutouts, attribute)
		}
	}

	return
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

	x2048 := imaging.Resize(rgba, 2048, 2048, imaging.NearestNeighbor)
	err = png.Encode(f, x2048)
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
	for _, cutout := range box.Cutouts {
		grid = append(grid, &gim.Grid{
			ImageFilePath: cutout.ImagePath,
		})
	}
	for _, binding := range box.Bindings {
		grid = append(grid, &gim.Grid{
			ImageFilePath: binding.ImagePath,
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

	extraChoicesMap := make(map[string]wr.Choice)
	for key, val := range ExtrasConfig {
		extraChoicesMap[key] = val
	}

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
			if extras[i] != "binding" {
				delete(extraChoicesMap, extras[i])
			}
		}

		extras = append(extras, chooserInterface.(string))
	}
	return extras, err
}

func (box Box) CreateHash() string {
	bindingNames := attributesToNames(box.Bindings)
	cutoutNames := attributesToNames(box.Cutouts)

	boxFormat := "Background:%v|Color:%v|Cutouts:%v|Bindings:%v"

	hash := md5.Sum([]byte(fmt.Sprintf(boxFormat, toRaw(box.Background.Name), toRaw(box.Color.Name), strings.Join(cutoutNames, ","), strings.Join(bindingNames, ","))))
	return hex.EncodeToString(hash[:])
}
