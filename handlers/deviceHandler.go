package handlers

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/Gonnekone/onlineStore/DTO"
	"github.com/Gonnekone/onlineStore/initializers"
	"github.com/Gonnekone/onlineStore/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func CreateDevice(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(10 << 20) // 10 MB max size
	if err != nil {
		http.Error(w, "error parsing form-data: "+err.Error(), http.StatusBadRequest)
		return
	}

	files := r.MultipartForm.File["files"]

	var deviceDTOs []DTO.DeviceDTO

	data := r.MultipartForm.Value["data"]

	if err := json.Unmarshal([]byte(data[0]), &deviceDTOs); err != nil {
		http.Error(w, "error parsing JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	for i, deviceDTO := range deviceDTOs {
		var device models.Device
		device = DTO.DeviceDTOToDevice(deviceDTO)

		if len(files) != 0 {
			file := files[i]
			fileName := uuid.NewString() + ".jpg"

			filePath := filepath.Join("static", fileName)
			dst, err := os.Create(filePath)
			if err != nil {
				http.Error(w, "Failed to save image: "+err.Error(), http.StatusInternalServerError)
				return
			}
			defer dst.Close()

			src, err := file.Open()
			if err != nil {
				http.Error(w, "Failed to open uploaded file: "+err.Error(), http.StatusInternalServerError)
				return
			}
			defer src.Close()

			if _, err := io.Copy(dst, src); err != nil {
				http.Error(w, "Failed to save image: "+err.Error(), http.StatusInternalServerError)
				return
			}

			device.Img = fileName
		} else {
			device.Img = "default.jpg"
		}

		if errors.Is(initializers.DB.Where("name = ?", device.Type.Name).First(&device.Type).Error, gorm.ErrRecordNotFound) {
			if err := initializers.DB.Create(&device.Type).Error; err != nil {
				http.Error(w, "Failed to create types: "+err.Error(), http.StatusBadRequest)
				return
			}
		}

		if errors.Is(initializers.DB.Where("name = ?", device.Brand.Name).First(&device.Brand).Error, gorm.ErrRecordNotFound) {
			if err := initializers.DB.Create(&device.Brand).Error; err != nil {
				http.Error(w, "Failed to create brands: "+err.Error(), http.StatusBadRequest)
				return
			}
		}

		if errors.Is(initializers.DB.Where("name = ?", device.Name).First(&device).Error, gorm.ErrRecordNotFound) {
			if err := initializers.DB.Create(&device).Error; err != nil {
				http.Error(w, "Failed to create devices: "+err.Error(), http.StatusBadRequest)
				return
			}
		}

		for _, i := range device.Info {
			i.DeviceID = device.ID
		}

		if errors.Is(initializers.DB.Where("id = ?", device.ID).First(&device.Info).Error, gorm.ErrRecordNotFound) {
			if err := initializers.DB.Create(&device.Info).Error; err != nil {
				http.Error(w, "Failed to create info: "+err.Error(), http.StatusBadRequest)
				return
			}
		}
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Devices created successfully"})
}

func ReadAllDevices(w http.ResponseWriter, r *http.Request) {
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))

	if page == 0 {
		page = 1
	}

	if limit == 0 {
		limit = 9
	}

	offset := (page - 1) * limit

	var typeIdPtr *int
	var brandIdPtr *int
	typeId, _ := strconv.Atoi(r.URL.Query().Get("typeId"))
	brandId, _ := strconv.Atoi(r.URL.Query().Get("brandId"))

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
		http.Error(w, "Failed to retrieve devices: "+err.Error(), http.StatusInternalServerError)
		return
	}

	var deviceDTOs []DTO.DeviceDTOImage
	for _, device := range devices {
		if err := initializers.DB.First(&device.Type, device.TypeID).Error; err != nil {
			http.Error(w, "Failed to retrieve type: "+err.Error(), http.StatusInternalServerError)
			return
		}

		if err := initializers.DB.First(&device.Brand, device.BrandID).Error; err != nil {
			http.Error(w, "Failed to retrieve brand: "+err.Error(), http.StatusInternalServerError)
			return
		}

		if err := initializers.DB.Where("device_id = ?", device.ID).Find(&device.Info).Error; err != nil {
			http.Error(w, "Failed to retrieve info: "+err.Error(), http.StatusInternalServerError)
			return
		}

		deviceDTOs = append(deviceDTOs, DTO.DeviceToDeviceDTOImage(device))
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{"count": count, "rows": deviceDTOs})
}

