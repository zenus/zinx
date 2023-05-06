package znet

import (
	"crypto/rand"
	"crypto/tls"
	"errors"
	"fmt"
	"github.com/zenus/zinx/logo"
	"github.com/zenus/zinx/zconf"
	//"github.com/zenus/zinx/zdecoder"
	"github.com/zenus/zinx/zlog"
	"github.com/zenus/zinx/zmetrics"
	"net"
	"os"
	"os/signal"
	"sync/atomic"
	"time"

	"github.com/zenus/zinx/ziface"
	//"github.com/zenus/zinx/zpack"
)

// Server 接口实现，定义一个Server服务类
type Server struct {
	//服务器的名称
	Name string
	//tcp4 or other
	IPVersion string
	//服务绑定的IP地址
	IP string
	//服务绑定的端口
	Port int
	// 服务绑定的websocket 端口
	WsPort int
	//当前Server的消息管理模块，用来绑定MsgID和对应的处理方法
	msgHandler ziface.IMsgHandle
	//路由模式
	RouterSlicesMode bool
	//当前Server的链接管理器
	ConnMgr ziface.IConnManager
	//该Server的连接创建时Hook函数
	onConnStart func(conn ziface.IConnection)
	//该Server的连接断开时的Hook函数
	onConnStop func(conn ziface.IConnection)
	//数据报文封包方式
	//	packet ziface.IDataPack
	//异步捕获链接关闭状态
	exitChan chan struct{}
	//断粘包解码器
	//	decoder ziface.IDecoder
	//心跳检测器
	//	hc ziface.IHeartbeatChecker

	// websocket
	//	upgrader *websocket.Upgrader
	// websocket 连接认证
	//	websocketAuth func(r *http.Request) error
	// connection id
	cID uint64
}

// NewServer 创建一个服务器句柄
func NewServer() ziface.IServer {
	logo.PrintLogo()

	s := &Server{
		Name:             zconf.GlobalObject.Name,
		IPVersion:        "tcp",
		IP:               zconf.GlobalObject.Host,
		Port:             zconf.GlobalObject.TCPPort,
		WsPort:           zconf.GlobalObject.WsPort,
		msgHandler:       newMsgHandle(),
		RouterSlicesMode: zconf.GlobalObject.RouterSlicesMode,
		ConnMgr:          newConnManager(),
		exitChan:         nil,
		//默认使用zinx的TLV封包方式
		//	packet:  zpack.Factory().NewPack(ziface.JsonDataPack),
		//decoder: zdecoder.NewJsonDecoder(), //默认使用TLV的解码方式
		/*
			upgrader: &websocket.Upgrader{
				ReadBufferSize: int(zconf.GlobalObject.IOReadBuffSize),
				CheckOrigin: func(r *http.Request) bool {
					return true
				},
			},
		*/
	}

	//提示当前配置信息
	zconf.GlobalObject.Show()

	return s
}

// NewServer 创建一个服务器句柄
func NewUserConfServer(config *zconf.Config) ziface.IServer {

	//刷新用户配置到全局配置变量
	zconf.UserConfToGlobal(config)

	//提示当前配置信息
	zconf.GlobalObject.Show()

	//打印logo
	logo.PrintLogo()

	s := &Server{
		Name:             config.Name,
		IPVersion:        "tcp",
		IP:               config.Host,
		Port:             config.TCPPort,
		WsPort:           config.WsPort,
		msgHandler:       newMsgHandle(),
		RouterSlicesMode: config.RouterSlicesMode,
		ConnMgr:          newConnManager(),
		exitChan:         nil,
		//	packet:           zpack.Factory().NewPack(ziface.ZinxDataPack),
		//	decoder:          zdecoder.NewTLVDecoder(), //默认使用TLV的解码方式
		//		upgrader: &websocket.Upgrader{
		//			ReadBufferSize: int(zconf.GlobalObject.IOReadBuffSize),
		//			CheckOrigin: func(r *http.Request) bool {
		//				return true
		//			},
		//		},
	}
	//	//更替打包方式
	//	for _, opt := range opts {
	//		opt(s)
	//	}
	//
	return s
}

