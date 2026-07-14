package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"jinya.de/goji/database"
)

func getAllMenus(w http.ResponseWriter, r *http.Request) {
	encoder := json.NewEncoder(w)

	menus, err := database.FindAllMenus()
	if err != nil {
		sendError("Failed to find all menus", http.StatusInternalServerError, encoder, w, err)
		return
	}

	encoder.Encode(menus)
}

func getMenuById(w http.ResponseWriter, r *http.Request) {
	encoder := json.NewEncoder(w)
	vars := mux.Vars(r)

	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		sendError("Id is invalid", http.StatusBadRequest, encoder, w, err)
		return
	}

	menu, err := database.FindMenuById(id)
	if err != nil {
		sendError("Menu not found", http.StatusNotFound, encoder, w, err)
		return
	}

	encoder.Encode(menu)
}

func createMenu(w http.ResponseWriter, r *http.Request) {
	encoder := json.NewEncoder(w)

	body, err := BodyFromRequest(r, "name")
	if err != nil {
		sendError("Invalid body", http.StatusBadRequest, encoder, w, err)
		return
	}

	menu := database.Menu{}
	if name, ok := GetFromBody[string](body, "name"); ok {
		menu.Name = *name
	}
	if logo, ok := GetIntFromBody(body, "logo"); ok {
		menu.Logo = &logo
	}

	err = menu.Create()
	if err != nil {
		switch detectErrorType(err) {
		case errorTypeUnique:
			sendError("Name is not unique", http.StatusConflict, encoder, w, err)
			return
		case errorTypeReference:
			sendError("Logo was not found", http.StatusConflict, encoder, w, err)
			return
		default:
			sendError("Failed to create menu", http.StatusInternalServerError, encoder, w, err)
			return
		}
	}

	w.WriteHeader(http.StatusCreated)
	encoder.Encode(&menu)
}

func updateMenu(w http.ResponseWriter, r *http.Request) {
	encoder := json.NewEncoder(w)
	vars := mux.Vars(r)

	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		sendError("Id is invalid", http.StatusBadRequest, encoder, w, err)
		return
	}

	menu, err := database.FindMenuById(id)
	if err != nil {
		sendError("Menu not found", http.StatusNotFound, encoder, w, err)
		return
	}

	body, err := BodyFromRequest(r, "name")
	if err != nil {
		sendError("Invalid body", http.StatusBadRequest, encoder, w, err)
		return
	}

	if name, ok := GetFromBody[string](body, "name"); ok {
		menu.Name = *name
	}
	if logo, ok := GetIntFromBody(body, "logo"); ok {
		menu.Logo = &logo
	}

	err = menu.Create()
	if err != nil {
		switch detectErrorType(err) {
		case errorTypeUnique:
			sendError("Name is not unique", http.StatusConflict, encoder, w, err)
			return
		case errorTypeReference:
			sendError("Logo was not found", http.StatusConflict, encoder, w, err)
			return
		default:
			sendError("Failed to create menu", http.StatusInternalServerError, encoder, w, err)
			return
		}
	}

	w.WriteHeader(http.StatusCreated)
	encoder.Encode(&menu)
}

func deleteMenu(w http.ResponseWriter, r *http.Request) {
	encoder := json.NewEncoder(w)
	vars := mux.Vars(r)

	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		sendError("Id is invalid", http.StatusBadRequest, encoder, w, err)
		return
	}

	menu, err := database.FindMenuById(id)
	if err != nil {
		sendError("Menu not found", http.StatusNotFound, encoder, w, err)
		return
	}

	err = menu.Delete()
	if err != nil {
		sendError("Failed to delete menu", http.StatusInternalServerError, encoder, w, err)
		return
	}
}

func getMenuItems(w http.ResponseWriter, r *http.Request) {
	encoder := json.NewEncoder(w)
	vars := mux.Vars(r)

	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		sendError("Id is invalid", http.StatusBadRequest, encoder, w, err)
		return
	}

	menu, err := database.FindMenuById(id)
	if err != nil {
		sendError("Menu not found", http.StatusNotFound, encoder, w, err)
		return
	}

	items, err := menu.GetItems()
	if err != nil {
		sendError("Failed to get menu items", http.StatusInternalServerError, encoder, w, err)
		return
	}

	encoder.Encode(items)
}

type menuItemBody struct {
	Highlighted    bool           `json:"highlighted"`
	Title          string         `json:"title"`
	Route          string         `json:"route"`
	BlogHomePage   *bool          `json:"blogHomePage"`
	ArtistId       *int           `json:"artistId"`
	ClassicPageId  *int           `json:"classicPageId"`
	ModernPageId   *int           `json:"modernPageId"`
	GalleryId      *int           `json:"galleryId"`
	BlogCategoryId *int           `json:"blogCategoryId"`
	Items          []menuItemBody `json:"items"`
}

func menuItemBodyToMenuItem(menuId int, isMenu bool, data []menuItemBody) []database.MenuItem {
	res := make([]database.MenuItem, len(data))
	for i, item := range data {
		res[i] = database.MenuItem{
			Position:      i,
			Title:         item.Title,
			Route:         &item.Route,
			Highlighted:   item.Highlighted,
			ArtistId:      item.ArtistId,
			GalleryId:     item.GalleryId,
			ModernPageId:  item.ModernPageId,
			ClassicPageId: item.ClassicPageId,
			CategoryId:    item.BlogCategoryId,
			BlogHomePage:  item.BlogHomePage,
		}
		if isMenu {
			res[i].MenuId = &menuId
		}
		if len(item.Items) > 0 {
			res[i].Children = menuItemBodyToMenuItem(0, false, item.Items)
		}
	}

	return res
}

func updateMenuItems(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	encoder := json.NewEncoder(w)
	vars := mux.Vars(r)

	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		sendError("Id is invalid", http.StatusBadRequest, encoder, w, err)
		return
	}

	menu, err := database.FindMenuById(id)
	if err != nil {
		sendError("Menu not found", http.StatusNotFound, encoder, w, err)
		return
	}

	var body []menuItemBody
	err = decoder.Decode(&body)
	if err != nil {
		sendError("Invalid body", http.StatusBadRequest, encoder, w, err)
		return
	}

	mappedItems := menuItemBodyToMenuItem(menu.Id, true, body)

	err = menu.ReplaceItems(mappedItems)
	if err != nil {
		switch detectErrorType(err) {
		case errorTypeReference:
			sendError("One of the galleries, artists, modern pages, classic pages or blog categories does not exist", http.StatusConflict, encoder, w, err)
			return
		default:
			sendError("Failed to replace menu items", http.StatusInternalServerError, encoder, w, err)
			return
		}
	}

	w.WriteHeader(http.StatusNoContent)
}
