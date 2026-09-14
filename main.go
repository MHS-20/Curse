package main

import (
	"io"
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
	mask := make([]string, glyphRows)
	for _, r := range strings.ToUpper(word) {
		g := glyph(r)
		for i := range mask {
			mask[i] += g[i] + " "
		}
	}
	return shade(mask)
}

// shade turns a '#' mask into terminal-parrot style art: stroke cells with
// more filled neighbours get denser characters.
func shade(mask []string) []string {
	filled := func(y, x int) bool {
		return y >= 0 && y < len(mask) && x >= 0 && x < len(mask[y]) && mask[y][x] == '#'
	}
	dense := []byte("0OKX")

	out := make([]string, len(mask))
	for y, row := range mask {
		line := []byte(row)
		for x := range line {
			if !filled(y, x) {
				continue
			}
			n := 0
			for _, d := range [][2]int{{-1, 0}, {1, 0}, {0, -1}, {0, 1}} {
				if filled(y+d[0], x+d[1]) {
					n++
				}
			}
			line[x] = dense[max(n-1, 0)]
		}
		out[y] = string(line)
	}
	return out
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
	if st, err := os.Stdin.Stat(); word == "" && err == nil && st.Mode()&os.ModeCharDevice == 0 {
		in, _ := io.ReadAll(os.Stdin)
		word = strings.Join(strings.Fields(string(in)), " ")
	}
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
