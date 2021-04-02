package main

import (
	"image/color"
	"log"

	_ "image/png"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/raveltan/chip-fa/cpu"
)

var splash *ebiten.Image

func init() {
	var err error
	splash, _, err = ebitenutil.NewImageFromFile("assets/splash.png")
	if err != nil {
		log.Fatal(err)
	}
}

type Emulator struct {
	cpu               *cpu.CPU
	splashScreenTimer int
}

func (e *Emulator) Update() error {
	if e.splashScreenTimer == 0 {
		e.cpu.DoCycle()
		// TODO: find a way to do conditional rerendering
		// if e.cpu.ShouldDraw {
		// 	log.Println("Should redraw")
		// }

		// TODO: Handle user input
	} else {
		e.splashScreenTimer--
	}
	return nil
}

func (e *Emulator) Draw(s *ebiten.Image) {
	if e.splashScreenTimer == 0 {
		for i, v := range e.cpu.Screen {
			drawColor := color.White
			if v == 0 {
				drawColor = color.Black
			}
			s.Set(i%64, i/64, drawColor)
		}
	} else {
		drawOptions := &ebiten.DrawImageOptions{}
		// 640 x 320 image size
		// canvas size 64 x 32
		drawOptions.GeoM.Scale(0.1, 0.1)
		s.DrawImage(splash, drawOptions)
	}
}

func (e *Emulator) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return outsideWidth / 24, outsideHeight / 24
}

func main() {
	// Initialize CPU
	cpu := new(cpu.CPU)
	cpu.Boot()
	// TODO: Enable load ROM from flag
	cpu.LoadROM("./roms/ibm_logo.ch8")
	// TODO: add cli utilities
	// TODO: add scaling functionalities
	ebiten.SetWindowSize(64*12, 32*12)
	ebiten.SetWindowTitle("Chip-Fa: Chip8's emulator")
	if err := ebiten.RunGame(&Emulator{cpu: cpu, splashScreenTimer: 3 * 60}); err != nil {
		log.Fatal(err)
	}
}