// NewDefaultRouterSlicesServer 创建一个默认自带一个Recover处理器的服务器句柄
//func NewDefaultRouterSlicesServer(opts ...Option) ziface.IServer {
//	logo.PrintLogo()
//	zconf.GlobalObject.RouterSlicesMode = true
//	s := &Server{
//		Name:             zconf.GlobalObject.Name,
//		IPVersion:        "tcp",
//		IP:               zconf.GlobalObject.Host,
//		Port:             zconf.GlobalObject.TCPPort,
//		WsPort:           zconf.GlobalObject.WsPort,
//		msgHandler:       newMsgHandle(),
//		RouterSlicesMode: zconf.GlobalObject.RouterSlicesMode,
//		ConnMgr:          newConnManager(),
//		exitChan:         nil,
//		//默认使用zinx的TLV封包方式
//		//	packet:  zpack.Factory().NewPack(ziface.ZinxDataPack),
//		//	decoder: zdecoder.NewTLVDecoder(), //默认使用TLV的解码方式
//		upgrader: &websocket.Upgrader{
//			ReadBufferSize: int(zconf.GlobalObject.IOReadBuffSize),
//			CheckOrigin: func(r *http.Request) bool {
//				return true
//			},
//		},
//	}
//
//	for _, opt := range opts {
//		opt(s)
//	}
//	s.Use(RouterRecovery)
//	//提示当前配置信息
//	zconf.GlobalObject.Show()
//
//	return s
//}

// NewUserRouterSlicesServer 创建一个用户配置的自带一个Recover处理器的服务器句柄，如果用户不希望Use这个方法，那么应该使用NewUserConfServer
//func NewUserConfDefaultRouterSlicesServer(config *zconf.Config, opts ...Option) ziface.IServer {
//
//	if !config.RouterSlicesMode {
//		panic("RouterSlicesMode is false")
//	}
//
//	//刷新用户配置到全局配置变量
//	zconf.UserConfToGlobal(config)
//
//	//提示当前配置信息
//	zconf.GlobalObject.Show()
//
//	//打印logo
//	logo.PrintLogo()
//
//	s := &Server{
//		Name:             config.Name,
//		IPVersion:        "tcp4",
//		IP:               config.Host,
//		Port:             config.TCPPort,
//		WsPort:           config.WsPort,
//		msgHandler:       newMsgHandle(),
//		RouterSlicesMode: config.RouterSlicesMode,
//		ConnMgr:          newConnManager(),
//		exitChan:         nil,
//		//	packet:           zpack.Factory().NewPack(ziface.ZinxDataPack),
//		//	decoder:          zdecoder.NewTLVDecoder(), //默认使用TLV的解码方式
//		upgrader: &websocket.Upgrader{
//			ReadBufferSize: int(zconf.GlobalObject.IOReadBuffSize),
//			CheckOrigin: func(r *http.Request) bool {
//				return true
//			},
//		},
//	}
//	//更替打包方式
//	for _, opt := range opts {
//		opt(s)
//	}
//	s.Use(RouterRecovery)
//	return s
//}

// ============== 实现 ziface.IServer 里的全部接口方法 ========
func (s *Server) StartConn(conn ziface.IConnection) {
	// HeartBeat 心跳检测
	/*	if s.hc != nil {
		//从Server端克隆一个心跳检测器
		heartBeatChecker := s.hc.Clone()

		//绑定当前链接
		heartBeatChecker.BindConn(conn)
	}*/

	//3.4 启动当前链接的处理业务
	conn.Start()
}

func (s *Server) ListenTcpConn() {
	//1 获取一个TCP的Addr
	addr, err := net.ResolveTCPAddr(s.IPVersion, fmt.Sprintf("%s:%d", s.IP, s.Port))
	if err != nil {
		zlog.Ins().ErrorF("[START] resolve tcp addr err: %v\n", err)
		return
	}
	// 2 监听服务器地址
	var listener net.Listener
	if zconf.GlobalObject.CertFile != "" && zconf.GlobalObject.PrivateKeyFile != "" {
		// 读取证书和密钥
		crt, err := tls.LoadX509KeyPair(zconf.GlobalObject.CertFile, zconf.GlobalObject.PrivateKeyFile)
		if err != nil {
			panic(err)
		}

		// TLS连接
		tlsConfig := &tls.Config{}
		tlsConfig.Certificates = []tls.Certificate{crt}
		tlsConfig.Time = time.Now
		tlsConfig.Rand = rand.Reader
		listener, err = tls.Listen(s.IPVersion, fmt.Sprintf("%s:%d", s.IP, s.Port), tlsConfig)
		if err != nil {
			panic(err)
		}
	} else {
		listener, err = net.ListenTCP(s.IPVersion, addr)
		if err != nil {
			panic(err)
		}
	}

	//3 启动server网络连接业务
	go func() {
		for {
			//3.1 设置服务器最大连接控制,如果超过最大连接，则等待
			if s.ConnMgr.Len() >= zconf.GlobalObject.MaxConn {
				zlog.Ins().InfoF("Exceeded the maxConnNum:%d, Wait:%d", zconf.GlobalObject.MaxConn, AcceptDelay.duration)
				AcceptDelay.Delay()
				continue
			}
			//3.2 阻塞等待客户端建立连接请求
			conn, err := listener.Accept()
			if err != nil {
				//Go 1.16+
				if errors.Is(err, net.ErrClosed) {
					zlog.Ins().ErrorF("Listener closed")
					return
				}
				zlog.Ins().ErrorF("Accept err: %v", err)
				AcceptDelay.Delay()
				continue
			}

			AcceptDelay.Reset()
			//3.4 处理该新连接请求的 业务 方法， 此时应该有 handler 和 conn是绑定的

			newCid := atomic.AddUint64(&s.cID, 1)
			dealConn := newServerConn(s, conn, newCid)

			go s.StartConn(dealConn)

		}
	}()
	select {
	case <-s.exitChan:
		err := listener.Close()
		if err != nil {
			zlog.Ins().ErrorF("listener close err: %v", err)
		}
	}
}

