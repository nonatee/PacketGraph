package main

import (
	"fmt"
	"log"

	"github.com/google/gopacket"
	"github.com/google/gopacket/pcap"
)

func main() {
	handle, err := pcap.OpenLive(`\Device\NPF_{CC3C6DC2-EA3C-4A08-AFBA-D3EB3E811FD5}`, 1600, true, pcap.BlockForever) // replace "en0" with your interface
	if err != nil {
		log.Fatal(err)
	}
	defer handle.Close()
	var counter = 0
	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
	for packet := range packetSource.Packets() {
		fmt.Println(packet) // prints raw packet info
		counter++
		if counter > 5 {
			close(packetSource.Packets())
		}
	}
}
