// Package dvd implements the classic bouncing logo screensaver.
package dvd

import (
	"math/rand"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/phlx0/drift/internal/config"
	"github.com/phlx0/drift/internal/scene"
)

const cornerFlashDur = 0.6 // seconds the bright flash lasts on a corner hit

// DVD renders a rounded box that bounces around the terminal, changing
// palette color on each wall bounce and flashing bright on corner hits.
type DVD struct {
	w, h  int
	theme scene.Theme
	rng   *rand.Rand

	// logo dimensions derived from the label at Init time
	logoW, logoH int
	label        []rune

	x, y   float64 // top-left position in cells (floating-point for smooth motion)
	vx, vy float64 // velocity in cells/second

	colorIdx    int     // index into theme.Palette
	cornerFlash float64 // counts down to 0 after a corner hit

	cfgSpeed float64
	cfgLabel string
}

func New(cfg config.DVDConfig) *DVD {
	return &DVD{
		cfgSpeed: cfg.Speed,
		cfgLabel: cfg.Label,
	}
}

func (d *DVD) Name() string { return "dvd" }

func (d *DVD) Init(w, h int, t scene.Theme) {
	d.w, d.h = w, h
	d.theme = t
	d.rng = rand.New(rand.NewSource(time.Now().UnixNano()))

	d.label = []rune(d.cfgLabel)

	// Logo layout (e.g. label = "drift"):
	//   ╭───────╮   logoW = len(label) + 4  (│·label·│)
	//   │       │   logoH = 5
	//   │ drift │
	//   │       │
	//   ╰───────╯
	d.logoW = len(d.label) + 4
	d.logoH = 5

	maxX := float64(w - d.logoW)
	maxY := float64(h - d.logoH)
	if maxX < 1 {
		maxX = 1
	}
	if maxY < 1 {
		maxY = 1
	}

	d.x = d.rng.Float64() * maxX
	d.y = d.rng.Float64() * maxY

	// Horizontal speed ~10 cells/s; vertical ~5 cells/s to look ~45° in most
	// fonts where cells are roughly twice as tall as they are wide.
	spd := 10.0 * d.cfgSpeed
	d.vx = spd
	if d.rng.Intn(2) == 0 {
		d.vx = -spd
	}
	d.vy = spd * 0.5
	if d.rng.Intn(2) == 0 {
		d.vy = -d.vy
	}

	d.colorIdx = d.rng.Intn(len(d.theme.Palette))
	d.cornerFlash = 0
}

func (d *DVD) Resize(w, h int) {
	d.w, d.h = w, h
	maxX := float64(w - d.logoW)
	maxY := float64(h - d.logoH)
	if maxX < 0 {
		maxX = 0
	}
	if maxY < 0 {
		maxY = 0
	}
	if d.x > maxX {
		d.x = maxX
	}
	if d.y > maxY {
		d.y = maxY
	}
	if d.x < 0 {
		d.x = 0
	}
	if d.y < 0 {
		d.y = 0
	}
}

func (d *DVD) Update(dt float64) {
	if d.cornerFlash > 0 {
		d.cornerFlash -= dt
		if d.cornerFlash < 0 {
			d.cornerFlash = 0
		}
	}

	d.x += d.vx * dt
	d.y += d.vy * dt

	maxX := float64(d.w - d.logoW)
	maxY := float64(d.h - d.logoH)
	if maxX < 0 {
		maxX = 0
	}
	if maxY < 0 {
		maxY = 0
	}

	hitH, hitV := false, false

	if d.x <= 0 {
		d.x = 0
		d.vx = -d.vx
		hitH = true
	} else if d.x >= maxX {
		d.x = maxX
		d.vx = -d.vx
		hitH = true
	}

	if d.y <= 0 {
		d.y = 0
		d.vy = -d.vy
		hitV = true
	} else if d.y >= maxY {
		d.y = maxY
		d.vy = -d.vy
		hitV = true
	}

	if hitH || hitV {
		d.colorIdx = (d.colorIdx + 1) % len(d.theme.Palette)
	}
	if hitH && hitV {
		d.cornerFlash = cornerFlashDur
	}
}

func (d *DVD) Draw(screen tcell.Screen) {
	if d.w < d.logoW || d.h < d.logoH {
		return
	}

	ox := int(d.x + 0.5)
	oy := int(d.y + 0.5)

	color := d.theme.Palette[d.colorIdx]
	if d.cornerFlash > 0 {
		// Fade from Bright back to the palette color as the flash decays.
		color = scene.Lerp(d.theme.Palette[d.colorIdx], d.theme.Bright, d.cornerFlash/cornerFlashDur)
	}
	st := color.Style()

	// Top:   ╭───...───╮
	screen.SetContent(ox, oy, '╭', nil, st)
	for i := 1; i < d.logoW-1; i++ {
		screen.SetContent(ox+i, oy, '─', nil, st)
	}
	screen.SetContent(ox+d.logoW-1, oy, '╮', nil, st)

	// Blank: │         │
	d.hline(screen, ox, oy+1, st)

	// Label: │ <label> │
	screen.SetContent(ox, oy+2, '│', nil, st)
	screen.SetContent(ox+1, oy+2, ' ', nil, st)
	for i, r := range d.label {
		screen.SetContent(ox+2+i, oy+2, r, nil, st)
	}
	screen.SetContent(ox+2+len(d.label), oy+2, ' ', nil, st)
	screen.SetContent(ox+d.logoW-1, oy+2, '│', nil, st)

	// Blank: │         │
	d.hline(screen, ox, oy+3, st)

	// Bottom: ╰───...───╯
	screen.SetContent(ox, oy+4, '╰', nil, st)
	for i := 1; i < d.logoW-1; i++ {
		screen.SetContent(ox+i, oy+4, '─', nil, st)
	}
	screen.SetContent(ox+d.logoW-1, oy+4, '╯', nil, st)
}

// hline draws a blank interior row: │<spaces>│
func (d *DVD) hline(screen tcell.Screen, ox, oy int, st tcell.Style) {
	screen.SetContent(ox, oy, '│', nil, st)
	for i := 1; i < d.logoW-1; i++ {
		screen.SetContent(ox+i, oy, ' ', nil, st)
	}
	screen.SetContent(ox+d.logoW-1, oy, '│', nil, st)
}
