package wspool

import (
	"context"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type Client interface {
	Get(ctx context.Context) (*websocket.Conn, error)
	Putback(c *websocket.Conn)
}

func New(
	addr string,
	usingUUID bool,
	poolSize int,
	header http.Header,
	timeout time.Duration,
) Client {
	wp := &wsPool{
		addr:    addr,
		timeout: timeout,
		header:  header,
		tokens:  make(chan struct{}, poolSize),
		pool:    make(chan *wsConn, poolSize),
	}
	for i := 0; i < poolSize; i++ {
		wp.tokens <- struct{}{}
	}
	return wp
}

type wsConn struct {
	conn      *websocket.Conn
	timestamp time.Time
}

type wsPool struct {
	addr    string
	timeout time.Duration
	header  http.Header

	tokens chan struct{}

	pool chan *wsConn
}

func (wp *wsPool) Get(ctx context.Context) (*websocket.Conn, error) {
	<-wp.tokens

	for {
		select {
		case ws := <-wp.pool:
			if time.Since(ws.timestamp) < wp.timeout {
				return ws.conn, nil
			} else {
				ws.conn.Close()
			}
		default:
			uid := uuid.New().String()
			dl := &websocket.Dialer{
				EnableCompression: true,
			}
			conn, _, err := dl.DialContext(ctx, wp.addr+uid, wp.header)
			if err != nil {
				return nil, err
			}
			return conn, nil
		}
	}
}

func (wp *wsPool) Putback(c *websocket.Conn) {
	if c != nil {
		select {
		case wp.pool <- &wsConn{conn: c, timestamp: time.Now()}:
		default:
			c.Close()
		}
	}
	wp.tokens <- struct{}{}
}
