package ijwt

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"crypto/sha512"
	"crypto/subtle"
	"encoding/base64"
	"encoding/json"
	"errors"
	"exam_api/internal/pkg/iclient"
	"fmt"
	"hash"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
)

// 算法类型
type Algorithm string

const (
	HS256 Algorithm = "HS256"
	HS384 Algorithm = "HS384"
	HS512 Algorithm = "HS512"
)

// 错误定义
var (
	ErrInvalidToken      = errors.New("invalid token")
	ErrInvalidAlgorithm  = errors.New("invalid algorithm")
	ErrTokenExpired      = errors.New("token has expired")
	ErrTokenNotValidYet  = errors.New("token is not valid yet")
	ErrSignatureInvalid  = errors.New("signature is invalid")
	ErrMalformedToken    = errors.New("malformed token")
	ErrSessionNotFound   = errors.New("session not found")
	ErrSequenceMismatch  = errors.New("invalid sequence number")
	ErrFingerprintChange = errors.New("client fingerprint changed")
)

// 会话存储结构
type SessionStore struct {
	sync.RWMutex
	sessions map[string]ExamSessionClaims
}

// 考试会话声明
type ExamSessionClaims struct {
	SessionID     string `json:"sid"` // 会话ID
	AssociationId string `json:"eid"` // 考试ID
	UserID        string `json:"uid"` // 用户ID
	ClientFP      string `json:"cfp"` // 客户端指纹
	IssuedAt      int64  `json:"iat"` // 签发时间
	ExpiresAt     int64  `json:"exp"` // 过期时间
	NotBefore     int64  `json:"nbf"` // 生效时间
}

// 主访问令牌声明
type AccessClaims struct {
	UserID   string `json:"uid"`   // 用户ID
	Username string `json:"uname"` // 用户名
	Role     string `json:"role"`  // 用户角色
	Exp      int64  `json:"exp"`   // 过期时间
	Iat      int64  `json:"iat"`   // 签发时间
}

// SecureJWT 安全JWT管理器
type SecureJWT struct {
	accessSecret []byte        // 主令牌密钥
	examSecret   []byte        // 考试令牌密钥
	algorithm    Algorithm     // 签名算法
	accessExpiry time.Duration // 主令牌有效期
	examExpiry   time.Duration // 考试令牌默认有效期
	issuer       string        // 签发者
	audience     string        // 受众
	sessionStore *SessionStore // 会话存储
}

// 安全配置选项
type SecurityOption func(*SecureJWT)

// 创建SecureJWT实例
func NewSecureJWT(accessSecret, examSecret string, options ...SecurityOption) *SecureJWT {
	j := &SecureJWT{
		accessSecret: []byte(accessSecret),
		examSecret:   []byte(examSecret),
		algorithm:    HS256,
		accessExpiry: 24 * time.Hour, // 主令牌默认24小时
		examExpiry:   2 * time.Hour,  // 考试令牌默认2小时
		issuer:       "exam-system",
		audience:     "exam-client",
		sessionStore: &SessionStore{sessions: make(map[string]ExamSessionClaims)},
	}

	for _, option := range options {
		option(j)
	}

	return j
}

// 配置选项
func WithAlgorithm(alg Algorithm) SecurityOption {
	return func(j *SecureJWT) {
		j.algorithm = alg
	}
}

func WithAccessExpiry(duration time.Duration) SecurityOption {
	return func(j *SecureJWT) {
		j.accessExpiry = duration
	}
}

func WithExamExpiry(duration time.Duration) SecurityOption {
	return func(j *SecureJWT) {
		j.examExpiry = duration
	}
}

func WithIssuer(issuer string) SecurityOption {
	return func(j *SecureJWT) {
		j.issuer = issuer
	}
}

func WithAudience(audience string) SecurityOption {
	return func(j *SecureJWT) {
		j.audience = audience
	}
}

// =====================
// 主令牌功能
// =====================

// 生成主访问令牌
func (j *SecureJWT) GenerateAccessToken(userID, username, role string) (string, error) {
	now := time.Now()
	claims := AccessClaims{
		UserID:   userID,
		Username: username,
		Role:     role,
		Iat:      now.Unix(),
		Exp:      now.Add(j.accessExpiry).Unix(),
	}

	return j.generateToken(claims, j.accessSecret)
}

