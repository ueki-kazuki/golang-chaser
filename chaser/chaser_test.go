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
