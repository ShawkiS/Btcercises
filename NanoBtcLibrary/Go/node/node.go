package node

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"reflect"

	rpc "github.com/Btcercises/NanoBtcLibrary/network/rpc"
	"github.com/Btcercises/NanoBtcLibrary/network/util"
)

type Node struct {
	Connection *net.TCPConn
	Testnet    bool
	Logging    bool
}

type NodeConnectOption func(*Node) *net.TCPConn

func WithHostName(host string, ports ...int) NodeConnectOption {
	return func(node *Node) *net.TCPConn {
		var port int
		if len(ports) == 0 {
			port = 8333
			if node.Testnet {
				port = 18333
			}
		} else {
			port = ports[0]
		}
		result, err := net.Dial("tcp", fmt.Sprintf("%s:%d", host, port))
		if err != nil {
			panic(err)
		}
		return result.(*net.TCPConn)
	}
}

func NewNode(option NodeConnectOption, testnet bool, logging bool) *Node {
	result := &Node{
		Testnet: testnet,
		Logging: logging,
	}
	result.Connection = option(result)
	return result
}

func (node *Node) Close() error {
	return node.Connection.Close()
}

func (node *Node) Handshake() (bool, error) {
	if ok, err := node.Send(rpc.NewVersionMessage(nil)); !ok {
		return ok, err
	}
	verack, err := node.WaitFor(rpc.VerackMessageOption())
	if err != nil {
		return false, err
	}
	if verack == nil {
		return false, errors.New("no response received")
	}
	return true, nil
}

func (node *Node) Send(message rpc.Message) (bool, error) {
	envelope := rpc.NewEnvelope(message.Command(), message.Serialize(), node.Testnet)
	if node.Logging {
		fmt.Fprintf(os.Stdout, "sending: %v\n", envelope)
	}
	_, err := node.Connection.Write(envelope.Serialize())
	if err != nil {
		return false, err
	}
	return true, nil
}

func (node *Node) Read() (*rpc.NetworkEnvelope, error) {
	bufCh := make(chan []byte)
	errCh := make(chan error)
	go func(conn *net.TCPConn, bufCh chan []byte, errCh chan error) {
		header := make([]byte, 24)
		_, err := io.ReadFull(conn, header)
		if err != nil {
			errCh <- err
			return
		}
		bufCh <- header
		payloadLength := int(util.LittleEndianToInt32(header[16:20]))
		payload := make([]byte, payloadLength)
		_, err = io.ReadFull(conn, payload)
		if err != nil {
			errCh <- err
			return
		}
		bufCh <- payload
		close(bufCh)
	}(node.Connection, bufCh, errCh)
	response := make([]byte, 0)
	for buf := range bufCh {
		select {
		case err := <-errCh:
			return nil, err
		default:
		}
		response = append(response, buf...)
	}
	return rpc.ParseEnvelope(bytes.NewReader(response), node.Testnet), nil
}

func (node *Node) WaitFor(messageTypes ...rpc.ReceiveMessageTypeOption) (rpc.Message, error) {
	commands := make(map[string]rpc.Message)
	for _, option := range messageTypes {
		messageType := option()
		message, ok := reflect.New(messageType.Elem()).Interface().(rpc.Message)
		if !ok {
			panic("Failed to cast to Message type!")
		}
		commands[string(message.Command())] = message
	}
	for {
		envelope, err := node.Read()
		if err != nil {
			return nil, err
		}
		command := string(envelope.Command)
		if node.Logging {
			fmt.Fprintf(os.Stdout, "received: %s\n", command)
		}
		switch command {
		case "version":
			node.Send(rpc.NewVerackMessage())
		case "ping":
			node.Send(rpc.NewPongMessage(envelope.Payload))
		}
		if result, ok := commands[command]; ok {
			result.Parse(envelope.Stream())
			return result, nil
		}
	}
}
