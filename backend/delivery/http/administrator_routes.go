package http

import (
	"strconv"
	"time"

	"trustdonations.org/m/v2/authorization/jwt"
	"trustdonations.org/m/v2/domain/administrator"
	"trustdonations.org/m/v2/domain/church"
	"trustdonations.org/m/v2/domain/donator"
	"trustdonations.org/m/v2/domain/transaction"
	"trustdonations.org/m/v2/paymentgateway/stripe"
	"trustdonations.org/m/v2/user"
	"trustdonations.org/m/v2/validation"

	"log"
	"net/http"
	"os"

	"trustdonations.org/m/v2/delivery/middleware"

	"github.com/labstack/echo"
)

// AdministratorHandler represents an HTTP API handler for admins.
type AdministratorHandler struct {
	*echo.Echo

	Usecase administrator.Usecase

	Authorization jwt.Services

	PaymentGateway stripe.Services

	Logger *log.Logger
}

// NewAdministratorHandler returns a new instance of AdministratorHandler.
func NewAdministratorHandler() *AdministratorHandler {
	h := &AdministratorHandler{
		Echo:   echo.New(),
		Logger: log.New(os.Stderr, "", log.LstdFlags),
	}

	public := h.Group("/api")
	public.POST("/administrator", h.handleRegisterAdministrator)
	public.POST("/administrator/login", h.handleLoginAdministrator)
	public.POST("/administrator/authorization", h.handleAuthorizationAdministrator)

	restricted := h.Group("/api")
	restricted.Use(middleware.JwtMiddleware)
	restricted.GET("/administrator", h.handleReadAdministrator)
	restricted.GET("/administrator/administrators/:churchID", h.handleReadAdministrators)
	restricted.PATCH("/administrator", h.handlePatchAdministrator)
	restricted.PATCH("/administrator/administrators/:churchID/:administratorUUID", h.handlePatchAdministratorAccess)
	// restricted.DELETE("/administrator", h.handleDeleteAdministrator)

	restricted.POST("/administrator/church", h.handleCreateChurch)
	restricted.POST("/administrator/church/add", h.handleAddChurch)
	restricted.PATCH("/administrator/church", h.handlePatchChurch)
	restricted.DELETE("/administrator/church/:churchID", h.handleDeleteChurch)
	restricted.GET("/administrator/churches", h.handleReadChurches)

	restricted.POST("/administrator/donator", h.handleCreateDonator)
	restricted.PATCH("/administrator/donator", h.handlePatchDonator)
	restricted.DELETE("/administrator/donator/:churchID/:donatorID", h.handleDeleteDonator)
	restricted.GET("/administrator/donators/:churchID", h.handleReadDonators)

	restricted.POST("/administrator/donation", h.handleCreateDonation)
	restricted.GET("/administrator/donations/:churchID", h.handleReadDonations)
	restricted.PATCH("/administrator/donation", h.handlePatchDonation)
	restricted.GET("/administrator/donations/report/:churchID/:lowerRange/:upperRange/:timePeriod/:multiplier", h.handleGetChurchDonationReport)
	restricted.GET("/administrator/donations/statement/:churchID/:lowerRange/:upperRange", h.handleGetDonationStatements)
	// restricted.DELETE("/administrator/donation", h.handleDeleteDonation)

	return h
}

func (h *AdministratorHandler) handleRegisterAdministrator(ctx echo.Context) error {
	var req administratorRequest

	// Decode the request.
	if err := ctx.Bind(&req); err != nil || req.Administrator == nil {
		return ResponseError(ctx.Response().Writer, ErrInvalidJSON, http.StatusBadRequest, h.Logger, "Administrator")
	}

	a := req.Administrator
	switch err := h.Usecase.Register(a); err {
	case nil:
		tokenPair, err := middleware.GenerateTokenPair(*a.UUID, middleware.AccestokenLimit, middleware.RefreshtokenLimit)
		if err != nil {
			return ResponseError(ctx.Response().Writer, err, http.StatusInternalServerError, h.Logger, "Administrator")
		}
		// Do not send password to avoid exploits.
		a.Password = nil
		encodeJSON(ctx.Response().Writer, &administratorResponse{Administrator: a, Token: tokenPair}, h.Logger, "Administrator")
	case administrator.ErrAdministratorExists:
		return ResponseError(ctx.Response().Writer, err, http.StatusConflict, h.Logger, "Administrator")
	case administrator.ErrAdministratorIncorrectCredentials:
		return ResponseError(ctx.Response().Writer, err, http.StatusUnauthorized, h.Logger, "Administrator")
	case administrator.ErrAdministratorIncompleteDetails, validation.ErrValidationUserEmail, validation.ErrValidationPassword:
		return ResponseError(ctx.Response().Writer, err, http.StatusBadRequest, h.Logger, "Administrator")
	default:
		return ResponseError(ctx.Response().Writer, err, http.StatusInternalServerError, h.Logger, "Administrator")
	}

	return nil
}

