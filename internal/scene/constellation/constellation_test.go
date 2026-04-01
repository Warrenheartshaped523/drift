package constellation

import (
	"testing"

	"github.com/phlx0/drift/internal/config"
	"github.com/phlx0/drift/internal/scene"
)

func TestConstellationInitPlacesStarsInBounds(t *testing.T) {
	c := New(config.Default().Scene.Constellation)
	c.Init(80, 24, scene.Themes["cosmic"])

	for i, s := range c.stars {
		if s.x < 0 || s.x >= float64(c.w) {
			t.Errorf("star %d x=%f out of [0, %d)", i, s.x, c.w)
		}
		if s.y < 0 || s.y >= float64(c.h) {
			t.Errorf("star %d y=%f out of [0, %d)", i, s.y, c.h)
		}
	}
}

func TestConstellationStarCountMatchesConfig(t *testing.T) {
	cfg := config.Default().Scene.Constellation
	cfg.StarCount = 50
	c := New(cfg)
	c.Init(80, 24, scene.Themes["cosmic"])
	if len(c.stars) != 50 {
		t.Errorf("expected 50 stars, got %d", len(c.stars))
	}
}

func TestConstellationTwinkleDisabledZerosFreq(t *testing.T) {
	cfg := config.Default().Scene.Constellation
	cfg.Twinkle = false
	c := New(cfg)
	c.Init(80, 24, scene.Themes["cosmic"])
	for i, s := range c.stars {
		if s.twinkleFreq != 0 {
			t.Errorf("star %d has non-zero twinkleFreq %f with twinkle disabled", i, s.twinkleFreq)
		}
	}
}

func TestConstellationTwinkleEnabledNonZeroFreq(t *testing.T) {
	cfg := config.Default().Scene.Constellation
	cfg.Twinkle = true
	c := New(cfg)
	c.Init(80, 24, scene.Themes["cosmic"])
	anyNonZero := false
	for _, s := range c.stars {
		if s.twinkleFreq != 0 {
			anyNonZero = true
			break
		}
	}
	if !anyNonZero {
		t.Error("expected at least one star with non-zero twinkleFreq when twinkle is enabled")
	}
}

func TestConstellationResizeRepositionsOutOfBoundsStars(t *testing.T) {
	c := New(config.Default().Scene.Constellation)
	c.Init(80, 24, scene.Themes["cosmic"])

	// Force a star out of bounds then shrink the terminal.
	c.stars[0].x = 200
	c.stars[0].y = 100
	c.Resize(60, 20)

	s := c.stars[0]
	if s.x >= float64(c.w) || s.y >= float64(c.h) {
		t.Errorf("out-of-bounds star not repositioned after Resize: x=%f y=%f w=%d h=%d", s.x, s.y, c.w, c.h)
	}
}

func TestConstellationUpdateDoesNotPanic(t *testing.T) {
	c := New(config.Default().Scene.Constellation)
	c.Init(80, 24, scene.Themes["cosmic"])
	for i := 0; i < 10; i++ {
		c.Update(0.033)
	}
}

func TestConstellationSmallTerminalDoesNotPanic(t *testing.T) {
	c := New(config.Default().Scene.Constellation)
	c.Init(5, 3, scene.Themes["cosmic"])
	c.Update(0.033)
}
