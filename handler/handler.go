package handler

import (
	"encoding/json"
	"event/repo"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"
)

const dateFormat = "2006-01-02"

type UpdateEventReq struct {
	UserID int
	Date   time.Time
	Title  string
}

type Handler struct {
	server *repo.Server
}

func New(server *repo.Server) *Handler {
	return &Handler{
		server: server,
	}
}

func (h *Handler) CreateEventHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Unable to parse form", http.StatusBadRequest)
		return
	}

	user_id := r.FormValue("user_id")
	dateStr := r.FormValue("date")
	title := r.FormValue("title")

	date, err := time.Parse(dateFormat, dateStr)
	if err != nil {
		fmt.Println("Error while parsing time: ", err)
		return
	}

	userId, err := strconv.Atoi(user_id)
	if err != nil {
		log.Println(err)
	}

	err = h.server.CreateEvent(r.Context(), userId, date, title)
	if err != nil {
		http.Error(w, "Event not found", http.StatusNotFound)
		log.Println(err)
		return
	}

	str := fmt.Sprintf(`"result": {
		date: %s
		}`, date)

	response, err := json.Marshal(str)
	if err != nil {
		fmt.Println("Error while marshaling response CreateEventHandler: ", err)
		return
	}
	w.Write(response)
}

func (h *Handler) GetEventsForDayHandler(w http.ResponseWriter, r *http.Request) {
	dateStr := r.URL.Query().Get("date")
	user_id := r.URL.Query().Get("user_id")

	userId, err := strconv.Atoi(user_id)
	if err != nil {
		fmt.Println("getEventForDayHandler atoi: ", err)
		return
	}

	date, err := time.Parse(dateFormat, dateStr)
	if err != nil {
		fmt.Println("Error while parsing time: ", err)
		return
	}

	res, err := h.server.GetEventsForDay(r.Context(), userId, date)
	if err != nil {
		fmt.Println("getEventForDayHandler getEventForDay: ", err)
		errMsg, err := json.Marshal(fmt.Sprintf(`"error": "%s"`, err))
		if err != nil {
			fmt.Println("Error while marshla error")
		}
		w.Write(errMsg)
		return
	}

	jsonBytes, err := json.MarshalIndent(res, "", "\t") // Use "\t" for tabs
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonBytes)
}

func (h *Handler) GetEventsForWeekHandler(w http.ResponseWriter, r *http.Request) {
	dateStr := r.URL.Query().Get("date")
	user_id := r.URL.Query().Get("user_id")

	userId, err := strconv.Atoi(user_id)
	if err != nil {
		fmt.Println("getEventForDayHandler atoi: ", err)
		return
	}

	date, err := time.Parse(dateFormat, dateStr)
	if err != nil {
		fmt.Println("Error while parsing time: ", err)
		return
	}

	endDate := date.AddDate(0, 0, 7)

	res, err := h.server.GetEventsForWeek(r.Context(), userId, date, endDate)
	if err != nil {
		fmt.Println("getEventForDayHandler getEventForDay: ", err)
		errMsg, err := json.Marshal(fmt.Sprintf(`"error": "%s"`, err))
		if err != nil {
			fmt.Println("Error while marshla error")
		}
		w.Write(errMsg)
		return
	}

	jsonBytes, err := json.MarshalIndent(res, "", "\t") // Use "\t" for tabs
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonBytes)
}

func (h *Handler) GetEventsForMonthHandler(w http.ResponseWriter, r *http.Request) {
	dateStr := r.URL.Query().Get("date")
	user_id := r.URL.Query().Get("user_id")

	userId, err := strconv.Atoi(user_id)
	if err != nil {
		fmt.Println("getEventForDayHandler atoi: ", err)
		return
	}

	date, err := time.Parse(dateFormat, dateStr)
	if err != nil {
		fmt.Println("Error while parsing time: ", err)
		return
	}

	endDate := date.AddDate(0, 1, 0)

	res, err := h.server.GetEventsForWeek(r.Context(), userId, date, endDate)
	if err != nil {
		fmt.Println("getEventForDayHandler getEventForDay: ", err)
		errMsg, err := json.Marshal(fmt.Sprintf(`"error": "%s"`, err))
		if err != nil {
			fmt.Println("Error while marshla error")
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(errMsg)
		return
	}

	jsonBytes, err := json.MarshalIndent(res, "", "\t")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonBytes)
}

func (h *Handler) RemoveEventHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Unable to parse form", http.StatusBadRequest)
		return
	}

	dateStr := r.FormValue("date")
	idStr := r.FormValue("id")

	date, err := time.Parse(dateFormat, dateStr)
	if err != nil {
		fmt.Println("Error while parsing time: ", err)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Println(err)
	}

	err = h.server.DeleteEvent(r.Context(), date, id)
	if err != nil {
		fmt.Println("error while deleting event", err)

		w.Header().Set("Content-Type", "application/json")
		jsonBytes, err := json.MarshalIndent(err, "", "\t")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonBytes)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	jsonBytes, err := json.MarshalIndent("Событие удалено", "", "\t")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(jsonBytes)
}

func (h *Handler) UpdateEventHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Unable to parse form", http.StatusBadRequest)
		return
	}

	idStr := r.FormValue("id")
	dateStr := r.FormValue("date")
	title := r.FormValue("title")

	date, err := time.Parse(dateFormat, dateStr)
	if err != nil {
		fmt.Println("Error while parsing time: ", err)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Println(err)
	}

	res, err := h.server.UpdateEvent(r.Context(), id, date, title)
	if err != nil {
		fmt.Println("Error while updating event")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	jsonBytes, err := json.MarshalIndent(res, "", "\t")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(jsonBytes)
}
