package handlers

import (
	"github.com/Gonnekone/onlineStore/DTO"
	"github.com/Gonnekone/onlineStore/initializers"
	"github.com/Gonnekone/onlineStore/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

func CreateBrand(c *gin.Context) {
	var brands []string

	if err := c.ShouldBindJSON(&brands); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	for _, name := range brands {
		if result := initializers.DB.Create(&models.Brand{Name: name}); result.Error != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to create brands"})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "Brands created successfully"})
}

func ReadAllBrands(c *gin.Context) {
	var brands []models.Brand

	if result := initializers.DB.Find(&brands); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve brands"})
		return
	}

	var brandDTOs []DTO.BrandRequest

	for _, brand := range brands {
		brandDTOs = append(brandDTOs, DTO.BrandToBrandRequest(brand))
	}

	c.JSON(http.StatusOK, brandDTOs)
}

func UpdateBrand(c *gin.Context) {
	var brandDTOs []DTO.BrandRequest
	if err := c.ShouldBindJSON(&brandDTOs); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	for _, brandDTO := range brandDTOs {
		var brand models.Brand
		if result := initializers.DB.First(&brand, brandDTO.ID); result.Error != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Brand not found"})
			return
		}

		brand.Name = brandDTO.Name

		if result := initializers.DB.Save(&brand); result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update brand"})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "Brands updated successfully"})
}

func DeleteBrand(c *gin.Context) {
	var indexes []uint

	if err := c.ShouldBindJSON(&indexes); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	for _, i := range indexes {
		if result := initializers.DB.Unscoped().Delete(&models.Brand{}, i); result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete brand"})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "Brands deleted successfully"})
}