//func (s *Server) ListenWebsocketConn() {
//
//	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
//		//1. 设置服务器最大连接控制,如果超过最大连接，则等待
//		if s.ConnMgr.Len() >= zconf.GlobalObject.MaxConn {
//			zlog.Ins().InfoF("Exceeded the maxConnNum:%d, Wait:%d", zconf.GlobalObject.MaxConn, AcceptDelay.duration)
//			AcceptDelay.Delay()
//			return
//		}
//		// 2. 如果需要 websocket 认证请设置认证信息
//		if s.websocketAuth != nil {
//			err := s.websocketAuth(r)
//			if err != nil {
//				zlog.Ins().ErrorF(" websocket auth err:%v", err)
//				w.WriteHeader(401)
//				AcceptDelay.Delay()
//				return
//			}
//
//		}
//		// 判断 header 里面是有子协议
//		if len(r.Header.Get("Sec-Websocket-Protocol")) > 0 {
//			s.upgrader.Subprotocols = websocket.Subprotocols(r)
//		}
//		// 4. 升级成 websocket 连接
//		conn, err := s.upgrader.Upgrade(w, r, nil)
//		if err != nil {
//			zlog.Ins().ErrorF("new websocket err:%v", err)
//			w.WriteHeader(500)
//			AcceptDelay.Delay()
//			return
//		}
//		// 5. 处理该新连接请求的 业务 方法， 此时应该有 handler 和 conn是绑定的
//		newCid := atomic.AddUint64(&s.cID, 1)
//		wsConn := newWebsocketConn(s, conn, newCid)
//		go s.StartConn(wsConn)
//
//	})
//
//	err := http.ListenAndServe(fmt.Sprintf("%s:%d", s.IP, s.WsPort), nil)
//	if err != nil {
//		panic(err)
//	}
//}

// Start 开启网络服务
func (s *Server) Start() {
	zlog.Ins().InfoF("[START] Server name: %s,listener at IP: %s, Port %d is starting", s.Name, s.IP, s.Port)
	s.exitChan = make(chan struct{})

	// 将解码器添加到拦截器
	/*	if s.decoder != nil {
		s.msgHandler.AddInterceptor(s.decoder)
	}*/
	// 启动worker工作池机制
	s.msgHandler.StartWorkerPool()

	//开启一个go去做服务端Listener业务
	switch zconf.GlobalObject.Mode {
	case zconf.ServerModeTcp:
		go s.ListenTcpConn()
	case zconf.ServerModeWebsocket:
		//go s.ListenWebsocketConn()
	default:
		go s.ListenTcpConn()
		//go s.ListenWebsocketConn()
	}

	// Prometheus Metrics 指标统计指标初始化
	zmetrics.InitZinxMetrics()

	// 启动Metrics Prometheus服务
	if zconf.GlobalObject.PrometheusMetricsEnable == true && zconf.GlobalObject.PrometheusServer == true {
		if zmetrics.RunMetricsService(zconf.GlobalObject) != nil {
			zlog.Ins().ErrorF("RunMetricsService err")
		}
	}
}

// Stop 停止服务
func (s *Server) Stop() {
	zlog.Ins().InfoF("[STOP] Zinx server , name %s", s.Name)

	//将其他需要清理的连接信息或者其他信息 也要一并停止或者清理
	s.ConnMgr.ClearConn()
	s.exitChan <- struct{}{}
	close(s.exitChan)
}

