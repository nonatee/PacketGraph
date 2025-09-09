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
	game := &Game{
		ipMap:      make(map[string]string),
		packetChan: make(chan gopacket.Packet, 100),
		nodeMap:    make(map[string]Node),
	}
	for _, iface := range ifs {
		log.Println("Interface:", iface.Name, iface.Description)
		handle, err := pcap.OpenLive(iface.Name, 1600, true, pcap.BlockForever)
		if err != nil {
			log.Fatal(err)
		}
		packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
		go func() {
			for packet := range packetSource.Packets() {
				game.packetChan <- packet
			}
		}()
	}
	ebiten.SetWindowSize(1000, 1000)
	ebiten.SetWindowTitle("Circles & Lines")

	return game
}