func (h *AdministratorHandler) handleLoginAdministrator(ctx echo.Context) error {
	var req administratorRequest

	// Decode the request.
	if err := ctx.Bind(&req); err != nil || req.Administrator == nil {
		return ResponseError(ctx.Response().Writer, ErrInvalidJSON, http.StatusBadRequest, h.Logger, "Administrator")
	}
	a := req.Administrator

	switch exists, err := h.Usecase.Login(a); err {
	case nil:
		tokenPair, err := middleware.GenerateTokenPair(*exists.UUID, middleware.AccestokenLimit, middleware.RefreshtokenLimit)
		if err != nil {
			return ResponseError(ctx.Response().Writer, err, http.StatusInternalServerError, h.Logger, "Administrator")
		}
		// Do not send password to avoid exploits.
		exists.Password = nil
		encodeJSON(ctx.Response().Writer, &administratorResponse{Administrator: exists, Token: tokenPair}, h.Logger, "Administrator")
	case administrator.ErrAdministratorIncorrectCredentials:
		return ResponseError(ctx.Response().Writer, err, http.StatusUnauthorized, h.Logger, "Administrator")
	case user.ErrUserNotFound, administrator.ErrAdministratorNotFound:
		return ResponseError(ctx.Response().Writer, administrator.ErrAdministratorNotFound, http.StatusNotFound, h.Logger, "Administrator")
	}
	return nil
}

func (h *AdministratorHandler) handleAuthorizationAdministrator(ctx echo.Context) error {
	var req authorizationAdministratorRequest

	// Decode the request.
	if err := ctx.Bind(&req); err != nil || req.Token == nil {
		return ResponseError(ctx.Response().Writer, ErrInvalidJSON, http.StatusBadRequest, h.Logger, "Administrator")
	}

	token := req.Token
	switch newToken, err := h.Authorization.Administrator.Authorize(token); err {
	case nil:
		encodeJSON(ctx.Response().Writer, &authorizationAdministratorResponse{Token: newToken}, h.Logger, "Administrator")
	case jwt.ErrJWTAuth:
		return ResponseError(ctx.Response().Writer, err, http.StatusUnauthorized, h.Logger, "Administrator")
	default:
		return ResponseError(ctx.Response().Writer, err, http.StatusInternalServerError, h.Logger, "Administrator")
	}

	return nil
}

func (h *AdministratorHandler) handleReadAdministrator(ctx echo.Context) error {
	//extract administrator uuid from token stored by JwtMiddleware handler func
	token := ctx.Get("user")
	uuid, err := h.Authorization.Administrator.ExtractUUIDFromToken(token)
	if err != nil {
		return ResponseError(ctx.Response().Writer, err, http.StatusInternalServerError, h.Logger, "Administrator")
	}

	// Find administrator.
	switch exists, err := h.Usecase.Services.Administrator.Read(&administrator.Administrator{UUID: uuid}); exists {
	case nil:
		if err != nil {
			return ResponseError(ctx.Response().Writer, err, http.StatusInternalServerError, h.Logger, "Administrator")
		}
		return NotFound(ctx.Response().Writer)
	default:
		byEmail := false
		uExists, err := h.Usecase.Services.User.Read(&user.User{UUID: exists.UUID}, byEmail)
		if err != nil {
			return ResponseError(ctx.Response().Writer, err, http.StatusInternalServerError, h.Logger, "Administrator")
		}
		// Do not send password to avoid exploits
		exists.Email = uExists.Email
		exists.Password = nil
		encodeJSON(ctx.Response().Writer, &administratorResponse{Administrator: exists}, h.Logger, "Administrator")
	}
	return nil
}

