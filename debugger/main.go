package debugger

import (
	"strings"

	"gopkg.in/abiosoft/ishell.v2"
)

type Debugger struct {
	shell                   *ishell.Shell
	ResumeEmulationCallback func()
	PauseEmulationCallback  func()
	ExitCallback            func()
}

func (d *Debugger) StartDebugShell() {
	if d.shell != nil {
		return
	}
	d.shell = ishell.New()

	// display welcome info.
	d.shell.Println("Chip-fa Interactive Debugger Shell")

	d.shell.AddCmd(&ishell.Cmd{
		Name: "greet",
		Help: "greet user",
		Func: func(c *ishell.Context) {
			c.Println("Hello", strings.Join(c.Args, " "))
		},
	})
	d.shell.AddCmd(&ishell.Cmd{
		Name: "resume",
		Help: "Resume the ongoing emulation",
		Func: func(c *ishell.Context) {
			c.Println("Resumed Emulation")
			d.ResumeEmulationCallback()
		},
	})
	d.shell.AddCmd(&ishell.Cmd{
		Name: "pause",
		Help: "Pause the ongoing emulation",
		Func: func(c *ishell.Context) {
			c.Println("Paused Emulation")
			d.PauseEmulationCallback()
		},
	})
	d.shell.AddCmd(&ishell.Cmd{
		Name: "exit",
		Help: "Exit chip-fa",
		Func: func(c *ishell.Context) {
			c.Println("Thank you for using Chip-fa")
			d.ExitCallback()
		},
	})

	// run shell
	d.shell.Run()
}
