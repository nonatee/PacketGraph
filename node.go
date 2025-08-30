package main

import (
	"image/color"
	"math"
)

type Node struct {
	position Point
	name     string
	color    color.RGBA
	counter  uint16
	radius   float32
}
type Point struct {
	x float32
	y float32
}

func (p Point) Add(other Point) Point {
	return Point{
		x: p.x + other.x,
		y: p.y + other.y,
	}
}
func (p Point) Mult(multiplier float32) Point {
	return Point{
		x: p.x * multiplier,
		y: p.y * multiplier,
	}
}

func Uint32ToRGBA(c uint32) color.RGBA {
	return color.RGBA{
		R: uint8((c >> 24) & 0xFF),
		G: uint8((c >> 16) & 0xFF),
		B: uint8((c >> 8) & 0xFF),
		A: uint8(c & 0xFF),
	}
}
func Distance(pos1 Point, pos2 Point) float64 {
	return math.Sqrt(math.Pow(float64(pos1.x-pos2.x), 2) + math.Pow(float64(pos1.y-pos2.y), 2))
}
func UpdateNodePos(g *Game, curNode Node) Node {
	for _, node := range g.nodeMap {
		if curNode.radius+node.radius > float32(Distance(node.position, curNode.position)) {
			curNode.position = curNode.position.Add(node.position.Mult(-1)).Mult(0.05).Add(curNode.position)
		}
	}
	return curNode
}
