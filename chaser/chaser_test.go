package chaser

import (
	"net"
	"net/textproto"
	"reflect"
	"testing"
)

func TestClient_GetReady(t *testing.T) {
	type fields struct {
		port int16
		host string
		name string
	}
	tests := []struct {
		name   string
		fields fields
		want   []int
	}{
		{
			name: "Test",
			fields: fields{
				name: "Test",
				host: "127.0.0.1",
				port: 2009,
			},
			want: []int{0, 0, 0, 0, 0, 0, 0, 0, 0},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s, c := net.Pipe()
			client := &Client{
				conn: textproto.NewConn(c),
			}
			defer client.Close()
			go func() {
				conn := textproto.NewConn(s)
				conn.PrintfLine("%s", "@")

				res, err := conn.ReadLine()
				if err != nil {
					t.Error(err)
				}
				switch res {
				case "gr":
					conn.PrintfLine("%v", "1000000000")
				default:
					t.Logf("%v", res)
				}
				conn.Close()
			}()

			got, err := client.GetReady()
			if err != nil {
				t.Error(err)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Client.GetReady() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestClient_Walk(t *testing.T) {
	type fields struct {
		port int16
		host string
		name string
	}
	tests := []struct {
		name   string
		fields fields
		want   []int
		f      func(client *Client) ([]int, error)
	}{
		{
			name: "Test",
			fields: fields{
				name: "Test",
				host: "127.0.0.1",
				port: 2009,
			},
			want: []int{2, 2, 2, 0, 0, 0, 0, 0, 0},
			f:    func(client *Client) ([]int, error) { return client.WalkUp() },
		},
		{
			name: "Test",
			fields: fields{
				name: "Test",
				host: "127.0.0.1",
				port: 2009,
			},
			want: []int{2, 0, 0, 2, 0, 0, 2, 0, 0},
			f:    func(client *Client) ([]int, error) { return client.WalkLeft() },
		},
		{
			name: "Test",
			fields: fields{
				name: "Test",
				host: "127.0.0.1",
				port: 2009,
			},
			want: []int{0, 0, 2, 0, 0, 2, 0, 0, 2},
			f:    func(client *Client) ([]int, error) { return client.WalkRight() },
		},
		{
			name: "Test",
			fields: fields{
				name: "Test",
				host: "127.0.0.1",
				port: 2009,
			},
			want: []int{0, 0, 0, 0, 0, 0, 2, 2, 2},
			f:    func(client *Client) ([]int, error) { return client.WalkDown() },
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s, c := net.Pipe()
			client := &Client{
				conn: textproto.NewConn(c),
			}
			defer client.Close()
			go func() {
				conn := textproto.NewConn(s)
				res, err := conn.ReadLine()
				if err != nil {
					t.Error(err)
				}
				switch res {
				case "wu":
					conn.PrintfLine("%v", "1222000000")
				case "wl":
					conn.PrintfLine("%v", "1200200200")
				case "wr":
					conn.PrintfLine("%v", "1002002002")
				case "wd":
					conn.PrintfLine("%v", "1000000222")
				default:
					t.Logf("%v", res)
				}
				conn.Close()
			}()

			got, err := tt.f(client)
			if err != nil {
				t.Error(err)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Client.GetReady() = %v, want %v", got, tt.want)
			}
		})
	}
}
