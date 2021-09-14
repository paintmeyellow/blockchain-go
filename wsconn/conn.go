package wsconn

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"sync"
	"time"
)

const SubsChanLen = 8

const _EMPTY_ = ""

var (
	ErrInvalidConnection = errors.New("pubsub: invalid connection")
	ErrConnectionClosed  = errors.New("pubsub: connection closed")
	ErrTimeout           = errors.New("pubsub: timeout")
)

type Conn struct {
	mu       sync.RWMutex
	subsMu   sync.RWMutex
	subs     map[string][]chan *Msg
	closed   bool
	nextID   uint64
	conn     *websocket.Conn
	respMap  map[string]chan *Msg
	respRand *rand.Rand
}

type Msg struct {
	ID    string `json:"id"`
	Reply string `json:"reply"`
	Subj  string `json:"subject"`
	Data  []byte `json:"data"`
}

type MsgHandler func(msg *Msg)

func Upgrade(w http.ResponseWriter, r *http.Request, h http.Header) (*Conn, error) {
	u := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		WriteBufferPool: &sync.Pool{},
	}
	c, err := u.Upgrade(w, r, h)
	if err != nil {
		return nil, err
	}
	return newConn(c)
}

func Connect(url string) (*Conn, error) {
	c, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		return nil, err
	}
	return newConn(c)
}

func newConn(c *websocket.Conn) (*Conn, error) {
	if c == nil {
		return nil, ErrInvalidConnection
	}
	ps := Conn{
		subs:     make(map[string][]chan *Msg, SubsChanLen),
		respMap:  make(map[string]chan *Msg),
		respRand: rand.New(rand.NewSource(time.Now().UnixNano())),
		conn:     c,
	}
	go ps.readLoop()
	return &ps, nil
}

func (c *Conn) readLoop() {
	for {
		_, r, err := c.conn.NextReader()
		if err == nil {
			var msg Msg
			if err := json.NewDecoder(r).Decode(&msg); err != nil {
				log.Println(err)
				continue
			}
			if msg.Reply != _EMPTY_ {
				c.mu.Lock()
				ch, ok := c.respMap[msg.Reply]
				c.mu.Unlock()
				if ok {
					go func(ch chan *Msg) { ch <- &msg }(ch)
				}
				continue
			}
			c.subsMu.Lock()
			subs := c.subs[msg.Subj]
			c.subsMu.Unlock()
			for _, ch := range subs {
				go func(ch chan *Msg) { ch <- &msg }(ch)
			}
		}
		if err != nil {
			log.Println(err)
			break
		}
	}
}

func (c *Conn) publish(subj, reply, id string, data []byte) error {
	c.mu.RLock()
	defer c.mu.RUnlock()
	if c.closed {
		return ErrConnectionClosed
	}
	if id == _EMPTY_ {
		id = c.randomToken()
	}
	msg := Msg{
		ID:    id,
		Reply: reply,
		Subj:  subj,
		Data:  data,
	}
	fmt.Println("published", "id:", msg.ID, "reply" ,msg.Reply)
	return c.conn.WriteJSON(&msg)
}

func (c *Conn) subscribe(subj string, ch chan *Msg) {
	c.mu.Lock()
	if c.closed {
		return
	}
	c.mu.Unlock()
	c.subsMu.Lock()
	c.subs[subj] = append(c.subs[subj], ch)
	c.subsMu.Unlock()
}

const (
	tokenLen      = 8
	letterBytes   = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

func (c *Conn) randomToken() string {
	var sb strings.Builder
	sb.Grow(tokenLen)
	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := tokenLen-1, c.respRand.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = c.respRand.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			sb.WriteByte(letterBytes[idx])
			i--
		}
		cache >>= letterIdxBits
		remain--
	}
	return sb.String()
}

func (c *Conn) Subscribe(subj string, h MsgHandler) {
	ch := make(chan *Msg)
	c.subscribe(subj, ch)
	go func() {
		for {
			m, ok := <-ch
			if !ok {
				return
			}
			h(m)
		}
	}()
}

func (c *Conn) RespondOnMsg(m *Msg, data []byte) error {
	return c.publish(m.Subj, m.ID, _EMPTY_, data)
}

func (c *Conn) Publish(subj string, data []byte) error {
	return c.publish(subj, _EMPTY_, _EMPTY_, data)
}

func (c *Conn) Request(subj string, data []byte, timeout time.Duration) (*Msg, error) {
	id := c.randomToken()

	c.mu.Lock()
	mch := make(chan *Msg)
	c.respMap[id] = mch
	c.mu.Unlock()

	if err := c.publish(subj, _EMPTY_, id, data); err != nil {
		return nil, err
	}

	var msg *Msg
	var err error
	t := time.NewTimer(timeout)

	select {
	case msg = <-mch:
	case <-t.C:
		err = ErrTimeout
	}

	c.mu.Lock()
	delete(c.respMap, id)
	c.mu.Unlock()

	return msg, err
}

func (c *Conn) Close() {
	c.mu.Lock()
	defer c.mu.Unlock()
	closed := c.closed
	if closed {
		c.subsMu.Lock()
		for _, subs := range c.subs {
			for _, ch := range subs {
				close(ch)
			}
		}
		c.subsMu.Unlock()
		return
	}
	c.closed = true
}
