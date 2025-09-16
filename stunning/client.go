package stunning

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/pkg/errors"
)

type Client struct {
}

func NewClient() *Client {
	return &Client{}
}

func (c *Client) Send(ctx context.Context, destAddr string, message Message) error {

	addr, err := net.ResolveUDPAddr("udp", destAddr)
	if err != nil {
		return errors.Wrap(err, "client:send => couldn't resolve udp addr")
	}

	conn, err := net.DialUDP("udp", nil, addr)

	if err != nil {
		return errors.Wrap(err, "client:send => failed to perform connection to destAddr")
	}
	defer conn.Close()

	messageBytes := message.Encode()

	conn.Write(messageBytes)

	buffer := make([]byte, 1024)
	conn.SetReadDeadline(time.Now().Add(30 * time.Second))
	n, addr, err := conn.ReadFromUDP(buffer)
	if err != nil {
		return err
	}
	fmt.Printf("\nReceived from %s: %s\n", addr, string(buffer[:n]))

	fmt.Println("About to decode message")
	resMessage, err := DecodeMessageFromBytes(buffer)
	if err != nil {
		return errors.Wrap(err, "client:send => failed to decode message from bytes")
	}

	fmt.Println(resMessage)

	mappedAddress, err := resMessage.GetMappedAddress()
	if err != nil {
		return errors.Wrap(err, "client:send => failed to get mapped address")
	}

	fmt.Println(mappedAddress.Address, mappedAddress.Port, mappedAddress.Family)

	return nil
}
