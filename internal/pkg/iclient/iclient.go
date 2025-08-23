package iclient

// ClientInfo 表示客户端上下文信息
type ClientInfo struct {
	IP         string
	UserAgent  string
	Accept     string
	AcceptLang string
	Platform   string
	CanvasFP   string
	WebGLFP    string
	FontsFP    string
}