func ReadOneDevice(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(r.URL.Query().Get("id"))
	var device models.Device
	if err := initializers.DB.First(&device, id).Error; err != nil {
		http.Error(w, "Failed to retrieve device: "+err.Error(), http.StatusInternalServerError)
		return
	}

	var deviceDTO DTO.DeviceDTOImage
	if err := initializers.DB.First(&device.Type, device.TypeID).Error; err != nil {
		http.Error(w, "Failed to retrieve type: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if err := initializers.DB.First(&device.Brand, device.BrandID).Error; err != nil {
		http.Error(w, "Failed to retrieve brand: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if err := initializers.DB.Where("device_id = ?", device.ID).Find(&device.Info).Error; err != nil {
		http.Error(w, "Failed to retrieve info: "+err.Error(), http.StatusInternalServerError)
		return
	}

	deviceDTO = DTO.DeviceToDeviceDTOImage(device)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(deviceDTO)
}

func UpdateDevice(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(10 << 20) // 10 MB max size
	if err != nil {
		http.Error(w, "error parsing form-data: "+err.Error(), http.StatusBadRequest)
		return
	}

	files := r.MultipartForm.File["files"]

	var deviceRequests []DTO.DeviceRequest

	data := r.MultipartForm.Value["data"]

	if err := json.Unmarshal([]byte(data[0]), &deviceRequests); err != nil {
		http.Error(w, "error parsing JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	for i, deviceRequest := range deviceRequests {
		var device models.Device

		if len(files) != 0 {
			file := files[i]
			fileName := uuid.NewString() + ".jpg"

			filePath := filepath.Join("static", fileName)
			dst, err := os.Create(filePath)
			if err != nil {
				http.Error(w, "Failed to save image: "+err.Error(), http.StatusInternalServerError)
				return
			}
			defer dst.Close()

			src, err := file.Open()
			if err != nil {
				http.Error(w, "Failed to open uploaded file: "+err.Error(), http.StatusInternalServerError)
				return
			}
			defer src.Close()

			if _, err := io.Copy(dst, src); err != nil {
				http.Error(w, "Failed to save image: "+err.Error(), http.StatusInternalServerError)
				return
			}

			device.Img = fileName
		}

		if err := initializers.DB.First(&device, deviceRequest.ID).Error; err != nil {
			http.Error(w, "Failed to retrieve devices: "+err.Error(), http.StatusInternalServerError)
			return
		}

		if err := initializers.DB.First(&device.Type, deviceRequest.ID).Error; err != nil {
			http.Error(w, "Failed to retrieve type: "+err.Error(), http.StatusInternalServerError)
			return
		}

		if err := initializers.DB.Where("name = ?", device.Brand.Name).Error; err != nil {
			http.Error(w, "Failed to retrieve brand: "+err.Error(), http.StatusInternalServerError)
			return
		}

		if err := initializers.DB.Where("device_id = ?", deviceRequest.ID).Find(&device.Info).Error; err != nil {
			http.Error(w, "Failed to retrieve info: "+err.Error(), http.StatusInternalServerError)
			return
		}

		device = DTO.DeviceRequestToDevice(deviceRequest, device)

		device.Type.ID = device.TypeID
		device.Brand.ID = device.BrandID

		if err := initializers.DB.Unscoped().Where("device_id = ?", device.ID).Delete(&models.DeviceInfo{}).Error; err != nil {
			http.Error(w, "Failed to delete info: "+err.Error(), http.StatusInternalServerError)
			return
		}

		if err := initializers.DB.Save(&device).Error; err != nil {
			http.Error(w, "Failed to update device: "+err.Error(), http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Devices updated successfully"})
}

func DeleteDevice(w http.ResponseWriter, r *http.Request) {
	var indexes []uint

	if err := json.NewDecoder(r.Body).Decode(&indexes); err != nil {
		http.Error(w, "error parsing JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	for _, i := range indexes {
		var img string
		if err := initializers.DB.Raw("SELECT img FROM devices WHERE id = ?", i).Scan(&img).Error; err != nil {
			http.Error(w, "Failed to retrieve img: "+err.Error(), http.StatusInternalServerError)
			return
		}

		if img != "default.jpg" {
			imgPath := filepath.Join("static", img)
			if err := os.Remove(imgPath); err != nil {
				http.Error(w, "Failed to delete img file: "+err.Error(), http.StatusInternalServerError)
				return
			}
		}

		if err := initializers.DB.Unscoped().Where("device_id = ?", i).Delete(&models.DeviceInfo{}).Error; err != nil {
			http.Error(w, "Failed to delete info: "+err.Error(), http.StatusInternalServerError)
			return
		}

		if err := initializers.DB.Unscoped().Delete(&models.Device{}, i).Error; err != nil {
			http.Error(w, "Failed to delete device: "+err.Error(), http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Devices deleted successfully"})
}
