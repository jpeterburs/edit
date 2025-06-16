package main

import (
	"fmt"
	"os"

	"github.com/nsf/termbox-go"
)

var (
	lines            []string
	cursorX, cursorY int
	filename         string
)

func loadFile(name string) {
	data, err := os.ReadFile(name)
	if err != nil {
		lines = []string{""}
		return
	}
	lines = []string{}
	curr := ""
	for _, b := range data {
		if b == '\n' {
			lines = append(lines, curr)
			curr = ""
		} else {
			curr += string(b)
		}
	}
	lines = append(lines, curr)
}

func saveFile(name string) {
	output := ""
	for i, l := range lines {
		output += l
		if i != len(lines)-1 {
			output += "\n"
		}
	}
	os.WriteFile(name, []byte(output), 0644)
}

func draw() {
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
	for y, line := range lines {
		for x, ch := range line {
			termbox.SetCell(x, y, ch, termbox.ColorDefault, termbox.ColorDefault)
		}
	}
	termbox.SetCursor(cursorX, cursorY)
	termbox.Flush()
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run editor.go <filename>")
		return
	}
	filename = os.Args[1]
	loadFile(filename)

	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	defer termbox.Close()

	draw()
	for {
		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			switch ev.Key {
			case termbox.KeyCtrlQ:
				return
			case termbox.KeyCtrlS:
				saveFile(filename)
			case termbox.KeyArrowLeft:
				if cursorX > 0 {
					cursorX--
				}
			case termbox.KeyArrowRight:
				if cursorX < len(lines[cursorY]) {
					cursorX++
				}
			case termbox.KeyArrowUp:
				if cursorY > 0 {
					cursorY--
					if cursorX > len(lines[cursorY]) {
						cursorX = len(lines[cursorY])
					}
				}
			case termbox.KeyArrowDown:
				if cursorY < len(lines)-1 {
					cursorY++
					if cursorX > len(lines[cursorY]) {
						cursorX = len(lines[cursorY])
					}
				}
			case termbox.KeySpace:
				line := lines[cursorY]
				lines[cursorY] = line[:cursorX] + " " + line[cursorX:]
				cursorX++
			case termbox.KeyBackspace, termbox.KeyBackspace2:
				if cursorX > 0 {
					line := lines[cursorY]
					lines[cursorY] = line[:cursorX-1] + line[cursorX:]
					cursorX--
				} else if cursorY > 0 {
					prev := lines[cursorY-1]
					curr := lines[cursorY]
					cursorX = len(prev)
					lines[cursorY-1] = prev + curr
					lines = append(lines[:cursorY], lines[cursorY+1:]...)
					cursorY--
				}
			case termbox.KeyEnter:
				line := lines[cursorY]
				newLine := line[cursorX:]
				lines[cursorY] = line[:cursorX]
				lines = append(lines[:cursorY+1], append([]string{newLine}, lines[cursorY+1:]...)...)
				cursorY++
				cursorX = 0
			default:
				if ev.Ch != 0 {
					line := lines[cursorY]
					lines[cursorY] = line[:cursorX] + string(ev.Ch) + line[cursorX:]
					cursorX++
				}
			}
			draw()
		case termbox.EventError:
			panic(ev.Err)
		}
	}
}
