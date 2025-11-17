package main

import (
	"backend/ent"
	"backend/ent/event"
	"backend/ent/listing"
	"backend/ent/nftaccessory"
	"backend/ent/nftmoment"
	"backend/ent/user"
	"backend/transactions"
	"backend/utils"
	"fmt"
	"io"
	"log"
	"math"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
)

// Handler adalah struct kustom yang akan kita gunakan
// untuk "menyuntikkan" (inject) koneksi database 'ent' kita
// ke dalam fungsi-fungsi API kita.
type Handler struct {
	DB *ent.Client
}

type Pagination struct {
	TotalItems  int `json:"totalItems"`
	TotalPages  int `json:"totalPages"`
	CurrentPage int `json:"currentPage"`
	PageSize    int `json:"pageSize"`
}

// APIResponse adalah "bungkusan" standar kita.
// 'omitempty' akan menyembunyikan field jika nilainya kosong/nil
type APIResponse struct {
	Data       interface{} `json:"data,omitempty"`
	Pagination *Pagination `json:"pagination,omitempty"`
	Error      string      `json:"error,omitempty"`
}

func getPagination(c echo.Context) (limit, offset, page, pageSize int) {
	// Nilai default
	const defaultPageSize = 10
	const defaultPage = 1

	page, err := strconv.Atoi(c.QueryParam("page"))
	if err != nil || page < 1 {
		page = defaultPage
	}

	pageSize, err = strconv.Atoi(c.QueryParam("pageSize"))
	if err != nil || pageSize < 1 {
		pageSize = defaultPageSize
	}

	limit = pageSize
	offset = (page - 1) * pageSize

	// Kembalikan semua nilai
	return limit, offset, page, pageSize
}

// GET /moments -> Mengambil SEMUA (untuk 'Explore')
// GET /moments?owner_address=0x123 -> Mengambil HANYA milik '0x123'
// GET /moments?owner_address=0x123&page=2 -> Pagination
func (h *Handler) getMoments(c echo.Context) error {
	ctx := c.Request().Context()

	// 1. Dapatkan parameter pagination
	limit, offset, page, pageSize := getPagination(c)

	// 2. Siapkan query dasar
	query := h.DB.NFTMoment.Query()

	// 3. Terapkan Filter (jika ada)
	//    Cek apakah ada filter 'owner_address' di URL
	ownerAddress := c.QueryParam("owner_address")
	if ownerAddress != "" {
		// Jika ada, tambahkan 'Where' clause (filter) ke query
		query = query.Where(
			// Filter berdasarkan relasi 'owner'
			nftmoment.HasOwnerWith(user.AddressEQ(ownerAddress)),
		)
	}
	// ---

	// 4. Hitung total item (setelah filter diterapkan)
	totalItems, err := query.Count(ctx)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, APIResponse{Error: err.Error()})
	}

	// 5. Buat Metadata Pagination
	totalPages := int(math.Ceil(float64(totalItems) / float64(pageSize)))
	pagination := &Pagination{
		TotalItems:  totalItems,
		TotalPages:  totalPages,
		CurrentPage: page,
		PageSize:    pageSize,
	}

	// 6. Jalankan Query UTAMA dengan Limit/Offset
	moments, err := query.
		WithOwner().
		WithEquippedAccessories(). // <-- 'Preload' data aksesoris yang terpasang
		WithMintedWithPass().      // <-- 'Preload' data EventPass yang digunakan
		Limit(limit).
		Offset(offset).
		Order(ent.Desc("id")). // Urutkan dari yang terbaru
		All(ctx)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, APIResponse{Error: err.Error()})
	}

	// 7. Kembalikan Respon Standar (Terbungkus)
	response := APIResponse{
		Data:       moments,
		Pagination: pagination,
	}
	return c.JSON(http.StatusOK, response)
}

