package chaser

import (
	"log"
	"net"
	"net/textproto"
	"reflect"
	"testing"
)

func TestClient_NewClient(t *testing.T) {
	type fields struct {
		port int
		host string
		name string
	}
	tests := []struct {
		name   string
		fields fields
		want   []int
	}{
		{
			name: "NewClient1",
			fields: fields{
				name: "TestUser",
				host: "127.0.0.1",
				port: 2009,
			},
		},
		{
			name: "NewClient2",
			fields: fields{
				name: "TestUser",
				host: "127.0.0.1",
				port: 2009,
			},
		},
	}
	psock, e := net.Listen("tcp", ":2009")
	if e != nil {
		log.Fatal(e)
		return
	}
	go func() {
		for {
			// Run test server
			defer psock.Close()
			conn, e := psock.Accept()
			if e != nil {
				log.Fatal(e)
				return
			}
			go func(conn *textproto.Conn) {
				_, err := conn.ReadLine()
				if err != nil {
					log.Fatal(err)
				}
			}(textproto.NewConn(conn))
		}
	}()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, err := NewClient(tt.fields.name, tt.fields.host, tt.fields.port)
			if err != nil {
				t.Error(err)
				return
			}
			if !reflect.DeepEqual(tt.fields.name, client.name) {
				t.Errorf("NewClient() = %v, want %v", client.name, tt.fields.name)
			}
			defer client.Close()
		})
	}
}

func TestClient_InvalidHostname(t *testing.T) {
	type fields struct {
		port int
		host string
		name string
	}
	tests := []struct {
		name   string
		fields fields
		want   []int
	}{
		{
			name: "example.com",
			fields: fields{
				name: "TestUser",
				host: "example.com",
				port: 2009,
			},
		},
		{
			name: "Invalid Port",
			fields: fields{
				name: "TestUser",
				host: "127.0.0.1",
				port: 1,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewClient(tt.fields.name, tt.fields.host, tt.fields.port)
			if err == nil {
				t.Error(err)
			}
		})
	}
}

func TestClient_IllegalResponse(t *testing.T) {
	type fields struct {
		port int
		host string
		name string
	}
	tests := []struct {
		name    string
		fields  fields
		want    []int
		gameset bool
	}{
		{
			name: "IllegalResponse",
			fields: fields{
				name: "IllegalResponse",
				host: "127.0.0.1",
				port: 2009,
			},
			want:    nil,
			gameset: false,
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
				conn.PrintfLine("%s", "X")
				conn.Close()
			}()

			_, err := client.GetReady()
			if err == nil {
				t.Error(err)
			}
		})
	}
}
func TestClient_GetReady(t *testing.T) {
	type fields struct {
		port int
		host string
		name string
	}
	tests := []struct {
		name    string
		fields  fields
		want    []int
		gameset bool
	}{
		{
			name: "GetReady",
			fields: fields{
				name: "GetReady",
				host: "127.0.0.1",
				port: 2009,
			},
			want:    []int{0, 0, 0, 0, 0, 0, 0, 0, 0},
			gameset: false,
		},
		{
			name: "GameOver",
			fields: fields{
				name: "GameOver",
				host: "127.0.0.1",
				port: 2009,
			},
			want:    nil,
			gameset: true,
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
					switch tt.name {
					case "GetReady":
						conn.PrintfLine("%v", "1000000000")
					case "GameOver":
						conn.PrintfLine("%v", "0000000000")
					}
				default:
					t.Logf("%v", res)
				}
				conn.Close()
			}()

			got, err := client.GetReady()
			if err != nil && err.Error() != "GameSet" {
				t.Error(err)
			}
			if !reflect.DeepEqual(client.GameSet, tt.gameset) {
				t.Errorf("Client.GetReady() = %v, want %v", client.GameSet, tt.gameset)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Client.GetReady() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestClient_Walk(t *testing.T) {
	type fields struct {
		port int
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
			name: "WalkUp",
			fields: fields{
				name: "WalkUp",
				host: "127.0.0.1",
				port: 2009,
			},
			want: []int{2, 2, 2, 0, 0, 0, 0, 0, 0},
			f:    func(client *Client) ([]int, error) { return client.WalkUp() },
		},
		{
			name: "WalkLeft",
			fields: fields{
				name: "WalkLeft",
				host: "127.0.0.1",
				port: 2009,
			},
			want: []int{2, 0, 0, 2, 0, 0, 2, 0, 0},
			f:    func(client *Client) ([]int, error) { return client.WalkLeft() },
		},
		{
			name: "WalkRight",
			fields: fields{
				name: "WalkRight",
				host: "127.0.0.1",
				port: 2009,
			},
			want: []int{0, 0, 2, 0, 0, 2, 0, 0, 2},
			f:    func(client *Client) ([]int, error) { return client.WalkRight() },
		},
		{
			name: "WalkDown",
			fields: fields{
				name: "WalkDown",
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
				if _, err := conn.ReadLine(); err != nil {
					t.Logf("%v", err)
				}
				conn.Close()
			}()

			got, err := tt.f(client)
			if err != nil {
				t.Error(err)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Client.%s() = %v, want %v", tt.name, got, tt.want)
			}
		})
	}
}

