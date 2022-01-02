package chaser

import (
	"errors"
	"fmt"
	"net"
)

const (
	GameSet = 0
	TopLeft
	Top
	TopRight
	Left
	Center
	Right
	BottomLeft
	Bottom
	BottomRight
)

const (
	Normal = 0
	ENEMY
	BLOCK
	ITEM
)

type Client struct {
	conn    net.Conn
	port    int16
	address net.IP
	name    string
}

func NewClient(name string, host string, port int16) (*Client, error) {
	ip := net.ParseIP(host)
	if ip == nil {
		return nil, errors.New("ParseIP Error")
	}
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", ip, port))
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	if _, err := conn.Write([]byte(fmt.Sprintln(name))); err != nil {
		fmt.Println(err)
		return nil, err
	}

	return &Client{
		conn:    conn,
		name:    name,
		address: ip.To4(),
		port:    port,
	}, nil
}

func (client *Client) Close() {
	if client.conn != nil {
		if err := client.conn.Close(); err != nil {
			fmt.Println(err)
		}
	}
}

func (client *Client) order(command string) []int8 {
	return []int8{0, 0, 0, 0, 0, 0, 0, 0, 0}
}

func (client *Client) GetReady() ([]int8, error) {
	fmt.Println("GetReady")
	buf := make([]byte, 0, 32)
	n, err := client.conn.Read(buf)
	if err != nil {
		return nil, err
	}
	if len(buf) < 1 {
		fmt.Println("read empty buffer")
		return nil, errors.New("read empty buffer")
	}
	response := string(buf[:n])
	if response[0] != '@' {
		fmt.Println("connection failed")
		client.conn.Close()
		return nil, errors.New("connection failed")
	}
	if _, err := client.conn.Write([]byte("gr")); err != nil {
		return nil, err
	}
	if _, err := client.conn.Read(buf); err != nil {
		return nil, err
	}
	switch buf[0] {
	case '0':
		fmt.Println("GameSet")
		client.conn.Close()
		return nil, errors.New("GameSet")
	case '1':
		fmt.Printf("%v", buf)
	default:
		fmt.Println("responce error")
		return nil, errors.New("responce error")

	}
	return []int8{0, 0, 0, 0, 0, 0, 0, 0, 0}, nil
}

func (client *Client) WalkUp() []int8 {
	fmt.Println("WalkUp")
	return client.order("wu")
}

func (client *Client) WalkLeft() []int8 {
	fmt.Println("WalkLeft")
	return client.order("wl")
}

func (client *Client) WalkRight() []int8 {
	fmt.Println("WalkRight")
	return client.order("wr")
}

func (client *Client) WalkDown() []int8 {
	fmt.Println("WalkDown")
	return client.order("wd")
}

func (client *Client) PutUp() []int8 {
	fmt.Println("PutUp")
	return client.order("pu")
}

func (client *Client) PutLeft() []int8 {
	fmt.Println("PutLeft")
	return client.order("pl")
}

func (client *Client) PutRight() []int8 {
	fmt.Println("PutRight")
	return client.order("pr")
}

func (client *Client) PutDown() []int8 {
	fmt.Println("PutDown")
	return client.order("pd")
}

func (client *Client) LookUp() []int8 {
	fmt.Println("LookUp")
	return client.order("lu")
}

func (client *Client) LookLeft() []int8 {
	fmt.Println("LookLeft")
	return client.order("ll")
}

func (client *Client) LookRight() []int8 {
	fmt.Println("LookRight")
	return client.order("lr")
}

func (client *Client) LookDown() []int8 {
	fmt.Println("LookDown")
	return client.order("ld")
}

func (client *Client) SearchUp() []int8 {
	fmt.Println("SearchUp")
	return client.order("su")
}

func (client *Client) SearchLeft() []int8 {
	fmt.Println("SearchLeft")
	return client.order("sl")
}

func (client *Client) SearchRight() []int8 {
	fmt.Println("SearchRight")
	return client.order("sr")
}

func (client *Client) SearchDown() []int8 {
	fmt.Println("SearchDown")
	return client.order("sd")
}