func (h *AdministratorHandler) handleReadAdministrators(ctx echo.Context) error {
	id, _ := strconv.ParseInt(ctx.Param("churchID"), 10, 64)

	//extract administrator uuid from token stored by JwtMiddleware handler func
	token := ctx.Get("user")
	uuid, err := h.Authorization.Administrator.ExtractUUIDFromToken(token)
	if err != nil {
		return ResponseError(ctx.Response().Writer, err, http.StatusInternalServerError, h.Logger, "Administrator")
	}

	a := &administrator.Administrator{UUID: uuid}
	c := &church.Church{ID: &id}

	switch churches, err := h.Usecase.ReadAdministrators(a, c); err {
	case nil:
		encodeJSON(ctx.Response().Writer, churches, h.Logger, "Administrator")
	case administrator.ErrAdministratorDoesNotBelongToChurch:
		return ResponseError(ctx.Response().Writer, err, http.StatusUnauthorized, h.Logger, "Administrator")
	default:
		return ResponseError(ctx.Response().Writer, err, http.StatusInternalServerError, h.Logger, "Administrator")
	}

	return nil
}

func (h *AdministratorHandler) handlePatchAdministrator(ctx echo.Context) error {
	var req administratorRequest

	// Decode the request.
	if err := ctx.Bind(&req); err != nil || req.Administrator == nil {
		return ResponseError(ctx.Response().Writer, ErrInvalidJSON, http.StatusBadRequest, h.Logger, "Administrator")
	}

	// extract administrator uuid from token stored by JwtMiddleware handler func
	token := ctx.Get("user")
	uuid, err := h.Authorization.Administrator.ExtractUUIDFromToken(token)
	if err != nil {
		return ResponseError(ctx.Response().Writer, err, http.StatusInternalServerError, h.Logger, "Administrator")
	}

	a := req.Administrator
	a.UUID = uuid

	// Update administrator.
	switch updatedAdministrator, err := h.Usecase.EditAdministrator(a); err {
	case nil:
		// Do not send password in response.
		a.Password = nil
		encodeJSON(ctx.Response().Writer, &administratorResponse{Administrator: updatedAdministrator}, h.Logger, "Administrator")
	case administrator.ErrAdministratorNotFound:
		ResponseError(ctx.Response().Writer, err, http.StatusNotFound, h.Logger, "Administrator")
	case validation.ErrValidationPassword, validation.ErrValidationUserEmail:
		return ResponseError(ctx.Response().Writer, err, http.StatusBadRequest, h.Logger, "Administrator")
	default:
		ResponseError(ctx.Response().Writer, err, http.StatusInternalServerError, h.Logger, "Administrator")
	}

	return nil

}
func (h *AdministratorHandler) handlePatchAdministratorAccess(ctx echo.Context) error {
	administratorToEditUUID := ctx.Param("administratorUUID")
	churchID, _ := strconv.ParseInt(ctx.Param("churchID"), 10, 64)

	var req churchAdministratorRequest

	// Decode the request.
	if err := ctx.Bind(&req); err != nil {
		return ResponseError(ctx.Response().Writer, ErrInvalidJSON, http.StatusBadRequest, h.Logger, "Administrator")
	}

	// extract administrator uuid from token stored by JwtMiddleware handler func
	token := ctx.Get("user")
	uuid, err := h.Authorization.Administrator.ExtractUUIDFromToken(token)
	if err != nil {
		return ResponseError(ctx.Response().Writer, err, http.StatusInternalServerError, h.Logger, "Administrator")
	}

	// Update administrator.
	switch updatedAdministrator, err := h.Usecase.EditAdministrators(
		&church.Church{ID: &churchID},
		&administrator.Administrator{UUID: uuid},
		&administrator.Administrator{UUID: &administratorToEditUUID},
		req.Administrator,
	); err {
	case nil:
		encodeJSON(ctx.Response().Writer, &churchAdministratorResponse{Administrator: updatedAdministrator, AdministratorUUID: &administratorToEditUUID}, h.Logger, "Administrator")
	case administrator.ErrAdministratorNotFound:
		ResponseError(ctx.Response().Writer, err, http.StatusNotFound, h.Logger, "Administrator")
	case validation.ErrValidationPassword, validation.ErrValidationUserEmail:
		return ResponseError(ctx.Response().Writer, err, http.StatusBadRequest, h.Logger, "Administrator")
	case administrator.ErrAdministratorNotAuthorized, administrator.ErrAdministratorFieldNotEditable:
		return ResponseError(ctx.Response().Writer, err, http.StatusUnauthorized, h.Logger, "Administrator")
	default:
		ResponseError(ctx.Response().Writer, err, http.StatusInternalServerError, h.Logger, "Administrator")
	}
	return nil
}