// getAccessories adalah handler untuk GET /accessories
// Ini akan mengambil semua aksesoris dari database
// Endpoint ini sekarang mendukung:
//
//	GET /accessories -> Mengambil SEMUA (untuk halaman 'Explore')
//	GET /accessories?owner_address=0x123 -> Mengambil HANYA milik '0x123'
//	GET /accessories?owner_address=0x123&page=2 -> Pagination
func (h *Handler) getAccessories(c echo.Context) error {
	ctx := c.Request().Context()

	// 1. Dapatkan parameter pagination
	limit, offset, page, pageSize := getPagination(c)

	// 2. Siapkan query dasar
	query := h.DB.NFTAccessory.Query()

	// 3. --- INI ADALAH LOGIKA BARU ANDA ---
	//    Cek apakah ada filter 'owner_address' di URL
	ownerAddress := c.QueryParam("owner_address")
	if ownerAddress != "" {
		// Jika ada, tambahkan 'Where' clause (filter) ke query
		query = query.Where(
			// Filter berdasarkan relasi 'owner'
			nftaccessory.HasOwnerWith(user.AddressEQ(ownerAddress)),
		)
	}
	// --- AKHIR LOGIKA BARU ---

	// 4. Hitung total item (setelah filter diterapkan)
	totalItems, err := query.Count(ctx)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, APIResponse{Error: err.Error()})
	}

	// 5. Buat Metadata Pagination
	totalPages := int(math.Ceil(float64(totalItems) / float64(pageSize)))
	pagination := &Pagination{
		TotalItems:  totalItems,
		TotalPages:  totalPages,
		CurrentPage: page,
		PageSize:    pageSize,
	}

	// 6. Jalankan Query UTAMA dengan Limit/Offset
	accessories, err := query.
		WithOwner(). // (Opsional: 'preload' data owner)
		Limit(limit).
		Offset(offset).
		Order(ent.Desc("id")). // Urutkan dari yang terbaru
		All(ctx)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, APIResponse{Error: err.Error()})
	}

	// 7. Kembalikan Respon Standar (Terbungkus)
	response := APIResponse{
		Data:       accessories,
		Pagination: pagination,
	}
	return c.JSON(http.StatusOK, response)
}

func (h *Handler) handleUGCUpload(c echo.Context) (string, error) {
	// 1. Ambil data FILE dari form multipart
	file, err := c.FormFile("thumbnail")
	if err != nil {
		log.Println("Error mengambil form file 'thumbnail':", err)
		return "", fmt.Errorf("file 'thumbnail' wajib ada")
	}

	// 2. Buka file yang di-upload
	src, err := file.Open()
	if err != nil {
		return "", fmt.Errorf("gagal membuka file yang di-upload: %w", err)
	}
	defer src.Close()

	// 3. Buat nama file sementara yang unik
	ext := filepath.Ext(file.Filename)
	tempFilePath := fmt.Sprintf("temp_upload_%d%s", time.Now().UnixNano(), ext)

	dst, err := os.Create(tempFilePath)
	if err != nil {
		return "", fmt.Errorf("gagal membuat file sementara: %w", err)
	}

	// 4. Salin file
	if _, err = io.Copy(dst, src); err != nil {
		dst.Close()
		os.Remove(tempFilePath)
		return "", fmt.Errorf("gagal menyimpan file sementara: %w", err)
	}
	dst.Close()

	// 5. Jadwalkan penghapusan
	defer func() {
		log.Println("Menghapus file sementara:", tempFilePath)
		os.Remove(tempFilePath)
	}()

	// 6. --- TAHAP MODERASI (PENTING) ---
	// (Panggil Google Cloud Vision / AWS Rekognition Anda di sini)
	// isSafe, err := runModeration(tempFilePath)
	// if err != nil || !isSafe {
	//     return "", fmt.Errorf("konten foto dilarang")
	// }
	// log.Println("Moderasi AI lolos.")

	// 7. Upload ke Pinata
	log.Println("Mengunggah", tempFilePath, "ke Pinata...")
	pinataResp, err := utils.UploadToPinata(tempFilePath)
	if err != nil {
		log.Printf("Gagal upload ke Pinata: %v", err)
		return "", fmt.Errorf("gagal mengunggah ke IPFS: %w", err)
	}

	// Buat URL IPFS yang benar
	thumbnailUrl := fmt.Sprintf("https://white-lazy-marten-351.mypinata.cloud/ipfs/%s", pinataResp.IpfsHash)
	log.Println("Berhasil di-pin ke Pinata:", thumbnailUrl)

	return thumbnailUrl, nil
}

