package main

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/sullyh7/social/internal/store"
)

type postKey string

const postCtx postKey = "post"

type CreatePostPayload struct {
	Content string   `json:"content" validate:"required,max=100"`
	Title   string   `json:"title" validate:"required,max=1000"`
	Tags    []string `json:"tags"`
}

type UpdatePostPayload struct {
	Content *string `json:"content" validate:"required,max=100"`
	Title   *string `json:"title" validate:"required,max=1000"`
}

func (a *application) createPostHandler(w http.ResponseWriter, r *http.Request) {
	payload := new(CreatePostPayload)
	if err := readJson(w, r, payload); err != nil {
		a.badRequest(w, r, err)
		return
	}

	if err := Validator.Struct(payload); err != nil {
		a.badRequest(w, r, err)
		return
	}

	userId := 1
	ctx := r.Context()
	post := &store.Post{
		Content: payload.Content,
		Title:   payload.Title,
		Tags:    payload.Tags,
		UserID:  int64(userId),
	}
	if err := a.store.Posts.Create(ctx, post); err != nil {
		a.internalServerError(w, r, err)
		return
	}
	if err := writeJSON(w, http.StatusCreated, post); err != nil {
		a.internalServerError(w, r, err)
		return
	}
}

func (a *application) getPostHandler(w http.ResponseWriter, r *http.Request) {
	post := getPostFromCtx(r)

	comments, err := a.store.Comments.GetByPostID(r.Context(), post.ID)
	if err != nil {
		a.internalServerError(w, r, err)
		return
	}
	post.Comments = comments
	if err := writeJSON(w, http.StatusOK, post); err != nil {
		a.internalServerError(w, r, err)
		return
	}
}

func (a *application) patchPostHandler(w http.ResponseWriter, r *http.Request) {
	post := getPostFromCtx(r)
	payload := new(UpdatePostPayload)
	if err := readJson(w, r, payload); err != nil {
		a.badRequest(w, r, err)
		return
	}
	if err := Validator.Struct(payload); err != nil {
		a.badRequest(w, r, err)
		return
	}
	if payload.Content == nil || payload.Title == nil {
		a.badRequest(w, r, errors.New("bad request"))
		return
	}
	post.Content = *payload.Content
	post.Title = *payload.Title
	if err := writeJSON(w, http.StatusOK, post); err != nil {
		a.internalServerError(w, r, err)
		return
	}
}

func (a *application) deletePostHandler(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "postID")
	ctx := r.Context()
	postID, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		a.badRequest(w, r, err)
		return
	}
	err = a.store.Posts.DeleteByID(ctx, postID)
	if err != nil {
		switch {
		case errors.Is(store.ErrNotFound, err):
			a.notFound(w, r, err)
			return
		default:
			a.internalServerError(w, r, err)
			return
		}
	}
	if err := writeJSON(w, http.StatusOK, map[string]string{"message": "deleted"}); err != nil {
		a.internalServerError(w, r, err)
		return
	}
}

func (a *application) postCtxMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		idParam := chi.URLParam(r, "postID")
		ctx := r.Context()
		postID, err := strconv.ParseInt(idParam, 10, 64)
		if err != nil {
			a.badRequest(w, r, err)
			return
		}
		post, err := a.store.Posts.GetByID(ctx, postID)
		if err != nil {
			switch {
			case errors.Is(store.ErrNotFound, err):
				a.notFound(w, r, err)
				return
			default:
				a.internalServerError(w, r, err)
				return
			}
		}
		ctx = context.WithValue(ctx, postCtx, post)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func getPostFromCtx(r *http.Request) *store.Post {
	post, _ := r.Context().Value("post").(*store.Post)
	return post
}
