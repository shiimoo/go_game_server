package net

import (
	"context"
	"errors"
	"fmt"
	"net"
	"time"
)

/* mgr */

type TcpMgr struct {
	ip   string // 监听id
	port int    // 监听端口

	ctx    context.Context
	cancel context.CancelFunc

	listener *net.TCPListener // 监听器
	// 链接集合, 中央管理器
	count int // 临时计数器，将被新的id管理器代替
	fds   map[int]*Fd

	cFunc func(c net.Conn) // 链接回调
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
	fd := newFd(m.count, c)
	m.fds[fd.Id()] = fd
	m.count += 1
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
	id int      // 分配的唯一id
	c  net.Conn // 链接
}

func newFd(id int, c net.Conn) *Fd {
	return &Fd{id, c}
}

func (f *Fd) Id() int {
	return f.id
}
