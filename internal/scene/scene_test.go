package scene

import (
	"math"
	"testing"

	"github.com/phlx0/drift/internal/config"
)

func TestThemesAllHavePalette(t *testing.T) {
	for name, theme := range Themes {
		if len(theme.Palette) == 0 {
			t.Errorf("theme %q has empty Palette", name)
		}
		if len(theme.Dim) == 0 {
			t.Errorf("theme %q has empty Dim", name)
		}
		if len(theme.Palette) != len(theme.Dim) {
			t.Errorf("theme %q: Palette len %d != Dim len %d", name, len(theme.Palette), len(theme.Dim))
		}
	}
}

func TestByName(t *testing.T) {
	for _, name := range Names() {
		if s := ByName(name); s == nil {
			t.Errorf("ByName(%q) returned nil", name)
		}
	}
	if s := ByName("does-not-exist"); s != nil {
		t.Errorf("ByName(unknown) should return nil, got %v", s)
	}
	if s := ByName("orrery"); s == nil {
		t.Fatal("ByName(orrery) returned nil")
	}
}

func TestLerp(t *testing.T) {
	black := RGBColor{0, 0, 0}
	white := RGBColor{255, 255, 255}

	mid := Lerp(black, white, 0.5)
	if mid.R != 127 && mid.R != 128 {
		t.Errorf("Lerp(black, white, 0.5).R = %d, want ~127", mid.R)
	}

	at0 := Lerp(black, white, 0)
	if at0 != black {
		t.Errorf("Lerp at t=0 should equal a, got %v", at0)
	}

	at1 := Lerp(black, white, 1)
	if at1 != white {
		t.Errorf("Lerp at t=1 should equal b, got %v", at1)
	}

	// Clamp test: t > 1 should not exceed b
	over := Lerp(black, white, 2)
	if over != white {
		t.Errorf("Lerp with t>1 should clamp to b, got %v", over)
	}
}

func TestScenesInitDoNotPanic(t *testing.T) {
	theme := Themes["cosmic"]
	for _, s := range All() {
		t.Run(s.Name(), func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					t.Errorf("scene %q panicked on Init: %v", s.Name(), r)
				}
			}()
			s.Init(120, 40, theme)
			s.Update(0.033)
			s.Resize(80, 24)
		})
	}
}

func TestOrreryBuildStarsHandlesNarrowTerminal(t *testing.T) {
	o := NewOrrery(config.Default().Scene.Orrery)
	o.Init(8, 12, Themes["cosmic"])

	for _, star := range o.stars {
		if star.x < 0 || star.x >= float64(o.pw) {
			t.Fatalf("star x out of bounds for narrow terminal: %f not in [0,%d)", star.x, o.pw)
		}
		if star.y < 0 || star.y >= float64(o.ph) {
			t.Fatalf("star y out of bounds for narrow terminal: %f not in [0,%d)", star.y, o.ph)
		}
	}
}

func TestOrreryResizePreservesActiveFlybys(t *testing.T) {
	o := NewOrrery(config.Default().Scene.Orrery)
	o.Init(120, 40, Themes["cosmic"])

	o.asteroid = orreryAsteroid{
		active: true,
		x:      11.5,
		y:      17.25,
		vx:     3.5,
		vy:     -2.25,
		size:   1.2,
	}
	o.ufo = orreryUFO{
		active:    true,
		x:         73.0,
		y:         21.0,
		vx:        -5.0,
		vy:        1.5,
		targetX:   66.0,
		targetY:   18.0,
		hoverTime: 0.75,
	}

	o.Resize(100, 32)

	if !o.asteroid.active {
		t.Fatal("active asteroid was reset during resize")
	}
	if o.asteroid.x != 11.5 || o.asteroid.y != 17.25 || o.asteroid.vx != 3.5 || o.asteroid.vy != -2.25 {
		t.Fatalf("asteroid state changed during resize: %+v", o.asteroid)
	}

	if !o.ufo.active {
		t.Fatal("active UFO was reset during resize")
	}
	if o.ufo.x != 73.0 || o.ufo.y != 21.0 || o.ufo.targetX != 66.0 || o.ufo.targetY != 18.0 || o.ufo.hoverTime != 0.75 {
		t.Fatalf("ufo state changed during resize: %+v", o.ufo)
	}
}

func TestOrrerySpawnAsteroidProducesFlybyVelocity(t *testing.T) {
	o := NewOrrery(config.Default().Scene.Orrery)
	o.Init(120, 40, Themes["cosmic"])

	spawned := false
	for i := 0; i < 64; i++ {
		o.spawnAsteroid()

		if !o.asteroid.active {
			if o.asteroid.cooldown <= 0 {
				t.Fatalf("spawnAsteroid backed off without setting a retry cooldown on iteration %d", i)
			}
			continue
		}
		spawned = true

		speed := math.Hypot(o.asteroid.vx, o.asteroid.vy)
		if speed < 14.0 {
			t.Fatalf("spawnAsteroid produced degenerate velocity on iteration %d: vx=%f vy=%f", i, o.asteroid.vx, o.asteroid.vy)
		}
	}

	if !spawned {
		t.Fatal("spawnAsteroid never produced a valid active asteroid in 64 attempts")
	}
}

func TestOrreryRingedPlanetHaloClearsIntersectingCells(t *testing.T) {
	o := NewOrrery(config.Default().Scene.Orrery)
	o.Init(120, 40, Themes["cosmic"])

	ringedIndex := -1
	for i, body := range o.bodies {
		if body.hasRing {
			ringedIndex = i
			break
		}
	}
	if ringedIndex == -1 {
		t.Fatal("expected ringed planet in orrery bodies")
	}

	body := o.bodies[ringedIndex]
	haloRadius := body.size + 2.8
	cellX := int(body.x) / 2
	cellY := int(body.y) / 4

	for subCol := 0; subCol < 2; subCol++ {
		for subRow := 0; subRow < 4; subRow++ {
			px := cellX*2 + subCol
			py := cellY*4 + subRow
			o.pixels[px][py] = 0.5
			o.pixelOwner[px][py] = orreryOrbitOwner
		}
	}

	o.clearPlanetHalo(body.x, body.y, haloRadius)

	for subCol := 0; subCol < 2; subCol++ {
		for subRow := 0; subRow < 4; subRow++ {
			px := cellX*2 + subCol
			py := cellY*4 + subRow
			if o.pixels[px][py] != 0 || o.pixelOwner[px][py] != 0 {
				t.Fatalf("ringed planet halo did not clear intersecting cell at (%d,%d)", px, py)
			}
		}
	}
}
