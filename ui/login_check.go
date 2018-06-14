package ui

import (
	"net/http"

	"github.com/miniflux/miniflux/http/context"
	"github.com/miniflux/miniflux/http/cookie"
	"github.com/miniflux/miniflux/http/request"
	"github.com/miniflux/miniflux/http/response"
	"github.com/miniflux/miniflux/http/response/html"
	"github.com/miniflux/miniflux/http/route"
	"github.com/miniflux/miniflux/logger"
	"github.com/miniflux/miniflux/ui/form"
	"github.com/miniflux/miniflux/ui/session"
	"github.com/miniflux/miniflux/ui/view"
)

// CheckLogin validates the username/password and redirects the user to the unread page.
func (c *Controller) CheckLogin(w http.ResponseWriter, r *http.Request) {
	ctx := context.New(r)
	sess := session.New(c.store, ctx)

	authForm := form.NewAuthForm(r)

	view := view.New(c.tpl, ctx, sess)
	view.Set("errorMessage", "Invalid username or password.")
	view.Set("form", authForm)

	if err := authForm.Validate(); err != nil {
		logger.Error("[Controller:CheckLogin] %v", err)
		html.OK(w, view.Render("login"))
		return
	}

	if err := c.store.CheckPassword(authForm.Username, authForm.Password); err != nil {
		logger.Error("[Controller:CheckLogin] %v", err)
		html.OK(w, view.Render("login"))
		return
	}

	sessionToken, userID, err := c.store.CreateUserSession(authForm.Username, r.UserAgent(), request.RealIP(r))
	if err != nil {
		html.ServerError(w, err)
		return
	}

	logger.Info("[Controller:CheckLogin] username=%s just logged in", authForm.Username)
	c.store.SetLastLogin(userID)

	userLanguage, err := c.store.UserLanguage(userID)
	if err != nil {
		html.ServerError(w, err)
		return
	}

	sess.SetLanguage(userLanguage)

	http.SetCookie(w, cookie.New(
		cookie.CookieUserSessionID,
		sessionToken,
		c.cfg.IsHTTPS,
		c.cfg.BasePath(),
	))

	response.Redirect(w, r, route.Path(c.router, "unread"))
}
