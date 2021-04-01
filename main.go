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
	return outsideWidth / 12, outsideHeight / 12
}

func main() {
	// TODO: Enable load ROM
	// TODO: add cli utilities
	// TODO: add scaling functionalities
	ebiten.SetWindowSize(64*12, 32*12)
	ebiten.SetWindowTitle("Chiper")

	cpu := cpu.CPU{}
	cpu.Boot()
	if err := ebiten.RunGame(&Emulator{cpu: cpu}); err != nil {
		log.Fatal(err)
	}
}
