package main

import (
	"flag"
	"fmt"
	"log"
	"log/slog"
	"net"
	"reflect"
)

const defaultListenAddr = ":5001"

type Config struct {
	ListenAddr string
}

type Message struct {
	// data []byte
	cmd  Command
	peer *Peer
}

type Server struct {
	Config
	peers     map[*Peer]bool
	ln        net.Listener
	addPeerCh chan *Peer
	delPeerCh chan *Peer
	quitCh    chan struct{}
	msgCh     chan Message
	kv        *KV
}

func NewServer(cfg Config) *Server {
	if len(cfg.ListenAddr) == 0 {
		cfg.ListenAddr = defaultListenAddr
	}
	return &Server{
		Config:    cfg,
		peers:     make(map[*Peer]bool),
		addPeerCh: make(chan *Peer),
		delPeerCh: make(chan *Peer),
		quitCh:    make(chan struct{}),
		msgCh:     make(chan Message),
		kv:        NewKV(),
	}
}

func (s *Server) Start() error {
	ln, err := net.Listen("tcp", s.ListenAddr)
	if err != nil {
		return err
	}
	s.ln = ln
	go s.loop()
	slog.Info("goredis server running", "listenAddr", s.ListenAddr)
	return s.acceptLoop()
}

func (s *Server) handleMessage(msg Message) error {
	// cmd, err := parseCommand(string(msg.data))
	// if err != nil {
	// 	return err
	// }
	slog.Info("got message from client", "type", reflect.TypeOf(msg.cmd))
	switch v := msg.cmd.(type) {
	case SetCommand:
		if err := s.kv.Set(v.key, v.val); err != nil {
			return err
		}
		_, err := msg.peer.Send([]byte("OK\r\n"))
		if err != nil {
			return fmt.Errorf("peer send error: %s", err)
		}
	case GetCommand:
		val, ok := s.kv.Get(v.key)
		if !ok {
			return fmt.Errorf("key not found")
		}
		_, err := msg.peer.Send(val)
		if err != nil {
			// slog.Error("peer send error", "err", err)
			return fmt.Errorf("peer send error: %s", err)
		}
		// return s.set(v.key, v.val)
		// slog.Info("someone wants to set a key into the hash table", "key", v.key, "val", v.val)
	case HelloCommand:
		fmt.Println("hello handled on server")
		spec := map[string]string{
			"server": "redis",
			"role":   "master",
		}
		_, err := msg.peer.Send(respWriteMap(spec))
		if err != nil {
			// slog.Error("peer send error", "err", err)
			return fmt.Errorf("peer send error: %s", err)
		}
		// fmt.Println("This is the hello command from the client")
	}
	return nil
}

func (s *Server) loop() {
	for {
		select {
		case msg := <-s.msgCh:
			if err := s.handleMessage(msg); err != nil {
				slog.Error("message error", "err", err)
			}
		case <-s.quitCh:
			return
		case peer := <-s.addPeerCh:
			slog.Info("peer connected", "remoteAddr", peer.conn.RemoteAddr())
			s.peers[peer] = true
			// default:
			// 	fmt.Println("foo")
		case peer := <-s.delPeerCh:
			slog.Info("peer disconnected", "remoteAddr", peer.conn.RemoteAddr())
			delete(s.peers, peer)
		}

	}
}

func (s *Server) acceptLoop() error {
	for {
		conn, err := s.ln.Accept()
		if err != nil {
			slog.Error("accept error", "err", err)
			continue
		}
		go s.handleConn(conn)
	}
}

func (s *Server) handleConn(conn net.Conn) {
	peer := NewPeer(conn, s.msgCh, s.delPeerCh)
	s.addPeerCh <- peer
	slog.Info("peer connected", "remoteAddr", conn.RemoteAddr())
	if err := peer.readLoop(); err != nil {
		slog.Error("peer read error", "err", err, "remoteAddr", conn.RemoteAddr())
	}
}

func main() {
	listenAddr := flag.String("listenAddr", defaultListenAddr, "listen address of the goredis server")
	flag.Parse()
	server := NewServer(Config{
		ListenAddr: *listenAddr,
	})
	log.Fatal(server.Start())

	// fmt.Println(server.kv.data)

}