// 验证主访问令牌
func (j *SecureJWT) VerifyAccessToken(tokenString string) (*AccessClaims, error) {
	var claims AccessClaims
	if err := j.parseToken(tokenString, j.accessSecret, &claims); err != nil {
		return nil, err
	}

	// 检查是否过期
	if time.Now().Unix() > claims.Exp {
		return nil, ErrTokenExpired
	}

	return &claims, nil
}

// =====================
// 考试令牌功能
// =====================

// 生成考试令牌
func (j *SecureJWT) GenerateExamToken(
	accessToken string,
	AssociationId string,
	examDuration time.Duration,
	clientInfo *iclient.ClientInfo,
) (string, *ExamSessionClaims, error) {
	// 验证主令牌
	accessClaims, err := j.VerifyAccessToken(accessToken)
	if err != nil {
		return "", nil, fmt.Errorf("invalid access token: %w", err)
	}

	// 生成唯一会话ID
	sessionID := uuid.New().String()

	// 获取客户端指纹
	clientFP := j.generateClientFingerprintFromInfo(clientInfo)

	now := time.Now()
	expiry := examDuration
	if expiry == 0 {
		expiry = j.examExpiry
	}

	claims := ExamSessionClaims{
		SessionID:     sessionID,
		AssociationId: AssociationId,
		UserID:        accessClaims.UserID,
		ClientFP:      clientFP,
		IssuedAt:      now.Unix(),
		ExpiresAt:     now.Add(expiry).Unix(),
		NotBefore:     now.Unix(),
	}

	// 生成令牌
	token, err := j.generateToken(claims, j.examSecret)
	if err != nil {
		return "", nil, err
	}

	return token, &claims, nil
}

// 验证考试令牌
func (j *SecureJWT) VerifyExamToken(
	accessToken string,
	examToken string,
	clientInfo *iclient.ClientInfo) (*ExamSessionClaims, error) {
	// 1. 验证主访问令牌
	accessClaims, err := j.VerifyAccessToken(accessToken)
	if err != nil {
		return nil, fmt.Errorf("access token invalid: %w", err)
	}

	// 2. 解析考试令牌
	var examClaims ExamSessionClaims
	if err := j.parseToken(examToken, j.examSecret, &examClaims); err != nil {
		return nil, fmt.Errorf("exam token invalid: %w", err)
	}

	// 3. 检查是否过期
	if time.Now().Unix() > examClaims.ExpiresAt {
		return nil, ErrTokenExpired
	}
	// 4. 检查令牌关联性
	if examClaims.UserID != accessClaims.UserID {
		return nil, errors.New("token user mismatch")
	}
	return &examClaims, nil
}

// =====================
// 辅助方法
// =====================

// 生成令牌
func (j *SecureJWT) generateToken(claims interface{}, secret []byte) (string, error) {
	// 序列化 Header
	header := map[string]interface{}{
		"typ": "JWT",
		"alg": j.algorithm,
	}
	headerJSON, err := json.Marshal(header)
	if err != nil {
		return "", fmt.Errorf("failed to marshal header: %w", err)
	}
	headerBase64 := base64.RawURLEncoding.EncodeToString(headerJSON)

	// 序列化 Claims
	claimsJSON, err := json.Marshal(claims)
	if err != nil {
		return "", fmt.Errorf("failed to marshal claims: %w", err)
	}
	claimsBase64 := base64.RawURLEncoding.EncodeToString(claimsJSON)

	// 拼接签名内容
	signingInput := fmt.Sprintf("%s.%s", headerBase64, claimsBase64)

	// 生成签名
	signature, err := j.sign(signingInput, secret)
	if err != nil {
		return "", err
	}

	// 生成完整 Token
	return fmt.Sprintf("%s.%s", signingInput, signature), nil
}

// 解析令牌
func (j *SecureJWT) parseToken(tokenString string, secret []byte, claims interface{}) error {
	// 分割 Token
	parts := strings.Split(tokenString, ".")
	if len(parts) != 3 {
		return ErrMalformedToken
	}

	// 验证签名
	signingInput := fmt.Sprintf("%s.%s", parts[0], parts[1])
	if err := j.verifySignature(signingInput, parts[2], secret); err != nil {
		return err
	}

	// 解析Claims
	claimsJSON, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return fmt.Errorf("failed to decode claims: %w", err)
	}

	if err := json.Unmarshal(claimsJSON, claims); err != nil {
		return fmt.Errorf("failed to parse claims: %w", err)
	}

	return nil
}

