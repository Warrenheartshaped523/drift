package scene

import (
	"math"

	"github.com/gdamore/tcell/v2"
)

func (o *Orrery) Draw(screen tcell.Screen) {
	if o.w <= 0 || o.h <= 0 {
		return
	}

	o.clearScratch()
	o.drawStars()
	o.drawOrbits()

	for x := 0; x < o.pw; x++ {
		for y := 0; y < o.ph; y++ {
			if b := o.trail[x][y]; b > 0.05 {
				o.stampPixel(x, y, o.trailOwner[x][y], b)
			}
		}
	}

	o.drawSun()
	for _, body := range o.bodies {
		o.drawBody(body)
	}
	o.drawAsteroid()
	o.drawUFO()

	for cx := 0; cx < o.w; cx++ {
		for cy := 0; cy < o.h; cy++ {
			ch, style, ok := o.buildBrailleCell(cx, cy)
			if ok {
				screen.SetContent(cx, cy, ch, nil, style)
			}
		}
	}

}

func (o *Orrery) stampDisc(cx, cy, radius float64, owner uint8, brightness float64, trail bool) {
	minX := int(cx - radius - 1)
	maxX := int(cx + radius + 1)
	minY := int(cy - radius - 1)
	maxY := int(cy + radius + 1)

	for x := minX; x <= maxX; x++ {
		for y := minY; y <= maxY; y++ {
			dx := float64(x) - cx
			dy := float64(y) - cy
			dist := math.Sqrt(dx*dx + dy*dy)
			if dist > radius {
				continue
			}
			falloff := 1 - dist/math.Max(radius, 0.01)
			value := brightness * (0.45 + falloff*0.55)
			if trail {
				o.stampTrailPixel(x, y, owner, value)
			} else {
				o.stampPixel(x, y, owner, value)
			}
		}
	}
}

func (o *Orrery) stampPlanetDisc(cx, cy, radius float64, owner uint8, brightness float64) {
	minX := int(cx - radius - 1)
	maxX := int(cx + radius + 1)
	minY := int(cy - radius - 1)
	maxY := int(cy + radius + 1)

	softEdge := 0.28
	for x := minX; x <= maxX; x++ {
		for y := minY; y <= maxY; y++ {
			dx := float64(x) - cx
			dy := float64(y) - cy
			dist := math.Sqrt(dx*dx + dy*dy)
			if dist > radius {
				continue
			}

			value := brightness
			if dist > radius-softEdge {
				edgeT := (radius - dist) / math.Max(softEdge, 0.01)
				value = brightness * (0.55 + clamp64(edgeT, 0, 1)*0.45)
			}

			o.stampPixel(x, y, owner, value)
		}
	}
}

func (o *Orrery) stampRing(cx, cy, inner, outer float64, owner uint8, brightness float64) {
	minX := int(cx - outer - 1)
	maxX := int(cx + outer + 1)
	minY := int(cy - outer - 1)
	maxY := int(cy + outer + 1)

	for x := minX; x <= maxX; x++ {
		for y := minY; y <= maxY; y++ {
			dx := float64(x) - cx
			dy := float64(y) - cy
			dist := math.Sqrt(dx*dx + dy*dy)
			if dist < inner || dist > outer {
				continue
			}
			falloff := 1 - math.Abs(dist-(inner+outer)/2)/math.Max((outer-inner)/2, 0.01)
			o.stampPixel(x, y, owner, brightness*(0.5+falloff*0.5))
		}
	}
}

func (o *Orrery) stampEllipse(cx, cy, rx, ry float64, owner uint8, brightness float64, trail bool) {
	minX := int(cx - rx - 1)
	maxX := int(cx + rx + 1)
	minY := int(cy - ry - 1)
	maxY := int(cy + ry + 1)

	for x := minX; x <= maxX; x++ {
		for y := minY; y <= maxY; y++ {
			dx := (float64(x) - cx) / math.Max(rx, 0.01)
			dy := (float64(y) - cy) / math.Max(ry, 0.01)
			dist := math.Sqrt(dx*dx + dy*dy)
			if dist > 1 {
				continue
			}
			falloff := 1 - dist
			value := brightness * (0.45 + falloff*0.55)
			if trail {
				o.stampTrailPixel(x, y, owner, value)
			} else {
				o.stampPixel(x, y, owner, value)
			}
		}
	}
}

func (o *Orrery) stampEllipseRing(cx, cy, innerRx, innerRy, outerRx, outerRy float64, owner uint8, brightness float64) {
	minX := int(cx - outerRx - 1)
	maxX := int(cx + outerRx + 1)
	minY := int(cy - outerRy - 1)
	maxY := int(cy + outerRy + 1)

	for x := minX; x <= maxX; x++ {
		for y := minY; y <= maxY; y++ {
			outerDX := (float64(x) - cx) / math.Max(outerRx, 0.01)
			outerDY := (float64(y) - cy) / math.Max(outerRy, 0.01)
			outerDist := math.Sqrt(outerDX*outerDX + outerDY*outerDY)
			if outerDist > 1 {
				continue
			}

			innerDX := (float64(x) - cx) / math.Max(innerRx, 0.01)
			innerDY := (float64(y) - cy) / math.Max(innerRy, 0.01)
			innerDist := math.Sqrt(innerDX*innerDX + innerDY*innerDY)
			if innerDist < 1 {
				continue
			}

			falloff := 1 - outerDist
			o.stampPixel(x, y, owner, brightness*(0.5+falloff*0.5))
		}
	}
}

func (o *Orrery) stampTrailPixel(px, py int, owner uint8, brightness float64) {
	if px < 0 || px >= o.pw || py < 0 || py >= o.ph {
		return
	}
	if brightness > o.trail[px][py] {
		o.trail[px][py] = brightness
		o.trailOwner[px][py] = owner
	}
}

func (o *Orrery) stampPixel(px, py int, owner uint8, brightness float64) {
	if px < 0 || px >= o.pw || py < 0 || py >= o.ph {
		return
	}
	if brightness > o.pixels[px][py] {
		o.pixels[px][py] = brightness
		o.pixelOwner[px][py] = owner
	}
}
