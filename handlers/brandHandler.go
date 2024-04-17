package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/Gonnekone/onlineStore/DTO"
	"github.com/Gonnekone/onlineStore/initializers"
	"github.com/Gonnekone/onlineStore/models"
)

func CreateBrand(w http.ResponseWriter, r *http.Request) {
	var brands []string

	if err := json.NewDecoder(r.Body).Decode(&brands); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	for _, name := range brands {
		if result := initializers.DB.Create(&models.Brand{Name: name}); result.Error != nil {
			http.Error(w, "Failed to create brands", http.StatusBadRequest)
			return
		}
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "Brands created successfully"}`))
}

func ReadAllBrands(w http.ResponseWriter, r *http.Request) {
	var brands []models.Brand

	if result := initializers.DB.Find(&brands); result.Error != nil {
		http.Error(w, "Failed to retrieve brands", http.StatusInternalServerError)
		return
	}

	var brandDTOs []DTO.BrandRequest

	for _, brand := range brands {
		brandDTOs = append(brandDTOs, DTO.BrandToBrandRequest(brand))
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(brandDTOs)
}

func UpdateBrand(w http.ResponseWriter, r *http.Request) {
	var brandDTOs []DTO.BrandRequest
	if err := json.NewDecoder(r.Body).Decode(&brandDTOs); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	for _, brandDTO := range brandDTOs {
		var brand models.Brand
		if result := initializers.DB.First(&brand, brandDTO.ID); result.Error != nil {
			http.Error(w, "Brand not found", http.StatusNotFound)
			return
		}

		brand.Name = brandDTO.Name

		if result := initializers.DB.Save(&brand); result.Error != nil {
			http.Error(w, "Failed to update brand", http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "Brands updated successfully"}`))
}

func DeleteBrand(w http.ResponseWriter, r *http.Request) {
	var indexes []uint

	if err := json.NewDecoder(r.Body).Decode(&indexes); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	for _, i := range indexes {
		if result := initializers.DB.Unscoped().Delete(&models.Brand{}, i); result.Error != nil {
			http.Error(w, "Failed to delete brand", http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "Brands deleted successfully"}`))
}
