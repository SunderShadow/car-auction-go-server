package auction

import (
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"strconv"
)

type Server struct {
	connectionUpgrader *websocket.Upgrader
	connections        map[*websocket.Conn]bool
}

func NewServer() *Server {
	return &Server{
		connectionUpgrader: &websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
		},
		connections: make(map[*websocket.Conn]bool),
	}
}

func (server *Server) WriteJSON(event string, data interface{}) error {
	for conn, _ := range server.connections {
		return conn.WriteJSON(map[string]any{
			"event": event,
			"data":  data,
		})
	}

	return nil
}

func (server *Server) ServeWebsocket(w http.ResponseWriter, r *http.Request) {
	conn, _ := server.connectionUpgrader.Upgrade(w, r, nil)

	server.connections[conn] = true

	defer func() {
		conn.Close()
		delete(server.connections, conn)
	}()

	log.Println(`New user connected. Current auction connections: ` + strconv.FormatInt(int64(len(server.connections)), 10))

	for {
		data := make(map[string]any)
		err := conn.ReadJSON(data)

		if err != nil {
			break
		}

	}
}
