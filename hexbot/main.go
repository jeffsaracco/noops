package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"sync"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
	colors "gopkg.in/go-playground/colors.v1"
)

// Colors is stuff
type Colors struct {
	Values []Dot `json:"colors"`
}

//Coordinate is a point
type Coordinate struct {
	X int `json:"x"`
	Y int `json:"y"`
}

// Dot is like {"value":"#EDA9E9","coordinates":{"x":200,"y":186}}<Paste>
type Dot struct {
	Value       string     `json:"value"`
	Coordinates Coordinate `json:"coordinates"`
}

//Height is the height
const Height = 768

//Width is the width
const Width = 1024

func run() {
	cfg := pixelgl.WindowConfig{
		Title:  "Hexbot",
		Bounds: pixel.R(0, 0, Width, Height),
		VSync:  true,
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}
	imd := imdraw.New(nil)
	// draw(imd)

	for !win.Closed() {
		if win.JustPressed(pixelgl.MouseButtonLeft) {
			draw(imd)
		}
		if win.JustPressed(pixelgl.MouseButtonRight) {
			imd.Clear()
		}
		win.Clear(colornames.Aliceblue)
		imd.Draw(win)
		win.Update()
	}
}

func main() {
	pixelgl.Run(run)
}

func draw(imd *imdraw.IMDraw) {
	urls := []string{
		fmt.Sprintf("https://api.noopschallenge.com/hexbot?count=%d&width=%d&height=%d", rand.Intn(5000), Width, Height),
	}
	jsonResponses := make(chan Colors)

	var wg sync.WaitGroup

	wg.Add(len(urls))

	go func(url string) {
		defer wg.Done()
		res, err := http.Get(url)
		if err != nil {
			log.Fatal(err)
		} else {
			defer res.Body.Close()
			body, err := ioutil.ReadAll(res.Body)
			if err != nil {
				log.Fatal(err)
			} else {
				var dots Colors
				err = json.Unmarshal(body, &dots)
				if err != nil {
					log.Fatal(err)
				}
				jsonResponses <- dots
			}
		}
	}(urls[0])

	go func() {
		wg.Wait()
		close(jsonResponses)
	}()

	for response := range jsonResponses {
		for _, c := range response.Values {
			parsedColor, err := colors.Parse(c.Value)
			if err != nil {
				log.Fatal(err)
			}
			color := parsedColor.ToRGB()
			imd.Color = pixel.RGB(float64(color.R)/255, float64(color.G)/255, float64(color.B)/255)
			imd.Push(pixel.V(float64(c.Coordinates.X), float64(c.Coordinates.Y)))
			imd.Circle(float64(rand.Intn(20)), 0)
		}
	}
}