func (h *AdministratorHandler) handleCreateChurch(ctx echo.Context) error {
	var req churchRequest
	// Decode the request.
	if err := ctx.Bind(&req); err != nil || req.Church == nil {
		return ResponseError(ctx.Response().Writer, ErrInvalidJSON, http.StatusBadRequest, h.Logger, "Administrator")
	}

	//extract administrator uuid from token stored by JwtMiddleware handler func
	token := ctx.Get("user")
	uuid, err := h.Authorization.Administrator.ExtractUUIDFromToken(token)
	if err != nil {
		return ResponseError(ctx.Response().Writer, err, http.StatusInternalServerError, h.Logger, "Administrator")
	}

	a := &administrator.Administrator{UUID: uuid}
	c := req.Church
	switch newChurch, err := h.Usecase.CreateChurch(a, c); err {
	case nil:
		newChurch.Password = nil
		encodeJSON(ctx.Response().Writer, &churchResponse{Church: newChurch, ChurchID: newChurch.ID}, h.Logger, "Administrator")
	case church.ErrChurchExists:
		return ResponseError(ctx.Response().Writer, err, http.StatusConflict, h.Logger, "Administrator")
	case administrator.ErrAdministratorChurchIncompleteDetails, validation.ErrValidationUserEmail, validation.ErrValidationPassword:
		return ResponseError(ctx.Response().Writer, err, http.StatusBadRequest, h.Logger, "Administrator")
	default:
		return ResponseError(ctx.Response().Writer, err, http.StatusInternalServerError, h.Logger, "Administrator")
	}
	return nil
}

func (h *AdministratorHandler) handleAddChurch(ctx echo.Context) error {
	var req churchRequest
	// Decode the request.
	if err := ctx.Bind(&req); err != nil || req.Church == nil || req.Church.Email == nil {
		return ResponseError(ctx.Response().Writer, ErrInvalidJSON, http.StatusBadRequest, h.Logger, "Administrator")
	}
	//extract administrator uuid from token stored by JwtMiddleware handler func
	token := ctx.Get("user")
	uuid, err := h.Authorization.Administrator.ExtractUUIDFromToken(token)
	if err != nil {
		return ResponseError(ctx.Response().Writer, err, http.StatusInternalServerError, h.Logger, "Administrator")
	}

	a := &administrator.Administrator{UUID: uuid}
	a.UUID = uuid
	c := req.Church

	switch newChurch, err := h.Usecase.AddChurch(a, c); err {
	case nil:
		newChurch.Password = nil
		encodeJSON(ctx.Response().Writer, &churchResponse{Church: newChurch, ChurchID: newChurch.ID}, h.Logger, "Administrator")
	case administrator.ErrAdministratorChurchExists:
		return ResponseError(ctx.Response().Writer, err, http.StatusConflict, h.Logger, "Administrator")
	case administrator.ErrAdministratorChurchIncorrectCredentilas:
		return ResponseError(ctx.Response().Writer, err, http.StatusUnauthorized, h.Logger, "Administrator")
	default:
		return ResponseError(ctx.Response().Writer, err, http.StatusInternalServerError, h.Logger, "Administrator")
	}
	return nil
}
func (h *AdministratorHandler) handleReadChurches(ctx echo.Context) error {
	token := ctx.Get("user")
	uuid, err := h.Authorization.Administrator.ExtractUUIDFromToken(token)
	if err != nil {
		return ResponseError(ctx.Response().Writer, err, http.StatusInternalServerError, h.Logger, "Administrator")
	}

	a := &administrator.Administrator{UUID: uuid}

	switch churches, err := h.Usecase.ReadChurches(a); err {
	case nil:
		encodeJSON(ctx.Response().Writer, churches, h.Logger, "Administrator")
	case administrator.ErrAdministratorNoChurches:
		return ResponseError(ctx.Response().Writer, err, http.StatusNoContent, h.Logger, "Administrator")
	default:
		return ResponseError(ctx.Response().Writer, err, http.StatusInternalServerError, h.Logger, "Administrator")
	}

	return nil
}
func (h *AdministratorHandler) handlePatchChurch(ctx echo.Context) error {
	var req churchRequest
	// Decode the request.
	if err := ctx.Bind(&req); err != nil || req.Church == nil || req.Church.ID == nil {
		return ResponseError(ctx.Response().Writer, ErrInvalidJSON, http.StatusBadRequest, h.Logger, "Administrator")
	}
	//extract administrator uuid from token stored by JwtMiddleware handler func
	token := ctx.Get("user")
	uuid, err := h.Authorization.Administrator.ExtractUUIDFromToken(token)
	if err != nil {
		return ResponseError(ctx.Response().Writer, err, http.StatusInternalServerError, h.Logger, "Administrator")
	}

	a := &administrator.Administrator{UUID: uuid}
	c := req.Church

	switch editedChurch, err := h.Usecase.EditChurch(a, c, *c.ID); err {
	case nil:
		encodeJSON(ctx.Response().Writer, &churchResponse{Church: editedChurch, ChurchID: c.ID}, h.Logger, "Administrator")
	case church.ErrChurchFieldEditUnPriveledged, validation.ErrValidationPassword, validation.ErrValidationUserEmail:
		return ResponseError(ctx.Response().Writer, err, http.StatusBadRequest, h.Logger, "Administrator")
	default:
		return ResponseError(ctx.Response().Writer, err, http.StatusInternalServerError, h.Logger, "Administrator")
	}
	return nil
}
func (h *AdministratorHandler) handleDeleteChurch(ctx echo.Context) error {
	cID, _ := strconv.ParseInt(ctx.Param("churchID"), 10, 64)

	token := ctx.Get("user")
	uuid, err := h.Authorization.Administrator.ExtractUUIDFromToken(token)
	if err != nil {
		return ResponseError(ctx.Response().Writer, err, http.StatusInternalServerError, h.Logger, "Administrator")
	}

	a := &administrator.Administrator{UUID: uuid}

	// Delete administrator.
	switch _, err := h.Usecase.EditChurch(a, nil, cID); err {
	case nil:
		encodeJSON(ctx.Response().Writer, &churchResponse{ChurchID: &cID}, h.Logger, "Administrator")
	case administrator.ErrAdministratorNotFound, church.ErrChurchDonatorDoesNotExists, donator.ErrDonatorExists:
		ResponseError(ctx.Response().Writer, err, http.StatusNotFound, h.Logger, "Administrator")
	default:
		ResponseError(ctx.Response().Writer, err, http.StatusInternalServerError, h.Logger, "Administrator")
	}

	return nil
}

