package main

import (
	"testing"
)

func TestMp_GetFirstMp(t *testing.T) {
	mp := NewMp("知一码园", Token, Cookie)
	_, _ = mp.GetFirstMp()
}
