//go:build unit

package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProfileUpdate_SetNumber(t *testing.T) {
	upd := NewProfileUpdate().SetNumber("score", 42.5)

	assert.Equal(t, map[string]float64{"score": 42.5}, upd.NumberSets)
	assert.Nil(t, upd.NumberIncr)
	assert.Nil(t, upd.StringSets)
}

func TestProfileUpdate_IncrNumber(t *testing.T) {
	upd := NewProfileUpdate().IncrNumber("login_count", 1)

	assert.Equal(t, map[string]float64{"login_count": 1}, upd.NumberIncr)
	assert.Nil(t, upd.NumberSets)
	assert.Nil(t, upd.StringSets)
}

func TestProfileUpdate_SetString(t *testing.T) {
	upd := NewProfileUpdate().SetString("nickname", "alice")

	assert.Equal(t, map[string]string{"nickname": "alice"}, upd.StringSets)
	assert.Nil(t, upd.NumberSets)
	assert.Nil(t, upd.NumberIncr)
}

func TestProfileUpdate_Chaining(t *testing.T) {
	upd := NewProfileUpdate().
		SetNumber("score", 100).
		IncrNumber("login_count", 1).
		SetString("nickname", "bob").
		SetNumber("level", 5).
		IncrNumber("xp", 250)

	assert.Equal(t, map[string]float64{"score": 100, "level": 5}, upd.NumberSets)
	assert.Equal(t, map[string]float64{"login_count": 1, "xp": 250}, upd.NumberIncr)
	assert.Equal(t, map[string]string{"nickname": "bob"}, upd.StringSets)
}

func TestProfileUpdate_OverwriteSameKey(t *testing.T) {
	upd := NewProfileUpdate().
		SetNumber("score", 10).
		SetNumber("score", 20)

	assert.Equal(t, 20.0, upd.NumberSets["score"])
}

func TestProfileUpdate_Empty(t *testing.T) {
	upd := NewProfileUpdate()

	assert.Nil(t, upd.NumberSets)
	assert.Nil(t, upd.NumberIncr)
	assert.Nil(t, upd.StringSets)
}
