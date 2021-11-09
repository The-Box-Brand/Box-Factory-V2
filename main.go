package main

import (
	"fmt"
	"io/fs"
	"log"
	"math/rand"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/The-Box-Brand/Box-Factory-V2/boxes"
)

// Loads attributes before main is ran
func init() {
	if err := loadAttributes(); err != nil {
		log.Fatal("Failed to load all attributes: " + err.Error())
	}
}

var attributesToNumber = make(map[string]map[string]int64)

func main() {
	rand.Seed(time.Now().UnixNano())

	attributesToNumber["background"] = make(map[string]int64)
	attributesToNumber["color"] = make(map[string]int64)
	attributesToNumber["cutout"] = make(map[string]int64)
	attributesToNumber["binding"] = make(map[string]int64)

	var strs []string
	var boxese []boxes.Box

retry:
	for x := 1; ; x++ {
		box, err := boxes.CreateBox()
		if err != nil {
			log.Fatal(err)
		}

		attributesToNumber["background"][box.Background.Name]++
		attributesToNumber["color"][box.Color.Name]++

		for _, cutout := range box.Cutouts {
			attributesToNumber["cutout"][cutout.Name]++
		}
		for _, binding := range box.Bindings {
			attributesToNumber["binding"][binding.Name]++
		}

		hash := box.CreateHash()
		for i := 0; i < len(strs); i++ {
			if strs[i] == hash {
				fmt.Println("got same box: here are all the unique boxes - " + fmt.Sprint(len(strs)))
				x--
				continue retry
				/* 	jsonBytes, _ := json.MarshalIndent(attributesToNumber, "", "	")
				log.Fatal(string(jsonBytes)) */
			}
		}
		err = box.SaveBox(x)
		if err != nil {
			log.Fatal(err)
		}

		strs = append(strs, hash)
		boxese = append(boxese, box)

		return
	}

}

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
