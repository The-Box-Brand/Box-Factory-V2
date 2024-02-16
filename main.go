package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"image"
	"image/draw"
	"image/gif"
	"image/png"
	"io/fs"
	"log"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/The-Box-Brand/Box-Factory-V2/boxes"
)

type miniFactory struct {
	sync.RWMutex
	attributesToNumber map[string]map[string]int64
	factory            boxes.Factory
	secrets            map[int]bool
	uniques            map[string]bool
	duration           time.Duration
}

// Loads attributes before main is ran
func init() {
	if err := loadAttributes(); err != nil {
		log.Fatal("Failed to load all attributes: " + err.Error())
	}

	// Uncomment to get constants for making custom boxes
	for key, attributesArray := range boxes.Traits {
		for _, attribute := range attributesArray {
			fmt.Printf(`%v = "%v"`, strings.ToUpper(strings.ReplaceAll(http.CanonicalHeaderKey(strings.ReplaceAll(key+"-"+attribute.Name, " ", "-")), "-", "_")), attribute.ImagePath)
			fmt.Println()
		}
	}

}

func main() {
	boxes.CreateCustom([]string{
		boxes.BACKGROUND_BG_BRIGHT_BLUE,
		boxes.BACKGROUND_LINES,
		boxes.COLOR_BRIGHT_GREEN,
		boxes.BOX_LINES,
		boxes.STRAP_DOUBLE_STRAPPED_RIGHT,
		boxes.STRAP_DOUBLE_STRAPPED_LEFT,
		boxes.STATE_WET,
		//boxes.STATE_BEACON,
	})

	createGIF(10)

	createCanvas()

	mf := miniFactory{}

	mf.createManyUnique(100)

	fmt.Println(mf.duration)
}

func createTest() {
	factory := boxes.CreateFactory()
	box, err := factory.CreateBox()
	if err != nil {
		log.Fatal(err)
	}

	box.SaveAs("./TBB/test.png", false)
}

func createGIF(framesNum int) {
	factory := boxes.CreateFactory()
	var frames []*image.Paletted

	wg := sync.WaitGroup{}
	wg.Add(framesNum)
	for i := 0; i < framesNum; i++ {
		go func() {
			defer wg.Done()

			box, err := factory.CreateBox()
			if err != nil {
				log.Fatal(err)
			}

			img, err := box.GetIMG(1080)
			if err != nil {
				log.Fatal(err)
			}

			t1 := time.Now()
			buf := bytes.Buffer{}
			err = gif.Encode(&buf, img, nil)
			if err != nil {
				log.Fatal(err)
			}

			tmpimg, err := gif.Decode(&buf)
			if err != nil {
				log.Fatal(err)
			}

			/* palettedImage := image.NewPaletted(img.Bounds(), palette.Plan9)
			draw.Draw(palettedImage, palettedImage.Rect, img, img.Bounds().Min, draw.Over)
			*/
			frames = append(frames, tmpimg.(*image.Paletted))
			fmt.Println(time.Since(t1))
		}()
	}

	delays := make([]int, framesNum)
	for j := range delays {
		delays[j] = 25
	}

	f, err := os.Create("./TBB/test.gif")
	if err != nil {
		log.Fatal(err)
	}

	wg.Wait()
	err = gif.EncodeAll(f, &gif.GIF{Image: frames, Delay: delays, LoopCount: 0})
	if err != nil {
		log.Fatal(err)
	}

}

