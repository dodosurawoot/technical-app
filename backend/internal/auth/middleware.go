package auth

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"airclean-tracker/backend/internal/config"
	"airclean-tracker/backend/internal/models"
	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

const userContextKey = "currentUser"

type Middleware struct {
	cfg      config.Config
	db       *gorm.DB
	verifier *oidc.IDTokenVerifier
}

type Claims struct {
	Subject           string   `json:"sub"`
	Email             string   `json:"email"`
	Name              string   `json:"name"`
	PreferredUsername string   `json:"preferred_username"`
	Role              string   `json:"role"`
	Groups            []string `json:"groups"`
}

func New(ctx context.Context, cfg config.Config, db *gorm.DB) (*Middleware, error) {
	m := &Middleware{cfg: cfg, db: db}
	if cfg.AuthentikIssuerURL == "" || cfg.AuthentikClientID == "" {
		return m, nil
	}
	provider, err := oidc.NewProvider(ctx, cfg.AuthentikIssuerURL)
	if err != nil {
		return nil, err
	}
	m.verifier = provider.Verifier(&oidc.Config{ClientID: cfg.AuthentikClientID})
	return m, nil
}

func (m *Middleware) RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		user, err := m.currentUser(c.Request.Context(), c.Request.Header)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}
		c.Set(userContextKey, user)
		c.Next()
	}
}

func (m *Middleware) RequireRole(roles ...string) gin.HandlerFunc {
	allowed := map[string]bool{}
	for _, role := range roles {
		allowed[role] = true
	}
	return func(c *gin.Context) {
		user := CurrentUser(c)
		if user == nil || !allowed[user.Role] {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "ไม่มีสิทธิ์ดำเนินการ"})
			return
		}
		c.Next()
	}
}

func CurrentUser(c *gin.Context) *models.User {
	v, ok := c.Get(userContextKey)
	if !ok {
		return nil
	}
	user, _ := v.(*models.User)
	return user
}

func (m *Middleware) currentUser(ctx context.Context, header http.Header) (*models.User, error) {
	if m.cfg.DevAuth && m.verifier == nil {
		return m.upsertUser(models.User{
			ProviderSubject: "dev-user",
			Email:           "admin@airclean.local",
			Name:            "ผู้ดูแลระบบ",
			Username:        "admin",
			Role:            models.RoleAdmin,
		})
	}
	token := bearer(header.Get("Authorization"))
	if token == "" {
		return nil, errors.New("ต้องเข้าสู่ระบบ")
	}
	idToken, err := m.verifier.Verify(ctx, token)
	if err != nil {
		return nil, errors.New("token ไม่ถูกต้อง")
	}
	var claims Claims
	if err := idToken.Claims(&claims); err != nil {
		return nil, errors.New("อ่านข้อมูลผู้ใช้ไม่ได้")
	}
	role := roleFromClaims(claims)
	if role == "" {
		role = models.RoleViewer
	}
	email := strings.TrimSpace(claims.Email)
	if email == "" {
		email = claims.Subject + "@authentik.local"
	}
	return m.upsertUser(models.User{
		ProviderSubject: claims.Subject,
		Email:           email,
		Name:            claims.Name,
		Username:        claims.PreferredUsername,
		Role:            role,
	})
}

func (m *Middleware) upsertUser(in models.User) (*models.User, error) {
	var user models.User
	err := m.db.Where("email = ?", in.Email).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		if in.Role == "" {
			in.Role = models.RoleViewer
		}
		if err := m.db.Create(&in).Error; err != nil {
			return nil, err
		}
		return &in, nil
	}
	if err != nil {
		return nil, err
	}
	user.ProviderSubject = in.ProviderSubject
	user.Name = in.Name
	user.Username = in.Username
	if user.Role == "" {
		user.Role = in.Role
	}
	if err := m.db.Save(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func bearer(raw string) string {
	if strings.HasPrefix(strings.ToLower(raw), "bearer ") {
		return strings.TrimSpace(raw[7:])
	}
	return ""
}

func roleFromClaims(claims Claims) string {
	if isRole(claims.Role) {
		return claims.Role
	}
	for _, group := range claims.Groups {
		g := strings.ToLower(group)
		if strings.Contains(g, "admin") {
			return models.RoleAdmin
		}
		if strings.Contains(g, "team") || strings.Contains(g, "technician") {
			return models.RoleTeam
		}
		if strings.Contains(g, "viewer") {
			return models.RoleViewer
		}
	}
	return ""
}

func isRole(role string) bool {
	return role == models.RoleAdmin || role == models.RoleTeam || role == models.RoleViewer
}
