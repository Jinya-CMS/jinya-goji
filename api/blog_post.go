package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"jinya.de/goji/api/middleware"
	"jinya.de/goji/database"
)

func getAllBlogPosts(w http.ResponseWriter, r *http.Request) {
	encoder := json.NewEncoder(w)

	posts, err := database.FindAllBlogPosts()
	if err != nil {
		sendError("Failed to find all blog posts", http.StatusInternalServerError, encoder, w, err)
		return
	}

	encoder.Encode(posts)
}

func getBlogPostById(w http.ResponseWriter, r *http.Request) {
	encoder := json.NewEncoder(w)
	vars := mux.Vars(r)

	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		sendError("Id is invalid", http.StatusBadRequest, encoder, w, err)
		return
	}

	post, err := database.FindBlogPostById(id)
	if err != nil {
		sendError("Blog post not found", http.StatusNotFound, encoder, w, err)
		return
	}

	encoder.Encode(post)
}

func executeBlogPostHook(post *database.BlogPost, r *http.Request) {
	category, err := post.GetCategory()
	if err != nil {
		slog.Error("Failed to get category for hook execution", "error", err, "post", post.Id)
		return
	}

	if post.Public && category != nil && category.WebhookEnabled && category.WebhookUrl != nil {
		slog.Debug("Category has webhook, trigger it", "category", category.Id, "post", post.Id)
		url := r.URL.Scheme + "://" + r.URL.Host + "/" + post.CreatedAt.Format("2006/01/02") + "/" + post.Slug

		body := map[string]any{
			"post": post,
			"url":  url,
		}
		data, err := json.Marshal(body)
		if err != nil {
			slog.Error("Failed to marshal post for hook execution", "error", err, "post", post.Id)
			return
		}

		req, err := http.NewRequest(http.MethodPost, *category.WebhookUrl, bytes.NewBuffer(data))
		if err != nil {
			slog.Error("Failed to create request for hook execution", "error", err, "post", post.Id)
			return
		}

		req.Header.Add("Content-Type", "application/json")
		req.Header.Add("Referer", r.RequestURI)

		res, err := http.DefaultClient.Do(req)
		if err != nil {
			slog.Error("Failed to send request for hook execution", "error", err, "post", post.Id)
			return
		}

		defer res.Body.Close()
		if res.StatusCode > 299 {
			resBody, err := io.ReadAll(res.Body)
			if err == nil {
				slog.Warn("Hook target returned non 200 status code", "post", post.Id, "responseBody", string(resBody))
			} else {
				slog.Warn("Hook target returned non 200 status code", "post", post.Id)
			}
		}
	}
}

func createBlogPost(w http.ResponseWriter, r *http.Request) {
	encoder := json.NewEncoder(w)

	body, err := BodyFromRequest(r, "title", "slug")
	if err != nil {
		sendError("Invalid body", http.StatusBadRequest, encoder, w, err)
		return
	}

	artist := middleware.GetUserFromRequest(r)
	post := database.BlogPost{
		CreatorId:   artist.Id,
		UpdatedById: artist.Id,
	}

	if title, ok := GetFromBody[string](body, "title"); ok {
		post.Title = *title
	}
	if slug, ok := GetFromBody[string](body, "slug"); ok {
		post.Slug = *slug
	}
	if category, ok := GetIntFromBody(body, "categoryId"); ok {
		if _, err = database.FindBlogCategoryById(category); err != nil {
			sendError("Category not found", http.StatusConflict, encoder, w, err)
			return
		}

		post.CategoryId = &category
	}
	if headerImage, ok := GetIntFromBody(body, "headerImageId"); ok {
		if _, err = database.FindFileById(headerImage); err != nil {
			sendError("Header image not found", http.StatusConflict, encoder, w, err)
			return
		}

		post.HeaderImageId = &headerImage
	}

	err = post.Create()
	if err != nil {
		switch detectErrorType(err) {
		case errorTypeUnique:
			sendError("Title is not unique", http.StatusConflict, encoder, w, err)
			return
		default:
			sendError("Failed to create blog post", http.StatusInternalServerError, encoder, w, err)
			return
		}
	}

	w.WriteHeader(http.StatusCreated)
	encoder.Encode(&post)

	go executeBlogPostHook(&post, r)
}

