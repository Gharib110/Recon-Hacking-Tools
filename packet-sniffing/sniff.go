package packet_sniffing

import (
	"bytes"
	"fmt"
	"github.com/google/gopacket"
	"github.com/google/gopacket/pcap"
	"log"
)

// SniffTCPPort get network interface and tcp port number to sniff
func sniffTCPPort(port int, iface string) {
	var snapshotLen = int32(1600)
	var promiscuous = false
	var timeout = pcap.BlockForever
	var filter = fmt.Sprintf("tcp and port %d", port)
	var devFound = false

	devices, err := pcap.FindAllDevs()
	if err != nil {
		log.Panicln(err)
	}

	for _, device := range devices {
		if device.Name == iface {
			devFound = true
		}
	}
	if !devFound {
		log.Panicf("Device named '%s' does not exist\n", iface)
	}

	handle, err := pcap.OpenLive(iface, snapshotLen, promiscuous, timeout)
	if err != nil {
		log.Panicln(err)
	}
	defer handle.Close()

	if err = handle.SetBPFFilter(filter); err != nil {
		log.Panicln(err)
	}

	source := gopacket.NewPacketSource(handle, handle.LinkType())
	for packet := range source.Packets() {
		fmt.Println(packet)
	}

}

// SniffFTP try to get USER and PASS from raw FTP proto packets
func SniffFTP(iface string) {
	var (
		snapshotLen = int32(1600)
		promiscuous = false
		timeout     = pcap.BlockForever
		filter      = "tcp and dst port 21"
		devFound    = false
	)

	devices, err := pcap.FindAllDevs()
	if err != nil {
		log.Panicln(err)
	}

	for _, device := range devices {
		if device.Name == iface {
			devFound = true
		}
	}
	if !devFound {
		log.Panicf("Device named '%s' does not exist\n", iface)
	}

	handle, err := pcap.OpenLive(iface, snapshotLen, promiscuous, timeout)
	if err != nil {
		log.Panicln(err)
	}
	defer handle.Close()

	if err := handle.SetBPFFilter(filter); err != nil {
		log.Panicln(err)
	}

	source := gopacket.NewPacketSource(handle, handle.LinkType())
	for packet := range source.Packets() {
		appLayer := packet.ApplicationLayer()
		if appLayer == nil {
			continue
		}
		payload := appLayer.Payload()
		if bytes.Contains(payload, []byte("USER")) {
			fmt.Print(string(payload))
		} else if bytes.Contains(payload, []byte("PASS")) {
			fmt.Print(string(payload))
		}
	}
}
