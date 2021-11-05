package main

import (
	"fmt"
	"io/fs"
	"log"
	"math/rand"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/The-Box-Brand/Box-Factory-V2/boxes"
)

// Loads attributes before main is ran
func init() {
	loadAttributes()
}

func main() {
	rand.Seed(time.Now().UnixNano())

	for {
		box, err := boxes.CreateBox()
		if err != nil {
			log.Fatal(err)
		}

		box.SaveBox()
		time.Sleep(100 * time.Millisecond)
	}

}

func loadAttributes() {
	// Walk through every single file in this directory
	filepath.WalkDir("./", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// Make sure we are inside one of the trait folders
		pathSplit := strings.Split(path, `\`)
		if len(pathSplit) < 2 {
			return nil
		}

		// Folder name put to lowercase
		traitName := strings.ToLower(pathSplit[0])
		if traitName == "boxes" {
			return nil
		}

		// Name of image with .png removed
		artwork := strings.ReplaceAll(pathSplit[1], ".png", "")

		// Name of image splitted by ~
		artworkSplit := strings.Split(artwork, "~")
		fmt.Println(artworkSplit)
		artworkName := strings.ReplaceAll(artworkSplit[0], "_", " ")

		// Converting rarity string to integer
		artworkRarity, err := strconv.Atoi(artworkSplit[1])
		if err != nil {
			return err
		}

		fmt.Printf("Trait Type: %v\nArtwork Name: %v\nArtwork Rarity: %v\n\n", traitName, artworkName, artworkRarity)

		// Adding the attribute to the corresponding trait inside the map
		boxes.Traits[traitName] = append(boxes.Traits[traitName], boxes.Attribute{
			Name:   artworkName,
			Rarity: artworkRarity,

			// Cleaning path
			ImagePath: "./" + strings.ReplaceAll(path, `\`, "/"),
		})

		return nil
	})
}
