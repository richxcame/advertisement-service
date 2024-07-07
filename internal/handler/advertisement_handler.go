package handler

import (
	"net/http"
	"strconv"

	"advertisement-service/internal/domain"
	"advertisement-service/internal/service"

	"github.com/labstack/echo/v4"
)

var adService service.Advertisement

// SetService initializes the service used by the handlers
func SetService(s service.Advertisement) {
	adService = s
}

// CreateAdvertisement creates a new advertisement
func CreateAdvertisement(c echo.Context) error {
	var ad domain.Advertisement
	if err := c.Bind(&ad); err != nil {
		return err
	}

	insertedAd, err := adService.Create(c.Request().Context(), ad)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusCreated, insertedAd)
}

// GetAdvertisement retrieves an advertisement by its ID
func GetAdvertisement(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, "Invalid ID")
	}

	ad, err := adService.Get(c.Request().Context(), int64(id))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, ad)
}

// UpdateAdvertisement updates an existing advertisement
func UpdateAdvertisement(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, "Invalid ID")
	}

	var ad domain.Advertisement
	if err := c.Bind(&ad); err != nil {
		return err
	}

	updatedAd, err := adService.Update(c.Request().Context(), int64(id), &ad)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, updatedAd)
}

// DeleteAdvertisement deletes an advertisement
func DeleteAdvertisement(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, "Invalid ID")
	}

	deletedID, err := adService.Delete(c.Request().Context(), int64(id))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, deletedID)
}

// ListAdvertisements lists all advertisements
func ListAdvertisements(c echo.Context) error {
	offset := 0
	limit := 10
	offset, err := strconv.Atoi(c.QueryParam("offset"))
	if err != nil {
		offset = 0
	}
	limit, err = strconv.Atoi(c.QueryParam("limit"))
	if err != nil {
		limit = 10
	}
	sortBy := c.QueryParam("sort_by")
	sortDir := c.QueryParam("sort_dir")

	ads, err := adService.List(c.Request().Context(), int64(offset), int64(limit), sortBy, sortDir)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, ads)
}
