package controller

import (
	"emersonargueta/m/v1/modules/communitygoaltracker/domain/achiever"
	"emersonargueta/m/v1/modules/communitygoaltracker/dto"
	"emersonargueta/m/v1/modules/communitygoaltracker/infrastructure/persistence"
	"emersonargueta/m/v1/modules/communitygoaltracker/usecase"
	"emersonargueta/m/v1/shared/infrastructure/http/authorization"
	"emersonargueta/m/v1/shared/infrastructure/http/response"
	"log"
	"net/http"

	"github.com/labstack/echo"
)

var _ Controller = &retrieveAchieverController{}

type retrieveAchieverController struct {
	Usecase       *usecase.RetrieveAchieverUsecase
	Logger        *log.Logger
	Authorization authorization.JwtService
}

// NewRetrieveAchieverController for retrieveAchiever achiever usecase
func NewRetrieveAchieverController(cgtRepos *persistence.Services, logger *log.Logger, authorizationService authorization.JwtService) Controller {
	retrieveAchieverUsecase := usecase.NewRetrieveAchieverUsecase(&cgtRepos.Achiever, authorizationService)

	ctrl := &retrieveAchieverController{
		Usecase:       retrieveAchieverUsecase,
		Logger:        logger,
		Authorization: authorizationService,
	}
	return ctrl
}

// Execute the usecase
func (ctrl *retrieveAchieverController) Execute(ctx echo.Context) (e error) {
	// extract user id from authKey stored by JwtMiddleware handler func
	authKey := ctx.Get("user").(string)
	userID, err := ctrl.Authorization.VerifyTokenAndExtractIDClaim(authKey)
	if err != nil {
		return response.ErrorResponse(ctx.Response().Writer, err, http.StatusInternalServerError, ctrl.Logger)
	}

	rDTO := &usecase.RetrieveAchieverDTO{AchieverUserID: userID}
	switch retrievedAchiever, e := ctrl.Usecase.Execute(rDTO); e {
	case nil:
		achieverDTO := dto.ToDTO(retrievedAchiever)
		response.EncodeJSON(ctx.Response().Writer, &achieverResponse{Achiever: achieverDTO}, ctrl.Logger)
	case achiever.ErrAchieverNotFound:
		return response.ErrorResponse(ctx.Response().Writer, e, http.StatusNotFound, ctrl.Logger)
	default:
		return response.ErrorResponse(ctx.Response().Writer, e, http.StatusInternalServerError, ctrl.Logger)
	}

	return nil
}