func (mf *miniFactory) createUnique(num int, wg *sync.WaitGroup) {
	mf.Lock()
	defer mf.Unlock()

retry:
	var box boxes.Box
	var err error
	if mf.secrets[num] {
		box = mf.factory.CreateSecretBox()
		fmt.Println(num, box.Secret.Name)
		mf.attributesToNumber["secret"][box.Secret.Name]++
	} else {
		t1 := time.Now()
		box, err = mf.factory.CreateBox()
		if err != nil {
			log.Fatal(err)
		}
		mf.duration += time.Since(t1)
		mf.attributesToNumber["background"][box.Background.Name]++
		mf.attributesToNumber["color"][box.Color.Name]++

		for _, cutout := range box.Cutouts {
			mf.attributesToNumber["cutout"][cutout.Name]++
		}
		for _, strap := range box.Straps {
			mf.attributesToNumber["strap"][strap.Name]++
		}
		for _, adhesive := range box.Adhesives {
			mf.attributesToNumber["adhesive"][adhesive.Name]++
		}
		if box.Label.ImagePath != "" {
			mf.attributesToNumber["label"][box.Label.Name]++
		}
		if box.State.ImagePath != "" {
			fmt.Println(num, box.State.Name)
			mf.attributesToNumber["state"][box.State.Name]++
		}

		hash := box.CreateHash()

		if _, ok := mf.uniques[hash]; ok {
			goto retry
		}

		mf.uniques[hash] = true

	}

	go func() {
		defer wg.Done()
		box.SaveAs("./TBB/"+fmt.Sprint(num)+".png", false)
	}()

	/* 	if err != nil {
		log.Fatal(err)
	} */

}
func (mf *miniFactory) createManyUnique(amount int) {
	mf.attributesToNumber = make(map[string]map[string]int64)
	mf.attributesToNumber["background"] = make(map[string]int64)
	mf.attributesToNumber["color"] = make(map[string]int64)
	mf.attributesToNumber["cutout"] = make(map[string]int64)
	mf.attributesToNumber["strap"] = make(map[string]int64)
	mf.attributesToNumber["adhesive"] = make(map[string]int64)
	mf.attributesToNumber["label"] = make(map[string]int64)
	mf.attributesToNumber["secret"] = make(map[string]int64)
	mf.attributesToNumber["state"] = make(map[string]int64)
	mf.secrets = make(map[int]bool)
	mf.uniques = make(map[string]bool)

	for i := 0; i < len(boxes.Traits["secret"]); {
	retry:
		num := rand.Intn(amount)
		if _, ok := mf.secrets[num]; ok {
			goto retry
		}
		mf.secrets[num] = true
		i++
	}

	wg := &sync.WaitGroup{}
	wg.Add(amount)

	mf.factory = boxes.CreateFactory()

	for x := 1; x <= amount; x++ {
		go mf.createUnique(x, wg)
	}

	wg.Wait()

	jsonBytes, err := json.MarshalIndent(mf.attributesToNumber, "", "	")
	if err != nil {
		log.Fatal(err)
	}

	f, err := os.Create("./TBB/boxes.json")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	fmt.Fprint(f, string(jsonBytes))
}
func createCanvas() {
	factory := boxes.CreateFactory()
	uniques := make(map[string]bool)

	width := 400
	height := 400
	m := image.NewRGBA(image.Rect(0, 0, width, height))
	maxBoxesOnX := width/54 + 2
	maxBoxesOnY := height/45 + 2

	fmt.Println(maxBoxesOnX)
	x := 0
	y := 0

	var newBoxes = make([]boxes.Box, maxBoxesOnX*maxBoxesOnY)

	for i := 0; i < maxBoxesOnX*maxBoxesOnY; i++ {
	retry:
		box, err := factory.CreateBox()
		if err != nil {
			log.Fatal(err)
		}
		box.Background = boxes.Attribute{}

		hash := box.CreateHash()
		if _, ok := uniques[hash]; ok {
			goto retry
		}

		locX := i % maxBoxesOnX
		locY := i / maxBoxesOnX

		if i != 0 {
			if newBoxes[locX-1+locY*maxBoxesOnX].Color.ImagePath == box.Color.ImagePath {
				goto retry
			}

		}
		if locY != 0 {
			if locX != 0 {
				if newBoxes[locX+(locY-1)*maxBoxesOnX].Color.ImagePath == box.Color.ImagePath {
					goto retry
				}
				if newBoxes[locX-1+(locY-1)*maxBoxesOnX].Color.ImagePath == box.Color.ImagePath {
					goto retry
				}
			}
			if newBoxes[locX+1+(locY-1)*maxBoxesOnX].Color.ImagePath == box.Color.ImagePath {
				goto retry
			}

		}
		newBoxes[i] = box
		uniques[hash] = true
	}

	for i := 0; i < maxBoxesOnX*maxBoxesOnY; i++ {
		boxPNG, err := newBoxes[i].GetPNG()
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

	jsonMap := make(map[int]map[string]string)

	for i, box := range newBoxes {
		jsonMap[i+1] = make(map[string]string)
		jsonMap[i+1]["color"] = box.Color.Name
		jsonMap[i+1]["cutouts"] = strings.Join(attributesToNames(box.Cutouts), ",")
		jsonMap[i+1]["adhesives"] = strings.Join(attributesToNames(box.Adhesives), ",")
		jsonMap[i+1]["label"] = box.Label.Name
		jsonMap[i+1]["straps"] = strings.Join(attributesToNames(box.Straps), ",")

	}

	f, err := os.Create("./TBB/canvas.png")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	//	x2048 := imaging.Resize(m, 2048, 2048, imaging.NearestNeighbor)

	err = png.Encode(f, m)
	if err != nil {
		log.Fatal(err)
	}
}

func toRaw(attributeName string) string {
	return strings.ReplaceAll(attributeName, " ", "_")
}

func attributesToNames(attributes []boxes.Attribute) (strs []string) {
	for i := range attributes {
		strs = append(strs, toRaw(attributes[i].Name))
	}
	// Sorts the array, doesn't return anything and just modifies the original array
	sort.Strings(strs)
	return
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