func (h *Handler) freeMintMoment(c echo.Context) error {
	// 1. Ambil data TEKS
	recipient := c.FormValue("recipient")
	name := c.FormValue("name")
	description := c.FormValue("description")

	// 2. Validasi
	if recipient == "" || name == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "recipient dan name adalah field wajib"})
	}

	// 3. Panggil helper untuk 'pekerjaan kotor' (upload)
	thumbnailUrl, err := h.handleUGCUpload(c)
	if err != nil {
		// handleUGCUpload sudah mem-format error-nya
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	// 4. Panggil transaksi
	err = transactions.FreeMintNFTMoment(
		recipient,
		name,
		description,
		thumbnailUrl,
	)
	if err != nil {
		log.Printf("Gagal menjalankan transaksi mint: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	// 5. Kirim respon sukses
	return c.JSON(http.StatusCreated, map[string]string{
		"message":   "NFT minted successfully!",
		"recipient": recipient,
		"name":      name,
		"thumbnail": thumbnailUrl,
	})
}

func (h *Handler) mintMomentWithEventPass(c echo.Context) error {
	// 1. Ambil data TEKS
	recipient := c.FormValue("recipient")
	eventPassID := c.FormValue("eventPassID")
	tier := c.FormValue("tier")
	name := c.FormValue("name")
	description := c.FormValue("description")

	// 2. Validasi
	if recipient == "" || name == "" || eventPassID == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "recipient, name, dan eventPassID adalah field wajib"})
	}

	// 3. Panggil helper untuk 'pekerjaan kotor' (upload)
	thumbnailUrl, err := h.handleUGCUpload(c)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	// 4. Panggil transaksi
	err = transactions.MintNFTMomentWithEventPass(
		recipient,
		eventPassID,
		name,
		description,
		thumbnailUrl,
		tier,
	)
	if err != nil {
		log.Printf("Gagal menjalankan transaksi mint-with-pass: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	// 5. Kirim respon sukses
	return c.JSON(http.StatusCreated, map[string]string{
		"message":   "NFT minted with Event Pass successfully!",
		"recipient": recipient,
		"name":      name,
		"thumbnail": thumbnailUrl,
	})
}

// --- HANDLER BARU: GET /listings ---
// Mengambil daftar penjualan (listings) dari marketplace
// Mendukung Pagination: ?page=1&pageSize=30
// Mendukung Filter: ?seller_address=0x...
func (h *Handler) getListings(c echo.Context) error {
	ctx := c.Request().Context()

	// 1. Dapatkan parameter pagination
	limit, offset, page, pageSize := getPagination(c)

	// 2. Siapkan query dasar
	query := h.DB.Listing.Query()

	// 3. Terapkan Filter (jika ada)
	sellerAddress := c.QueryParam("seller_address")
	if sellerAddress != "" {
		query = query.Where(
			listing.HasSellerWith(user.AddressEQ(sellerAddress)),
		)
	}

	// 4. HITUNG TOTAL ITEM (PENTING!)
	// Jalankan query COUNT() SEBELUM Limit/Offset
	totalItems, err := query.Count(ctx)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, APIResponse{Error: err.Error()})
	}

	// 5. Buat Metadata Pagination
	totalPages := int(math.Ceil(float64(totalItems) / float64(pageSize)))
	pagination := &Pagination{
		TotalItems:  totalItems,
		TotalPages:  totalPages,
		CurrentPage: page,
		PageSize:    pageSize,
	}

	// 6. Jalankan Query UTAMA dengan Limit/Offset
	listings, err := query.
		WithSeller().
		WithNftAccessory().
		Limit(limit).
		Offset(offset).
		Order(ent.Desc("id")).
		All(ctx)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, APIResponse{Error: err.Error()})
	}

	// 7. Kembalikan Respon Standar (Terbungkus)
	response := APIResponse{
		Data:       listings,
		Pagination: pagination,
	}
	return c.JSON(http.StatusOK, response)
}

