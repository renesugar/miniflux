// Copyright 2018 Frédéric Guillot. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package ui

import (
	"net/http"

	"github.com/miniflux/miniflux/http/context"
	"github.com/miniflux/miniflux/http/request"
	"github.com/miniflux/miniflux/http/response/html"
	"github.com/miniflux/miniflux/http/route"
	"github.com/miniflux/miniflux/logger"
	"github.com/miniflux/miniflux/model"
	"github.com/miniflux/miniflux/storage"
	"github.com/miniflux/miniflux/ui/session"
	"github.com/miniflux/miniflux/ui/view"
)

// ShowUnreadEntry shows a single feed entry in "unread" mode.
func (c *Controller) ShowUnreadEntry(w http.ResponseWriter, r *http.Request) {
	ctx := context.New(r)

	user, err := c.store.UserByID(ctx.UserID())
	if err != nil {
		html.ServerError(w, err)
		return
	}

	entryID, err := request.IntParam(r, "entryID")
	if err != nil {
		html.BadRequest(w, err)
		return
	}

	builder := c.store.NewEntryQueryBuilder(user.ID)
	builder.WithEntryID(entryID)
	builder.WithoutStatus(model.EntryStatusRemoved)

	entry, err := builder.GetEntry()
	if err != nil {
		html.ServerError(w, err)
		return
	}

	if entry == nil {
		html.NotFound(w)
		return
	}

	entryPaginationBuilder := storage.NewEntryPaginationBuilder(c.store, user.ID, entry.ID, user.EntryDirection)
	entryPaginationBuilder.WithStatus(model.EntryStatusUnread)
	prevEntry, nextEntry, err := entryPaginationBuilder.Entries()
	if err != nil {
		html.ServerError(w, err)
		return
	}

	nextEntryRoute := ""
	if nextEntry != nil {
		nextEntryRoute = route.Path(c.router, "unreadEntry", "entryID", nextEntry.ID)
	}

	prevEntryRoute := ""
	if prevEntry != nil {
		prevEntryRoute = route.Path(c.router, "unreadEntry", "entryID", prevEntry.ID)
	}

	// We change the status here, otherwise we cannot get the pagination for unread items.
	if entry.Status == model.EntryStatusUnread {
		err = c.store.SetEntriesStatus(user.ID, []int64{entry.ID}, model.EntryStatusRead)
		if err != nil {
			logger.Error("[Controller:ShowUnreadEntry] %v", err)
			html.ServerError(w, nil)
			return
		}
	}

	sess := session.New(c.store, ctx)
	view := view.New(c.tpl, ctx, sess)
	view.Set("entry", entry)
	view.Set("prevEntry", prevEntry)
	view.Set("nextEntry", nextEntry)
	view.Set("nextEntryRoute", nextEntryRoute)
	view.Set("prevEntryRoute", prevEntryRoute)
	view.Set("menu", "unread")
	view.Set("user", user)
	view.Set("countUnread", c.store.CountUnreadEntries(user.ID))
	view.Set("hasSaveEntry", c.store.HasSaveEntry(user.ID))

	html.OK(w, view.Render("entry"))
}