func (h *AdministratorHandler) handleCreateDonator(ctx echo.Context) error {
	var req donatorRequest
	// Decode the request.
	if err := ctx.Bind(&req); err != nil || req.Donator == nil || req.Church == nil {
		return ResponseError(ctx.Response().Writer, ErrInvalidJSON, http.StatusBadRequest, h.Logger, "Administrator")
	}
	//extract administrator uuid from token stored by JwtMiddleware handler func
	token := ctx.Get("user")
	uuid, err := h.Authorization.Administrator.ExtractUUIDFromToken(token)
	if err != nil {
		return ResponseError(ctx.Response().Writer, err, http.StatusInternalServerError, h.Logger, "Administrator")
	}

	a := &administrator.Administrator{UUID: uuid}
	c := req.Church
	d := req.Donator

	switch newDonator, err := h.Usecase.CreateDonator(a, c, d); err {
	case nil:
		encodeJSON(ctx.Response().Writer, &donatorResponse{Donator: newDonator, DonatorID: newDonator.ID}, h.Logger, "Administrator")
	case administrator.ErrAdministratorDonatorIncompleteDetails:
		return ResponseError(ctx.Response().Writer, err, http.StatusBadRequest, h.Logger, "Administrator")
	case church.ErrChurchDonatorExists:
		return ResponseError(ctx.Response().Writer, err, http.StatusConflict, h.Logger, "Administrator")
	default:
		return ResponseError(ctx.Response().Writer, err, http.StatusInternalServerError, h.Logger, "Administrator")
	}
	return nil
}

func (h *AdministratorHandler) handleReadDonators(ctx echo.Context) error {
	id, _ := strconv.ParseInt(ctx.Param("churchID"), 10, 64)

	//extract administrator uuid from token stored by JwtMiddleware handler func
	token := ctx.Get("user")
	uuid, err := h.Authorization.Administrator.ExtractUUIDFromToken(token)
	if err != nil {
		return ResponseError(ctx.Response().Writer, err, http.StatusInternalServerError, h.Logger, "Administrator")
	}

	a := &administrator.Administrator{UUID: uuid}
	c := &church.Church{ID: &id}

	switch donators, err := h.Usecase.ReadDonators(a, c); err {
	case nil:
		encodeJSON(ctx.Response().Writer, donators, h.Logger, "Administrator")
	case administrator.ErrAdministratorDoesNotBelongToChurch:
		return ResponseError(ctx.Response().Writer, err, http.StatusUnauthorized, h.Logger, "Administrator")
	default:
		return ResponseError(ctx.Response().Writer, err, http.StatusInternalServerError, h.Logger, "Administrator")
	}

	return nil
}

