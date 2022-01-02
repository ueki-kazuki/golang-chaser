package chaser

import (
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
		want   []int8
	}{
		{
			name: "Test",
			fields: fields{
				name: "Test",
				host: "127.0.0.1",
				port: 2009,
			},
			want: []int8{0, 0, 0, 0, 0, 0, 0, 0, 0},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, err := NewClient(
				tt.fields.name,
				tt.fields.host,
				tt.fields.port)
			if err != nil {
				t.Error(err)
				return
			}
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
