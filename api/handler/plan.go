package handler

import (
	"github.com/compliance-framework/configuration-service/api"
	"github.com/compliance-framework/configuration-service/domain"
	"github.com/compliance-framework/configuration-service/event"
	"github.com/compliance-framework/configuration-service/store"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"net/http"
)

// TODO: Publishing the events from the handler is not a good idea. We should
//  publish the events from the domain services, following the business logic.

type PlanHandler struct {
	store     store.PlanStore
	publisher event.Publisher
	sugar     *zap.SugaredLogger
}

func (h *PlanHandler) Register(api *echo.Group) {
	api.POST("/plan", h.CreatePlan)
	api.POST("/plan/:id/assets", h.AddAsset)
}

func NewPlanHandler(l *zap.SugaredLogger, s store.PlanStore, p event.Publisher) *PlanHandler {
	return &PlanHandler{
		sugar:     l,
		store:     s,
		publisher: p,
	}
}

// CreatePlan godoc
// @Summary 		Create a plan
// @Description 	Creates a new plan in the system
// @Accept  		json
// @Produce  		json
// @Param   		plan body createPlanRequest true "Plan to add"
// @Success 		201 {object} planIdResponse
// @Failure 		401 {object} api.Error
// @Failure 		422 {object} api.Error
// @Failure 		500 {object} api.Error
// @Router 			/api/plan [post]
func (h *PlanHandler) CreatePlan(ctx echo.Context) error {
	// Initialize a new plan object
	p := domain.NewPlan()

	// Initialize a new createPlanRequest object
	req := createPlanRequest{}

	// Bind the incoming request to the plan object
	// If there's an error, return a 422 status code with the error message
	if err := req.bind(ctx, p); err != nil {
		return ctx.JSON(http.StatusUnprocessableEntity, api.NewError(err))
	}

	// Attempt to create the plan in the store
	// If there's an error, return a 500 status code with the error message
	id, err := h.store.Create(p)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, api.NewError(err))
	}

	// Publish an event indicating that a plan was created
	// If there's an error, log it
	// TODO: Should only publish when the Timing and the Subjects are set
	err = h.publisher(event.PlanCreated{Uuid: p.Uuid}, event.TopicTypePlan)
	if err != nil {
		h.sugar.Errorf("error publishing event: %v", err)
	}

	// If everything went well, return a 201 status code with the ID of the created plan
	return ctx.JSON(http.StatusCreated, planIdResponse{
		Id: id.(string),
	})
}

// AddAsset godoc
// @Summary Add asset to a plan
// @Description This method adds an asset to a specific plan by its ID.
// @Tags Plan
// @Accept  json
// @Produce  json
// @Param id path string true "Plan ID"
// @Param asset body addAssetRequest true "Asset to add"
// @Success 200 {object} api.Response "Successfully added the asset to the plan"
// @Failure 404 {object} api.Response "Plan not found"
// @Failure 422 {object} api.Response "Unprocessable Entity: Error binding the request"
// @Failure 500 {object} api.Response "Internal Server Error"
// @Router /plans/{id}/assets [post]
func (h *PlanHandler) AddAsset(ctx echo.Context) error {
	plan, err := h.store.GetById(ctx.Param("id"))
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, api.NewError(err))
	} else if plan == nil {
		return ctx.JSON(http.StatusNotFound, api.NotFound())
	}

	req := &addAssetRequest{}
	if err := ctx.Bind(req); err != nil {
		return err
	}

	plan.AddAsset(domain.Uuid(req.AssetUuid), req.Type)
	err = h.store.Update(plan)
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, nil)
}