func (h *AdministratorHandler) handlePatchDonator(ctx echo.Context) error {
	var req donatorRequest

	// Decode the request.
	if err := ctx.Bind(&req); err != nil || req.Donator == nil || req.Church == nil || req.Donator.ID == nil || req.Church.ID == nil {
		return ResponseError(ctx.Response().Writer, ErrInvalidJSON, http.StatusBadRequest, h.Logger, "Administrator")
	}

	// extract administrator uuid from token stored by JwtMiddleware handler func
	token := ctx.Get("user")
	uuid, err := h.Authorization.Administrator.ExtractUUIDFromToken(token)
	if err != nil {
		return ResponseError(ctx.Response().Writer, err, http.StatusInternalServerError, h.Logger, "Administrator")
	}

	a := &administrator.Administrator{UUID: uuid}
	c := req.Church
	d := req.Donator
	// Update donator.
	switch updatedDonator, err := h.Usecase.EditDonator(a, c, d, *d.ID); err {
	case nil:
		// Do not send password in response.
		a.Password = nil
		encodeJSON(ctx.Response().Writer, &donatorResponse{Donator: updatedDonator, DonatorID: d.ID}, h.Logger, "Administrator")
	case administrator.ErrAdministratorNotFound, church.ErrChurchDonatorDoesNotExists, donator.ErrDonatorExists:
		ResponseError(ctx.Response().Writer, err, http.StatusNotFound, h.Logger, "Administrator")
	default:
		ResponseError(ctx.Response().Writer, err, http.StatusInternalServerError, h.Logger, "Administrator")
	}

	return nil

}

func (h *AdministratorHandler) handleDeleteDonator(ctx echo.Context) error {
	cID, _ := strconv.ParseInt(ctx.Param("churchID"), 10, 64)
	dID, _ := strconv.ParseInt(ctx.Param("donatorID"), 10, 64)

	token := ctx.Get("user")
	uuid, err := h.Authorization.Administrator.ExtractUUIDFromToken(token)
	if err != nil {
		return ResponseError(ctx.Response().Writer, err, http.StatusInternalServerError, h.Logger, "Administrator")
	}

	a := &administrator.Administrator{UUID: uuid}
	c := &church.Church{ID: &cID}

	// Delete administrator.
	switch _, err := h.Usecase.EditDonator(a, c, nil, dID); err {
	case nil:
		encodeJSON(ctx.Response().Writer, &donatorResponse{DonatorID: &dID}, h.Logger, "Administrator")
	case administrator.ErrAdministratorNotFound, church.ErrChurchDonatorDoesNotExists, donator.ErrDonatorExists:
		ResponseError(ctx.Response().Writer, err, http.StatusNotFound, h.Logger, "Administrator")
	default:
		ResponseError(ctx.Response().Writer, err, http.StatusInternalServerError, h.Logger, "Administrator")
	}

	return nil
}

func (h *AdministratorHandler) handleCreateDonation(ctx echo.Context) error {
	var req donationRequest
	// Decode the request.
	err := ctx.Bind(&req)

	if err != nil ||
		req.Donation.DonatorID == nil ||
		req.Donation.ChurchID == nil ||
		req.Donation == nil ||
		req.Donation.Amount == nil ||
		req.Donation.Date == nil ||
		(req.Donation.Date != nil && *req.Donation.Date == "") {
		return ResponseError(ctx.Response().Writer, ErrInvalidJSON, http.StatusBadRequest, h.Logger, "Administrator")
	}
	//extract administrator uuid from token stored by JwtMiddleware handler func
	token := ctx.Get("user")
	uuid, err := h.Authorization.Administrator.ExtractUUIDFromToken(token)
	if err != nil {
		return ResponseError(ctx.Response().Writer, err, http.StatusInternalServerError, h.Logger, "Administrator")
	}

	a := &administrator.Administrator{UUID: uuid}
	donation := req.Donation
	if _, err := time.Parse("2006-01-02T15:04:05.999Z", *donation.Date); err != nil {
		return ResponseError(ctx.Response().Writer, err, http.StatusBadRequest, h.Logger, "Administrator")
	}

	switch newDonation, err := h.Usecase.CreateDonation(a, donation); err {
	case nil:
		encodeJSON(ctx.Response().Writer, &donationResponse{Donation: newDonation}, h.Logger, "Administrator")
	case administrator.ErrAdministratorDonatorIncompleteDetails:
		return ResponseError(ctx.Response().Writer, err, http.StatusBadRequest, h.Logger, "Administrator")
	case church.ErrChurchDonatorExists:
		return ResponseError(ctx.Response().Writer, err, http.StatusConflict, h.Logger, "Administrator")
	default:
		return ResponseError(ctx.Response().Writer, err, http.StatusInternalServerError, h.Logger, "Administrator")
	}
	return nil
}

