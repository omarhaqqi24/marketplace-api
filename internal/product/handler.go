package product

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/omarhaqqi24/marketplace-api/internal/auth"
)

type Handler interface {
	Create(c *gin.Context)
	List(c *gin.Context)
	Get(c *gin.Context)
	Update(c *gin.Context)
	Delete(c *gin.Context)
}

type handler struct {
	service Service
}

func NewHandler(service Service) Handler {
	return &handler{
		service: service,
	}
}

func (h *handler) Create(c *gin.Context) {
	var req CreateProductRequest

	err := c.ShouldBindJSON(&req)

	// log.Printf("%+v\n", req)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	claims := c.MustGet("claims").(*auth.Claims)
	sellerID, err := uuid.Parse(claims.UserID)

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "invalid user id in token",
		})
		return
	}

	product, err := h.service.Create(req, sellerID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": true,
			"message": "failed to create pruduct",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "product created successfully",
		"data": gin.H{
			"id":          product.ID,
			"name":        product.Name,
			"description": product.Description,
			"stock":       product.Stock,
			"price":       product.Price,
		},
	})
}

func (h *handler) List(c *gin.Context) {
	products, err := h.service.List()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "products fetch failed",
		})
		return
	}

	response := make([]ProductResponse, 0, len(products))

	for _, p := range products {
		response = append(response, ProductResponse{
			ID:          p.ID,
			Name:        p.Name,
			Description: p.Description,
			Price:       p.Price,
			Stock:       p.Stock,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "product list fetched successfully",
		"data":    response,
	})
}

func (h *handler) Get(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "invalid product id",
		})
		return
	}

	product, err := h.service.GetByID(id)

	if err != nil {
		if errors.Is(err, ErrProductNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"message": err.Error(),
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "failed to fetch product",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "product fetched successfully",
		"data": ProductResponse{
			ID:          product.ID,
			Name:        product.Name,
			Description: product.Description,
			Price:       product.Price,
			Stock:       product.Stock,
		},
	})
}

func (h *handler) Update(c *gin.Context) {
	var req UpdateProductRequest

	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "invalid product id",
		})
		return
	}

	sellerID, err := uuid.Parse(c.MustGet("claims").(*auth.Claims).UserID)

	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{
			"success": false,
			"message": "invalid user id in token",
		})
		return
	}

	err = c.ShouldBindJSON(&req)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	product, err := h.service.Update(id, sellerID, req)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "product updated successfully",
		"data": ProductResponse{
			ID:          product.ID,
			Name:        product.Name,
			Description: product.Description,
			Price:       product.Price,
			Stock:       product.Stock,
		},
	})
}

func (h *handler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "invalid product id",
		})
		return
	}

	userID, err := uuid.Parse(c.MustGet("claims").(*auth.Claims).UserID)

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "invalid user id in token",
		})
		return
	}

	err = h.service.Delete(id, userID)

	if err != nil {
		if errors.Is(err, ErrForbiddenAccess) {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"message": err.Error(),
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "product delete successfully",
	})
}
