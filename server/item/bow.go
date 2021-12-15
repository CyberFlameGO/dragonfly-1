package item

import (
	"github.com/df-mc/dragonfly/server/world"
	"github.com/df-mc/dragonfly/server/world/sound"
	"github.com/go-gl/mathgl/mgl64"
	"math"
	"time"
)

// Bow is a ranged weapon that fires arrows.
type Bow struct{}

// MaxCount always returns 1.
func (Bow) MaxCount() int {
	return 1
}

// DurabilityInfo ...
func (Bow) DurabilityInfo() DurabilityInfo {
	return DurabilityInfo{
		MaxDurability: 385,
		BrokenItem:    simpleItem(Stack{}),
	}
}

// Release ...
func (Bow) Release(releaser Releaser, duration time.Duration, ctx *UseContext) {
	if !releaser.GameMode().Visible() {
		return
	}

	arrow := NewStack(Arrow{}, 1)
	inCreative := releaser.GameMode().CreativeInventory()
	if !inCreative {
		if !ctx.Require(arrow) {
			return
		}
		ctx.DamageItem(1)
		ctx.Consume(arrow)
	}

	ticks := duration.Milliseconds() / 50
	if ticks < 3 {
		return
	}

	t := float64(ticks) / 20
	force := math.Min((t*t+t*2)/3, 1)
	if force < 0.1 {
		return
	}

	rYaw, rPitch := releaser.Rotation()
	yaw, pitch := -rYaw, -rPitch
	if rYaw > 180 {
		yaw = 360 - rYaw
	}

	proj, ok := world.EntityByName("minecraft:arrow")
	if !ok {
		return
	}

	p, ok := proj.(interface {
		New(pos, vel mgl64.Vec3, yaw, pitch float64, critical, canPickup, creativePickup bool, baseDamage float64) world.Entity
	})
	if !ok {
		return
	}
	e := p.New(eyePosition(releaser), directionVector(releaser).Mul(force*3), yaw, pitch, force >= 1, true, inCreative, 2)
	if o, ok := e.(owned); ok {
		o.Own(releaser)
	}

	releaser.PlaySound(sound.BowShoot{})
	releaser.World().AddEntity(e)
}

// EncodeItem ...
func (Bow) EncodeItem() (name string, meta int16) {
	return "minecraft:bow", 0
}
