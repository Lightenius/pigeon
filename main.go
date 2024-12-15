package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"

	"github.com/pion/stun"
)

const ownPort int = 5515
const targetPort string = ":5515"

func main() {
	println("Your IP Address is -> ", getMyIP())
	println("Please Enter Target IP Address -> ")
	var targetIp string
	fmt.Scanln(&targetIp)

	go recieve()

	send(strings.TrimSuffix(targetIp, "\n"))

}

func recieve() {
	c, err := net.ListenUDP("udp", &net.UDPAddr{IP: net.ParseIP("0.0.0.0"), Port: ownPort})
	if err != nil {
		panic(err)
	}
	defer c.Close()

	buffer := make([]byte, 1024)

	for {
		n, sender, err := c.ReadFromUDP(buffer)
		if err != nil {
			panic(err)
		}

		fmt.Println(sender.IP.String(), " ::> ", string(buffer[:n]))
	}
}

func send(adres string) {
	adres = adres + targetPort

	c, err := net.Dial("udp", adres)
	if err != nil {
		panic(err)
	}
	defer c.Close()

	reader := bufio.NewReader(os.Stdin)
	for {
		print("\nEnter Message -> ")
		inp, err := reader.ReadString('\n')
		if err != nil {
			println(err)
			continue
		}

		c.Write([]byte(inp))

	}

}

func getMyIP() string {
	u, err := stun.ParseURI("stun:stun.l.google.com:19302")
	if err != nil {
		panic(err)
	}

	var ipstr string
	// Creating a "connection" to STUN server.
	c, err := stun.DialURI(u, &stun.DialConfig{})
	if err != nil {
		panic(err)
	}
	// Building binding request with random transaction id.
	message := stun.MustBuild(stun.TransactionID, stun.BindingRequest)
	// Sending request to STUN server, waiting for response message.
	if err := c.Do(message, func(res stun.Event) {
		if res.Error != nil {
			panic(res.Error)
		}
		// Decoding XOR-MAPPED-ADDRESS attribute from message.
		var xorAddr stun.XORMappedAddress
		if err := xorAddr.GetFrom(res.Message); err != nil {
			panic(err)
		}
		ipstr = xorAddr.IP.String()
	}); err != nil {
		panic(err)
	}
	return ipstr
}