func TestClient_Put(t *testing.T) {
	type fields struct {
		port int
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
			name: "PutUp",
			fields: fields{
				name: "PutUp",
				host: "127.0.0.1",
				port: 2009,
			},
			want: []int{0, 2, 0, 0, 0, 0, 0, 0, 0},
			f:    func(client *Client) ([]int, error) { return client.PutUp() },
		},
		{
			name: "PutLeft",
			fields: fields{
				name: "PutLeft",
				host: "127.0.0.1",
				port: 2009,
			},
			want: []int{0, 0, 0, 2, 0, 0, 0, 0, 0},
			f:    func(client *Client) ([]int, error) { return client.PutLeft() },
		},
		{
			name: "PutRight",
			fields: fields{
				name: "PutRight",
				host: "127.0.0.1",
				port: 2009,
			},
			want: []int{0, 0, 0, 0, 0, 2, 0, 0, 0},
			f:    func(client *Client) ([]int, error) { return client.PutRight() },
		},
		{
			name: "PutDown",
			fields: fields{
				name: "PutDown",
				host: "127.0.0.1",
				port: 2009,
			},
			want: []int{0, 0, 0, 0, 0, 0, 0, 2, 0},
			f:    func(client *Client) ([]int, error) { return client.PutDown() },
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
				case "pu":
					conn.PrintfLine("%v", "1020000000")
				case "pl":
					conn.PrintfLine("%v", "1000200000")
				case "pr":
					conn.PrintfLine("%v", "1000002000")
				case "pd":
					conn.PrintfLine("%v", "1000000020")
				default:
					t.Logf("%v", res)
				}
				if _, err := conn.ReadLine(); err != nil {
					t.Logf("%v", err)
				}
				conn.Close()
			}()

			got, err := tt.f(client)
			if err != nil {
				t.Error(err)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Client.%s() = %v, want %v", tt.name, got, tt.want)
			}
		})
	}
}

func TestClient_Look(t *testing.T) {
	type fields struct {
		port int
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
			name: "LookUp",
			fields: fields{
				name: "LookUp",
				host: "127.0.0.1",
				port: 2009,
			},
			want: []int{0, 2, 0, 0, 0, 0, 0, 0, 0},
			f:    func(client *Client) ([]int, error) { return client.LookUp() },
		},
		{
			name: "LookLeft",
			fields: fields{
				name: "LookLeft",
				host: "127.0.0.1",
				port: 2009,
			},
			want: []int{0, 0, 0, 2, 0, 0, 0, 0, 0},
			f:    func(client *Client) ([]int, error) { return client.LookLeft() },
		},
		{
			name: "LookRight",
			fields: fields{
				name: "LookRight",
				host: "127.0.0.1",
				port: 2009,
			},
			want: []int{0, 0, 0, 0, 0, 2, 0, 0, 0},
			f:    func(client *Client) ([]int, error) { return client.LookRight() },
		},
		{
			name: "LookDown",
			fields: fields{
				name: "LookDown",
				host: "127.0.0.1",
				port: 2009,
			},
			want: []int{0, 0, 0, 0, 0, 0, 0, 2, 0},
			f:    func(client *Client) ([]int, error) { return client.LookDown() },
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
				case "lu":
					conn.PrintfLine("%v", "1020000000")
				case "ll":
					conn.PrintfLine("%v", "1000200000")
				case "lr":
					conn.PrintfLine("%v", "1000002000")
				case "ld":
					conn.PrintfLine("%v", "1000000020")
				default:
					t.Logf("%v", res)
				}
				if _, err := conn.ReadLine(); err != nil {
					t.Logf("%v", err)
				}
				conn.Close()
			}()

			got, err := tt.f(client)
			if err != nil {
				t.Error(err)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Client.%s() = %v, want %v", tt.name, got, tt.want)
			}
		})
	}
}

func TestClient_Search(t *testing.T) {
	type fields struct {
		port int
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
			name: "SearchUp",
			fields: fields{
				name: "SearchUp",
				host: "127.0.0.1",
				port: 2009,
			},
			want: []int{0, 2, 0, 0, 0, 0, 0, 0, 0},
			f:    func(client *Client) ([]int, error) { return client.SearchUp() },
		},
		{
			name: "SearchLeft",
			fields: fields{
				name: "SearchLeft",
				host: "127.0.0.1",
				port: 2009,
			},
			want: []int{0, 0, 0, 2, 0, 0, 0, 0, 0},
			f:    func(client *Client) ([]int, error) { return client.SearchLeft() },
		},
		{
			name: "SearchRight",
			fields: fields{
				name: "SearchRight",
				host: "127.0.0.1",
				port: 2009,
			},
			want: []int{0, 0, 0, 0, 0, 2, 0, 0, 0},
			f:    func(client *Client) ([]int, error) { return client.SearchRight() },
		},
		{
			name: "SearchDown",
			fields: fields{
				name: "SearchDown",
				host: "127.0.0.1",
				port: 2009,
			},
			want: []int{0, 0, 0, 0, 0, 0, 0, 2, 0},
			f:    func(client *Client) ([]int, error) { return client.SearchDown() },
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
				case "su":
					conn.PrintfLine("%v", "1020000000")
				case "sl":
					conn.PrintfLine("%v", "1000200000")
				case "sr":
					conn.PrintfLine("%v", "1000002000")
				case "sd":
					conn.PrintfLine("%v", "1000000020")
				default:
					t.Logf("%v", res)
				}
				if _, err := conn.ReadLine(); err != nil {
					t.Logf("%v", err)
				}
				conn.Close()
			}()

			got, err := tt.f(client)
			if err != nil {
				t.Error(err)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Client.%s() = %v, want %v", tt.name, got, tt.want)
			}
		})
	}
}
