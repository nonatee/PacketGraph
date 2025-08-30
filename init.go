package main

import (
	"log"

	"github.com/google/gopacket"
	"github.com/google/gopacket/pcap"
	"github.com/hajimehoshi/ebiten/v2"
)

func Init() *Game {
	ifs, err := pcap.FindAllDevs()
	if err != nil {
		log.Fatal(err)
	}

	for _, iface := range ifs {
		log.Println("Interface:", iface.Name, iface.Description)
	}
	ebiten.SetWindowSize(1000, 1000)
	ebiten.SetWindowTitle("Circles & Lines")
	handle, err := pcap.OpenLive(`\Device\NPF_{CC3C6DC2-EA3C-4A08-AFBA-D3EB3E811FD5}`, 1600, true, pcap.BlockForever)
	if err != nil {
		log.Fatal(err)
	}
	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
	game := &Game{
		ipMap:      make(map[string]string),
		packetChan: make(chan gopacket.Packet, 100),
		nodeMap:    make(map[string]Node),
	}
	game.handle = handle
	go func() {
		for packet := range packetSource.Packets() {
			game.packetChan <- packet
		}
	}()
	return game
}
