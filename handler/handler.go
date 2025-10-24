package handler

import (
	"event/repo"
	"event/utils"
	"net/http"
)

type Handler struct {
	server *repo.Server
	utils  utils.Utils
}

func New(server *repo.Server, utils utils.Utils) *Handler {
	return &Handler{
		server: server,
		utils:  utils,
	}
}

func (h *Handler) CreateEventHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		h.utils.WriteError(w, err, http.StatusBadRequest)
		return
	}

	userId, err := h.utils.ParseUserID(r.FormValue("user_id"))
	if err != nil {
		h.utils.WriteError(w, err, http.StatusBadRequest)
		return
	}

	date, err := h.utils.ParseDate(r.FormValue("date"))
	if err != nil {
		h.utils.WriteError(w, err, http.StatusBadRequest)
		return
	}

	title := r.FormValue("title")

	if err := h.server.CreateEvent(r.Context(), userId, date, title); err != nil {
		h.utils.WriteError(w, err, http.StatusInternalServerError)
		return
	}

	h.utils.WriteJSON(w, "Событие добавлено")
}

func (h *Handler) GetEventsForDayHandler(w http.ResponseWriter, r *http.Request) {
	userId, err := h.utils.ParseUserID(r.URL.Query().Get("user_id"))
	if err != nil {
		h.utils.WriteError(w, err, http.StatusBadRequest)
		return
	}

	date, err := h.utils.ParseDate(r.URL.Query().Get("date"))
	if err != nil {
		h.utils.WriteError(w, err, http.StatusBadRequest)
		return
	}

	res, err := h.server.GetEventsForDay(r.Context(), userId, date)
	if err != nil {
		h.utils.WriteError(w, err, http.StatusInternalServerError)
		return
	}

	h.utils.WriteJSON(w, res)
}

func (h *Handler) GetEventsForWeekHandler(w http.ResponseWriter, r *http.Request) {
	userId, err := h.utils.ParseUserID(r.URL.Query().Get("user_id"))
	if err != nil {
		h.utils.WriteError(w, err, http.StatusBadRequest)
		return
	}

	date, err := h.utils.ParseDate(r.URL.Query().Get("date"))
	if err != nil {
		h.utils.WriteError(w, err, http.StatusBadRequest)
		return
	}

	endDate := date.AddDate(0, 0, 7)
	res, err := h.server.GetEventsForWeek(r.Context(), userId, date, endDate)
	if err != nil {
		h.utils.WriteError(w, err, http.StatusInternalServerError)
		return
	}

	h.utils.WriteJSON(w, res)
}

func (h *Handler) GetEventsForMonthHandler(w http.ResponseWriter, r *http.Request) {
	userId, err := h.utils.ParseUserID(r.URL.Query().Get("user_id"))
	if err != nil {
		h.utils.WriteError(w, err, http.StatusBadRequest)
		return
	}

	date, err := h.utils.ParseDate(r.URL.Query().Get("date"))
	if err != nil {
		h.utils.WriteError(w, err, http.StatusBadRequest)
		return
	}

	endDate := date.AddDate(0, 1, 0)
	res, err := h.server.GetEventsForWeek(r.Context(), userId, date, endDate)
	if err != nil {
		h.utils.WriteError(w, err, http.StatusInternalServerError)
		return
	}

	h.utils.WriteJSON(w, res)
}

func (h *Handler) RemoveEventHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		h.utils.WriteError(w, err, http.StatusBadRequest)
		return
	}

	date, err := h.utils.ParseDate(r.FormValue("date"))
	if err != nil {
		h.utils.WriteError(w, err, http.StatusBadRequest)
		return
	}

	id, err := h.utils.ParseUserID(r.FormValue("id"))
	if err != nil {
		h.utils.WriteError(w, err, http.StatusBadRequest)
		return
	}

	if err := h.server.DeleteEvent(r.Context(), date, id); err != nil {
		h.utils.WriteError(w, err, http.StatusInternalServerError)
		return
	}

	h.utils.WriteJSON(w, "Событие удалено")
}

func (h *Handler) UpdateEventHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		h.utils.WriteError(w, err, http.StatusBadRequest)
		return
	}

	id, err := h.utils.ParseUserID(r.FormValue("id"))
	if err != nil {
		h.utils.WriteError(w, err, http.StatusBadRequest)
		return
	}

	userId, err := h.utils.ParseUserID(r.FormValue("user_id"))
	if err != nil {
		h.utils.WriteError(w, err, http.StatusBadRequest)
		return
	}

	title := r.FormValue("title")

	res, err := h.server.UpdateEvent(r.Context(), id, userId, title)
	if err != nil {
		h.utils.WriteError(w, err, http.StatusInternalServerError)
		return
	}

	h.utils.WriteJSON(w, res)
}
