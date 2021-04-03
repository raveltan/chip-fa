package emulator

import (
	"fmt"
	"log"

	"image/color"
	_ "image/png"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/raveltan/chip-fa/cpu"
	"github.com/raveltan/chip-fa/wavegen"
)

type Emulator struct {
	cpu             *cpu.CPU
	beepAudioPlayer *audio.Player
	beepAudioTimer  int
	scaleFactor     float64
	debug           bool
	pause           bool
}

func (e *Emulator) Update() error {
	// Reset keypad state
	for i := range e.cpu.KeypadStates {
		e.cpu.KeypadStates[i] = 0
	}
	// Set keypad states
	for k := ebiten.Key(0); k <= ebiten.KeyMax; k++ {
		if ebiten.IsKeyPressed(k) {
			switch k.String() {
			case "1":
				e.cpu.KeypadStates[0] = 1
			case "2":
				e.cpu.KeypadStates[1] = 1
			case "3":
				e.cpu.KeypadStates[2] = 1
			case "4":
				e.cpu.KeypadStates[3] = 1
			case "Q":
				e.cpu.KeypadStates[4] = 1
			case "W":
				e.cpu.KeypadStates[5] = 1
			case "E":
				e.cpu.KeypadStates[6] = 1
			case "R":
				e.cpu.KeypadStates[7] = 1
			case "A":
				e.cpu.KeypadStates[8] = 1
			case "S":
				e.cpu.KeypadStates[9] = 1
			case "D":
				e.cpu.KeypadStates[10] = 1
			case "F":
				e.cpu.KeypadStates[11] = 1
			case "Z":
				e.cpu.KeypadStates[12] = 1
			case "X":
				e.cpu.KeypadStates[13] = 1
			case "C":
				e.cpu.KeypadStates[14] = 1
			case "V":
				e.cpu.KeypadStates[15] = 1

			case "0":
				// add to cli help
				if e.debug {
					e.pause = true
				}
			case "9":
				if e.debug {
					e.pause = false
				}
			}

		}
	}
	if !e.pause {
		if e.beepAudioPlayer == nil {
			// Pass the (infinite) stream to audio.NewPlayer.
			// After calling Play, the stream never ends as long as the player object lives.
			var err error
			e.beepAudioPlayer, err = audio.NewPlayer(audio.NewContext(44100), &wavegen.Stream{})
			if err != nil {
				return err
			}
		}

		// FIXME: find a way to do conditional rerendering
		// Update sound timer
		if e.cpu.SoundTimer > 0 || e.beepAudioTimer > 0 {
			if e.cpu.SoundTimer > 0 {
				if e.cpu.SoundTimer == 1 {
					// Beep the buzzer
					e.beepAudioPlayer.Play()
					e.beepAudioTimer = 25
				}
				e.cpu.SoundTimer--
			}
			if e.beepAudioTimer > 0 {
				if e.beepAudioTimer == 1 {
					// Stop wave generator
					e.beepAudioPlayer.Pause()
				}
				e.beepAudioTimer--
			}
		}
		e.cpu.DoCycle()
	}
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
	return (outsideWidth) / (12 * int(e.scaleFactor)), (outsideHeight) / int((12 * e.scaleFactor))
}

func StartEmulation(rom string, DPIscale float64, displayScale float64, cyclePerSecond int, debug bool) {
	if DPIscale == 0 {
		DPIscale = (ebiten.DeviceScaleFactor())
	}
	// Initialize CPU
	cpu := new(cpu.CPU)
	cpu.Boot()
	if err := cpu.LoadROM(rom); err != nil {
		log.Fatal(fmt.Sprintf("error: Unable to open ROM, %v", err))
	}
	ebiten.SetWindowSize(64*12*int(displayScale), 32*12*int(displayScale))
	ebiten.SetWindowTitle("Chip-Fa")
	ebiten.SetMaxTPS(cyclePerSecond)
	if err := ebiten.RunGame(&Emulator{cpu: cpu, scaleFactor: DPIscale, debug: debug}); err != nil {
		log.Fatal(err)
	}

}
