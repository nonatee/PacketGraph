package main

import (
	"fmt"
	"image/color"
	"log"
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
}
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

func (g *Game) Update() error {
	select {
	case packet := <-g.packetChan:
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

					} else {
						g.nodeMap[ans.IP.String()] = Node{Point{rand.Float32(), rand.Float32()}, ans.IP.String(), Uint32ToRGBA(rand.Uint32()), 100}
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
func (g *Game) Layout(w, h int) (int, int) { return 400, 400 }
func main() {
	ifs, err := pcap.FindAllDevs()
	if err != nil {
		log.Fatal(err)
	}

	for _, iface := range ifs {
		log.Println("Interface:", iface.Name, iface.Description)
	}
	ebiten.SetWindowSize(400, 400)
	ebiten.SetWindowTitle("Circles & Lines")
	handle, err := pcap.OpenLive(`\Device\NPF_{CC3C6DC2-EA3C-4A08-AFBA-D3EB3E811FD5}`, 1600, true, pcap.BlockForever)
	if err != nil {
		log.Fatal(err)
	}
	defer handle.Close()
	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
	game := &Game{
		ipMap:      make(map[string]string),
		packetChan: make(chan gopacket.Packet, 100),
		nodeMap:    make(map[string]Node),
	}
	go func() {
		for packet := range packetSource.Packets() {
			game.packetChan <- packet
		}
	}()
	ebiten.RunGame(game)
}
