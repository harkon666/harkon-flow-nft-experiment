package main

import (
	"backend/ent"
	"backend/transactions"
	"backend/utils"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

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

	// 1. Ambil data TEKS dari form multipart
	recipient := c.FormValue("recipient")
	eventPassID := c.FormValue("eventPassID")
	useFreeMint := c.FormValue("useFreeMint")
	tier := c.FormValue("tier")
	name := c.FormValue("name")
	description := c.FormValue("description")

	// 2. Validasi input teks
	if recipient == "" || name == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "recipient dan name adalah field wajib",
		})
	}

	// 3. Ambil data FILE dari form multipart
	// "thumbnail" adalah 'name' dari input file di frontend
	file, err := c.FormFile("thumbnail")
	if err != nil {
		log.Println("Error mengambil form file 'thumbnail':", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "file 'thumbnail' wajib ada"})
	}

	// 4. Buka file yang di-upload
	src, err := file.Open()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Gagal membuka file yang di-upload"})
	}
	defer src.Close()

	// 5. Simpan file ke path sementara (temp)
	// (Di aplikasi produksi, gunakan nama unik/random)
	tempFilePath := "moment_" + fmt.Sprint(time.Now().UnixNano()) + "_" + file.Filename
	dst, err := os.Create(tempFilePath)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Gagal membuat file sementara"})
	}

	// Salin file yang di-upload ke file sementara
	if _, err = io.Copy(dst, src); err != nil {
		dst.Close()             // Tutup dulu sebelum hapus
		os.Remove(tempFilePath) // Hapus file temp jika copy gagal
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Gagal menyimpan file sementara"})
	}
	dst.Close() // Tutup file setelah selesai 'copy'

	// --- Jadwalkan PENGHAPUSAN file sementara ---
	// 'defer' akan berjalan di akhir fungsi 'mintMoment'
	defer func() {
		log.Println("Menghapus file sementara:", tempFilePath)
		os.Remove(tempFilePath)
	}()

	// 6. --- TAHAP MODERASI (PENTING) ---
	// Di sinilah Anda memanggil Google Cloud Vision / AWS Rekognition
	// pada 'tempFilePath'
	//
	// isSafe, err := runModeration(tempFilePath)
	// if err != nil || !isSafe {
	// 	  return c.JSON(http.StatusForbidden, map[string]string{"error": "Konten foto dilarang"})
	// }
	// log.Println("Moderasi AI lolos.")

	// 7. Upload file yang aman ke Pinata
	log.Println("Mengunggah", tempFilePath, "ke Pinata...")
	pinataResp, err := utils.UploadToPinata(tempFilePath)
	if err != nil {
		log.Printf("Gagal upload ke Pinata: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Gagal mengunggah ke IPFS"})
	}

	// Buat URL IPFS yang benar
	thumbnailUrl := fmt.Sprintf("ipfs://%s", pinataResp.IpfsHash)
	log.Println("Berhasil di-pin ke Pinata:", thumbnailUrl)

	// 8. Panggil transaksi blockchain dengan URL IPFS baru
	//    (Kita ganti 'req.Thumbnail' dengan 'thumbnailUrl')
	err = transactions.MintNFTMoment(
		recipient,
		eventPassID,
		name,
		description,
		thumbnailUrl,
		useFreeMint,
		tier, // <-- Menggunakan URL IPFS yang baru!
	)

	// 9. Tangani hasilnya
	if err != nil {
		log.Printf("Gagal menjalankan transaksi mint: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}

	// 10. Jika sukses, kirim respon 201 Created
	return c.JSON(http.StatusCreated, map[string]string{
		"message":   "NFT minted successfully!",
		"recipient": recipient,
		"name":      name,
		"thumbnail": thumbnailUrl, // Kirimkan URL IPFS baru ke frontend
	})
}
