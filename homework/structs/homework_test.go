package main

import (
	"math"
	"testing"
	"unsafe"

	"github.com/stretchr/testify/assert"
)

type Option func(*GamePerson)

func WithName(name string) func(*GamePerson) {
	return func(person *GamePerson) {
		copy(person.name[:], []byte(name))
	}
}

func WithCoordinates(x, y, z int) func(*GamePerson) {
	return func(person *GamePerson) {
		person.x = int32(x)
		person.y = int32(y)
		person.z = int32(z)
	}
}

func WithGold(gold int) func(*GamePerson) {
	return func(person *GamePerson) {
		person.goldHouse = person.goldHouse&(^goldMask) | (uint32(gold) << 1)
	}
}

func WithMana(mana int) func(*GamePerson) {
	return func(person *GamePerson) {
		person.specs2 = person.specs2&(^manaMask) | (uint32(mana) << manaShift)
	}
}

func WithHealth(health int) func(*GamePerson) {
	return func(person *GamePerson) {
		person.specs2 = person.specs2&(^healthMask) | (uint32(health) << healthShift)
	}
}

func WithType(personType int) func(*GamePerson) {
	return func(person *GamePerson) {
		person.specs1 = person.specs1&(^typeMask) | (byte(personType) << typeShift)
	}
}

func WithRespect(respect int) func(*GamePerson) {
	return func(person *GamePerson) {
		person.specs2 = person.specs2&(^respectMask) | (uint32(respect) << respectShift)
	}
}

func WithStrength(strength int) func(*GamePerson) {
	return func(person *GamePerson) {
		person.specs2 = person.specs2&(^strengthMask) | (uint32(strength) << strengthShift)
	}
}

func WithExperience(experience int) func(*GamePerson) {
	return func(person *GamePerson) {
		person.specs2 = person.specs2&(^experienceMask) | (uint32(experience) << experienceShift)
	}
}

func WithLevel(level int) func(*GamePerson) {
	return func(person *GamePerson) {
		person.specs1 = person.specs1&(^levelMask) | (byte(level) << levelShift)
	}
}

func WithHouse() func(*GamePerson) {
	return func(person *GamePerson) {
		person.goldHouse = person.goldHouse&(^houseMask) + 1
	}
}

func WithGun() func(*GamePerson) {
	return func(person *GamePerson) {
		person.specs1 = person.specs1&(^gunMask) | (1 << gunShift)
	}
}

func WithFamily() func(*GamePerson) {
	return func(person *GamePerson) {
		person.specs1 = person.specs1&(^familyMask) | (1 << familyShift)
	}
}

const (
	BuilderGamePersonType = iota
	BlacksmithGamePersonType
	WarriorGamePersonType
)

const (
	houseMask uint32 = 0x00000001
	goldMask  uint32 = ^houseMask

	manaMask        uint32 = 0xFFC00000
	healthMask      uint32 = 0x003FF000
	respectMask     uint32 = 0x00000F00
	strengthMask    uint32 = 0x000000F0
	experienceMask  uint32 = 0x0000000F
	manaShift              = 22
	healthShift            = 12
	respectShift           = 8
	strengthShift          = 4
	experienceShift        = 0

	levelMask   byte = 0xF0
	typeMask    byte = 0x0C
	gunMask     byte = 0x02
	familyMask  byte = 0b01
	levelShift       = 4
	typeShift        = 2
	gunShift         = 1
	familyShift      = 0
)

type GamePerson struct {
	name [42]byte

	// level, type, gun, family
	specs1 byte

	x, y, z int32

	goldHouse uint32

	// mana, health, respect, strength, experience
	specs2 uint32
}

func NewGamePerson(options ...Option) GamePerson {
	// need to implement
	p := GamePerson{}

	for _, opt := range options {
		opt(&p)
	}

	return p
}

func (p *GamePerson) Name() string {
	for i := range p.name {
		if p.name[i] == 0 {
			return string(p.name[:i])
		}
	}

	return string(p.name[:])
}

func (p *GamePerson) X() int {
	return int(p.x)
}

func (p *GamePerson) Y() int {
	return int(p.y)
}

func (p *GamePerson) Z() int {
	return int(p.z)
}

func (p *GamePerson) Gold() int {
	return int(p.goldHouse >> 1)
}

func (p *GamePerson) Mana() int {
	return int(p.specs2 & manaMask >> manaShift)
}

func (p *GamePerson) Health() int {
	return int(p.specs2 & healthMask >> healthShift)
}

func (p *GamePerson) Type() int {
	return int(p.specs1 & typeMask >> typeShift)
}

func (p *GamePerson) Respect() int {
	return int(p.specs2 & respectMask >> respectShift)
}

func (p *GamePerson) Strength() int {
	return int(p.specs2 & strengthMask >> strengthShift)
}

func (p *GamePerson) Experience() int {
	return int(p.specs2 & experienceMask >> experienceShift)
}

func (p *GamePerson) Level() int {
	return int(p.specs1 & levelMask >> levelShift)
}

func (p *GamePerson) HasHouse() bool {
	return p.goldHouse&houseMask == 1
}

func (p *GamePerson) HasGun() bool {
	return (p.specs1 & gunMask >> gunShift) == 1
}

func (p *GamePerson) HasFamily() bool {
	return (p.specs1 & familyMask >> familyShift) == 1
}

func TestGamePerson(t *testing.T) {
	assert.LessOrEqual(t, unsafe.Sizeof(GamePerson{}), uintptr(64))

	const x, y, z = math.MinInt32, math.MaxInt32, 0
	const name = "aaaaaaaaaaaaa_bbbbbbbbbbbbb_cccccccccccccc"
	const personType = BuilderGamePersonType
	const gold = math.MaxInt32
	const mana = 1000
	const health = 1000
	const respect = 10
	const strength = 10
	const experience = 10
	const level = 10

	options := []Option{
		WithName(name),
		WithCoordinates(x, y, z),
		WithGold(gold),
		WithMana(mana),
		WithHealth(health),
		WithRespect(respect),
		WithStrength(strength),
		WithExperience(experience),
		WithLevel(level),
		WithHouse(),
		WithFamily(),
		WithType(personType),
	}

	person := NewGamePerson(options...)
	assert.Equal(t, name, person.Name())
	assert.Equal(t, x, person.X())
	assert.Equal(t, y, person.Y())
	assert.Equal(t, z, person.Z())
	assert.Equal(t, gold, person.Gold())
	assert.Equal(t, mana, person.Mana())
	assert.Equal(t, health, person.Health())
	assert.Equal(t, respect, person.Respect())
	assert.Equal(t, strength, person.Strength())
	assert.Equal(t, experience, person.Experience())
	assert.Equal(t, level, person.Level())
	assert.True(t, person.HasHouse())
	assert.True(t, person.HasFamily())
	assert.False(t, person.HasGun())
	assert.Equal(t, personType, person.Type())
}
