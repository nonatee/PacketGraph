package main

import (
	"fmt"
	"image/color"
	"math/rand"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"golang.org/x/image/font/basicfont"
)

type Game struct {
	packetChan chan gopacket.Packet
	ipMap      map[string]string
	nodeMap    map[string]Node
}

func (g *Game) Update() error {
	select {
	case packet := <-g.packetChan:
		fmt.Println(packet)
		if dnsLayer := packet.Layer(layers.LayerTypeDNS); dnsLayer != nil {
			dns, _ := dnsLayer.(*layers.DNS)
			if !dns.QR { // query
				for _, q := range dns.Questions {
					fmt.Println("DNS Query:", string(q.Name))
				}
			} else { // response
				for _, ans := range dns.Answers {
					fmt.Println("DNS Response:", string(ans.Name), ans.IP)
					g.ipMap[ans.IP.String()] = string(ans.Name)
					_, ok := g.nodeMap[ans.IP.String()]
					if ok {
						new := g.nodeMap[string(ans.Name)]
						new.counter++
						g.nodeMap[string(ans.Name)] = new
					} else {
						g.nodeMap[string(ans.Name)] = Node{
							Point{500 + rand.Float32()*100, 500 + rand.Float32()*100},
							string(ans.Name),
							Uint32ToRGBA(rand.Uint32()),
							100,
							100}
					}
				}
			}
		}
	default:

	}

	for key := range g.nodeMap {
		newNode := UpdateNodePos(g, g.nodeMap[key])
		g.nodeMap[key] = newNode
	}

	return nil
}
func (g *Game) Draw(screen *ebiten.Image) {
	for key, node := range g.nodeMap {
		vector.DrawFilledCircle(screen, node.position.x, node.position.y, float32(node.counter), node.color, true)
		text.Draw(screen, key, basicfont.Face7x13, int(node.position.x-node.radius), int(node.position.y), color.White)
	}
}
func (g *Game) Layout(w, h int) (int, int) { return 1000, 1000 }
