package main

import "image/color"

type Node struct {
	position Point
	name     string
	color    color.RGBA
	counter  uint16
}
type Point struct {
	x float32
	y float32
}

func Uint32ToRGBA(c uint32) color.RGBA {
	return color.RGBA{
		R: uint8((c >> 24) & 0xFF),
		G: uint8((c >> 16) & 0xFF),
		B: uint8((c >> 8) & 0xFF),
		A: uint8(c & 0xFF),
	}
}
