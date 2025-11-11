package main

import (
	"backend/ent"
	"net/http"

	"github.com/labstack/echo/v4"
)

// Handler adalah struct kustom yang akan kita gunakan
// untuk "menyuntikkan" (inject) koneksi database 'ent' kita
// ke dalam fungsi-fungsi API kita.
type Handler struct {
	DB *ent.Client
}

// getAccessories adalah handler untuk GET /accessories
// Ini akan mengambil semua aksesoris dari database
func (h *Handler) getAccessories(c echo.Context) error {
	ctx := c.Request().Context()

	// Gunakan 'h.DB' (klien 'ent' kita) untuk query ke database
	accessories, err := h.DB.NFTAccessory.Query().
		// Opsional: Anda bisa 'preload' data owner
		// WithOwner(). 
		All(ctx)
	
	if err != nil {
		// Jika ada error, kirim respon error 500
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}

	// Jika sukses, kirim data 'accessories' sebagai JSON
	return c.JSON(http.StatusOK, accessories)
}

// getMoments adalah handler untuk GET /moments
func (h *Handler) getMoments(c echo.Context) error {
	ctx := c.Request().Context()

	moments, err := h.DB.NFTMoment.Query().
		// Anda bisa 'preload' data relasi 'equip' di sini
		// WithEquippedAccessories().
		All(ctx)
	
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, moments)
}

// Anda bisa menambahkan handler lain di sini...
// (misal: getAccessoryByID, createListing, dll.)