package main

import (
	"fmt"
	"io"
	"log"
	"net"

	"github.com/tidwall/resp"
)

type Peer struct {
	conn  net.Conn
	msgCh chan Message
	delCh chan *Peer
}

func NewPeer(conn net.Conn, msgCh chan Message, delCh chan *Peer) *Peer {
	return &Peer{
		conn:  conn,
		msgCh: msgCh,
		delCh: delCh,
	}
}

func (p *Peer) Send(msg []byte) (int, error) {
	return p.conn.Write(msg)
}

func (p *Peer) readLoop() error {
	rd := resp.NewReader(p.conn)
	for {
		v, _, err := rd.ReadValue()
		if err == io.EOF {
			p.delCh <- p
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		if v.Type() == resp.Array {
			for _, value := range v.Array() {
				var cmd Command
				switch value.String() {
				case CommandSET:
					if len(v.Array()) != 3 {
						return fmt.Errorf("invalid number of variables for SET command")
					}
					cmd = SetCommand{
						key: v.Array()[1].Bytes(),
						val: v.Array()[2].Bytes(),
					}

					// fmt.Printf("got SET cmd %+v\n", cmd)
				case CommandGET:
					if len(v.Array()) != 2 {
						return fmt.Errorf("invalid number of variables for GET command")
					}
					cmd = GetCommand{
						key: v.Array()[1].Bytes(),
					}

					// fmt.Printf("got GET cmd %+v\n", cmd)
				case CommandHELLO:
					cmd = HelloCommand{
						value: v.Array()[1].String(),
					}
				default:
					// fmt.Println(v.Array())
					panic("this command is not being handled")
				}
				p.msgCh <- Message{
					cmd:  cmd,
					peer: p,
				}
			}
		}
	}
	return nil
}
