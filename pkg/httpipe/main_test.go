package httpipe

import (
	"io/ioutil"
	"net/http"
	"reflect"
	"strconv"
	"testing"
	"time"

	"github.com/gorilla/websocket"
)

func TestPages(t *testing.T) {
	htmlMessage := "HTML Page"
	jsMessage := "Javascript Code"
	jsonMessage := "JSON Data"

	tests := []struct {
		url  string
		want string
	}{
		{"/whatever.html", htmlMessage},
		{"/whatever.js", jsMessage},
		{"/whatever.json", jsonMessage},
		{"/whatever", htmlMessage},
		{"/dir/", htmlMessage},
		{"/dir/whatever.html", htmlMessage},
		{"/dir/whatever.js", jsMessage},
		{"/dir/whatever.json", jsonMessage},
	}

	var s = new(Server)
	s.StartBackground(":9999")
	defer s.StopBackground()

	s.ChangeDefault(htmlMessage)
	s.ChangeJS(jsMessage)
	s.ChangeJSON(jsonMessage)

	for _, tt := range tests {
		t.Run(tt.url, func(t *testing.T) {
			url := "http://127.0.0.1:9999" + tt.url
			r, e := http.Get(url)
			if e != nil {
				t.Errorf("%s = ERROR!", tt.url)
				return
			}
			defer r.Body.Close()
			body, _ := ioutil.ReadAll(r.Body)
			got := string(body)
			if got != tt.want {
				t.Errorf("%s = %s, want %s", tt.url, got, tt.want)
				return
			}
		})
	}
}

func TestChanges(t *testing.T) {
	msg1 := "Help me, Obi-Wan Kenobi. You're my only hope"
	msg2 := "I find your lack of faith disturbing"
	msg3 := "Do. Or do not. There is no try"

	url := "http://127.0.0.1:9999/whatever.html"

	tests := []struct {
		name string
		want string
	}{
		{"First change", msg1},
		{"Second change", msg2},
		{"Third change", msg3},
	}

	var s = new(Server)
	s.StartBackground(":9999")
	defer s.StopBackground()
	s.ChangeDefault(msg1)

	for _, tt := range tests {
		t.Run(url, func(t *testing.T) {
			s.ChangeDefault(tt.want)
			r, e := http.Get(url)
			if e != nil {
				t.Errorf("%s = ERROR!", url)
				return
			}
			defer r.Body.Close()
			body, _ := ioutil.ReadAll(r.Body)
			got := string(body)
			if got != tt.want {
				t.Errorf("%s = %s, want %s", tt.name, got, tt.want)
				return
			}
		})
	}
}

func TestReadWS(t *testing.T) {
	msg1 := []string{"1", "10", "100", "1000", "10000", "100000", "1000000", "10000000", "100000000", "1000000000"}
	url := "ws://127.0.0.1:9998/ws/whatever"

	tests := []struct {
		name    string
		howmany int
		want    []string
	}{
		{"Read Message", 1, msg1},
		{"Read Message with 10 simultaneous clients", 10, msg1},
	}

	for _, tt := range tests {
		t.Run(url, func(t *testing.T) {
			// Prepare Server
			var s = new(Server)
			s.StartBackground(":9998")
			defer s.StopBackground()

			// Create connections
			for i := 0; i < tt.howmany; i++ {
				go func() {
					c, _, e := websocket.DefaultDialer.Dial(url, nil)
					if e != nil {
						t.Errorf("Error connecting to server, ERROR = %v", e)
					}
					defer c.Close()

					// Send messages
					for _, s := range tt.want {
						msg := []byte(s)
						e = c.WriteMessage(websocket.TextMessage, msg)
						if e != nil {
							t.Errorf("Error sending message: %v", e)
						}
					}
				}()
			}

			// Wait a little bit
			time.Sleep(time.Millisecond * 100)

			// Check if the amount of connections match
			l := len(s.readWS)
			if l != tt.howmany {
				t.Errorf("Simultaneous connections = %d, want %d", l, tt.howmany)
			}

			// Check if they were received properly
			for i := 0; i < l; i++ {
				var got []string

				for len(s.readWS[i]) > 0 {
					v := s.ReadFromConn(i)
					got = append(got, v)
				}

				if !reflect.DeepEqual(tt.want, got) {
					t.Errorf("%s (%d) = %s, want %s", tt.name, i, got, tt.want)
				}
			}
		})
	}
}

func TestWriteWS(t *testing.T) {
	msg1 := []string{"1", "10", "100", "1000", "10000", "100000", "1000000", "10000000", "100000000", "1000000000"}
	url := "ws://127.0.0.1:9997/ws/whatever"

	tests := []struct {
		name    string
		howmany int
		want    []string
	}{
		{"Write Message", 1, msg1},
		{"Write Message with 10 simultaneous clients", 10, msg1},
	}

	for _, tt := range tests {
		t.Run(url, func(t *testing.T) {
			// Prepare Server
			var s = new(Server)
			s.StartBackground(":9997")
			defer s.StopBackground()

			// Create connections
			var threads []*websocket.Conn
			for i := 0; i < tt.howmany; i++ {
				c, _, e := websocket.DefaultDialer.Dial(url, nil)
				if e != nil {
					t.Errorf("Error connecting to server, ERROR = %v", e)
				}
				defer c.Close()
				threads = append(threads, c)
			}

			// Check if the amount of connections match
			l := len(s.readWS)
			if l != tt.howmany {
				t.Errorf("Simultaneous connections = %d, want %d", l, tt.howmany)
			}

			// Send data for each of them
			for i := 0; i < l; i++ {
				msg1[0] = strconv.Itoa(i)

				// Write messages
				for _, m := range msg1 {
					s.WriteToConn(i, m)
				}
			}

			// Wait a little bit
			time.Sleep(time.Millisecond * 100)

			// Check if they were sent properly
			l = len(threads)
			for i := 0; i < l; i++ {
				var got []string

				for len(tt.want) > len(got) {
					_, msg, e := threads[i].ReadMessage()
					if e != nil {
						t.Errorf("Error reading thread %d: %v", i, e)
					}
					data := s.byteToWsData(msg)
					got = append(got, data)
				}
				tt.want[0] = strconv.Itoa(i)

				if !reflect.DeepEqual(tt.want, got) {
					t.Errorf("%s (%d) = %s, want %s", tt.name, i, got, tt.want)
				}
			}
		})
	}
}
