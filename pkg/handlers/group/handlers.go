package group

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/wanderfusion/nomadcore/pkg/errors"
	"github.com/wanderfusion/nomadcore/pkg/handlers"
	"github.com/wanderfusion/nomadcore/pkg/services/group"

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

func getUserIDFromContext(r *http.Request) (uuid.UUID, bool) {
	userID, ok := r.Context().Value(handlers.UserIdContextKey).(uuid.UUID)
	return userID, ok
}

func handleError(w http.ResponseWriter, r *http.Request, err errors.Error) {
	handlers.RespondWithError(w, r, err)
}

func (h *Handlers) CreateGroup(w http.ResponseWriter, r *http.Request) {
	userID, ok := getUserIDFromContext(r)
	if !ok {
		handleError(w, r, ErrInternalServerError)
		return
	}
	var req CreateGroupReq
	if err := handlers.FromRequest(r, &req); err != nil {
		handleError(w, r, ErrBadRequest.Wrap(err))
		return
	}

	msg, err := h.Service.CreateGroup(userID, req.Name, req.Description)
	if err != nil {
		handleError(w, r, ErrInternalServerError.Wrap(err))
		return
	}

	handlers.RespondWithData(w, r, msg)
}

func (h *Handlers) GetGroups(w http.ResponseWriter, r *http.Request) {
	userID, ok := getUserIDFromContext(r)
	if !ok {
		handleError(w, r, ErrInvalidContext)
		return
	}

	groups, err := h.Service.GetGroups(userID)
	if err != nil {
		handleError(w, r, ErrInternalServerError.Wrap(err))
		return
	}

	groupIDs := make([]uuid.UUID, 0)
	for _, calendar := range groups {
		groupIDs = append(groupIDs, calendar.ID)
	}

	dates, err := h.Service.GetDates(groupIDs)
	if err != nil {
		handleError(w, r, ErrInternalServerError.Wrap(err))
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
	userID, ok := getUserIDFromContext(r)
	if !ok {
		handleError(w, r, ErrInvalidContext)
		return
	}
	var req AddDatesToGroupReq
	if err := handlers.FromRequest(r, &req); err != nil {
		handleError(w, r, ErrBadRequest.Wrap(err))
		return
	}

	dates := group.Dates{
		From: req.Dates.From,
		To:   req.Dates.To,
	}
	msg, err := h.Service.AddDatesToGroup(userID, req.GroupID, dates)
	if err != nil {
		handleError(w, r, ErrInternalServerError.Wrap(err))
		return
	}

	handlers.RespondWithData(w, r, msg)
}

func (h *Handlers) AddUsersToGroup(w http.ResponseWriter, r *http.Request) {
	// todo: check if user is in group
	// userID, ok := getUserIDFromContext(r)
	// if !ok {
	// 	handleError(w, r, ErrContextInvalid, http.StatusInternalServerError)
	// 	return
	// }
	var req AddUsersToGroupReq
	if err := handlers.FromRequest(r, &req); err != nil {
		handleError(w, r, ErrBadRequest.Wrap(err))
		return
	}

	msg, err := h.Service.AddUsersToGroup(req.Usernames, req.GroupID)
	if err != nil {
		handleError(w, r, ErrInternalServerError.Wrap(err))
		return
	}

	handlers.RespondWithData(w, r, msg)
}

func (h *Handlers) GetGroupDetails(w http.ResponseWriter, r *http.Request) {
	// Extract Group ID from URL using Chi router
	groupIDStr := chi.URLParam(r, "groupID")
	groupID, err := uuid.Parse(groupIDStr)
	if err != nil {
		handleError(w, r, ErrBadRequest.Wrap(err))
		return
	}

	userID, ok := getUserIDFromContext(r)
	if !ok {
		handleError(w, r, ErrInvalidContext)
		return
	}

	groupDetails, groupDates, groupUsers, svcErr := h.Service.GetGroupDetails(userID, groupID)
	if svcErr != nil {
		handleError(w, r, ErrInternalServerError.Wrap(err))
		return
	}

	// Convert to DTOs
	groupDTO := GroupDTO{
		ID:          groupDetails.ID,
		CreatedAt:   groupDetails.CreatedAt,
		Name:        groupDetails.Name,
		Description: groupDetails.Description,
	}

	var dateDTOs []DateDTO
	for _, date := range groupDates {
		dateDTO := DateDTO{
			ID:   date.ID,
			From: date.FromDate,
			To:   date.ToDate,
		}
		dateDTOs = append(dateDTOs, dateDTO)
	}

	var userDTOs []UserDTO
	for _, user := range groupUsers {
		userDTO := UserDTO{
			ID: user.UserID,
		}
		userDTOs = append(userDTOs, userDTO)
	}

	// Prepare response
	response := GetGroupDetailsRes{
		Group:     groupDTO,
		GroupDate: dateDTOs,
		GroupUser: userDTOs,
	}

	handlers.RespondWithData(w, r, response)
}