// --- HANDLER BARU: GET /events ---
// Mengambil daftar event (seperti Luma)
// Mendukung Pagination: ?page=1&pageSize=10
// Mendukung Filter: ?type=0 (0=online, 1=offline)
func (h *Handler) getEvents(c echo.Context) error {
	ctx := c.Request().Context()

	// 1. Dapatkan parameter pagination lengkap
	limit, offset, page, pageSize := getPagination(c)

	// 2. Siapkan query dasar
	query := h.DB.Event.Query()

	// 3. Terapkan Filter (kode Anda sudah benar)
	eventTypeParam := c.QueryParam("type")
	if eventTypeParam != "" {
		eventType, err := strconv.Atoi(eventTypeParam)
		if err == nil { // Abaikan jika 'type' bukan angka
			query = query.Where(event.EventTypeEQ(uint8(eventType)))
		}
	}
	// (Anda bisa menambahkan filter lain di sini)

	// 4. HITUNG TOTAL ITEM (PENTING!)
	// Jalankan query COUNT() SEBELUM Limit/Offset
	totalItems, err := query.Count(ctx)
	if err != nil {
		// Gunakan 'APIResponse' untuk error
		return c.JSON(http.StatusInternalServerError, APIResponse{Error: err.Error()})
	}

	// 5. Buat Metadata Pagination
	totalPages := int(math.Ceil(float64(totalItems) / float64(pageSize)))
	pagination := &Pagination{
		TotalItems:  totalItems,
		TotalPages:  totalPages,
		CurrentPage: page,
		PageSize:    pageSize,
	}

	// 6. Jalankan Query UTAMA dengan Limit/Offset
	events, err := query.
		WithHost(). // Ambil data 'User' (host)
		// WithAttendances(). // Hati-hati: 'Eager loading' ini bisa sangat berat jika ada 1000 peserta
		Limit(limit).
		Offset(offset).
		Order(ent.Desc("start_date")). // Urutkan dari yang paling baru
		All(ctx)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, APIResponse{Error: err.Error()})
	}

	// 7. Kembalikan Respon Standar (Terbungkus)
	response := APIResponse{
		Data:       events,
		Pagination: pagination,
	}
	return c.JSON(http.StatusOK, response)
}

// --- HANDLER BARU: GET /profiles/:address ---
// Mengambil semua data untuk satu halaman profil pengguna
func (h *Handler) getUserProfile(c echo.Context) error {
	ctx := c.Request().Context()
	address := c.Param("address") // Ambil ':address' dari URL

	if address == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "address is required"})
	}

	// Ambil 'User' dan SEMUA relasinya dalam satu query
	user, err := h.DB.User.Query().
		Where(user.AddressEQ(address)).
		// Eager load semua data yang terkait dengan User ini
		WithMoments().      // Ambil 10 momen terakhir (contoh pagination)
		WithAccessories().  // Ambil 10 aksesoris terakhir
		WithEventPasses().  // Ambil 10 pass terakhir
		WithHostedEvents(). // Ambil 10 event yang di-host
		WithListings().     // Ambil 10 listing terakhir
		Only(ctx)

	if err != nil {
		if ent.IsNotFound(err) {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "User profile not found"})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	// (Anda bisa menambahkan pagination kustom untuk 'moments', 'accessories', dll.
	// di sini jika Anda tidak ingin 'eager load' semuanya)

	return c.JSON(http.StatusOK, APIResponse{Data: user})
}
