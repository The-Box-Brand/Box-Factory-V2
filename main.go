package main

import (
	"encoding/json"
	"fmt"
	"image"
	"image/draw"
	"image/png"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/The-Box-Brand/Box-Factory-V2/boxes"
	"github.com/disintegration/imaging"
)

// Loads attributes before main is ran
func init() {
	rand.Seed(time.Now().UnixNano())
	if err := loadAttributes(); err != nil {
		log.Fatal("Failed to load all attributes: " + err.Error())
	}
}

var attributesToNumber = make(map[string]map[string]int64)

func main() {
	createManyUnique(15)
	createCanvas()
	createTest()
}

func createTest() {
	box, err := boxes.CreateBox()
	if err != nil {
		log.Fatal(err)
	}

	box.SaveBox("test.png")
}

func createManyUnique(amount int) {

	attributesToNumber["background"] = make(map[string]int64)
	attributesToNumber["color"] = make(map[string]int64)
	attributesToNumber["cutout"] = make(map[string]int64)
	attributesToNumber["strap"] = make(map[string]int64)
	attributesToNumber["adhesive"] = make(map[string]int64)
	attributesToNumber["label"] = make(map[string]int64)

	var strs []string

retry:
	for x := 1; x < amount; x++ {
		box, err := boxes.CreateBox()
		if err != nil {
			log.Fatal(err)
		}

		attributesToNumber["background"][box.Background.Name]++
		attributesToNumber["color"][box.Color.Name]++

		for _, cutout := range box.Cutouts {
			attributesToNumber["cutout"][cutout.Name]++
		}
		for _, strap := range box.Straps {
			attributesToNumber["strap"][strap.Name]++
		}
		for _, adhesive := range box.Adhesives {
			attributesToNumber["adhesive"][adhesive.Name]++
		}
		if box.Label.ImagePath != "" {
			attributesToNumber["label"][box.Label.Name]++
		}

		hash := box.CreateHash()
		for i := 0; i < len(strs); i++ {
			if strs[i] == hash {
				fmt.Println("got same box: here are all the unique boxes - " + fmt.Sprint(len(strs)))
				x--
				continue retry

			}
		}
		err = box.SaveBox("./TBB/" + fmt.Sprint(x) + ".png")
		if err != nil {
			log.Fatal(err)
		}

		strs = append(strs, hash)
	}

	jsonBytes, _ := json.MarshalIndent(attributesToNumber, "", "	")
	log.Println(string(jsonBytes))
}
func createCanvas() {

	width := 1027
	height := 1027
	m := image.NewRGBA(image.Rect(0, 0, width, height))
	maxBoxesOnX := width/54 + 1
	maxBoxesOnY := height/45 + 1

	x := 0
	y := 0
	for i := 0; i < maxBoxesOnX*maxBoxesOnY; i++ {
		box, err := boxes.CreateBox()
		if err != nil {
			log.Fatal(err)
		}

		boxPNG, err := box.GetPNG()
		if err != nil {
			log.Fatal(err)
		}

		shifter := 36
		if y%2 != 0 {
			shifter = 64
		}

		draw.Draw(m, m.Bounds(), boxPNG, image.Point{-(x * 54) + shifter, -(y * 45) + 47}, draw.Over)

		x++
		if x == maxBoxesOnX {
			y++
			x = 0

		}
	}

	f, err := os.Create("canvas.png")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	x2048 := imaging.Resize(m, 2048, 2048, imaging.NearestNeighbor)
	err = png.Encode(f, x2048)
	if err != nil {
		log.Fatal(err)
	}
}