// Serve 运行服务
func (s *Server) Serve() {
	s.Start()
	//阻塞,否则主Go退出， listenner的go将会退出
	c := make(chan os.Signal, 1)
	//监听指定信号 ctrl+c kill信号
	signal.Notify(c, os.Interrupt, os.Kill)
	sig := <-c
	zlog.Ins().InfoF("[SERVE] Zinx server , name %s, Serve Interrupt, signal = %v", s.Name, sig)
}

// AddRouter 路由功能：给当前服务注册一个路由业务方法，供客户端链接处理使用
func (s *Server) AddRouter(cmd string, router ziface.IRouter) {
	if s.RouterSlicesMode {
		panic("Server RouterSlicesMode is true ")
	}
	s.msgHandler.AddRouter(cmd, router)
}
func (s *Server) AddRouterSlices(cmd string, router ...ziface.RouterHandler) ziface.IRouterSlices {
	if !s.RouterSlicesMode {
		panic("Server RouterSlicesMode is false ")
	}
	return s.msgHandler.AddRouterSlices(cmd, router...)
}

//func (s *Server) Group(start, end uint32, Handlers ...ziface.RouterHandler) ziface.IGroupRouterSlices {
//	if !s.RouterSlicesMode {
//		panic("Server RouterSlicesMode is false")
//	}
//	return s.msgHandler.Group(start, end, Handlers...)
//}

func (s *Server) Use(Handlers ...ziface.RouterHandler) ziface.IRouterSlices {
	if !s.RouterSlicesMode {
		panic("Server RouterSlicesMode is false")
	}
	return s.msgHandler.Use(Handlers...)
}

// GetConnMgr 得到链接管理
func (s *Server) GetConnMgr() ziface.IConnManager {
	return s.ConnMgr
}

// SetOnConnStart 设置该Server的连接创建时Hook函数
func (s *Server) SetOnConnStart(hookFunc func(ziface.IConnection)) {
	s.onConnStart = hookFunc
}

// SetOnConnStop 设置该Server的连接断开时的Hook函数
func (s *Server) SetOnConnStop(hookFunc func(ziface.IConnection)) {
	s.onConnStop = hookFunc
}

// GetOnConnStart 得到该Server的连接创建时Hook函数
func (s *Server) GetOnConnStart() func(ziface.IConnection) {
	return s.onConnStart
}

// 得到该Server的连接断开时的Hook函数
func (s *Server) GetOnConnStop() func(ziface.IConnection) {
	return s.onConnStop
}

/*func (s *Server) GetPacket() ziface.IDataPack {
	return s.packet
}*/

/*func (s *Server) SetPacket(packet ziface.IDataPack) {
	s.packet = packet
}*/

func (s *Server) GetMsgHandler() ziface.IMsgHandle {
	return s.msgHandler
}

// StartHeartBeat 启动心跳检测
// interval 每次发送心跳的时间间隔
/*func (s *Server) StartHeartBeat(interval time.Duration) {
	checker := NewHeartbeatChecker(interval)

	//添加心跳检测的路由
	s.AddRouter(checker.MsgID(), checker.Router())

	//server绑定心跳检测器
	s.hc = checker
}*/

// StartHeartBeatWithFunc 启动心跳检测
// option 心跳检测的配置
/*func (s *Server) StartHeartBeatWithOption(interval time.Duration, option *ziface.HeartBeatOption) {
	checker := NewHeartbeatChecker(interval)

	if option != nil {
		checker.SetHeartbeatMsgFunc(option.MakeMsg)
		checker.SetOnRemoteNotAlive(option.OnRemoteNotAlive)
		checker.BindRouter(option.HeadBeatMsgID, option.Router)
	}

	//添加心跳检测的路由
	s.AddRouter(checker.MsgID(), checker.Router())

	//server绑定心跳检测器
	s.hc = checker
}*/

/*func (s *Server) GetHeartBeat() ziface.IHeartbeatChecker {
	return s.hc
}*/

/*func (s *Server) SetDecoder(decoder ziface.IDecoder) {
	s.decoder = decoder
}*/

/*func (s *Server) GetLengthField() *ziface.LengthField {
	if s.decoder != nil {
		return s.decoder.GetLengthField()
	}
	return nil
}*/

func (s *Server) AddInterceptor(interceptor ziface.IInterceptor) {
	s.msgHandler.AddInterceptor(interceptor)
}

//func (s *Server) SetWebsocketAuth(f func(r *http.Request) error) {
//	s.websocketAuth = f
//}

func (s *Server) ServerName() string {
	return s.Name
}

func init() {
}
