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

		var cmd Command
		if v.Type() == resp.Array {
			rawCmd := v.Array()[0]
			switch rawCmd.String() {
			case CommandCLIENT:
				cmd = ClientCommand{
					value: v.Array()[1].String(),
				}
			case CommandGET:
				cmd = GetCommand{
					key: v.Array()[1].Bytes(),
				}
			case CommandSET:
				cmd = SetCommand{
					key: v.Array()[1].Bytes(),
					val: v.Array()[2].Bytes(),
				}

			case CommandHELLO:
				cmd = HelloCommand{
					value: v.Array()[1].String(),
				}
			default:
				fmt.Println("got this unhandled command", rawCmd)
			}
			p.msgCh <- Message{
				cmd:  cmd,
				peer: p,
			}
			// fmt.Println("this is cmd", v.Array()[0])
		}
	}
	return nil
}
