package main

import (
	"math"
	"testing"
	"unsafe"

	"github.com/stretchr/testify/assert"
)

type Option func(*GamePerson)

func WithName(name string) Option {
	return func(person *GamePerson) {
		// Копируем символы строки прямо в массив байт
		copy(person.Name[:], name)
		person.nameLen = uint8(len(name))
	}
}

func WithCoordinates(x, y, z int) Option {
	return func(person *GamePerson) {
		person.coordX = int16(x)
		person.coordY = int16(y)
		person.coordZ = int16(z)
	}
}

func WithGold(gold int) Option {
	return func(person *GamePerson) {
		person.goldVal = int16(gold)
	}
}

func WithMana(mana int) Option {
	return func(person *GamePerson) {
		person.manaVal = int16(mana)
	}
}

func WithHealth(health int) Option {
	return func(person *GamePerson) {
		person.healthVal = int16(health)
	}
}

func WithRespect(respect int) Option {
	return func(person *GamePerson) {
		person.respectVal = int16(respect)
	}
}

func WithStrength(strength int) Option {
	return func(person *GamePerson) {
		person.strengthVal = int16(strength)
	}
}

func WithExperience(experience int) Option {
	return func(person *GamePerson) {
		person.expVal = int32(experience)
	}
}

func WithLevel(level int) Option {
	return func(person *GamePerson) {
		person.levelVal = int16(level)
	}
}

func WithHouse() Option {
	return func(person *GamePerson) {
		person.flags |= 1
	}
}

func WithGun() Option {
	return func(person *GamePerson) {
		person.flags |= 2
	}
}

func WithFamily() Option {
	return func(person *GamePerson) {
		person.flags |= 4
	}
}

func WithType(personType int) Option {
	return func(person *GamePerson) {
		person.pType = uint8(personType)
	}
}

const (
	BuilderGamePersonType = iota
	BlacksmithGamePersonType
	WarriorGamePersonType
)

// Итоговый размер структуры: РОВНО 64 БАЙТА
type GamePerson struct {
	name        [42]byte // 42 байта — символы латиницы (inline)
	nameLen     uint8    // 1 байт  — реальная длина строки
	flags       uint8    // 1 байт  — битовая маска (House, Gun, Family)
	pType       uint8    // 1 байт  — тип персонажа
	_           uint8    // 1 байт  — padding для выравнивания следующего int32 на границу 4 байт
	expVal      int32    // 4 байта — опыт (оставляем большим, так как опыта обычно много)
	coordX      int16    // 2 байта
	coordY      int16    // 2 байта
	coordZ      int16    // 2 байта
	goldVal     int16    // 2 байта
	manaVal     int16    // 2 байта
	healthVal   int16    // 2 байта
	respectVal  int16    // 2 байта
	strengthVal int16    // 2 байта
	levelVal    int16    // 2 байта
}

func NewGamePerson(Options ...Option) GamePerson {
	p := GamePerson{}
	for _, opt := range Options {
		opt(&p)
	}
	return p
}

// Геттеры

func (p *GamePerson) Name() string {
	// Безопасно преобразуем часть массива в строку без сохранения указателя на внешние данные
	return string(p.name[:p.nameLen])
}

func (p *GamePerson) X() int {
	return int(p.coordX)
}

func (p *GamePerson) Y() int {
	return int(p.coordY)
}

func (p *GamePerson) Z() int {
	return int(p.coordZ)
}

func (p *GamePerson) Gold() int {
	return int(p.goldVal)
}

func (p *GamePerson) Mana() int {
	return int(p.manaVal)
}

func (p *GamePerson) Health() int {
	return int(p.healthVal)
}

func (p *GamePerson) Respect() int {
	return int(p.respectVal)
}

func (p *GamePerson) Strength() int {
	return int(p.strengthVal)
}

func (p *GamePerson) Experience() int {
	return int(p.expVal)
}

func (p *GamePerson) Level() int {
	return int(p.levelVal)
}

func (p *GamePerson) HasHouse() bool {
	return (p.flags & 1) != 0
}

func (p *GamePerson) HasGun() bool {
	return (p.flags & 2) != 0
}

func (p *GamePerson) HasFamily() bool {
	return (p.flags & 4) != 0
}

func (p *GamePerson) PersonType() int {
	return int(p.pType)
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

	Options := []Option{
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

	person := NewGamePerson(Options...)
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
	assert.True(t, person.HasFamilty())
	assert.False(t, person.HasGun())
	assert.Equal(t, personType, person.Type())
}
