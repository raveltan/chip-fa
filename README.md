# Chip-Fa
![Chip-fa icon](./icon.png)


A fully-featured [Chip8](https://en.wikipedia.org/wiki/CHIP-8) emulator written in GO.

## Current Status
The following is the planned features of Chip-fa:
- Screen emulation
- Sound emulation
- Keypad emulation
- Debugger
- Memory view (in debugger)

## Installation
Chip-fa's executable is currently available in these platforms:
- Windows x86
- Windows x64
- Linux x64


To download the executable, head to [release section](https://github.com/raveltan/chip-fa/releases).

For best experiece of using Chip-fa, please add it to the system PATH.
```bash
# Linux/Mac:
echo 'export PATH="$PATH:/path/to/chip-fa"' >> ~/.profile
source ~/.profile
# Windows:
set PATH=%PATH%;C:\path\to\chip-fa\
```

> All NON-LINUX releases is created with [XGO](https://github.com/karalabe/xgo) cross compile tool. compatibility not guaranteed.
## Usage
Running a ROM with default configuration
```bash
chip-fa -r roms/tetris.ch8
```
You can also set the window scaling using the -s flag
```bash
chip-fa -r roms/tetris.ch8 -s 2.0
```
If you are using HIDPI screen, please make sure to set the -x flag according to your system scaling if not done automatically.
```bash
chip-fa -r roms/tetris.ch8 -x 2.0
```
If you want to alter the amount of cycle per second, it can be done by the -c flag. 
```bash
# Run the emulator at 600/60 of the normal clock speed
chip-fa -r roms/tetris.ch8 -c 600
```
You can also enable debug mode for developing ROMS. (more about the debugger at the next section)
```bash
chip-fa -r roms/tetris.ch8 -d
```
## Debugging ROMS

Chip-fa includes a very powerful debugger for ROM developer. This debugger can be activated by pressing 0 when you are running a ROM the -d flag and will also be automatically activated when running a ROM with breakpoint (0x0001) instruction (this instruction is only available to this emulator). Some of the common use cases is available here:

Gather CPU related information (I,PC,stack,SP)
```bash
cpu
```

Gather instructions surrounding the current program counter
```bash
instruction-view
```
or
```bash
iv
```

Set a specific value to a register, in this case VF with value of 0xFF
```bash
sv 0xF 0xFF
```

more information about the command available in the debuger can be accessed from the help menu.
```bash
help
```


## Additional Operation Codes

- 0x0001: Debug breakpoint, if a ROM contains this operation code, and debuggin mode is on, the application while stop and activate the debugger shell.

## Screenshots
![Chip8 Logo ROM.](./ss/chip8.png)
![Particle Demo ROM.](./ss/particle.png)
![Tetris ROM.](./ss/tetris.png)

## Bug Reporting
Please create a new issue with the bugs detail.

## Support Us
Show your support to this project by submitting pull request or starring this project :))