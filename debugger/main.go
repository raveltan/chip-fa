package debugger

import (
	"fmt"
	"strconv"
	"strings"

	"gopkg.in/abiosoft/ishell.v2"
)

type Debugger struct {
	shell                   *ishell.Shell
	ResumeEmulationCallback func() bool
	PauseEmulationCallback  func() bool
	GetRegisterCallback     func() [16]uint8
	GetTimerCallback        func() [2]uint8
	ExitCallback            func()
	SetRegisterCallback     func(uint8, uint8)
	SetTimerCallback        func(bool, uint8)
	GetSpecialCallback      func() string
	SetICallback            func(uint16)
	SetPcCallback           func(uint16)
	GetMemoryViewCallback   func() []uint8
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
		Name:    "resume",
		Aliases: []string{"r"},
		Help:    "[r] Resume the ongoing emulation",
		Func: func(c *ishell.Context) {
			if d.ResumeEmulationCallback() {
				c.Println("Resumed Emulation")
				return
			}
			c.Println("Unable to resume as emulation is not paused")
		},
	})
	d.shell.AddCmd(&ishell.Cmd{
		Name:    "pause",
		Aliases: []string{"p"},
		Help:    "[p] Pause the ongoing emulation",
		Func: func(c *ishell.Context) {
			if d.PauseEmulationCallback() {
				c.Println("Paused Emulation")
				return
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
		Name:    "register",
		Aliases: []string{"v"},
		Help:    "[v] Get list of current registers states",
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
	d.shell.AddCmd(&ishell.Cmd{
		Name: "sv",
		Help: "Set a value to a specific register (ex: setv 0xF 0xFF)",
		Func: func(c *ishell.Context) {
			addressString := strings.Replace(c.Args[0], "0x", "", -1)
			address, err := strconv.ParseUint(addressString, 16, 8)
			if err != nil {
				c.Println(fmt.Sprintf("Unable to parse register location data: [0x%v], make sure that it is (0x0-0xF)", addressString))
				return
			}
			if address > 0xF {
				c.Println("Register location must be between 0x0 - 0xF (V0-VF)")
				return
			}
			dataString := strings.Replace(c.Args[1], "0x", "", -1)
			data, err := strconv.ParseUint(dataString, 16, 8)
			if err != nil {
				c.Println(fmt.Sprintf("Unable to parse new value data [0x%v], make sure that it is (0x00-0xFF)", dataString))
				return
			}
			d.SetRegisterCallback(uint8(address), uint8(data))
		},
	})

	d.shell.AddCmd(&ishell.Cmd{
		Name: "st",
		Help: "Set a value to a specific timer (ex: setv 0x0 0xFF, where 0x0 means delay timer and 0x1 means sound timer)",
		Func: func(c *ishell.Context) {
			addressString := strings.Replace(c.Args[0], "0x", "", -1)
			address, err := strconv.ParseUint(addressString, 16, 8)
			if err != nil {
				c.Println(fmt.Sprintf("Unable to parse register location data: [0x%v], make sure that it is either 0x0 or 0x1", addressString))
				return
			}
			if address > 0x1 {
				c.Println("Timer should be either 0x0 (delay timer) or 0x1 (sound timer)")
				return
			}
			dataString := strings.Replace(c.Args[1], "0x", "", -1)
			data, err := strconv.ParseUint(dataString, 16, 8)
			if err != nil {
				c.Println(fmt.Sprintf("Unable to parse new value data [0x%v], make sure that it is (0x00-0xFF)", dataString))
				return
			}
			d.SetTimerCallback(uint8(address) == 0, uint8(data))
		},
	})

	d.shell.AddCmd(&ishell.Cmd{
		Name:    "timer",
		Aliases: []string{"t"},
		Help:    "[t] Get list of current timers states",
		Func: func(c *ishell.Context) {
			timers := d.GetTimerCallback()
			data := [][]string{}
			for i, v := range timers {
				if i == 0 {
					data = append(data, []string{
						"Delay", fmt.Sprintf("0x%x", v),
					})
				} else {
					data = append(data, []string{
						"Sound", fmt.Sprintf("0x%x", v),
					})
				}
			}
			head, body := buildHorizontalTable(data)
			c.Println(head)
			c.Println(body)
		},
	})

	d.shell.AddCmd(&ishell.Cmd{
		Name: "cpu",
		Help: "Get CPU datas such as I (IndexRegister), pc(Program Counter), stack and stack pointer",
		Func: func(c *ishell.Context) {
			c.Println(d.GetSpecialCallback())
		},
	})

	d.shell.AddCmd(&ishell.Cmd{
		Name: "si",
		Help: "Set 16 bit unsigned value to the I (IndexRegister) (ex: si 0xFFF)",
		Func: func(c *ishell.Context) {
			valueString := strings.Replace(c.Args[0], "0x", "", -1)
			value, err := strconv.ParseUint(valueString, 16, 16)
			if err != nil {
				c.Println(fmt.Sprintf("Unable to parse value data: [0x%v], make sure that it is an unsigned 16bit integer", valueString))
				return
			}
			d.SetICallback(uint16(value))
		},
	})

	d.shell.AddCmd(&ishell.Cmd{
		Name: "spc",
		Help: "Set a 16 bit unsigned value to the PC (ProgarmCounter) (ex: spc 0xFFF) ",
		Func: func(c *ishell.Context) {
			valueString := strings.Replace(c.Args[0], "0x", "", -1)
			value, err := strconv.ParseUint(valueString, 16, 16)
			if err != nil {
				c.Println(fmt.Sprintf("Unable to parse value data: [0x%v], make sure that it is an unsigned 16bit integer", valueString))
				return
			}
			d.SetPcCallback(uint16(value))
		},
	})

	d.shell.AddCmd(&ishell.Cmd{
		Name:    "instruction-view",
		Aliases: []string{"iv"},
		Help:    "View current instruction data of Program counter +- 30 entries",
		Func: func(c *ishell.Context) {
			result := d.GetMemoryViewCallback()
			text := ""
			prevData := uint8(0)
			for i, v := range result {
				if i == 121 {
					text += fmt.Sprintf("{[0x%x}]\t", uint16(prevData)<<8|uint16(v))
				} else if i%2 == 0 {
					prevData = v
				} else {
					text += fmt.Sprintf("0x%x\t", uint16(prevData)<<8|uint16(v))
				}
			}
			c.Println(text)
		},
	})

	// run shell
	d.shell.Run()
}
