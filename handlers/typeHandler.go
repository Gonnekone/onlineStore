package handlers

import (
	"github.com/Gonnekone/onlineStore/DTO"
	"github.com/Gonnekone/onlineStore/initializers"
	"github.com/Gonnekone/onlineStore/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

func CreateType(c *gin.Context) {
	var types []string

	if err := c.ShouldBindJSON(&types); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	for _, name := range types {
		if err := initializers.DB.Create(&models.Type{Name: name}).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"Failed to create types": err.Error()})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "Types created successfully"})
}

func ReadAllTypes(c *gin.Context) {
	var types []models.Type

	if err := initializers.DB.Find(&types).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var typeDTOs []DTO.TypeRequest

	for _, t := range types {
		typeDTOs = append(typeDTOs, DTO.TypeToTypeRequest(t))
	}

	c.JSON(http.StatusOK, typeDTOs)
}

func UpdateType(c *gin.Context) {
	var typeDTOs []DTO.TypeRequest
	if err := c.ShouldBindJSON(&typeDTOs); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	for _, typeDTO := range typeDTOs {
		var t models.Type
		if err := initializers.DB.First(&t, typeDTO.ID).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"Type not found": err.Error()})
			return
		}

		t.Name = typeDTO.Name

		if err := initializers.DB.Save(&t).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Failed to update type": err.Error()})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "Types updated successfully"})
}

func DeleteType(c *gin.Context) {
	var indexes []uint

	if err := c.ShouldBindJSON(&indexes); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	for _, i := range indexes {
		if err := initializers.DB.Unscoped().Delete(&models.Type{}, i).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Failed to delete type": err.Error()})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "Types deleted successfully"})
}
