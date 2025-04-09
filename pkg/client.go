package chat

// clientWriter listens for outbound messages and writes them to the client connection.
func (s *Server) clientWriter(client *Client) {
	for msg := range client.out {
		_, err := client.conn.Write([]byte(msg + "\n"))
		if err != nil {
			break
		}
	}
}
