package main

import (
	"backend/ent"
	"backend/transactions"
	"net/http"
	"log"

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

func (h *Handler) mintMoment(c echo.Context) error {
	// 1. Siapkan variabel untuk menampung request body
	var req MintMomentRequest

	// 2. 'Bind' (ikat) JSON body yang masuk ke 'req' struct
	if err := c.Bind(&req); err != nil {
		log.Println("Error binding request:", err)
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request body: " + err.Error(),
		})
	}

	// 3. (Opsional) Validasi input
	if req.Recipient == "" || req.Name == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "recipient and name are required fields",
		})
	}

	// 4. Panggil fungsi 'MintNFTMoment' Anda (dari file 'mint_moment.go')
	//    Fungsi ini akan melakukan semua pekerjaan berat di blockchain
	err := transactions.MintNFTMoment(
		req.Recipient,
		req.Name,
		req.Description,
		req.Thumbnail,
	)

	// 5. Tangani hasilnya
	if err != nil {
		// Jika transaksi gagal, kirim error ke frontend
		log.Printf("Gagal menjalankan transaksi mint: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}

	// 6. Jika sukses, kirim respon 201 Created
	return c.JSON(http.StatusCreated, map[string]string{
		"message":   "NFT minted successfully!",
		"recipient": req.Recipient,
		"name":      req.Name,
	})
}