func (h *AdministratorHandler) handleReadDonations(ctx echo.Context) error {
	cID, _ := strconv.ParseInt(ctx.Param("churchID"), 10, 64)

	//extract administrator uuid from token stored by JwtMiddleware handler func
	token := ctx.Get("user")
	uuid, err := h.Authorization.Administrator.ExtractUUIDFromToken(token)
	if err != nil {
		return ResponseError(ctx.Response().Writer, err, http.StatusInternalServerError, h.Logger, "Administrator")
	}

	a := &administrator.Administrator{UUID: uuid}
	c := &church.Church{ID: &cID}

	switch donators, err := h.Usecase.ReadDonations(a, c); err {
	case nil:
		encodeJSON(ctx.Response().Writer, donators, h.Logger, "Administrator")
	case administrator.ErrAdministratorDoesNotBelongToChurch:
		return ResponseError(ctx.Response().Writer, err, http.StatusUnauthorized, h.Logger, "Administrator")
	case church.ErrChurchDonatorDoesNotExists:
		return ResponseError(ctx.Response().Writer, err, http.StatusBadRequest, h.Logger, "Administrator")
	default:
		return ResponseError(ctx.Response().Writer, err, http.StatusInternalServerError, h.Logger, "Administrator")
	}

	return nil
}

func (h *AdministratorHandler) handlePatchDonation(ctx echo.Context) error {
	var req donationRequest
	// Decode the request.
	err := ctx.Bind(&req)

	if err != nil ||
		req.Donation.DonatorID == nil ||
		req.Donation.ChurchID == nil ||
		req.Donation == nil ||
		req.Donation.Amount == nil ||
		req.Donation.Date == nil ||
		(req.Donation.Date != nil && *req.Donation.Date == "") {
		return ResponseError(ctx.Response().Writer, ErrInvalidJSON, http.StatusBadRequest, h.Logger, "Administrator")
	}
	//extract administrator uuid from token stored by JwtMiddleware handler func
	token := ctx.Get("user")
	uuid, err := h.Authorization.Administrator.ExtractUUIDFromToken(token)
	if err != nil {
		return ResponseError(ctx.Response().Writer, err, http.StatusInternalServerError, h.Logger, "Administrator")
	}

	a := &administrator.Administrator{UUID: uuid}
	donation := req.Donation
	if _, err := time.Parse("2006-01-02T15:04:05.999Z", *donation.Date); err != nil {
		return ResponseError(ctx.Response().Writer, err, http.StatusBadRequest, h.Logger, "Administrator")
	}

	switch newDonation, err := h.Usecase.EditDonation(a, donation); err {
	case nil:
		encodeJSON(ctx.Response().Writer, &donationResponse{Donation: newDonation}, h.Logger, "Administrator")
	case administrator.ErrAdministratorDonatorIncompleteDetails:
		return ResponseError(ctx.Response().Writer, err, http.StatusBadRequest, h.Logger, "Administrator")
	case church.ErrChurchDonatorExists:
		return ResponseError(ctx.Response().Writer, err, http.StatusConflict, h.Logger, "Administrator")
	default:
		return ResponseError(ctx.Response().Writer, err, http.StatusInternalServerError, h.Logger, "Administrator")
	}

	return nil
}

