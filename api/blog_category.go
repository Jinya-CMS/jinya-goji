package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"jinya.de/goji/database"
)

func getAllBlogCategories(w http.ResponseWriter, r *http.Request) {
	encoder := json.NewEncoder(w)

	categories, err := database.FindAllBlogCategories()
	if err != nil {
		sendError("Failed to find all blog categories", http.StatusInternalServerError, encoder, w, err)
		return
	}

	encoder.Encode(categories)
}

func getBlogCategoryById(w http.ResponseWriter, r *http.Request) {
	encoder := json.NewEncoder(w)
	vars := mux.Vars(r)

	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		sendError("Id is invalid", http.StatusBadRequest, encoder, w, err)
		return
	}

	category, err := database.FindBlogCategoryById(id)
	if err != nil {
		sendError("Blog category not found", http.StatusNotFound, encoder, w, err)
		return
	}

	encoder.Encode(category)
}

func createBlogCategory(w http.ResponseWriter, r *http.Request) {
	encoder := json.NewEncoder(w)
	category := database.BlogCategory{}

	body, err := BodyFromRequest(r, "name")
	if err != nil {
		sendError("Invalid body", http.StatusBadRequest, encoder, w, err)
		return
	}

	if name, ok := GetFromBody[string](body, "name"); ok {
		category.Name = *name
	}
	if description, ok := GetFromBody[string](body, "description"); ok {
		category.Description = description
	}
	if parentId, ok := GetIntFromBody(body, "parentId"); ok {
		if _, err = database.FindBlogCategoryById(parentId); err != nil {
			sendError("Parent category not found", http.StatusConflict, encoder, w, err)
			return
		}

		category.ParentId = &parentId
	}
	if webhookEnabled, ok := GetFromBody[bool](body, "webhookEnabled"); ok {
		category.WebhookEnabled = *webhookEnabled
	}
	if webhookUrl, ok := GetFromBody[string](body, "webhookUrl"); ok {
		category.WebhookUrl = webhookUrl
	}

	err = category.Create()
	if err != nil {
		switch detectErrorType(err) {
		case errorTypeUnique:
			sendError("Name is not unique", http.StatusConflict, encoder, w, err)
			return
		default:
			sendError("Failed to update blog category", http.StatusInternalServerError, encoder, w, err)
			return
		}
	}

	w.WriteHeader(http.StatusCreated)
	encoder.Encode(&category)
}

func updateBlogCategory(w http.ResponseWriter, r *http.Request) {
	encoder := json.NewEncoder(w)
	vars := mux.Vars(r)

	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		sendError("Id is invalid", http.StatusBadRequest, encoder, w, err)
		return
	}

	category, err := database.FindBlogCategoryById(id)
	if err != nil {
		sendError("Blog category not found", http.StatusNotFound, encoder, w, err)
		return
	}

	body, err := BodyFromRequest(r)
	if err != nil {
		sendError("Invalid body", http.StatusBadRequest, encoder, w, err)
		return
	}

	if name, ok := GetFromBody[string](body, "name"); ok {
		category.Name = *name
	}
	if description, ok := GetFromBody[string](body, "description"); ok {
		category.Description = description
	}
	if parentId, ok := GetIntFromBody(body, "parentId"); ok {
		if _, err = database.FindBlogCategoryById(parentId); err != nil {
			sendError("Parent category not found", http.StatusConflict, encoder, w, err)
			return
		}

		category.ParentId = &parentId
	}
	if webhookEnabled, ok := GetFromBody[bool](body, "webhookEnabled"); ok {
		category.WebhookEnabled = *webhookEnabled
	}
	if webhookUrl, ok := GetFromBody[string](body, "webhookUrl"); ok {
		category.WebhookUrl = webhookUrl
	}

	err = category.Update()
	if err != nil {
		switch detectErrorType(err) {
		case errorTypeUnique:
			sendError("Name is not unique", http.StatusConflict, encoder, w, err)
			return
		default:
			sendError("Failed to update blog category", http.StatusInternalServerError, encoder, w, err)
			return
		}
	}

	w.WriteHeader(http.StatusNoContent)
}

func deleteBlogCategory(w http.ResponseWriter, r *http.Request) {
	encoder := json.NewEncoder(w)
	vars := mux.Vars(r)

	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		sendError("Id is invalid", http.StatusBadRequest, encoder, w, err)
		return
	}

	category, err := database.FindBlogCategoryById(id)
	if err != nil {
		sendError("Blog category not found", http.StatusNotFound, encoder, w, err)
		return
	}

	err = category.Delete()
	if err != nil {
		sendError("Failed to delete blog category", http.StatusInternalServerError, encoder, w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func getPostsByCategory(w http.ResponseWriter, r *http.Request) {
	encoder := json.NewEncoder(w)
	vars := mux.Vars(r)

	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		sendError("Id is invalid", http.StatusBadRequest, encoder, w, err)
		return
	}

	category, err := database.FindBlogCategoryById(id)
	if err != nil {
		sendError("Blog category not found", http.StatusNotFound, encoder, w, err)
		return
	}

	posts, err := category.GetBlogPosts(true, false)
	if err != nil {
		sendError("Failed to get blog posts", http.StatusInternalServerError, encoder, w, err)
		return
	}

	encoder.Encode(posts)
}
