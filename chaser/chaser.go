package chaser

import (
	"errors"
	"fmt"
	"log"
	"net"
	"net/textproto"
	"strconv"
	"strings"
)

const (
	GameSet = iota
	UpLeft
	Up
	UpRight
	Left
	Center
	Right
	DownLeft
	Down
	DownRight
)

const (
	Normal = iota
	ENEMY
	BLOCK
	ITEM
)

type Client struct {
	conn    *textproto.Conn
	port    int
	address net.IP
	name    string
	GameSet bool
}

func NewClient(name string, host string, port int) (*Client, error) {
	ip := net.ParseIP(host)
	if ip == nil {
		return nil, errors.New("ParseIP Error")
	}
	conn, err := textproto.Dial("tcp", fmt.Sprintf("%s:%d", ip, port))
	if err != nil {
		log.Println(err)
		return nil, err
	}

	if err := conn.PrintfLine("%v", name); err != nil {
		log.Println(err)
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
			log.Println(err)
		}
	}
}

func strToIntArray(str string) []int {
	sa := strings.Split(str, "")
	ia := make([]int, 0, len(sa))
	for _, s := range sa {
		num, err := strconv.Atoi(s)
		if err != nil {
			log.Println(err)
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
		log.Println("GameSet!!")
		client.conn.Close()
		client.GameSet = true
		return nil, errors.New("GameSet")
	case '1':
		log.Printf("%v\n", response)
	default:
		log.Println("responce error")
		return nil, errors.New("responce error")
	}
	if command != "gr" {
		if err := client.conn.PrintfLine("#"); err != nil {
			return nil, err
		}
	}
	return strToIntArray(response[1:]), nil
}

func (client *Client) GetReady() ([]int, error) {
	log.Println("GetReady")
	res, err := client.conn.ReadLine()
	if err != nil {
		return nil, err
	}
	if res[0] != '@' {
		log.Println("connection failed")
		client.conn.Close()
		return nil, errors.New("connection failed")
	}
	return client.order("gr")
}

func (client *Client) WalkUp() ([]int, error) {
	log.Println("WalkUp")
	return client.order("wu")
}

func (client *Client) WalkLeft() ([]int, error) {
	log.Println("WalkLeft")
	return client.order("wl")
}

func (client *Client) WalkRight() ([]int, error) {
	log.Println("WalkRight")
	return client.order("wr")
}

func (client *Client) WalkDown() ([]int, error) {
	log.Println("WalkDown")
	return client.order("wd")
}

func (client *Client) PutUp() ([]int, error) {
	log.Println("PutUp")
	return client.order("pu")
}

func (client *Client) PutLeft() ([]int, error) {
	log.Println("PutLeft")
	return client.order("pl")
}

func (client *Client) PutRight() ([]int, error) {
	log.Println("PutRight")
	return client.order("pr")
}

func (client *Client) PutDown() ([]int, error) {
	log.Println("PutDown")
	return client.order("pd")
}

func (client *Client) LookUp() ([]int, error) {
	log.Println("LookUp")
	return client.order("lu")
}

func (client *Client) LookLeft() ([]int, error) {
	log.Println("LookLeft")
	return client.order("ll")
}

func (client *Client) LookRight() ([]int, error) {
	log.Println("LookRight")
	return client.order("lr")
}

func (client *Client) LookDown() ([]int, error) {
	log.Println("LookDown")
	return client.order("ld")
}

func (client *Client) SearchUp() ([]int, error) {
	log.Println("SearchUp")
	return client.order("su")
}

func (client *Client) SearchLeft() ([]int, error) {
	log.Println("SearchLeft")
	return client.order("sl")
}

func (client *Client) SearchRight() ([]int, error) {
	log.Println("SearchRight")
	return client.order("sr")
}

func (client *Client) SearchDown() ([]int, error) {
	log.Println("SearchDown")
	return client.order("sd")
}
