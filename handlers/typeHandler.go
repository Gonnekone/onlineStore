package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/Gonnekone/onlineStore/DTO"
	"github.com/Gonnekone/onlineStore/initializers"
	"github.com/Gonnekone/onlineStore/models"
)

func CreateType(w http.ResponseWriter, r *http.Request) {
	var types []string

	if err := json.NewDecoder(r.Body).Decode(&types); err != nil {
		http.Error(w, "error decoding JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	for _, name := range types {
		if err := initializers.DB.Create(&models.Type{Name: name}).Error; err != nil {
			http.Error(w, "Failed to create types: "+err.Error(), http.StatusBadRequest)
			return
		}
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Types created successfully"})
}

func ReadAllTypes(w http.ResponseWriter, r *http.Request) {
	var types []models.Type

	if err := initializers.DB.Find(&types).Error; err != nil {
		http.Error(w, "Failed to retrieve types: "+err.Error(), http.StatusInternalServerError)
		return
	}

	var typeDTOs []DTO.TypeRequest

	for _, t := range types {
		typeDTOs = append(typeDTOs, DTO.TypeToTypeRequest(t))
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(typeDTOs)
}

func UpdateType(w http.ResponseWriter, r *http.Request) {
	var typeDTOs []DTO.TypeRequest
	if err := json.NewDecoder(r.Body).Decode(&typeDTOs); err != nil {
		http.Error(w, "error decoding JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	for _, typeDTO := range typeDTOs {
		var t models.Type
		if err := initializers.DB.First(&t, typeDTO.ID).Error; err != nil {
			http.Error(w, "Type not found: "+err.Error(), http.StatusNotFound)
			return
		}

		t.Name = typeDTO.Name

		if err := initializers.DB.Save(&t).Error; err != nil {
			http.Error(w, "Failed to update type: "+err.Error(), http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Types updated successfully"})
}

func DeleteType(w http.ResponseWriter, r *http.Request) {
	var indexes []uint

	if err := json.NewDecoder(r.Body).Decode(&indexes); err != nil {
		http.Error(w, "error decoding JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	for _, i := range indexes {
		if err := initializers.DB.Unscoped().Delete(&models.Type{}, i).Error; err != nil {
			http.Error(w, "Failed to delete type: "+err.Error(), http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Types deleted successfully"})
}
