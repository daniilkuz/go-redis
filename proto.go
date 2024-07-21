package main

import (
	"bytes"
	"fmt"
)

const (
	CommandSET   = "SET"
	CommandGET   = "GET"
	CommandHELLO = "hello"
)

type Command interface {
}

type SetCommand struct {
	key, val []byte
}

type HelloCommand struct {
	value string
}

type GetCommand struct {
	key []byte
}

// func parseCommand(raw string) (Command, error) {
// 	rd := resp.NewReader(bytes.NewBufferString(raw))
// 	for {
// 		v, _, err := rd.ReadValue()
// 		if err == io.EOF {
// 			break
// 		}
// 		if err != nil {
// 			log.Fatal(err)
// 		}
// 		if v.Type() == resp.Array {
// 			for _, value := range v.Array() {
// 				switch value.String() {
// 				case CommandSET:
// 					if len(v.Array()) != 3 {
// 						return nil, fmt.Errorf("invalid number of variables for SET command")
// 					}
// 					cmd := SetCommand{
// 						key: v.Array()[1].Bytes(),
// 						val: v.Array()[2].Bytes(),
// 					}
// 					return cmd, nil
// 				case CommandGET:
// 					if len(v.Array()) != 2 {
// 						return nil, fmt.Errorf("invalid number of variables for GET command")
// 					}
// 					cmd := GetCommand{
// 						key: v.Array()[1].Bytes(),
// 					}
// 					return cmd, nil
// 				case CommandHELLO:
// 					cmd := HelloCommand{
// 						value: v.Array()[1].String(),
// 					}
// 					return cmd, nil

// 				}
// 			}
// 		}
// 	}
// 	return nil, fmt.Errorf("invalid or unknown command received: %s", raw)
// }

func respWriteMap(m map[string]string) string {
	buf := bytes.Buffer{}
	buf.WriteString("%" + fmt.Sprintf("%d\r\n", len(m)))
	for k, v := range m {
		buf.WriteString(fmt.Sprintf("+%s\r\n", k))
		buf.WriteString(fmt.Sprintf(":%s\r\n", v))
	}
	return buf.String()
}