func (h *AdministratorHandler) handleGetChurchDonationReport(ctx echo.Context) error {
	q := ctx.Request().URL.Query()
	categories := q["categories"]

	//extract administrator uuid from token stored by JwtMiddleware handler func
	token := ctx.Get("user")
	uuid, err := h.Authorization.Administrator.ExtractUUIDFromToken(token)
	if err != nil {
		return ResponseError(ctx.Response().Writer, err, http.StatusInternalServerError, h.Logger, "Administrator")
	}
	a := &administrator.Administrator{UUID: uuid}

	upperRange, err := time.Parse("2006-01-02T15:04:05.999Z", ctx.Param("upperRange"))
	if err != nil {
		return ResponseError(ctx.Response().Writer, err, http.StatusBadRequest, h.Logger, "Administrator")
	}
	lowerRange, err := time.Parse("2006-01-02T15:04:05.999Z", ctx.Param("lowerRange"))
	if err != nil {
		return ResponseError(ctx.Response().Writer, err, http.StatusBadRequest, h.Logger, "Administrator")

	}
	timePeriod, err := transaction.ParseTimePeriod(ctx.Param("timePeriod"))
	if err != nil {
		return ResponseError(ctx.Response().Writer, err, http.StatusInternalServerError, h.Logger, "Administrator")
	}
	multiplier, err := strconv.ParseInt(ctx.Param("multiplier"), 10, 64)
	if err != nil {
		return ResponseError(ctx.Response().Writer, err, http.StatusInternalServerError, h.Logger, "Administrator")
	}
	cID, err := strconv.ParseInt(ctx.Param("churchID"), 10, 64)
	if err != nil {
		return ResponseError(ctx.Response().Writer, err, http.StatusInternalServerError, h.Logger, "Administrator")
	}

	donationReport := transaction.DonationReport{
		ChurchID: &cID,
		SumFilter: &transaction.SumFilter{
			TimePeriod: &timePeriod,
			Multiplier: &multiplier,
		},
		TimeRange: &transaction.TimeRange{
			Upper: &upperRange,
			Lower: &lowerRange,
		},
		DonationCategories: categories,
	}

	switch donationReport, err := h.Usecase.ChurchDonationReport(a, &donationReport); err {
	case nil:
		encodeJSON(ctx.Response().Writer, donationReport, h.Logger, "Administrator")
	default:
		return ResponseError(ctx.Response().Writer, err, http.StatusInternalServerError, h.Logger, "Administrator")
	}

	return nil
}

func (h *AdministratorHandler) handleGetDonationStatements(ctx echo.Context) error {

	q := ctx.Request().URL.Query()
	categories := q["categories"]
	donatorIDs := make([]int64, len(q["donatorIDs"]))
	for i, donatorIDStr := range q["donatorIDs"] {
		donatorID, err := strconv.ParseInt(donatorIDStr, 10, 64)
		if err != nil {
			return ResponseError(ctx.Response().Writer, err, http.StatusInternalServerError, h.Logger, "Administrator")
		}
		donatorIDs[i] = donatorID
	}

	//extract administrator uuid from token stored by JwtMiddleware handler func
	token := ctx.Get("user")
	uuid, err := h.Authorization.Administrator.ExtractUUIDFromToken(token)
	if err != nil {
		return ResponseError(ctx.Response().Writer, err, http.StatusInternalServerError, h.Logger, "Administrator")
	}
	a := &administrator.Administrator{UUID: uuid}

	upperRange, err := time.Parse("2006-01-02T15:04:05.999Z", ctx.Param("upperRange"))
	if err != nil {
		return ResponseError(ctx.Response().Writer, err, http.StatusBadRequest, h.Logger, "Administrator")
	}
	lowerRange, err := time.Parse("2006-01-02T15:04:05.999Z", ctx.Param("lowerRange"))
	if err != nil {
		return ResponseError(ctx.Response().Writer, err, http.StatusBadRequest, h.Logger, "Administrator")

	}

	cID, err := strconv.ParseInt(ctx.Param("churchID"), 10, 64)
	if err != nil {
		return ResponseError(ctx.Response().Writer, err, http.StatusInternalServerError, h.Logger, "Administrator")
	}

	donationReport := transaction.DonationReport{
		DonatorIDs: donatorIDs,
		ChurchID:   &cID,
		TimeRange: &transaction.TimeRange{
			Upper: &upperRange,
			Lower: &lowerRange,
		},
		DonationCategories: categories,
	}

	switch donationReport, err := h.Usecase.DonationStatementReport(a, &donationReport); err {
	case nil:
		encodeJSON(ctx.Response().Writer, donationReport, h.Logger, "Administrator")
	default:
		return ResponseError(ctx.Response().Writer, err, http.StatusInternalServerError, h.Logger, "Administrator")
	}

	return nil
}
