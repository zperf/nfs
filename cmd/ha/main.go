package main

import (
	"crypto/rand"
	"encoding/binary"
	"net"
	"sync"

	"go.uber.org/zap"

	"github.com/zperf/nfs/v3/logging"
)

func coin() bool {
	var n uint64
	err := binary.Read(rand.Reader, binary.BigEndian, &n)
	if err != nil {
		zap.S().Fatal("coin failed", err)
	}

	return n%2 == 0
}

type ConnHolder struct {
	c    net.Conn
	lock sync.RWMutex
}

func (h *ConnHolder) Conn() net.Conn {
	var r net.Conn
	h.lock.RLock()
	r = h.c
	h.lock.RUnlock()
	return r
}

func (h *ConnHolder) Roll(c net.Conn) {
	if c == h.c {
		return
	}
	zap.S().Info("Connection roll")
	h.lock.Lock()
	h.c = c
	h.lock.Unlock()
}

func onAccept(conn net.Conn, dst1 net.Conn, dst2 net.Conn) {
	log := zap.S()
	holder := ConnHolder{c: dst1}

	log.Info("Accepted..")

	var wg sync.WaitGroup
	wg.Add(2)
	exit := make(chan struct{})

	go func(wg *sync.WaitGroup, h *ConnHolder) {
		defer wg.Done()
		// read conn, write to dst
		buf := make([]byte, 8192)
	p:
		for {
			select {
			case <-exit:
				return
			default:
				n, err := conn.Read(buf)
				if err != nil || n == 0 {
					break p
				}

				w, err := h.Conn().Write(buf[:n])
				if err != nil || w != n {
					break p
				}
			}
		}
		exit <- struct{}{}
		log.Info("routine1 exited")
	}(&wg, &holder)

	go func(wg *sync.WaitGroup, h *ConnHolder) {
		defer wg.Done()
		// read dst, write to conn
		buf := make([]byte, 8192)
	p:
		for {
			select {
			case <-exit:
				return
			default:
				n, err := h.Conn().Read(buf)
				if err != nil || n == 0 {
					break p
				}

				w, err := conn.Write(buf[:n])
				if err != nil || w != n {
					break p
				}
			}
		}
		exit <- struct{}{}
		log.Info("routine2 exited")
	}(&wg, &holder)

	wg.Wait()
	close(exit)
	err := conn.Close()
	if err != nil {
		log.Warn("close client conn failed", err)
	}
	log.Info("Connection closed")
}

type setNoDelay interface {
	SetNoDelay(bool) error
}

func main() {
	log := logging.New().Sugar()

	dst1, err := net.Dial("tcp", "192.168.90.221:80")
	if err != nil {
		log.Fatalf("dial failed, %v", err)
	}

	dst2, err := net.Dial("tcp", "192.168.126.16:3456")
	if err != nil {
		log.Fatalf("dial failed, %v", err)
	}

	listener, err := net.Listen("tcp", ":3456")
	if err != nil {
		log.Fatalf("listen failed, %v", err)
	}

	_ = dst1.(setNoDelay).SetNoDelay(true)
	_ = dst2.(setNoDelay).SetNoDelay(true)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatalf("accept failed, %v", err)
		}

		go onAccept(conn, dst1, dst2)
	}
}
