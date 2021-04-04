package emulator

import (
	"fmt"
	"log"
	"os"

	"image/color"
	_ "image/png"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/raveltan/chip-fa/cpu"
	"github.com/raveltan/chip-fa/debugger"
	"github.com/raveltan/chip-fa/wavegen"
)

type Emulator struct {
	Cpu             *cpu.CPU
	beepAudioPlayer *audio.Player
	beepAudioTimer  int
	scaleFactor     float64
	debug           *debugger.Debugger
	Pause           bool
}

func (e *Emulator) Update() error {
	// Reset keypad state
	for i := range e.Cpu.KeypadStates {
		e.Cpu.KeypadStates[i] = 0
	}
	// Set keypad states
	for k := ebiten.Key(0); k <= ebiten.KeyMax; k++ {
		if ebiten.IsKeyPressed(k) {
			switch k.String() {
			case "1":
				e.Cpu.KeypadStates[0] = 1
			case "2":
				e.Cpu.KeypadStates[1] = 1
			case "3":
				e.Cpu.KeypadStates[2] = 1
			case "4":
				e.Cpu.KeypadStates[3] = 1
			case "Q":
				e.Cpu.KeypadStates[4] = 1
			case "W":
				e.Cpu.KeypadStates[5] = 1
			case "E":
				e.Cpu.KeypadStates[6] = 1
			case "R":
				e.Cpu.KeypadStates[7] = 1
			case "A":
				e.Cpu.KeypadStates[8] = 1
			case "S":
				e.Cpu.KeypadStates[9] = 1
			case "D":
				e.Cpu.KeypadStates[10] = 1
			case "F":
				e.Cpu.KeypadStates[11] = 1
			case "Z":
				e.Cpu.KeypadStates[12] = 1
			case "X":
				e.Cpu.KeypadStates[13] = 1
			case "C":
				e.Cpu.KeypadStates[14] = 1
			case "V":
				e.Cpu.KeypadStates[15] = 1
			case "0":
				if e.debug != nil {
					e.Pause = true
					go e.debug.StartDebugShell()
				}
			}
		}
	}
	if !e.Pause {
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
		if e.Cpu.SoundTimer > 0 || e.beepAudioTimer > 0 {
			if e.Cpu.SoundTimer > 0 {
				if e.Cpu.SoundTimer == 1 {
					// Beep the buzzer
					e.beepAudioPlayer.Play()
					e.beepAudioTimer = 25
				}
				e.Cpu.SoundTimer--
			}
			if e.beepAudioTimer > 0 {
				if e.beepAudioTimer == 1 {
					// Stop wave generator
					e.beepAudioPlayer.Pause()
				}
				e.beepAudioTimer--
			}
		}
		e.Cpu.DoCycle()
	}
	return nil
}

func (e *Emulator) Draw(s *ebiten.Image) {
	for i, v := range e.Cpu.Screen {
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
	// Initialize CPU
	cpu := new(cpu.CPU)
	cpu.Boot()
	if err := cpu.LoadROM(rom); err != nil {
		log.Fatal(fmt.Sprintf("error: Unable to open ROM, %v", err))
	}
	ebiten.SetWindowSize(64*12*int(displayScale), 32*12*int(displayScale))
	ebiten.SetWindowTitle("Chip-Fa")
	ebiten.SetMaxTPS(cyclePerSecond)

	// Setup emulator and debugger
	emulator := &Emulator{Cpu: cpu, scaleFactor: DPIscale}
	emulator.debug = &debugger.Debugger{ResumeEmulationCallback: func() bool {
		if !emulator.Pause {
			return false
		}
		emulator.Pause = false
		return true
	}, PauseEmulationCallback: func() bool {
		if emulator.Pause {
			return false
		}
		emulator.Pause = false
		return true
	}, ExitCallback: func() {
		os.Exit(0)
	}, GetRegisterCallback: func() [16]uint8 {
		return emulator.Cpu.Register
	}}

	// Start emulation
	if err := ebiten.RunGame(emulator); err != nil {
		log.Fatal(err)
	}

}
