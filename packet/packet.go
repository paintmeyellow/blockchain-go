package packet

import "errors"

type Type string

const (
	ConnectionInit Type = "connection_init"
	ConnectionAck  Type = "connection_ack"
	Request        Type = "request"
	Next           Type = "next"
)

type Packet struct {
	ID      uint64  `json:"id"`
	Type    Type    `json:"type"`
	Payload Payload `json:"payload"`
}

type Payload map[string]interface{}

func (p Payload) StringValue(key string) (string, error) {
	v, ok := p[key]
	if !ok {
		return "", errors.New("key " + key + " is not found")
	}
	res, ok := v.(string)
	if !ok {
		return "", errors.New("error type assertion")
	}
	return res, nil
}
