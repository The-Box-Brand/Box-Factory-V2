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
	"sync"
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

func main() {
	mf := miniFactory{}
	mf.createManyUnique(15)

	createCanvas()
	createTest()
}

func createTest() {
	factory := boxes.CreateFactory()
	box, err := factory.CreateBox()
	if err != nil {
		log.Fatal(err)
	}

	box.SaveBox("test.png")
}

type miniFactory struct {
	sync.RWMutex
	hashes             []string
	attributesToNumber map[string]map[string]int64
	factory            boxes.Factory
	secrets            map[int]bool
}

func (mf *miniFactory) createUnique(num int, wg *sync.WaitGroup) {
	defer wg.Done()
	mf.Lock()
	defer mf.Unlock()
retry:
	fmt.Println("Creating box: ", num)

	var box boxes.Box
	var err error
	if mf.secrets[num] {
		box = mf.factory.CreateSecretBox()
		mf.attributesToNumber["secret"][box.Secret.Name]++
	} else {
		box, err = mf.factory.CreateBox()
		if err != nil {
			log.Fatal(err)
		}

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

		hash := box.CreateHash()
		for i := 0; i < len(mf.hashes); i++ {
			if mf.hashes[i] == hash {
				fmt.Println("got same box: here are all the unique boxes - " + fmt.Sprint(len(mf.hashes)))
				goto retry

			}
		}
		mf.hashes = append(mf.hashes, hash)
	}

	err = box.SaveBox("./TBB/" + fmt.Sprint(num) + ".png")
	if err != nil {
		log.Fatal(err)
	}

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
	mf.secrets = make(map[int]bool)

	for i := 0; i < len(boxes.Traits["secret"]); {
		num := rand.Intn(amount)
		if _, ok := mf.secrets[num]; ok {
			continue
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

	jsonBytes, _ := json.MarshalIndent(mf.attributesToNumber, "", "	")
	log.Println(string(jsonBytes))
}
func createCanvas() {
	factory := boxes.CreateFactory()

	width := 1027
	height := 1027
	m := image.NewRGBA(image.Rect(0, 0, width, height))
	maxBoxesOnX := width/54 + 1
	maxBoxesOnY := height/45 + 1

	x := 0
	y := 0
	for i := 0; i < maxBoxesOnX*maxBoxesOnY; i++ {
		box, err := factory.CreateBox()
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
