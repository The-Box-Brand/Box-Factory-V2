package main

import (
	"io/fs"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"

	"github.com/The-Box-Brand/Box-Factory-V2/boxes"
)

func loadAttributes() error {
	// Walk through every single file in this directory
	return filepath.WalkDir("./", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		char := `\`
		if runtime.GOOS == "darwin" {
			char = `/`
		}
		// Make sure we are inside one of the trait folders
		pathSplit := strings.Split(path, char)
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
		if len(artworkSplit) < 2 {
			return nil
		}

		artworkName := strings.ReplaceAll(artworkSplit[0], "_", " ")

		// Converting rarity string to integer
		artworkRarity, err := strconv.Atoi(artworkSplit[1])
		if err != nil {
			return err
		}

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
