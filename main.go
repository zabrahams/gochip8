package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/veandco/go-sdl2/sdl"
	"github.com/zabrahams/gochip8/beeper"
	"github.com/zabrahams/gochip8/chip8"
	"github.com/zabrahams/gochip8/screen"
)

const helpMsg = `
gochip8 is a chip 8 emulator! 
You can use it as follows:
./gochip8 mode rom
where mode is either:
	run - runs the rom
	dis - dissassembles the rom
	debug - runs the rom in debug mode
and rom is a path to the rom
`

func main() {

	if len(os.Args) < 3 {
		fmt.Println(helpMsg)
		os.Exit(1)
	}

	subcommand := os.Args[1]
	if subcommand == "" {
		panic("need subcommand: run, dis or debug")
	}
	programFile := os.Args[2]
	if programFile == "" {
		panic("no program file given")
	}

	switch subcommand {
	case "run":
		run(programFile)
	case "debug":
		debug(programFile)
	case "dis":
		dis(programFile)
	default:
		panic(fmt.Sprintf("unknown command: %s", subcommand))
	}

}

func dis(programFile string) {
	file, err := os.Open(programFile)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	program, err := ioutil.ReadAll(file)
	if err != nil {
		panic(err)
	}
	builder := chip8.Disassemble(program, 0x200)
	fmt.Println(builder.String())
}

func debug(programFile string) {
	var (
		command string
		s       struct{}
	)
	fmt.Println("Starting Chip8 Emulator")

	screen := screen.NewScreen()
	defer screen.Close()

	beeper := beeper.NewSDLBeeper()
	defer beeper.Close()

	c8 := chip8.NewChip8(beeper)
	c8.Load(programFile)
	c8.String()
	quit := false
	running := false
	for !quit {
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch event.(type) {
			case *sdl.QuitEvent:
				println("Quit")
				quit = true
			case *sdl.KeyboardEvent:
				kevent := event.(*sdl.KeyboardEvent)
				if kevent.Type == sdl.KEYUP && kevent.Keysym.Sym == sdl.K_PERIOD {
					c8.Stop <- s
					c8.String()
					running = false
				}
			}
		}
		if !running {
			fmt.Print("command: (h for help) ")
			fmt.Scanln(&command)
			switch command {
			case "s":
				c8.ExecInstr()
				c8.String()
			case "r":
				c8.Run()
				running = true
			case "q":
				quit = true
			case "h":
				fmt.Println("you can use the following commands (s)tep, (r)un, (q)uit.  you can also use '.' to stop the running program.")
			}
		}
		kbState := sdl.GetKeyboardState()
		newKBState := parseKbState(kbState)

		c8.Keyboard.Update(newKBState)
		screen.Update(c8.FrameBuffer)

	}
	fmt.Println("Closing Chip8 Emulator")
}

func run(programFile string) {
	fmt.Println("Starting Chip8 Emulator")

	screen := screen.NewScreen()
	defer screen.Close()

	beeper := beeper.NewSDLBeeper()
	defer beeper.Close()

	c8 := chip8.NewChip8(beeper)
	c8.Load(programFile)
	c8.Run()
	running := true
	for running {
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch event.(type) {
			case *sdl.QuitEvent:
				println("Quit")
				running = false
			}
		}
		kbState := sdl.GetKeyboardState()
		newKBState := parseKbState(kbState)

		c8.Keyboard.Update(newKBState)
		screen.Update(c8.FrameBuffer)
	}
	fmt.Println("Closing Chip8 Emulator")
}

func parseKbState(kbState []uint8) uint16 {
	var keys uint16
	keys = 0
	if kbState[sdl.SCANCODE_1] == 1 {
		keys = keys | (0x1 << 0)
	}
	if kbState[sdl.SCANCODE_2] == 1 {
		keys = keys | (0x1 << 1)
	}
	if kbState[sdl.SCANCODE_3] == 1 {
		keys = keys | (0x1 << 2)
	}
	if kbState[sdl.SCANCODE_4] == 1 {
		keys = keys | (0x1 << 3)
	}
	if kbState[sdl.SCANCODE_Q] == 1 {
		keys = keys | (0x1 << 4)
	}
	if kbState[sdl.SCANCODE_W] == 1 {
		keys = keys | (0x1 << 5)
	}
	if kbState[sdl.SCANCODE_E] == 1 {
		keys = keys | (0x1 << 6)
	}
	if kbState[sdl.SCANCODE_R] == 1 {
		keys = keys | (0x1 << 7)
	}
	if kbState[sdl.SCANCODE_A] == 1 {
		keys = keys | (0x1 << 8)
	}
	if kbState[sdl.SCANCODE_S] == 1 {
		keys = keys | (0x1 << 9)
	}
	if kbState[sdl.SCANCODE_D] == 1 {
		keys = keys | (0x1 << 10)
	}
	if kbState[sdl.SCANCODE_F] == 1 {
		keys = keys | (0x1 << 11)
	}
	if kbState[sdl.SCANCODE_Z] == 1 {
		keys = keys | (0x1 << 12)
	}
	if kbState[sdl.SCANCODE_X] == 1 {
		keys = keys | (0x1 << 13)
	}
	if kbState[sdl.SCANCODE_C] == 1 {
		keys = keys | (0x1 << 14)
	}
	if kbState[sdl.SCANCODE_V] == 1 {
		keys = keys | (0x1 << 15)
	}
	return keys
}
