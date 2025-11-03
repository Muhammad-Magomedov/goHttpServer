package handler

import (
	"event/repo"
	"event/utils"
	"net/http"
	"strconv"
	"time"
)

const dateFormat = "2006-01-02"

type Handler struct {
	server *repo.Server
}

func New(server *repo.Server) *Handler {
	return &Handler{
		server: server,
	}
}

func (h *Handler) CreateEventHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		utils.WriteError(w, err, http.StatusBadRequest)
		return
	}

	userID, err := strconv.Atoi(r.FormValue("user_id"))
	if err != nil {
		utils.WriteError(w, err, http.StatusBadRequest)
		return
	}

	date, err := time.Parse(dateFormat, r.FormValue("date"))
	if err != nil {
		utils.WriteError(w, err, http.StatusBadRequest)
		return
	}

	title := r.FormValue("title")

	event := repo.CreateEvent{
		UserID: userID,
		Date:   date,
		Title:  title,
	}

	if err := h.server.CreateEvent(r.Context(), event); err != nil {
		utils.WriteError(w, err, http.StatusInternalServerError)
		return
	}

	utils.WriteJSON(w, "Событие добавлено")
}

func (h *Handler) GetEventsForDayHandler(w http.ResponseWriter, r *http.Request) {
	userId, err := strconv.Atoi(r.URL.Query().Get("user_id"))
	if err != nil {
		utils.WriteError(w, err, http.StatusBadRequest)
		return
	}

	date, err := time.Parse(dateFormat, r.URL.Query().Get("date"))
	if err != nil {
		utils.WriteError(w, err, http.StatusBadRequest)
		return
	}

	res, err := h.server.GetEventsForDay(r.Context(), userId, date)
	if err != nil {
		utils.WriteError(w, err, http.StatusInternalServerError)
		return
	}

	utils.WriteJSON(w, res)
}

func (h *Handler) GetEventsForWeekHandler(w http.ResponseWriter, r *http.Request) {
	userId, err := strconv.Atoi(r.URL.Query().Get("user_id"))
	if err != nil {
		utils.WriteError(w, err, http.StatusBadRequest)
		return
	}

	date, err := time.Parse(dateFormat, r.URL.Query().Get("date"))
	if err != nil {
		utils.WriteError(w, err, http.StatusBadRequest)
		return
	}

	endDate := date.AddDate(0, 0, 7)
	res, err := h.server.GetEventsForDates(r.Context(), userId, date, endDate)
	if err != nil {
		utils.WriteError(w, err, http.StatusInternalServerError)
		return
	}

	utils.WriteJSON(w, res)
}

func (h *Handler) GetEventsForMonthHandler(w http.ResponseWriter, r *http.Request) {
	userId, err := strconv.Atoi(r.URL.Query().Get("user_id"))
	if err != nil {
		utils.WriteError(w, err, http.StatusBadRequest)
		return
	}

	date, err := time.Parse(dateFormat, r.URL.Query().Get("date"))
	if err != nil {
		utils.WriteError(w, err, http.StatusBadRequest)
		return
	}

	endDate := date.AddDate(0, 1, 0)
	res, err := h.server.GetEventsForDates(r.Context(), userId, date, endDate)
	if err != nil {
		utils.WriteError(w, err, http.StatusInternalServerError)
		return
	}

	utils.WriteJSON(w, res)
}

func (h *Handler) RemoveEventHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		utils.WriteError(w, err, http.StatusBadRequest)
		return
	}

	date, err := time.Parse(dateFormat, r.FormValue("date"))
	if err != nil {
		utils.WriteError(w, err, http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(r.FormValue("id"))
	if err != nil {
		utils.WriteError(w, err, http.StatusBadRequest)
		return
	}

	if err := h.server.DeleteEvent(r.Context(), date, id); err != nil {
		utils.WriteError(w, err, http.StatusInternalServerError)
		return
	}

	utils.WriteJSON(w, "Событие удалено")
}

func (h *Handler) UpdateEventHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		utils.WriteError(w, err, http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(r.FormValue("id"))
	if err != nil {
		utils.WriteError(w, err, http.StatusBadRequest)
		return
	}

	title := r.FormValue("title")

	res, err := h.server.UpdateEvent(r.Context(), id, title)
	if err != nil {
		utils.WriteError(w, err, http.StatusInternalServerError)
		return
	}

	utils.WriteJSON(w, res)
}
