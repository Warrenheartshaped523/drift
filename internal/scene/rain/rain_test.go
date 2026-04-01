package rain

import (
	"testing"

	"github.com/phlx0/drift/internal/config"
	"github.com/phlx0/drift/internal/scene"
)

func TestRainInitDoesNotPanic(t *testing.T) {
	r := New(config.Default().Scene.Rain)
	r.Init(80, 24, scene.Themes["cosmic"])
}

func TestRainDropCountScalesWithDensity(t *testing.T) {
	cfg := config.Default().Scene.Rain
	r := New(cfg)

	low := r.dropCount(80)

	cfg.Density = 0.8
	r2 := New(cfg)
	high := r2.dropCount(80)

	if high <= low {
		t.Errorf("higher density should produce more drops: low=%d high=%d", low, high)
	}
}

func TestRainDropCountZeroDensityFloor(t *testing.T) {
	cfg := config.Default().Scene.Rain
	cfg.Density = 0
	r := New(cfg)
	if n := r.dropCount(80); n < 1 {
		t.Errorf("dropCount should be >= 1, got %d", n)
	}
}

func TestRainInitPopulatesDrops(t *testing.T) {
	r := New(config.Default().Scene.Rain)
	r.Init(80, 24, scene.Themes["cosmic"])
	if len(r.drops) == 0 {
		t.Error("expected drops to be populated after Init")
	}
}

func TestRainResizeGrowsDrops(t *testing.T) {
	r := New(config.Default().Scene.Rain)
	r.Init(40, 24, scene.Themes["cosmic"])
	before := len(r.drops)
	r.Resize(120, 24)
	if len(r.drops) <= before {
		t.Errorf("expected more drops after widening terminal: before=%d after=%d", before, len(r.drops))
	}
}

func TestRainResizeShrinksDrop(t *testing.T) {
	r := New(config.Default().Scene.Rain)
	r.Init(120, 24, scene.Themes["cosmic"])
	before := len(r.drops)
	r.Resize(40, 24)
	if len(r.drops) >= before {
		t.Errorf("expected fewer drops after narrowing terminal: before=%d after=%d", before, len(r.drops))
	}
}

func TestRainSmallTerminalDoesNotPanic(t *testing.T) {
	r := New(config.Default().Scene.Rain)
	r.Init(1, 1, scene.Themes["cosmic"])
	r.Update(0.016)
}

func TestRainUpdateDoesNotPanic(t *testing.T) {
	r := New(config.Default().Scene.Rain)
	r.Init(80, 24, scene.Themes["cosmic"])
	for i := 0; i < 10; i++ {
		r.Update(0.033)
	}
}
