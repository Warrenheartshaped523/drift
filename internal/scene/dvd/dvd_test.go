package dvd

import (
	"testing"

	"github.com/gdamore/tcell/v2"
	"github.com/phlx0/drift/internal/config"
	"github.com/phlx0/drift/internal/scene"
)

func defaultDVD() *DVD {
	d := New(config.Default().Scene.DVD)
	d.Init(80, 24, scene.Themes["cosmic"])
	return d
}

func TestDVDInitDoesNotPanic(t *testing.T) {
	defaultDVD()
}

func TestDVDLogoDimensions(t *testing.T) {
	cfg := config.Default().Scene.DVD
	cfg.Label = "drift" // 5 chars
	d := New(cfg)
	d.Init(80, 24, scene.Themes["cosmic"])

	if d.logoW != 9 {
		t.Errorf("logoW: got %d, want 9 for label %q", d.logoW, cfg.Label)
	}
	if d.logoH != 5 {
		t.Errorf("logoH: got %d, want 5", d.logoH)
	}
}

func TestDVDInitPositionInBounds(t *testing.T) {
	d := defaultDVD()
	if d.x < 0 || d.x > float64(d.w-d.logoW) {
		t.Errorf("initial x=%f out of [0, %d]", d.x, d.w-d.logoW)
	}
	if d.y < 0 || d.y > float64(d.h-d.logoH) {
		t.Errorf("initial y=%f out of [0, %d]", d.y, d.h-d.logoH)
	}
}

func TestDVDBounceReverseVX(t *testing.T) {
	d := defaultDVD()
	d.x = float64(d.w - d.logoW) // push to right wall
	d.vx = 5.0
	vxBefore := d.vx

	d.Update(0.1)

	if d.vx == vxBefore {
		t.Error("vx should have reversed after hitting the right wall")
	}
}

func TestDVDBounceReverseVY(t *testing.T) {
	d := defaultDVD()
	d.y = float64(d.h - d.logoH) // push to bottom wall
	d.vy = 5.0
	vyBefore := d.vy

	d.Update(0.1)

	if d.vy == vyBefore {
		t.Error("vy should have reversed after hitting the bottom wall")
	}
}

func TestDVDColorChangesOnBounce(t *testing.T) {
	d := defaultDVD()
	d.x = float64(d.w - d.logoW)
	d.vx = 5.0
	colorBefore := d.colorIdx

	d.Update(0.1)

	if d.colorIdx == colorBefore {
		t.Error("colorIdx should change after a wall bounce")
	}
}

func TestDVDCornerFlashOnCornerHit(t *testing.T) {
	d := defaultDVD()
	// Position at bottom-right corner, moving toward it.
	d.x = float64(d.w - d.logoW)
	d.y = float64(d.h - d.logoH)
	d.vx = 5.0
	d.vy = 5.0

	d.Update(0.1)

	if d.cornerFlash <= 0 {
		t.Error("cornerFlash should be > 0 after hitting a corner")
	}
}

func TestDVDNoCornerFlashOnSingleWallHit(t *testing.T) {
	d := defaultDVD()
	// Hit only the right wall (not a corner).
	d.x = float64(d.w - d.logoW)
	d.y = 10
	d.vx = 5.0
	d.vy = 0.1

	d.Update(0.1)

	// Unless y also hits a wall, no corner flash.
	if d.y <= 0 || d.y >= float64(d.h-d.logoH) {
		t.Skip("y also hit a wall; corner flash is expected")
	}
	if d.cornerFlash > 0 {
		t.Error("cornerFlash should be 0 when only one wall was hit")
	}
}

func TestDVDCornerFlashDecays(t *testing.T) {
	d := defaultDVD()
	d.cornerFlash = cornerFlashDur

	d.Update(0.1)

	if d.cornerFlash >= cornerFlashDur {
		t.Error("cornerFlash should decrease over time")
	}
}

func TestDVDPositionStaysInBoundsAfterManyUpdates(t *testing.T) {
	d := defaultDVD()
	for i := 0; i < 1000; i++ {
		d.Update(0.033)
		if d.x < 0 || d.x > float64(d.w-d.logoW) {
			t.Fatalf("x=%f out of bounds after %d updates", d.x, i+1)
		}
		if d.y < 0 || d.y > float64(d.h-d.logoH) {
			t.Fatalf("y=%f out of bounds after %d updates", d.y, i+1)
		}
	}
}

func TestDVDDrawDoesNotPanic(t *testing.T) {
	d := defaultDVD()

	screen := tcell.NewSimulationScreen("")
	if err := screen.Init(); err != nil {
		t.Fatalf("screen init: %v", err)
	}
	defer screen.Fini()
	screen.SetSize(80, 24)

	d.Draw(screen)
}

func TestDVDSmallTerminalDoesNotPanic(t *testing.T) {
	d := New(config.Default().Scene.DVD)
	d.Init(5, 3, scene.Themes["cosmic"])

	screen := tcell.NewSimulationScreen("")
	if err := screen.Init(); err != nil {
		t.Fatalf("screen init: %v", err)
	}
	defer screen.Fini()
	screen.SetSize(5, 3)

	d.Update(0.033)
	d.Draw(screen)
}

func TestDVDResizeClampPosition(t *testing.T) {
	d := defaultDVD()
	d.x = 70
	d.y = 20

	d.Resize(20, 10)

	maxX := float64(d.w - d.logoW)
	maxY := float64(d.h - d.logoH)
	if maxX < 0 {
		maxX = 0
	}
	if maxY < 0 {
		maxY = 0
	}
	if d.x > maxX {
		t.Errorf("x=%f not clamped after Resize; max=%f", d.x, maxX)
	}
	if d.y > maxY {
		t.Errorf("y=%f not clamped after Resize; max=%f", d.y, maxY)
	}
}

func TestDVDEmptyLabelDoesNotPanic(t *testing.T) {
	cfg := config.Default().Scene.DVD
	cfg.Label = ""
	d := New(cfg)
	d.Init(80, 24, scene.Themes["cosmic"])

	screen := tcell.NewSimulationScreen("")
	if err := screen.Init(); err != nil {
		t.Fatalf("screen init: %v", err)
	}
	defer screen.Fini()
	screen.SetSize(80, 24)

	d.Update(0.033)
	d.Draw(screen)
}
