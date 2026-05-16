package controllers

import (
	"core-service/models"
	"core-service/models/dto"
	"core-service/pkg/response"
	"core-service/services"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/go-playground/validator/v10"
)

// CarController handles car-related HTTP requests
type CarController struct {
	carService services.ICarService
}

// ICarController defines the interface for car controller
type ICarController interface {
	GetAllCars(ctx *gin.Context)
	GetCarByID(ctx *gin.Context)
	CreateCar(ctx *gin.Context)
	UpdateCar(ctx *gin.Context)
	DeleteCar(ctx *gin.Context)
}

// NewCarController creates a new car controller
func NewCarController(carService services.ICarService) ICarController {
	return &CarController{
		carService: carService,
	}
}

// GetAllCars handles GET /api/v1/cars
// @Summary Get all cars
// @Description Retrieves all cars
// @Tags Cars
// @Produce json
// @Success 200 {object} response.Response
// @Router /api/v1/cars [get]
func (c *CarController) GetAllCars(ctx *gin.Context) {
	cars, err := c.carService.GetAllCars(ctx.Request.Context())
	if err != nil {
		response.InternalServerError(ctx, "Failed to fetch cars")
		return
	}

	response.OK(ctx, "Cars fetched successfully", cars)
}

// GetCarByID handles GET /api/v1/cars/:id
// @Summary Get car by ID
// @Description Retrieves a specific car by ID
// @Tags Cars
// @Produce json
// @Param id path string true "Car ID"
// @Success 200 {object} response.Response
// @Router /api/v1/cars/{id} [get]
func (c *CarController) GetCarByID(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		response.BadRequest(ctx, "Invalid car ID format")
		return
	}

	car, err := c.carService.GetCarByID(ctx.Request.Context(), id)
	if err != nil {
		response.InternalServerError(ctx, "Failed to fetch car")
		return
	}

	if car == nil {
		response.NotFound(ctx, "Car not found")
		return
	}

	response.OK(ctx, "Car fetched successfully", car)
}

// CreateCar handles POST /api/v1/cars
// @Summary Create a new car
// @Description Creates a new car
// @Tags Cars
// @Accept json
// @Produce json
// @Param request body dto.CreateCarRequest true "Car data"
// @Success 201 {object} response.Response
// @Router /api/v1/cars [post]
func (c *CarController) CreateCar(ctx *gin.Context) {
	var req dto.CreateCarRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.BadRequest(ctx, "Invalid request body: "+err.Error())
		return
	}

	if err := validator.New().Struct(req); err != nil {
		response.BadRequest(ctx, "Validation failed: "+err.Error())
		return
	}

	car := &models.Car{
		Brand:        req.Brand,
		Model:        req.Model,
		Year:         req.Year,
		LicensePlate: req.LicensePlate,
		Color:        req.Color,
		Transmission: req.Transmission,
		ImageURL:     req.ImageURL,
		Notes:        req.Notes,
		Status:       models.CarStatusAvailable,
	}

	if err := c.carService.CreateCar(ctx.Request.Context(), car); err != nil {
		response.InternalServerError(ctx, "Failed to create car")
		return
	}

	response.Created(ctx, "Car created successfully", car)
}

// UpdateCar handles PUT /api/v1/cars/:id
// @Summary Update a car
// @Description Updates an existing car
// @Tags Cars
// @Accept json
// @Produce json
// @Param id path string true "Car ID"
// @Param request body dto.UpdateCarRequest true "Car data"
// @Success 200 {object} response.Response
// @Router /api/v1/cars/{id} [put]
func (c *CarController) UpdateCar(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		response.BadRequest(ctx, "Invalid car ID format")
		return
	}

	// Get existing car
	car, err := c.carService.GetCarByID(ctx.Request.Context(), id)
	if err != nil {
		response.InternalServerError(ctx, "Failed to fetch car")
		return
	}

	if car == nil {
		response.NotFound(ctx, "Car not found")
		return
	}

	var req dto.UpdateCarRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.BadRequest(ctx, "Invalid request body: "+err.Error())
		return
	}

	// Update fields
	if req.Brand != "" {
		car.Brand = req.Brand
	}
	if req.Model != "" {
		car.Model = req.Model
	}
	if req.Year != 0 {
		car.Year = req.Year
	}
	if req.Color != "" {
		car.Color = req.Color
	}
	if req.Transmission != "" {
		car.Transmission = req.Transmission
	}
	if req.Status != "" {
		car.Status = req.Status
	}
	if req.Mileage != 0 {
		car.Mileage = req.Mileage
	}
	if req.ImageURL != "" {
		car.ImageURL = req.ImageURL
	}
	if req.Notes != "" {
		car.Notes = req.Notes
	}

	if err := c.carService.UpdateCar(ctx.Request.Context(), car); err != nil {
		response.InternalServerError(ctx, "Failed to update car")
		return
	}

	response.OK(ctx, "Car updated successfully", car)
}

// DeleteCar handles DELETE /api/v1/cars/:id
// @Summary Delete a car
// @Description Deletes a car by ID
// @Tags Cars
// @Produce json
// @Param id path string true "Car ID"
// @Success 200 {object} response.Response
// @Router /api/v1/cars/{id} [delete]
func (c *CarController) DeleteCar(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		response.BadRequest(ctx, "Invalid car ID format")
		return
	}

	car, err := c.carService.GetCarByID(ctx.Request.Context(), id)
	if err != nil {
		response.InternalServerError(ctx, "Failed to fetch car")
		return
	}

	if car == nil {
		response.NotFound(ctx, "Car not found")
		return
	}

	if err := c.carService.DeleteCar(ctx.Request.Context(), car); err != nil {
		response.InternalServerError(ctx, "Failed to delete car")
		return
	}

	response.OK(ctx, "Car deleted successfully", nil)
}