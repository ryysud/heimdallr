package connector

import (
	"context"
	"crypto/tls"
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"os"
	"sync"
	"time"

	"github.com/f110/lagrangian-proxy/pkg/database"
	"github.com/f110/lagrangian-proxy/pkg/logger"
	"go.uber.org/zap"
	"golang.org/x/xerrors"
)

type Relay struct {
	Addr     *net.TCPAddr
	name     string
	address  string
	server   *Server
	conn     *tls.Conn
	listener net.Listener

	mu       sync.Mutex
	accepted map[string]net.Conn
}

func NewRelay(ca database.CertificateAuthority, name string, server *Server, conn *tls.Conn) (*Relay, error) {
	hostname, err := os.Hostname()
	if err != nil {
		return nil, xerrors.Errorf(": %v", err)
	}

	l, err := net.Listen("tcp", ":0")
	if err != nil {
		return nil, xerrors.Errorf(": %v", err)
	}
	addr := l.Addr().(*net.TCPAddr)

	cert, privateKey, err := ca.NewServerCertificate(hostname)
	if err != nil {
		return nil, xerrors.Errorf(": %v", err)
	}
	listener := tls.NewListener(l, &tls.Config{
		Certificates: []tls.Certificate{
			{
				Certificate: [][]byte{cert.Raw},
				PrivateKey:  privateKey,
			},
		},
	})

	address := fmt.Sprintf("%s:%d", hostname, addr.Port)
	err = server.Locator.Set(
		context.Background(),
		&database.Relay{
			Name:      name,
			Addr:      address,
			UpdatedAt: time.Now(),
		},
	)
	if err != nil {
		l.Close()
		return nil, err
	}

	return &Relay{
		Addr:     addr,
		name:     name,
		address:  address,
		server:   server,
		conn:     conn,
		listener: listener,
		accepted: make(map[string]net.Conn),
	}, nil
}

func (r *Relay) Serve() error {
	logger.Log.Debug("Start relay", zap.String("for", r.name), zap.String("addr", r.Addr.String()))
	for {
		conn, err := r.listener.Accept()
		if err != nil {
			break
		}
		go r.acceptConn(conn)
	}
	logger.Log.Debug("Close relay listener", zap.String("name", r.name))
	return nil
}

func (r *Relay) Close() {
	if err := r.server.Locator.Delete(context.Background(), r.name, r.address); err != nil {
		logger.Log.Info("Failure delete relay record", zap.Error(err))
	}
	r.listener.Close()
	for _, c := range r.accepted {
		c.Close()
	}
}

func (r *Relay) acceptConn(childConn net.Conn) {
	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()
	defer childConn.Close()

	r.mu.Lock()
	r.accepted[childConn.RemoteAddr().String()] = childConn
	r.mu.Unlock()
	defer func() {
		r.mu.Lock()
		defer r.mu.Unlock()
		delete(r.accepted, childConn.RemoteAddr().String())
	}()

	header := make([]byte, 5)
	buf := make([]byte, 4*1024)
	for {
		_, err := io.ReadFull(childConn, header)
		if err != nil {
			return
		}

		bodySize := binary.BigEndian.Uint32(header[1:5])
		if cap(buf) < int(bodySize) {
			buf = make([]byte, bodySize)
		}
		var n int
		if bodySize > 0 {
			n, err = io.ReadAtLeast(childConn, buf, int(bodySize))
			if err != nil {
				return
			}
		}

		switch header[0] {
		case OpcodeDial:
			dialId := binary.BigEndian.Uint32(buf[:4])
			streamId, err := r.server.DialUpstreamForRelay(ctx, r.name, childConn, dialId)
			if err != nil {
				return
			}
			b := make([]byte, 8)
			binary.BigEndian.PutUint32(b[:4], dialId)
			binary.BigEndian.PutUint32(b[4:8], streamId)
			logger.Log.Debug("dial success", zap.Int32("Dial id", int32(dialId)), zap.Uint32("stream id", streamId))
			f := NewFrame()
			f.Write(b)
			f.EncodeTo(OpcodeDialSuccess, childConn)
		case OpcodeHeartbeat:
			childConn.SetWriteDeadline(time.Now().Add(heartbeatDuration * 2))
			childConn.Write(header)
			logger.Log.Debug("Write heartbeat of relay")
		default:
			r.conn.Write(append(header, buf[:n]...))
		}
	}
}