// 签名
func (j *SecureJWT) sign(signingInput string, key []byte) (string, error) {
	var h hash.Hash

	switch j.algorithm {
	case HS256:
		h = hmac.New(sha256.New, key)
	case HS384:
		h = hmac.New(sha512.New384, key)
	case HS512:
		h = hmac.New(sha512.New, key)
	default:
		return "", ErrInvalidAlgorithm
	}

	h.Write([]byte(signingInput))
	signature := h.Sum(nil)

	return base64.RawURLEncoding.EncodeToString(signature), nil
}

// 验证签名
func (j *SecureJWT) verifySignature(signingInput, signature string, key []byte) error {
	// 计算预期签名
	expectedSig, err := j.sign(signingInput, key)
	if err != nil {
		return err
	}

	// 使用恒定时间比较防止时序攻击
	if subtle.ConstantTimeCompare([]byte(expectedSig), []byte(signature)) != 1 {
		return ErrSignatureInvalid
	}

	return nil
}

// 生成客户端指纹 (网页版)
func (j *SecureJWT) generateClientFingerprint(r *http.Request) string {
	// 收集浏览器指纹信息
	fingerprint := fmt.Sprintf("%s|%s|%s|%s",
		r.UserAgent(),                      // 浏览器标识
		r.Header.Get("Accept"),             // 接受的内容类型
		r.Header.Get("Accept-Language"),    // 语言
		r.Header.Get("Sec-CH-UA-Platform"), // 浏览器平台
	)

	// 添加前端计算的指纹
	if canvasFP := r.Header.Get("X-Canvas-FP"); canvasFP != "" {
		fingerprint += "|" + canvasFP
	}

	if webglFP := r.Header.Get("X-WebGL-FP"); webglFP != "" {
		fingerprint += "|" + webglFP
	}

	if fontsFP := r.Header.Get("X-Fonts-FP"); fontsFP != "" {
		fingerprint += "|" + fontsFP
	}

	// 哈希生成最终指纹
	h := sha256.New()
	h.Write([]byte(fingerprint))
	return base64.RawURLEncoding.EncodeToString(h.Sum(nil))
}

// 从 ClientInfo 生成指纹
func (j *SecureJWT) generateClientFingerprintFromInfo(info *iclient.ClientInfo) string {
	fingerprint := fmt.Sprintf("%s|%s|%s|%s",
		info.UserAgent,
		info.Accept,
		info.AcceptLang,
		info.Platform,
	)
	if info.CanvasFP != "" {
		fingerprint += "|" + info.CanvasFP
	}
	if info.WebGLFP != "" {
		fingerprint += "|" + info.WebGLFP
	}
	if info.FontsFP != "" {
		fingerprint += "|" + info.FontsFP
	}

	h := sha256.New()
	h.Write([]byte(fingerprint))
	return base64.RawURLEncoding.EncodeToString(h.Sum(nil))
}

// 判断指纹变化是否允许
func (j *SecureJWT) isFingerprintChangeAllowed(original, current string) bool {
	// 实际应用中应有更复杂的逻辑
	// 这里简化为只检查前两部分是否一致
	origParts := strings.Split(original, "|")
	currParts := strings.Split(current, "|")

	if len(origParts) < 2 || len(currParts) < 2 {
		return false
	}

	// 检查UserAgent和平台
	return origParts[0] == currParts[0] && origParts[3] == currParts[3]
}

// =====================
// 工具函数
// =====================

// 获取客户端IP
func getClientIP(r *http.Request) string {
	// 检查代理头
	if ip := r.Header.Get("X-Forwarded-For"); ip != "" {
		return strings.Split(ip, ",")[0]
	}
	if ip := r.Header.Get("X-Real-Ip"); ip != "" {
		return ip
	}

	// 使用远程地址
	host, _, _ := strings.Cut(r.RemoteAddr, ":")
	return host
}

// 生成强随机JTI
func generateJTI() string {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return fmt.Sprintf("%d-%d", time.Now().UnixNano(), SecureInt63())
	}
	return base64.RawURLEncoding.EncodeToString(b)
}
