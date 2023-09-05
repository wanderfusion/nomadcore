package group

import (
	"net/http"

	"github.com/akxcix/nomadcore/pkg/handlers"
	"github.com/akxcix/nomadcore/pkg/services/group"

	"github.com/google/uuid"
)

type Handlers struct {
	Service *group.Service
}

func New(s *group.Service) *Handlers {
	h := Handlers{
		Service: s,
	}

	return &h
}

func (h *Handlers) CreateGroup(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(handlers.UserIdContextKey).(uuid.UUID)
	if !ok {
		handlers.RespondWithError(w, r, ErrContextInvalid, http.StatusInternalServerError)
		return
	}
	var req CreateGroupReq
	if err := handlers.FromRequest(r, &req); err != nil {
		handlers.RespondWithError(w, r, err, http.StatusBadRequest)
		return
	}

	msg, err := h.Service.CreateGroup(userID, req.Name, req.Description)
	if err != nil {
		handlers.RespondWithError(w, r, err, http.StatusInternalServerError)
		return
	}

	handlers.RespondWithData(w, r, msg)
}

func (h *Handlers) GetGroups(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(handlers.UserIdContextKey).(uuid.UUID)
	if !ok {
		handlers.RespondWithError(w, r, ErrContextInvalid, http.StatusInternalServerError)
		return
	}

	groups, err := h.Service.GetGroups(userID)
	if err != nil {
		handlers.RespondWithError(w, r, err, http.StatusInternalServerError)
		return
	}

	groupIDs := make([]uuid.UUID, 0)
	for _, calendar := range groups {
		groupIDs = append(groupIDs, calendar.ID)
	}

	dates, err := h.Service.GetDates(groupIDs)
	if err != nil {
		handlers.RespondWithError(w, r, err, http.StatusInternalServerError)
		return
	}

	dateMap := make(map[uuid.UUID][]DateDTO)
	for _, date := range dates {
		if _, exists := dateMap[date.GroupID]; !exists {
			dateMap[date.GroupID] = make([]DateDTO, 0)
		}

		dateSlice := dateMap[date.GroupID]
		dateSlice = append(dateSlice, DateDTO{
			ID:   date.ID,
			From: date.FromDate,
			To:   date.ToDate,
		})
		dateMap[date.GroupID] = dateSlice
	}

	groupDTOs := make([]GroupDTO, 0)
	for _, cal := range groups {
		dateDtos := dateMap[cal.ID]

		calDto := GroupDTO{
			ID:          cal.ID,
			Name:        cal.Name,
			Description: cal.Description,
			Dates:       dateDtos,
		}

		groupDTOs = append(groupDTOs, calDto)
	}

	res := GetGroupsRes{
		Groups: groupDTOs,
	}

	handlers.RespondWithData(w, r, res)
}

func (h *Handlers) AddDatesToGroup(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(handlers.UserIdContextKey).(uuid.UUID)
	if !ok {
		handlers.RespondWithError(w, r, ErrContextInvalid, http.StatusInternalServerError)
		return
	}
	var req AddDatesToGroupReq
	if err := handlers.FromRequest(r, &req); err != nil {
		handlers.RespondWithError(w, r, err, http.StatusBadRequest)
		return
	}

	dates := group.Dates{
		From: req.Dates.From,
		To:   req.Dates.To,
	}
	msg, err := h.Service.AddDatesToGroup(userID, req.GroupID, dates)
	if err != nil {
		handlers.RespondWithError(w, r, err, http.StatusInternalServerError)
		return
	}

	handlers.RespondWithData(w, r, msg)
}
