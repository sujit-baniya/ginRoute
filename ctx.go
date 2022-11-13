package ginRoute

import (
	"context"
	"github.com/sujit-baniya/framework/view"
	"mime/multipart"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	contracthttp "github.com/sujit-baniya/framework/contracts/http"
)

type Config struct {
	Mode        string `json:"mode"`
	ViewsLayout string `json:"views_layout"`
	Extension   string `json:"extension"`
	Path        string `json:"path"`
	View        *view.Engine
}

type Context struct {
	instance *gin.Context
	config   Config
}

func NewContext(ctx *gin.Context, config Config) contracthttp.Context {
	ct := &Context{instance: ctx, config: config}
	return ct
}

func (c *Context) Context() context.Context {
	ctx := context.Background()
	for key, value := range c.instance.Keys {
		ctx = context.WithValue(ctx, key, value)
	}

	return ctx
}

func (c *Context) Origin() *http.Request {
	return c.instance.Request
}

func (c *Context) String(format string, values ...any) error {
	c.instance.String(c.StatusCode(), format, values...)
	return nil
}

func (c *Context) Json(obj any) error {
	c.instance.JSON(c.StatusCode(), obj)
	return nil
}

func (c *Context) Render(name string, bind any, layouts ...string) error {
	return c.config.View.Render(c.instance.Writer, name, bind, layouts...)
}

func (c *Context) SendFile(filepath string, compress ...bool) error {
	c.instance.File(filepath)
	return nil
}

func (c *Context) Download(filepath, filename string) error {
	c.instance.FileAttachment(filepath, filename)
	return nil
}

func (c *Context) SetHeader(key, value string) contracthttp.Context {
	c.instance.Header(key, value)
	return c
}

func (c *Context) StatusCode() int {
	return c.instance.Writer.Status()
}

func (c *Context) Status(status int) contracthttp.Context {
	c.instance.Status(status)
	return c
}

func (c *Context) Vary(field string, values ...string) {
	c.Append(field)
}

func (c *Context) Append(field string, values ...string) {
	if len(values) == 0 {
		return
	}
	h := c.instance.GetHeader(field)
	originalH := h
	for _, value := range values {
		if len(h) == 0 {
			h = value
		} else if h != value && !strings.HasPrefix(h, value+",") && !strings.HasSuffix(h, " "+value) &&
			!strings.Contains(h, " "+value+",") {
			h += ", " + value
		}
	}
	if originalH != h {
		c.SetHeader(field, h)
	}
}

func (c *Context) WithValue(key string, value any) {
	c.instance.Set(key, value)
}

func (c *Context) Deadline() (deadline time.Time, ok bool) {
	return c.instance.Deadline()
}

func (c *Context) Done() <-chan struct{} {
	return c.instance.Done()
}

func (c *Context) Err() error {
	return c.instance.Err()
}

func (c *Context) Value(key any) any {
	return c.instance.Value(key)
}

func (c *Context) Params(key string) string {
	return c.instance.Param(key)
}

func (c *Context) Query(key, defaultValue string) string {
	return c.instance.DefaultQuery(key, defaultValue)
}

func (c *Context) Form(key, defaultValue string) string {
	return c.instance.DefaultPostForm(key, defaultValue)
}

func (c *Context) Bind(obj any) error {
	return c.instance.Bind(obj)
}

func (c *Context) SaveFile(name string, dst string) error {
	file, err := c.File(name)
	if err != nil {
		return err
	}
	return c.instance.SaveUploadedFile(file, dst)
}

func (c *Context) File(name string) (*multipart.FileHeader, error) {
	return c.instance.FormFile(name)
}

func (c *Context) Header(key, defaultValue string) string {
	header := c.instance.GetHeader(key)
	if header != "" {
		return header
	}

	return defaultValue
}

func (c *Context) Headers() http.Header {
	return c.instance.Request.Header
}

func (c *Context) Method() string {
	return c.instance.Request.Method
}

func (c *Context) Url() string {
	return c.instance.Request.RequestURI
}

func (c *Context) FullUrl() string {
	prefix := "https://"
	if c.instance.Request.TLS == nil {
		prefix = "http://"
	}

	if c.instance.Request.Host == "" {
		return ""
	}

	return prefix + c.instance.Request.Host + c.instance.Request.RequestURI
}

func (c *Context) AbortWithStatus(code int) {
	c.instance.AbortWithStatus(code)
}

func (c *Context) Next() error {
	c.instance.Next()
	return nil
}

func (c *Context) Cookies(key string, defaultValue ...string) string {
	str, _ := c.instance.Cookie(key)
	return str
}

func (c *Context) Cookie(co *contracthttp.Cookie) {
	switch co.SameSite {
	case "Lax":
		c.instance.SetSameSite(http.SameSiteLaxMode)
		break
	case "None":
		c.instance.SetSameSite(http.SameSiteNoneMode)
		break
	case "Strict":
		c.instance.SetSameSite(http.SameSiteStrictMode)
		break
	default:
		c.instance.SetSameSite(http.SameSiteDefaultMode)
	}
	c.instance.SetCookie(co.Name, co.Value, co.MaxAge, co.Path, co.Domain, co.Secure, co.HTTPOnly)
}

func (c *Context) Path() string {
	return c.instance.Request.URL.Path
}

func (c *Context) EngineContext() any {
	return c.instance
}

func (c *Context) Secure() bool {
	if c.instance.Request.Proto == "https" {
		return true
	}
	return false
}

func (c *Context) Ip() string {
	return c.instance.ClientIP()
}
