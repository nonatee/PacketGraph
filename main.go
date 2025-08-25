package main

import (
	"fmt"
	"image/color"
	"log"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type Game struct {
	packetChan chan gopacket.Packet
	ipMap      map[string]string
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
				}
			}
		}
	default:
		// no packet available right now → just continue
	}
	return nil
}
func (g *Game) Draw(screen *ebiten.Image) {
	vector.DrawFilledCircle(screen, 100, 100, 50, color.RGBA{255, 0, 0, 255}, true)
	vector.StrokeLine(screen, 50, 50, 200, 200, 2, color.RGBA{0, 0, 255, 255}, true)
}
func (g *Game) Layout(w, h int) (int, int) { return 400, 400 }
func main() {
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
		packetChan: make(chan gopacket.Packet, 100), // buffer
	}
	go func() {
		for packet := range packetSource.Packets() {
			game.packetChan <- packet
		}
	}()
	ebiten.RunGame(game)
}
