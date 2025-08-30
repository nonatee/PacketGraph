package main

import (
	"fmt"
	"math/rand"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type Game struct {
	packetChan chan gopacket.Packet
	ipMap      map[string]string
	nodeMap    map[string]Node
	handle     *pcap.Handle
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
						new := g.nodeMap[ans.IP.String()]
						new.counter++
						g.nodeMap[ans.IP.String()] = new
					} else {
						g.nodeMap[ans.IP.String()] = Node{Point{500 + rand.Float32(), 500 + rand.Float32()}, ans.IP.String(), Uint32ToRGBA(rand.Uint32()), 100}
					}
				}
			}
		}
	default:

	}
	return nil
}
func (g *Game) Draw(screen *ebiten.Image) {
	for _, node := range g.nodeMap {
		vector.DrawFilledCircle(screen, node.position.x, node.position.y, float32(node.counter), node.color, true)
	}
}
func (g *Game) Layout(w, h int) (int, int) { return 1000, 1000 }
