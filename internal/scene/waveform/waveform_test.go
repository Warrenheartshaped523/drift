package waveform

import (
	"testing"

	"github.com/phlx0/drift/internal/config"
	"github.com/phlx0/drift/internal/scene"
)

func TestWaveformInitDoesNotPanic(t *testing.T) {
	wf := New(config.Default().Scene.Waveform)
	wf.Init(80, 24, scene.Themes["cosmic"])
}

func TestWaveformLayersClampedToOne(t *testing.T) {
	cfg := config.Default().Scene.Waveform
	cfg.Layers = 0
	wf := New(cfg)
	wf.Init(80, 24, scene.Themes["cosmic"])
	if len(wf.layers) != 1 {
		t.Errorf("expected 1 layer when cfg.Layers=0, got %d", len(wf.layers))
	}
}

func TestWaveformLayersClampedToThree(t *testing.T) {
	cfg := config.Default().Scene.Waveform
	cfg.Layers = 10
	wf := New(cfg)
	wf.Init(80, 24, scene.Themes["cosmic"])
	if len(wf.layers) != 3 {
		t.Errorf("expected 3 layers when cfg.Layers=10, got %d", len(wf.layers))
	}
}

func TestWaveformResizeUpdatesPixelBuffer(t *testing.T) {
	wf := New(config.Default().Scene.Waveform)
	wf.Init(80, 24, scene.Themes["cosmic"])
	wf.Resize(40, 12)
	if wf.pw != 80 || wf.ph != 48 {
		t.Errorf("after Resize(40,12) expected pw=80 ph=48, got pw=%d ph=%d", wf.pw, wf.ph)
	}
	if len(wf.pixels) != 80 {
		t.Errorf("pixel buffer width mismatch: got %d, want 80", len(wf.pixels))
	}
}

func TestWaveformSmallTerminalDoesNotPanic(t *testing.T) {
	wf := New(config.Default().Scene.Waveform)
	wf.Init(1, 1, scene.Themes["cosmic"])
	wf.Update(0.016)
}

func TestWaveformUpdateAdvancesTime(t *testing.T) {
	wf := New(config.Default().Scene.Waveform)
	wf.Init(80, 24, scene.Themes["cosmic"])
	wf.Update(0.5)
	if wf.time != 0.5 {
		t.Errorf("expected time=0.5, got %f", wf.time)
	}
}
