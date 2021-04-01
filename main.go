package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"ravel.com/chiper/cpu"
)

type Emulator struct {
	cpu cpu.CPU
}

func (e *Emulator) Update() error {
	e.cpu.DoCycle()
	if e.cpu.ShouldDraw {
		log.Println("should draw")
	}

	// TODO: Handle user input

	return nil
}

func (e *Emulator) Draw(s *ebiten.Image) {

}

func (e *Emulator) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return outsideWidth / 2, outsideHeight / 2
}

func main() {
	// TODO: Enable load ROM
	// TODO: add cli utilities

	ebiten.SetWindowSize(600, 600)
	ebiten.SetWindowTitle("Chiper")

	cpu := cpu.CPU{}
	cpu.Boot()
	if err := ebiten.RunGame(&Emulator{cpu: cpu}); err != nil {
		log.Fatal(err)
	}
}
