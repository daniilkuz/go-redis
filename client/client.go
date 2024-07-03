package client

import (
	"bytes"
	"context"
	"io"
	"net"

	"github.com/tidwall/resp"
)

type Client struct {
	addr string
	conn net.Conn
}

func New(addr string) (*Client, error) {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return nil, err
	}
	return &Client{
		addr: addr,
		conn: conn,
	}, nil
}

func (c *Client) Set(ctx context.Context, key string, val string) error {

	buf := &bytes.Buffer{}
	wr := resp.NewWriter(buf)
	wr.WriteArray([]resp.Value{
		resp.StringValue("SET"),
		resp.StringValue(key),
		resp.StringValue(val),
	})
	// fmt.Printf("%s", buf.String())
	// _, err = conn.Write(buf.Bytes())
	_, err := io.Copy(c.conn, buf)
	return err
}

func (c *Client) Get(ctx context.Context, key string) (string, error) {

	buf := &bytes.Buffer{}
	wr := resp.NewWriter(buf)
	wr.WriteArray([]resp.Value{
		resp.StringValue("GET"),
		resp.StringValue(key),
	})
	// fmt.Printf("%s", buf.String())
	// _, err = conn.Write(buf.Bytes())
	_, err := io.Copy(c.conn, buf)

	if err != nil {
		return "", err
	}
	b := make([]byte, 1024)
	n, err := c.conn.Read(b)

	return string(b[:n]), err
}

func (c *Client) Close() error {
	return c.conn.Close()
}
