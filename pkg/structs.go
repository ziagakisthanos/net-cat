package chat

import (
	"net"
	"sync"
	"os"
)

// Maximum connections
const maxClients = 10

var welcomeBanner = "Welcome to TCP-Chat!\n" +
	"         _nnnn_\n" +
	"        dGGGGMMb\n" +
	"       @p~qp~~qMb\n" +
	"       M|@||@) M|\n" +
	"       @,----.JM|\n" +
	"      JS^\\__/  qKL\n" +
	"     dZP        qKRb\n" +
	"    dZP          qKKb\n" +
	"   fZP            SMMb\n" +
	"   HZM            MMMM\n" +
	"   FqM            MMMM\n" +
	" __| \".        |\\dS\"qML\n" +
	" |    `.       | `' \\Zq\n" +
	"_)      \\.___.,|     .'\n" +
	"\\____   )MMMMMP|   .'\n" +
	"     `-'       `--\n" +
	"[ENTER YOUR NAME]: "

// Client represents a connected chat user.
type Client struct {
	name string
	conn net.Conn
	out  chan string
}

// Server holds state for the chat server.
type Server struct {
	listenAddr string
	ln         net.Listener
	mu         sync.Mutex
	clients    map[string]*Client
	history    []string
	logFile     *os.File

}
