package main

import (
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"ravel.com/chiper/cpu"
)

type Emulator struct {
	cpu *cpu.CPU
}

func (e *Emulator) Update() error {
	e.cpu.DoCycle()
	// TODO: find a way to do conditional rerendering
	// if e.cpu.ShouldDraw {
	// 	log.Println("Should redraw")
	// }

	// TODO: Handle user input

	return nil
}

func (e *Emulator) Draw(s *ebiten.Image) {
	for i, v := range e.cpu.Screen {
		drawColor := color.White
		if v == 0 {
			drawColor = color.Black
		}
		s.Set(i%64, i/64, drawColor)
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
	ebiten.SetWindowTitle("Chiper")
	if err := ebiten.RunGame(&Emulator{cpu: cpu}); err != nil {
		log.Fatal(err)
	}
}
