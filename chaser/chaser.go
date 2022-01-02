package chaser

import (
	"errors"
	"fmt"
	"net"
	"net/textproto"
	"strconv"
	"strings"
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
	conn    *textproto.Conn
	port    int16
	address net.IP
	name    string
}

func NewClient(name string, host string, port int16) (*Client, error) {
	ip := net.ParseIP(host)
	if ip == nil {
		return nil, errors.New("ParseIP Error")
	}
	conn, err := textproto.Dial("tcp", fmt.Sprintf("%s:%d", ip, port))
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	if err := conn.PrintfLine("%v", name); err != nil {
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

func strToIntArray(str string) []int {
	sa := strings.Split(str, "")
	ia := make([]int, 0, len(sa))
	for _, s := range sa {
		num, err := strconv.Atoi(s)
		if err != nil {
			fmt.Println(err)
			return ia
		}
		ia = append(ia, num)
	}
	return ia
}

func (client *Client) order(command string) ([]int, error) {
	if err := client.conn.PrintfLine("%v", command); err != nil {
		return nil, err
	}
	response, err := client.conn.ReadLine()
	if err != nil {
		return nil, err
	}
	switch response[0] {
	case '0':
		fmt.Println("GameSet")
		client.conn.Close()
		return nil, errors.New("GameSet")
	case '1':
		fmt.Printf("%v", response)
	default:
		fmt.Println("responce error")
		return nil, errors.New("responce error")

	}
	return strToIntArray(response[1:]), nil
}

func (client *Client) GetReady() ([]int, error) {
	fmt.Println("GetReady")
	if response, err := client.conn.ReadLine(); err != nil {
		return nil, err
	} else if response[0] != '@' {
		fmt.Println("connection failed")
		client.conn.Close()
		return nil, errors.New("connection failed")
	}
	return client.order("gr")
}

func (client *Client) WalkUp() ([]int, error) {
	fmt.Println("WalkUp")
	return client.order("wu")
}

func (client *Client) WalkLeft() ([]int, error) {
	fmt.Println("WalkLeft")
	return client.order("wl")
}

func (client *Client) WalkRight() ([]int, error) {
	fmt.Println("WalkRight")
	return client.order("wr")
}

func (client *Client) WalkDown() ([]int, error) {
	fmt.Println("WalkDown")
	return client.order("wd")
}

func (client *Client) PutUp() ([]int, error) {
	fmt.Println("PutUp")
	return client.order("pu")
}

func (client *Client) PutLeft() ([]int, error) {
	fmt.Println("PutLeft")
	return client.order("pl")
}

func (client *Client) PutRight() ([]int, error) {
	fmt.Println("PutRight")
	return client.order("pr")
}

func (client *Client) PutDown() ([]int, error) {
	fmt.Println("PutDown")
	return client.order("pd")
}

func (client *Client) LookUp() ([]int, error) {
	fmt.Println("LookUp")
	return client.order("lu")
}

func (client *Client) LookLeft() ([]int, error) {
	fmt.Println("LookLeft")
	return client.order("ll")
}

func (client *Client) LookRight() ([]int, error) {
	fmt.Println("LookRight")
	return client.order("lr")
}

func (client *Client) LookDown() ([]int, error) {
	fmt.Println("LookDown")
	return client.order("ld")
}

func (client *Client) SearchUp() ([]int, error) {
	fmt.Println("SearchUp")
	return client.order("su")
}

func (client *Client) SearchLeft() ([]int, error) {
	fmt.Println("SearchLeft")
	return client.order("sl")
}

func (client *Client) SearchRight() ([]int, error) {
	fmt.Println("SearchRight")
	return client.order("sr")
}

func (client *Client) SearchDown() ([]int, error) {
	fmt.Println("SearchDown")
	return client.order("sd")
}
