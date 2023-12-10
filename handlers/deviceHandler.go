package handlers

import (
	"encoding/json"
	"errors"
	"github.com/Gonnekone/onlineStore/DTO"
	"github.com/Gonnekone/onlineStore/initializers"
	"github.com/Gonnekone/onlineStore/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
)

func CreateDevice(c *gin.Context) {
	form, err := c.MultipartForm()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error parsing form-data": err.Error()})
		return
	}

	files := form.File["files"]

	var deviceDTOs []DTO.DeviceDTO

	data := form.Value["data"]

	if err := json.Unmarshal([]byte(data[0]), &deviceDTOs); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error parsing JSON": err.Error()})
		return
	}

	for i, deviceDTO := range deviceDTOs {
		var device models.Device
		device = DTO.DeviceDTOToDevice(deviceDTO)

		if len(files) != 0 {
			file := files[i]
			fileName := uuid.NewString() + ".jpg"

			err := c.SaveUploadedFile(file, filepath.Join("static", fileName))
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save image"})
				return
			}
			device.Img = fileName
		} else {
			device.Img = "default.jpg"
		}

		if errors.Is(initializers.DB.Where("name = ?", device.Type.Name).First(&device.Type).Error, gorm.ErrRecordNotFound) {
			if err := initializers.DB.Create(&device.Type).Error; err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"Failed to create types": err.Error()})
				return
			}
		}

		if errors.Is(initializers.DB.Where("name = ?", device.Brand.Name).First(&device.Brand).Error, gorm.ErrRecordNotFound) {
			if err := initializers.DB.Create(&device.Brand).Error; err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"Failed to create brands": err.Error()})
				return
			}
		}

		if errors.Is(initializers.DB.Where("name = ?", device.Name).First(&device).Error, gorm.ErrRecordNotFound) {
			if err := initializers.DB.Create(&device).Error; err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"Failed to create devices": err.Error()})
				return
			}
		}

		for _, i := range device.Info {
			i.DeviceID = device.ID
		}

		if errors.Is(initializers.DB.Where("id = ?", device.ID).First(&device.Info).Error, gorm.ErrRecordNotFound) {
			if err := initializers.DB.Create(&device.Info).Error; err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"Failed to create info": err.Error()})
				return
			}
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "Devices created successfully"})
}

func ReadAllDevices(c *gin.Context) {
	limit, _ := strconv.Atoi(c.Query("limit"))
	page, _ := strconv.Atoi(c.Query("page"))

	if page == 0 {
		page = 1
	}

	if limit == 0 {
		limit = 9
	}

	offset := (page - 1) * limit

	var typeIdPtr *int
	var brandIdPtr *int
	typeId, _ := strconv.Atoi(c.Query("typeId"))
	brandId, _ := strconv.Atoi(c.Query("brandId"))

	if typeId != 0 {
		typeIdPtr = &typeId
	}

	if brandId != 0 {
		brandIdPtr = &brandId
	}

	var devices []models.Device
	var count int64
	if err := initializers.DB.Where("type_id = COALESCE(?, type_id) AND brand_id = COALESCE(?, brand_id)", typeIdPtr, brandIdPtr).
		Limit(limit).Offset(offset).Find(&devices).Count(&count).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Failed to retrieve devices": err.Error()})
		return
	}

	var deviceDTOs []DTO.DeviceDTOImage
	for _, device := range devices {
		if err := initializers.DB.First(&device.Type, device.TypeID).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Failed to retrieve type": err.Error()})
			return
		}

		if err := initializers.DB.First(&device.Brand, device.BrandID).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Failed to retrieve brand": err.Error()})
			return
		}

		if err := initializers.DB.Where("device_id = ?", device.ID).Find(&device.Info).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Failed to retrieve info": err.Error()})
			return
		}

		deviceDTOs = append(deviceDTOs, DTO.DeviceToDeviceDTOImage(device))
	}

	c.JSON(http.StatusOK, gin.H{"count": count, "rows": deviceDTOs})
}

func ReadOneDevice(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var device models.Device
	if err := initializers.DB.First(&device, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Failed to retrieve device": err.Error()})
		return
	}

	var deviceDTO DTO.DeviceDTOImage
	if err := initializers.DB.First(&device.Type, device.TypeID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Failed to retrieve type": err.Error()})
		return
	}

	if err := initializers.DB.First(&device.Brand, device.BrandID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Failed to retrieve brand": err.Error()})
		return
	}

	if err := initializers.DB.Where("device_id = ?", device.ID).Find(&device.Info).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Failed to retrieve info": err.Error()})
		return
	}

	deviceDTO = DTO.DeviceToDeviceDTOImage(device)

	c.JSON(http.StatusOK, deviceDTO)
}

func UpdateDevice(c *gin.Context) {
	form, err := c.MultipartForm()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error parsing form-data": err.Error()})
		return
	}

	files := form.File["files"]

	var deviceRequests []DTO.DeviceRequest

	data := form.Value["data"]

	if err := json.Unmarshal([]byte(data[0]), &deviceRequests); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error parsing JSON": err.Error()})
		return
	}

	for i, deviceRequest := range deviceRequests {
		var device models.Device

		if len(files) != 0 {
			file := files[i]
			fileName := uuid.NewString() + ".jpg"

			err := c.SaveUploadedFile(file, filepath.Join("static", fileName))
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"Failed to save image": err.Error()})
				return
			}
			device.Img = fileName
		}

		if err := initializers.DB.First(&device, deviceRequest.ID).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Failed to retrieve devices": err.Error()})
			return
		}

		if err := initializers.DB.First(&device.Type, deviceRequest.ID).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Failed to retrieve type": err.Error()})
			return
		}

		if err := initializers.DB.Where("name = ?", device.Brand.Name).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Failed to retrieve brand": err.Error()})
			return
		}

		if err := initializers.DB.Where("device_id = ?", deviceRequest.ID).Find(&device.Info).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Failed to retrieve info": err.Error()})
			return
		}

		device = DTO.DeviceRequestToDevice(deviceRequest, device)

		device.Type.ID = device.TypeID
		device.Brand.ID = device.BrandID

		if err := initializers.DB.Unscoped().Where("device_id = ?", device.ID).Delete(&models.DeviceInfo{}).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Failed to delete info": err.Error()})
			return
		}

		if err := initializers.DB.Save(&device).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Failed to update device": err.Error()})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "Devices updated successfully"})
}

func DeleteDevice(c *gin.Context) {
	var indexes []uint

	if err := c.ShouldBindJSON(&indexes); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error parsing JSON": err.Error()})
		return
	}

	for _, i := range indexes {
		var img string
		if err := initializers.DB.Raw("SELECT img FROM devices WHERE id = ?", i).Scan(&img).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Failed to retrieve img": err.Error()})
			return
		}

		if img != "default.jpg" {
			imgPath := filepath.Join("static", img)
			if err := os.Remove(imgPath); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"Failed to delete img file": err.Error()})
				return
			}
		}

		if err := initializers.DB.Unscoped().Where("device_id = ?", i).Delete(&models.DeviceInfo{}).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Failed to delete info": err.Error()})
			return
		}

		if err := initializers.DB.Unscoped().Delete(&models.Device{}, i).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Failed to delete device": err.Error()})
			return
		}

	}

	c.JSON(http.StatusOK, gin.H{"message": "Devices deleted successfully"})
}
