package main

import (
	"math"
	"math/rand"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"
)

var curses = []string{
	"DAMN",
	"SHIT",
	"CRAP",
	"CAZZO",
	"MINCHIA",
	"FICA",
	"OSTIA",
	"PUTTANA",
	"MIGNOTTA",
	"KITTEMOURT",
	"FUCK",
	"DIO PAGURO",
}

// One color per frame, like terminal-parrot, in a cyan-to-magenta neon range.
var palette = []int{51, 45, 39, 33, 63, 99, 135, 171, 207, 213}

const (
	glyphRows   = 5
	cycleFrames = 10
	sway        = 2 // max horizontal shift, in columns
	bob         = 1 // max vertical bounce, in rows
	frameDelay  = 75 * time.Millisecond
	height      = glyphRows + bob
)

func render(word string) []string {
	lines := make([]string, glyphRows)
	for _, r := range strings.ToUpper(word) {
		g := glyph(r)
		for i := range lines {
			lines[i] += g[i] + " "
		}
	}
	return lines
}

// frame builds one animation step: the whole word sways sideways and bobs.
func frame(lines []string, f int) string {
	phase := 2 * math.Pi * float64(f) / cycleFrames
	pad := strings.Repeat(" ", sway+int(math.Round(float64(sway)*math.Sin(phase))))
	dy := int(math.Round(float64(bob) * (1 - math.Cos(2*phase)) / 2))

	var b strings.Builder
	b.WriteString("\x1b[" + strconv.Itoa(height) + "A")
	b.WriteString("\x1b[38;5;" + strconv.Itoa(palette[f%len(palette)]) + "m")
	for row := 0; row < height; row++ {
		i := row - (bob - dy)
		if i >= 0 && i < glyphRows {
			b.WriteString(pad)
			b.WriteString(lines[i])
		}
		b.WriteString("\x1b[K\n")
	}
	b.WriteString("\x1b[0m")
	return b.String()
}

func main() {
	word := strings.Join(os.Args[1:], " ")
	if word == "" {
		word = curses[rand.Intn(len(curses))]
	}
	lines := render(word)

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)

	os.Stdout.WriteString("\x1b[?25l" + strings.Repeat("\n", height))
	defer os.Stdout.WriteString("\x1b[?25h")

	tick := time.NewTicker(frameDelay)
	defer tick.Stop()

	for f := 0; ; f++ {
		os.Stdout.WriteString(frame(lines, f))
		select {
		case <-sig:
			return
		case <-tick.C:
		}
	}
}
