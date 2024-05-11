package net

import (
	"context"
	"errors"
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/shiimoo/go_game_server/id"
)

/* mgr */

type TcpMgr struct {
	ip   string // 监听id
	port int    // 监听端口

	ctx    context.Context
	cancel context.CancelFunc

	listener *net.TCPListener // 监听器
	// 链接集合, 中央管理器
	fds    map[int]*Fd
	fdLock sync.RWMutex // fds 锁控制

	cFunc func(c net.Conn)
}

func NewTcpMgr(parent context.Context, ip string, port int) (*TcpMgr, error) {
	mgr := new(TcpMgr)
	// TODO: ip和port 合法性检查
	mgr.ip = ip
	mgr.port = port
	if err := mgr.initListen(); err != nil {
		return nil, err
	}
	mgr.fds = make(map[int]*Fd, 0)
	// 创建子context
	mgr.ctx, mgr.cancel = context.WithCancel(parent)

	return mgr, nil
}

func (m *TcpMgr) SetConnFunc(f func(c net.Conn)) {
	m.cFunc = f
}

func (m *TcpMgr) initListen() error {
	addr, err := net.ResolveTCPAddr("tcp", ":8080")
	if err != nil {
		return err
	}

	listerner, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return err
	}
	m.listener = listerner
	return nil
}

func (m *TcpMgr) AddFd(c net.Conn) {
	m.fdLock.Lock()
	defer m.fdLock.Unlock()
	fd := newFd(m.ctx, id.Gen(), c, 100) // 暂定100毫秒读取超时
	m.fds[fd.Id()] = fd
	go fd.start()
}

func (m *TcpMgr) _start() {
	for {
		select {
		case <-m.ctx.Done():
			_ = m.listener.Close()
			return
		default:
			err := m.listener.SetDeadline(time.Now().Add(time.Millisecond * 10))
			if err != nil {
				// fmt.Println("listen set timeout err:", err)
				continue
			}
			conn, err := m.listener.AcceptTCP()
			if err != nil {
				var netErr net.Error
				switch {
				case errors.As(err, &netErr):
					continue
				default:
					fmt.Println("listen AcceptTCP err:", err)
					continue
				}
			}
			// 添加链接
			if m.cFunc != nil {
				m.cFunc(conn)
			}
			fmt.Println("链接接入", conn.RemoteAddr(), m.cFunc == nil)
			m.AddFd(conn)
		}
	}
}

func (m *TcpMgr) Start() {
	go m._start()
}

func (m *TcpMgr) Close() {
	m.cancel()
}

/* connFd */

type Fd struct {
	id          int      // 分配的唯一id
	c           net.Conn // 链接
	readTimeout int      // 读取(毫秒) 超时

	ctx    context.Context
	cancel context.CancelFunc
}

func newFd(parent context.Context, id int, c net.Conn, timeout int) *Fd {
	if timeout < 0 {
		timeout = 0
	}
	f := new(Fd)
	f.id = id
	f.c = c
	f.readTimeout = timeout
	f.ctx, f.cancel = context.WithCancel(parent)
	return f
}

func (f *Fd) Id() int {
	return f.id
}

func (f *Fd) start() {
	for {
		select {
		case <-f.ctx.Done():
			_ = f.c.Close()
			// todo 通知链接关闭
			fmt.Println("conn close:", f.Id())
			return
		default:
			err := f.c.SetDeadline(time.Now().Add(time.Duration(f.readTimeout)))
			if err != nil {
				fmt.Println(" set read Time out err:", err)
			}
			bs := make([]byte, 128)
			n, err := f.c.Read(bs)
			if err != nil {
				if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
					continue
				} else {
					f.Close()
				}
			}

			bs = bs[:n]
		}
	}
}

func (f *Fd) Close() {
	f.cancel()
}