func updateBlogPost(w http.ResponseWriter, r *http.Request) {
	encoder := json.NewEncoder(w)
	vars := mux.Vars(r)

	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		sendError("Id is invalid", http.StatusBadRequest, encoder, w, err)
		return
	}

	post, err := database.FindBlogPostById(id)
	if err != nil {
		sendError("Blog post not found", http.StatusNotFound, encoder, w, err)
		return
	}

	body, err := BodyFromRequest(r)
	if err != nil {
		sendError("Invalid body", http.StatusBadRequest, encoder, w, err)
		return
	}

	if title, ok := GetFromBody[string](body, "title"); ok {
		post.Title = *title
	}
	if slug, ok := GetFromBody[string](body, "slug"); ok {
		post.Slug = *slug
	}
	if category, ok := GetIntFromBody(body, "categoryId"); ok {
		if _, err = database.FindBlogCategoryById(category); err != nil {
			sendError("Category not found", http.StatusConflict, encoder, w, err)
			return
		}

		post.CategoryId = &category
	}
	if headerImage, ok := GetIntFromBody(body, "headerImageId"); ok {
		if _, err = database.FindFileById(headerImage); err != nil {
			sendError("Header image not found", http.StatusConflict, encoder, w, err)
			return
		}

		post.HeaderImageId = &headerImage
	}

	artist := middleware.GetUserFromRequest(r)
	post.UpdatedById = artist.Id

	postWasPrivate := !post.Public

	err = post.Update()
	if err != nil {
		switch detectErrorType(err) {
		case errorTypeUnique:
			sendError("Title is not unique", http.StatusConflict, encoder, w, err)
			return
		default:
			sendError("Failed to create blog post", http.StatusInternalServerError, encoder, w, err)
			return
		}
	}

	w.WriteHeader(http.StatusNoContent)
	if postWasPrivate {
		go executeBlogPostHook(post, r)
	}
}

func deleteBlogPost(w http.ResponseWriter, r *http.Request) {
	encoder := json.NewEncoder(w)
	vars := mux.Vars(r)

	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		sendError("Id is invalid", http.StatusBadRequest, encoder, w, err)
		return
	}

	post, err := database.FindBlogPostById(id)
	if err != nil {
		sendError("Blog post not found", http.StatusNotFound, encoder, w, err)
		return
	}

	err = post.Delete()
	if err != nil {
		sendError("Failed to delete blog post", http.StatusInternalServerError, encoder, w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func getBlogPostSections(w http.ResponseWriter, r *http.Request) {
	encoder := json.NewEncoder(w)
	vars := mux.Vars(r)

	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		sendError("Id is invalid", http.StatusBadRequest, encoder, w, err)
		return
	}

	post, err := database.FindBlogPostById(id)
	if err != nil {
		sendError("Blog post not found", http.StatusNotFound, encoder, w, err)
		return
	}

	sections, err := post.GetSections()
	if err != nil {
		sendError("Failed to load sections", http.StatusInternalServerError, encoder, w, err)
		return
	}

	encoder.Encode(sections)
}

func updateBlogPostSections(w http.ResponseWriter, r *http.Request) {
	encoder := json.NewEncoder(w)
	decoder := json.NewDecoder(r.Body)
	vars := mux.Vars(r)

	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		sendError("Id is invalid", http.StatusBadRequest, encoder, w, err)
		return
	}

	post, err := database.FindBlogPostById(id)
	if err != nil {
		sendError("Blog post not found", http.StatusNotFound, encoder, w, err)
		return
	}

	var body []struct {
		Gallery *int    `json:"gallery"`
		File    *int    `json:"file"`
		Link    string  `json:"link"`
		Html    *string `json:"html"`
		Type    string  `json:"type"`
	}

	err = decoder.Decode(&body)
	if err != nil {
		sendError("Invalid body", http.StatusBadRequest, encoder, w, err)
		return
	}

	sections := make([]database.BlogPostSection, len(body))

	for i, section := range body {
		if (section.Type == "gallery" && section.Gallery == nil) || (section.Type == "file" && section.File == nil) || (section.Type == "html" && section.Html == nil) {
			sendError("Invalid body", http.StatusBadRequest, encoder, w, err)
			return
		}

		s := database.BlogPostSection{
			Html: section.Html,
		}
		if section.Link != "" {
			s.Link = &section.Link
		}

		if section.Gallery != nil {
			if _, err = database.FindGalleryById(*section.Gallery); err != nil {
				sendError(fmt.Sprintf("Gallery %d not found", *section.Gallery), http.StatusConflict, encoder, w, err)
				return
			}

			s.GalleryId = section.Gallery
		}
		if section.File != nil {
			if _, err = database.FindFileById(*section.File); err != nil {
				sendError(fmt.Sprintf("File %d not found", *section.File), http.StatusConflict, encoder, w, err)
				return
			}

			s.FileId = section.File
		}

		sections[i] = s
	}

	err = post.ReplaceSections(sections)
	if err != nil {
		sendError("Invalid body", http.StatusBadRequest, encoder, w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
