// Copyright 2018 Frédéric Guillot. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package ui

import (
	"net/http"

	"github.com/miniflux/miniflux/http/context"
	"github.com/miniflux/miniflux/http/request"
	"github.com/miniflux/miniflux/http/response/html"
	"github.com/miniflux/miniflux/ui/form"
	"github.com/miniflux/miniflux/ui/session"
	"github.com/miniflux/miniflux/ui/view"
)

// EditCategory shows the form to modify a category.
func (c *Controller) EditCategory(w http.ResponseWriter, r *http.Request) {
	ctx := context.New(r)
	sess := session.New(c.store, ctx)
	view := view.New(c.tpl, ctx, sess)

	user, err := c.store.UserByID(ctx.UserID())
	if err != nil {
		html.ServerError(w, err)
		return
	}

	categoryID, err := request.IntParam(r, "categoryID")
	if err != nil {
		html.BadRequest(w, err)
		return
	}

	category, err := c.store.Category(ctx.UserID(), categoryID)
	if err != nil {
		html.ServerError(w, err)
		return
	}

	if category == nil {
		html.NotFound(w)
		return
	}

	categoryForm := form.CategoryForm{
		Title: category.Title,
	}

	view.Set("form", categoryForm)
	view.Set("category", category)
	view.Set("menu", "categories")
	view.Set("user", user)
	view.Set("countUnread", c.store.CountUnreadEntries(user.ID))

	html.OK(w, view.Render("edit_category"))
}
