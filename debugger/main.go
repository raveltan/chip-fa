package debugger

import (
	"fmt"

	"gopkg.in/abiosoft/ishell.v2"
)

type Debugger struct {
	shell                   *ishell.Shell
	ResumeEmulationCallback func() bool
	PauseEmulationCallback  func() bool
	GetRegisterCallback     func() [16]uint8
	ExitCallback            func()
}

func buildHorizontalTable(data [][]string) (header string, content string) {
	for _, v := range data {
		maxLength := 0
		head := ""
		body := ""
		for i, v1 := range v {
			if len(v1) > maxLength {
				maxLength = len(v1)
			}
			if i == 0 {
				head += v1
			} else {
				body += v1
			}
		}
		if len(head) < maxLength {
			for i := 0; i < maxLength-len(head); i++ {
				head += " "
			}
		}

		if len(body) < maxLength {
			for i := 0; i < maxLength-len(body); i++ {
				body += " "
			}
		}
		header += head + "\t"
		content += body + "\t"
	}
	return
}

func (d *Debugger) StartDebugShell() {
	if d.shell != nil {
		return
	}
	d.shell = ishell.New()

	// display welcome info.
	d.shell.Println("Chip-fa Interactive Debugger Shell")

	d.shell.AddCmd(&ishell.Cmd{
		Name: "resume",
		Help: "Resume the ongoing emulation",
		Func: func(c *ishell.Context) {
			if d.ResumeEmulationCallback() {
				c.Println("Resumed Emulation")
				return
			}
			c.Println("Unable to resume as emulation is not paused")
		},
	})
	d.shell.AddCmd(&ishell.Cmd{
		Name: "pause",
		Help: "Pause the ongoing emulation",
		Func: func(c *ishell.Context) {
			if d.PauseEmulationCallback() {
				c.Println("Paused Emulation")
			}
			c.Println("Unable to pause emulation as it yet is not running")
		},
	})
	d.shell.AddCmd(&ishell.Cmd{
		Name: "exit",
		Help: "Exit Chip-Fa emulator",
		Func: func(c *ishell.Context) {
			c.Println("Thank you for using Chip-fa")
			d.ExitCallback()
		},
	})

	d.shell.AddCmd(&ishell.Cmd{
		Name: "register",
		Help: "Get list of current registers states",
		Func: func(c *ishell.Context) {
			registers := d.GetRegisterCallback()
			data := [][]string{}
			for i, v := range registers {
				data = append(data, []string{
					fmt.Sprintf("v%x", i), fmt.Sprintf("0x%x", v),
				})
			}
			head, body := buildHorizontalTable(data)
			c.Println(head)
			c.Println(body)
		},
	})

	// run shell
	d.shell.Run()
}
