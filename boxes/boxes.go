package boxes

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"image"
	"image/png"
	"math/rand"
	"os"
	"strings"
	"sync"

	"github.com/disintegration/imaging"
	wr "github.com/mroth/weightedrand"
	gim "github.com/ozankasikci/go-image-merge"
)

var Traits = make(map[string][]Attribute)

func CreateFactory() Factory {
	return Factory{RWMutex: sync.RWMutex{}}
}

func (factory *Factory) CreateSecretBox() Box {
	choice := rand.Intn(len(Traits["secret"]))

	var box Box
	box.Secret = Traits["secret"][choice]

	Traits["secret"] = append(Traits["secret"][:choice], Traits["secret"][choice+1:]...)

	return box
}

func (factory *Factory) CreateBox() (Box, error) {
	factory.Lock()
	defer factory.Unlock()

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

func CreateCustom(attrs []string) error {
	if len(attrs) < 1 {
		return fmt.Errorf("not enough attrs")
	}

	grids := []*gim.Grid{}
	for i := range attrs {
		grids = append(grids, &gim.Grid{ImageFilePath: attrs[i]})
	}

	rgba, _ := gim.New([]*gim.Grid{
		{
			ImageFilePath: attrs[0],
			Grids:         grids,
		},
	}, 1, 1).Merge()

	x2048 := imaging.Resize(rgba, 2048, 2048, imaging.NearestNeighbor)

	f, err := os.Create("./custom.png")
	if err != nil {
		return err
	}
	defer f.Close()

	err = png.Encode(f, x2048)
	if err != nil {
		return err
	}

	return nil
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

	// For each extra trait get a random attribute for it
	for i := 0; i < len(extras); i++ {
		attribute, err := generateRandomAttribute(Traits[extras[i]])
		if err != nil {
			return err
		}

		switch extras[i] {
		case "strap":
			if attribute.checkIfInside(box.Straps) {
				//i--
				continue
			}
			box.Straps = append(box.Straps, attribute)
		case "adhesive":
			if attribute.checkIfInside(box.Adhesives) {
				i--
				continue
			}
			box.Adhesives = append(box.Adhesives, attribute)
		case "cutout":
			box.Cutouts = append(box.Cutouts, attribute)
		case "label":
			box.Label = attribute
		case "state":
			box.State = attribute
		}
	}

	return
}

// Saves the box as a png, will need to do lots more in future
func (box Box) SaveAs(path string, isPNG bool) error {
	var img image.Image
	var err error

	if isPNG {
		img, err = box.GetPNG()
		if err != nil {
			return err
		}

		img = imaging.Resize(img, 2048, 2048, imaging.NearestNeighbor)
	} else {
		img, err = box.GetIMG(2048)
		if err != nil {
			return err
		}
	}

	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	err = png.Encode(f, img)
	if err != nil {
		return err
	}

	return err
}

func (box Box) save(path string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	err = png.Encode(f, box.IMG)
	if err != nil {
		return err
	}

	return err
}
func (box *Box) GetIMG(size int) (image.Image, error) {
	cond := box.Secret.ImagePath == ""
	rgba, err := gim.New([]*gim.Grid{
		{
			ImageFilePath: ternaryOperator(cond, box.Background.ImagePath, box.Secret.ImagePath).(string), // Switch to box.Color.ImagePath to have no background
			Grids:         ternaryOperator(cond, box.createGrids(false), []*gim.Grid{}).([]*gim.Grid),
		},
	}, 1, 1).Merge()

	if err != nil {
		return rgba, err
	}

	x2048 := imaging.Resize(rgba, size, size, imaging.NearestNeighbor)
	box.IMG = x2048

	return x2048, err
}

func (box *Box) GetPNG() (image.Image, error) {
	rgba, err := gim.New([]*gim.Grid{
		{
			ImageFilePath: box.Color.ImagePath,
			Grids:         box.createGrids(true),
		},
	}, 1, 1).Merge()
	if err != nil {
		return nil, err
	}

	box.IMG = rgba
	return rgba, err
}

// Create grid helper
func (box Box) createGrids(isPNG bool) (grid []*gim.Grid) {
	if !isPNG {
		grid = append(grid, &gim.Grid{
			ImageFilePath: "./Default/BACKGROUNDLINES.png",
		})
	}
	grid = append(grid, &gim.Grid{
		ImageFilePath: box.Color.ImagePath,
	}, &gim.Grid{
		ImageFilePath: "./Default/BOXLINES.png",
	})

	if box.State.ImagePath != "" {
		grid = append(grid, &gim.Grid{
			ImageFilePath: box.State.ImagePath,
		})
	}
	for _, cutout := range box.Cutouts {
		grid = append(grid, &gim.Grid{
			ImageFilePath: cutout.ImagePath,
		})
	}
	for _, adhesive := range box.Adhesives {
		grid = append(grid, &gim.Grid{
			ImageFilePath: adhesive.ImagePath,
		})
	}
	for _, strap := range box.Straps {
		grid = append(grid, &gim.Grid{
			ImageFilePath: strap.ImagePath,
		})
	}
	if box.Label.ImagePath != "" {
		grid = append(grid, &gim.Grid{
			ImageFilePath: box.Label.ImagePath,
		})
	}

	return
}

// Returns an array of extra traits that will then be used to add to the box
func generateExtras() (extras []string, err error) {
	// Gets the random number of extra traits the box will have
	chooserInterface, err := createChooser(NumberOfTraitsConfig)
	if err != nil {
		return extras, err
	}
	numberOfExtras := chooserInterface.(int)

	// Copies the config
	extraChoicesMap := make(map[string]wr.Choice)
	for key, val := range ExtrasConfig {
		extraChoicesMap[key] = val
	}

	for i := 0; i < numberOfExtras; i++ {

		// Don't allow duplicate extras
		for i := range extras {
			if extras[i] != "strap" && extras[i] != "adhesive" {
				delete(extraChoicesMap, extras[i])
			}
		}

		var extraChoices []wr.Choice
		for _, val := range extraChoicesMap {
			extraChoices = append(extraChoices, val)
		}

		chooserInterface, err := createChooser(extraChoices)
		if err != nil {
			return extras, err
		}

		extras = append(extras, chooserInterface.(string))
	}
	return extras, err
}

func (box Box) CreateHash() string {
	strapNames := attributesToNames(box.Straps)
	adhesiveNames := attributesToNames(box.Adhesives)
	cutoutNames := attributesToNames(box.Cutouts)

	boxFormat := "Background:%v|Color:%v|State:%v|Cutouts:%v|Adhesives:%v|Straps:%v|Label:%v"

	hash := md5.Sum([]byte(fmt.Sprintf(boxFormat, toRaw(box.Background.Name), toRaw(box.Color.Name), toRaw(box.State.Name), strings.Join(cutoutNames, ","), strings.Join(adhesiveNames, ","), strings.Join(strapNames, ","), toRaw(box.Label.Name))))
	return hex.EncodeToString(hash[:])
}
