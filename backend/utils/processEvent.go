package utils

import ( // Dibutuhkan jika Anda akan melakukan operasi DB
	"backend/ent"
	"backend/ent/attendance"
	"backend/ent/event"
	"backend/ent/eventpass"
	"backend/ent/listing"
	"backend/ent/nftaccessory"
	"backend/ent/nftmoment"
	"backend/ent/user"
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/onflow/cadence"
	"github.com/onflow/flow-go-sdk"
)

/**
 * getCadenceField[T cadence.Value] adalah fungsi generik
 * untuk mengambil dan memvalidasi tipe field dari event.
 *
 * T adalah tipe Cadence yang diharapkan (misal: cadence.Address)
 */
func getCadenceField[T cadence.Value](fields map[string]cadence.Value, key string) (T, error) {
	// 1. Ambil dari map
	fieldValue, ok := fields[key]
	if !ok {
		// 'var zero T' adalah cara untuk mendapatkan "tipe default"
		// dari T (misal: nil) agar kita bisa mengembalikannya.
		var zero T
		return zero, fmt.Errorf("field '%s' tidak ditemukan di event", key)
	}

	// 2. Lakukan type assertion ke tipe Generik 'T'
	typedValue, ok := fieldValue.(T)
	if !ok {
		var zero T
		// Error ini akan sangat jelas, misal:
		// "field 'brandAddress' bukan tipe cadence.Address (tipe: cadence.String)"
		return zero, fmt.Errorf("field '%s' bukan tipe yang diharapkan (tipe: %T)", key, fieldValue)
	}

	// 3. Kembalikan nilai yang sudah di-type-assert
	return typedValue, nil
}

func convertCadenceDictToGoMap(dict cadence.Dictionary) map[string]string {
	goMap := make(map[string]string)
	for _, pair := range dict.Pairs {
		// Asumsi Key dan Value adalah cadence.String
		key, okK := pair.Key.(cadence.String)
		value, okV := pair.Value.(cadence.String)
		if okK && okV {
			goMap[key.String()] = value.String()
		}
	}
	return goMap
}

// convertCadenceArrayToUint64Slice mengkonversi [(UInt64)?] Cadence ke []uint64 Go
func convertCadenceArrayToUint64Slice(arr cadence.Array) []uint64 {
	goSlice := make([]uint64, 0, len(arr.Values)) // Alokasikan 'slice' dengan kapasitas
	for _, val := range arr.Values {
		// 'val' adalah tipe 'cadence.Optional'
		if optVal, ok := val.(cadence.Optional); ok {
			// Cek apakah 'Optional'-nya tidak 'nil'
			if optVal.Value != nil {
				// Konversi isinya ke 'cadence.UInt64'
				if concreteVal, ok := optVal.Value.(cadence.UInt64); ok {
					goSlice = append(goSlice, uint64(concreteVal))
				}
			}
		}
	}
	return goSlice
}

func HandleCapabilityIssued(ctx context.Context, ev flow.Event, client *ent.Client) {

	// Dapatkan semua field dari event
	Fields := ev.Value.FieldsMappedByName()
	cadenceTypeString := Fields["type"].String()

	if !strings.Contains(cadenceTypeString, "&A.1bb6b1e0a5170088.UserProfile.Profile") {
		return
	}

	// 4. CEK ANDA: Apakah ini event untuk UserProfile?
	//    Kita cek apakah string-nya mengandung ".UserProfile."

	log.Println("Event UserProfile terdeteksi. Memproses...")

	// 5. Ambil & Parse 'address' field
	ownerAddressCadence, err := getCadenceField[cadence.Address](Fields, "address")
	if err != nil {
		log.Println("error taking address cadence")
	}

	// Dapatkan alamat sebagai string (misal: "0xf8d6e0586b0a20c7")
	userAddress := ownerAddressCadence.String()

	// 6. Pola "Get-or-Create" (Sangat Penting)
	// Coba cari user dulu
	existingUser, err := client.User.Query().
		Where(user.AddressEQ(userAddress)).
		Only(ctx)

	if err != nil {
		// Jika error-nya adalah "Not Found" (user baru)
		if ent.IsNotFound(err) {
			log.Printf("User baru terdeteksi: %s. Menyimpan ke database...", userAddress)

			// Buat user baru
			_, createErr := client.User.Create().
				SetAddress(userAddress).
				Save(ctx)

			if createErr != nil {
				log.Printf("Gagal menyimpan user baru %s: %v", userAddress, createErr)
			} else {
				log.Println("User baru berhasil disimpan.")
			}

			// Jika error-nya BUKAN "Not Found" (masalah DB lain)
		} else {
			log.Printf("Error saat query user %s: %v", userAddress, err)
		}

		// Jika tidak ada error (err == nil)
	} else {
		log.Printf("User %s sudah ada di database. (ID: %d)", existingUser.Address, existingUser.ID)
	}
}

func NFTMomentMinted(ctx context.Context, ev flow.Event, client *ent.Client) {
	var Fields = ev.Value.FieldsMappedByName()
	ownerAddressCadence, err := getCadenceField[cadence.Address](Fields, "recipient")
	idNftCadence, _ := getCadenceField[cadence.UInt64](Fields, "id")
	nameCadence, _ := getCadenceField[cadence.String](Fields, "name")
	descriptionCadence, _ := getCadenceField[cadence.String](Fields, "description")
	thumbnailCadence, _ := getCadenceField[cadence.String](Fields, "thumbnail")
	if err != nil {
		log.Println("Gagal parsing brandAddress:", err)
		return
	}
	ownerAddress := ownerAddressCadence.String()
	name := nameCadence.String()
	description := descriptionCadence.String()
	thumbnail := thumbnailCadence.String()

	isUserFound, err := client.User.Query().
		Where(
			user.AddressEQ(ownerAddress), // Gunakan predikat 'AddressEQ'
		).
		Only(ctx)
	if err != nil {
		log.Println("please setup moment collection")
	} else {
		log.Println("User found", isUserFound)
		nftMinted, err := client.NFTMoment.Create().
			SetName(name).
			SetDescription(description).
			SetThumbnail(thumbnail).
			SetNftID(uint64(idNftCadence)).
			SetOwnerID(isUserFound.ID).
			Save(ctx)

		if err != nil {
			log.Println("error when create insert NFT")
		} else {
			log.Println("nft minted", nftMinted)
		}
	}
}

func NFTAccessoryMinted(ctx context.Context, ev flow.Event, client *ent.Client) {
	var Fields = ev.Value.FieldsMappedByName()
	ownerAddressCadence, err := getCadenceField[cadence.Address](Fields, "recipient")
	idNftCadence, _ := getCadenceField[cadence.UInt64](Fields, "id")
	nameCadence, _ := getCadenceField[cadence.String](Fields, "name")
	descriptionCadence, _ := getCadenceField[cadence.String](Fields, "description")
	thumbnailCadence, _ := getCadenceField[cadence.String](Fields, "thumbnail")
	equipmentTypeCadence, _ := getCadenceField[cadence.String](Fields, "equipmentType")
	if err != nil {
		log.Println("failed parsing:", err)
		return
	}
	ownerAddress := ownerAddressCadence.String()
	name := nameCadence.String()
	description := descriptionCadence.String()
	thumbnail := thumbnailCadence.String()

	isUserFound, err := client.User.Query().
		Where(
			user.AddressEQ(ownerAddress), // Gunakan predikat 'AddressEQ'
		).
		Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			log.Println("please setup accessory collection")
		} else {
			log.Println("Error query")
		}

	} else {
		log.Println("User found", isUserFound)
		nftMinted, err := client.NFTAccessory.Create().
			SetName(name).
			SetDescription(description).
			SetThumbnail(thumbnail).
			SetNftID(uint64(idNftCadence)).
			SetOwnerID(isUserFound.ID).
			SetEquipmentType(equipmentTypeCadence.String()).
			Save(ctx)

		if err != nil {
			log.Println("error when create insert NFT")
		} else {
			log.Println("nft minted", nftMinted)
		}
	}

}

func NFTMomentEquipAccessory(ctx context.Context, ev flow.Event, client *ent.Client) {
	var Fields = ev.Value.FieldsMappedByName()
	nftAccessoryIdCadence, err := getCadenceField[cadence.Optional](Fields, "NftAccessoryId")
	nftMomentIdNftCadence, _ := getCadenceField[cadence.UInt64](Fields, "NftMomentId")
	prevNFTAccessoryIdCadence, _ := getCadenceField[cadence.Optional](Fields, "prevNFTAccessoryId")

	if err != nil {
		log.Println("failed parsing :", err)
		return
	}
	nftMoment, err := client.NFTMoment.Query().
		Where(
			nftmoment.NftIDEQ(uint64(nftMomentIdNftCadence)), // Gunakan predikat 'AddressEQ'
		).
		Only(ctx)
	if err != nil {
		log.Println("nftmoment not found")
	}

	client.NFTAccessory.Update().Where(
		nftaccessory.NftIDEQ(uint64(nftAccessoryIdCadence.Value.(cadence.UInt64))),
	).SetEquippedOnMoment(nftMoment).Save(ctx)

	if prevNFTAccessoryIdCadence.Value != nil {
		client.NFTAccessory.Update().Where(
			nftaccessory.NftIDEQ(uint64(prevNFTAccessoryIdCadence.Value.(cadence.UInt64))),
		).ClearEquippedOnMoment().Save(ctx)
		log.Println("success unequip accessory", uint64(prevNFTAccessoryIdCadence.Value.(cadence.UInt64)))
	}

	log.Println("success equip accessory", uint64(nftAccessoryIdCadence.Value.(cadence.UInt64)))
}

func NFTMomentUnequipAccessory(ctx context.Context, ev flow.Event, client *ent.Client) {
	var Fields = ev.Value.FieldsMappedByName()
	nftAccessoryIdCadence, err := getCadenceField[cadence.Optional](Fields, "NftAccessoryId")

	if err != nil {
		log.Println("nftmoment not found")
	}
	client.NFTAccessory.Update().Where(
		nftaccessory.NftIDEQ(uint64(nftAccessoryIdCadence.Value.(cadence.UInt64))),
	).ClearEquippedOnMoment().Save(ctx)
	log.Println("success unequip accessory")
}

func EventCreated(ctx context.Context, ev flow.Event, client *ent.Client) {

	// --- 1. Parsing Semua Field Event ---
	var Fields = ev.Value.FieldsMappedByName()

	hostAddressCadence, err := getCadenceField[cadence.Address](Fields, "hostAddress")
	eventIDCadence, _ := getCadenceField[cadence.UInt64](Fields, "eventID")
	if err != nil {
		log.Println("Gagal parsing hostAddress:", err)
		return
	}
	eventNameCadence, err := getCadenceField[cadence.String](Fields, "eventName")
	if err != nil {
		log.Println("Gagal parsing eventName:", err)
		return
	}
	descriptionCadence, _ := getCadenceField[cadence.String](Fields, "description")
	thumbnailURLCadence, _ := getCadenceField[cadence.String](Fields, "thumbnailURL")
	eventTypeCadence, _ := getCadenceField[cadence.UInt8](Fields, "eventType")
	locationCadence, _ := getCadenceField[cadence.String](Fields, "location")
	latCadence, _ := getCadenceField[cadence.Fix64](Fields, "lat")
	longCadence, _ := getCadenceField[cadence.Fix64](Fields, "long")
	quotaCadence, _ := getCadenceField[cadence.UInt64](Fields, "quota")
	startDateCadence, _ := getCadenceField[cadence.UFix64](Fields, "startDate")
	endDateCadence, _ := getCadenceField[cadence.UFix64](Fields, "endDate")

	// --- 2. Konversi Tipe Cadence ke Tipe Go ---

	// Alamat
	hostAddress := hostAddressCadence.String()
	eventID := uint64(eventIDCadence)
	// String
	eventName := eventNameCadence.String()
	description := descriptionCadence.String()
	thumbnailURL := thumbnailURLCadence.String()
	location := locationCadence.String()
	// Angka
	eventType := uint8(eventTypeCadence)
	quota := uint64(quotaCadence)
	// Fix64 ke Float64
	latFloat, _ := strconv.ParseFloat(latCadence.String(), 64)
	longFloat, _ := strconv.ParseFloat(longCadence.String(), 64)
	// UFix64 (Timestamp) ke time.Time
	startDateInt, _ := strconv.ParseInt(strings.Split(startDateCadence.String(), ".")[0], 10, 64)
	startDate := time.Unix(startDateInt, 0)
	endDateInt, _ := strconv.ParseInt(strings.Split(endDateCadence.String(), ".")[0], 10, 64)
	endDate := time.Unix(endDateInt, 0)

	// --- 3. Cari atau Buat Host (User) ---
	hostUser, err := client.User.Query().
		Where(
			user.AddressEQ(hostAddress),
		).
		Only(ctx)

	// Jika User (Host) tidak ditemukan, buat baru
	if err != nil {
		if ent.IsNotFound(err) {
			log.Println("Please create user profile")
		} else {
			// Error DB lain
			log.Printf("Error saat query host user %s: %v", hostAddress, err)
			return
		}
	}

	// --- 5. Simpan Event Baru ke Database ---

	// Cek dulu apakah event ini sudah kita indeks
	_, err = client.Event.Query().
		Where(event.EventIDEQ(eventID)).
		Only(ctx)

	// Jika 'err' BUKAN nil, berarti event belum ada (atau ada error lain)
	if err != nil {
		if ent.IsNotFound(err) {
			// Event belum ada, kita buat
			newEvent, createErr := client.Event.Create().
				SetEventID(eventID). // <-- Field unik Anda
				SetName(eventName).
				SetDescription(description).
				SetThumbnail(thumbnailURL).
				SetEventType(eventType).
				SetLocation(location).
				SetLat(latFloat).
				SetLong(longFloat).
				SetStartDate(startDate).
				SetEndDate(endDate).
				SetQuota(quota).
				SetHost(hostUser). // <-- Tautkan ke User (Host)
				Save(ctx)

			if createErr != nil {
				log.Printf("Gagal menyimpan event baru ID %d: %v", eventID, createErr)
			} else {
				log.Printf("Event baru berhasil di-indeks: %s (ID: %d)", newEvent.Name, newEvent.EventID)
			}

		} else {
			// Error DB lain
			log.Printf("Error saat query event ID %d: %v", eventID, err)
		}
	} else {
		// Jika err == nil, 'existingEvent' ditemukan
		log.Printf("Event ID %d sudah ada di database, dilewati.", eventID)
	}
}

// (Handler untuk event 'UserRegistered')
func UserRegistered(ctx context.Context, ev flow.Event, client *ent.Client) {
	var Fields = ev.Value.FieldsMappedByName()
	userAddress, err := getCadenceField[cadence.Address](Fields, "userAddress")
	eventID, _ := getCadenceField[cadence.UInt64](Fields, "eventID")
	// ... (parsing event untuk 'userAddress' dan 'eventID') ...
	if err != nil {
		log.Println("gagal parsing user address")
		return
	}
	// 1. Dapatkan 'User'
	user, err := client.User.Query().Where(user.AddressEQ(userAddress.String())).Only(ctx)
	// (handle 'IsNotFound')
	if ent.IsNotFound(err) {
		log.Println("please setup user profile")
		return
	}
	// 2. Dapatkan 'Event'
	event, err := client.Event.Query().Where(event.EventIDEQ(uint64(eventID))).Only(ctx)
	if err != nil {
		log.Println("Gagal menemukan event:", uint64(eventID))
		return
	}

	// 3. BUAT ENTRI 'ATTENDANCE' BARU
	// Ini adalah "lem" (perekat) yang menghubungkan keduanya
	_, err = client.Attendance.Create().
		SetUser(user).       // Tautkan ke User
		SetEvent(event).     // Tautkan ke Event
		SetCheckedIn(false). // Set status (sesuai kontrak Anda)
		Save(ctx)

	if err != nil {
		log.Println("Gagal menyimpan 'Attendance':", err)
	} else {
		log.Println("User", user.Address, "berhasil mendaftar ke", event.Name)
	}
}

func UserCheckedIn(ctx context.Context, ev flow.Event, client *ent.Client) {
	// --- 1. Parsing Event (Sama seperti 'Registered') ---
	var Fields = ev.Value.FieldsMappedByName()

	userAddressCadence, err := getCadenceField[cadence.Address](Fields, "userAddress")
	if err != nil {
		log.Println("gagal parsing user address:", err)
		return
	}
	eventIDCadence, err := getCadenceField[cadence.UInt64](Fields, "eventID")
	if err != nil {
		log.Println("gagal parsing eventID:", err)
		return
	}

	userAddress := userAddressCadence.String()
	eventID := uint64(eventIDCadence)

	// --- 2. Cari 'Attendance' Record yang Spesifik ---
	// Kita perlu mencari 'Attendance' yang menghubungkan User DAN Event ini.
	// Kita bisa menggunakan 'WhereHas' untuk memfilter berdasarkan relasi.

	attendanceRecord, err := client.Attendance.Query().
		Where(
			// Cari 'Attendance' yang...
			// ...memiliki 'event' di mana 'event_id' cocok
			attendance.HasEventWith(event.EventIDEQ(eventID)),
			// ...DAN memiliki 'user' di mana 'address' cocok
			attendance.HasUserWith(user.AddressEQ(userAddress)),
		).
		Only(ctx) // Kita harapkan hanya ada 1 hasil

	if err != nil {
		// Jika 'IsNotFound', berarti user ini tidak terdaftar
		// atau event/user tidak ada.
		if ent.IsNotFound(err) {
			log.Printf("Gagal check-in: 'Attendance' record tidak ditemukan untuk user %s di event %d. User harus register dulu.", userAddress, eventID)
		} else {
			// Error database lain
			log.Printf("Error query 'Attendance': %v", err)
		}
		return
	}

	// (Opsional) Cek apakah sudah check-in agar tidak kerja dua kali
	if attendanceRecord.CheckedIn {
		log.Printf("User %s sudah check-in ke event %d, dilewati.", userAddress, eventID)
		return
	}

	// --- 3. UPDATE 'Attendance' Record ---
	// Kita sudah dapat 'attendanceRecord', sekarang kita update
	_, err = attendanceRecord.Update().
		SetCheckedIn(true). // Set status menjadi 'true'
		Save(ctx)

	if err != nil {
		log.Printf("Gagal mengupdate 'Attendance' ke checked-in: %v", err)
	} else {
		log.Printf("User %s berhasil CHECK-IN ke event %d", userAddress, eventID)
	}
}

func EventPassMinted(ctx context.Context, ev flow.Event, client *ent.Client) {

	// --- 1. Parsing Event ---
	var Fields = ev.Value.FieldsMappedByName()

	// (ASUMSI ANDA SUDAH MEMPERBAIKI KONTRAK ANDA)
	recipientAddressCadence, err := getCadenceField[cadence.Address](Fields, "owner")
	if err != nil {
		log.Println("Gagal parsing 'owner' (recipient):", err)
		return
	}
	nameCadence, _ := getCadenceField[cadence.String](Fields, "name")
	descriptionCadence, _ := getCadenceField[cadence.String](Fields, "description")
	thumbnailCadence, _ := getCadenceField[cadence.String](Fields, "thumbnail")
	eventTypeCadence, _ := getCadenceField[cadence.UInt8](Fields, "event_type")

	// 'id' atau 'uuid' adalah ID unik dari pass SBT
	passIDCadence, err := getCadenceField[cadence.UInt64](Fields, "id")
	if err != nil {
		log.Println("Gagal parsing 'id' (passID):", err)
		return
	}

	// 'eventID' adalah ID dari 'EventManager'
	eventIDCadence, err := getCadenceField[cadence.UInt64](Fields, "eventID")
	if err != nil {
		log.Println("Gagal parsing 'eventID':", err)
		return
	}

	// --- 2. Konversi Tipe Go ---
	recipientAddress := recipientAddressCadence.String()
	passID := uint64(passIDCadence)
	eventID := uint64(eventIDCadence)

	// --- 3. Dapatkan Relasi (User & Event) ---

	// Dapatkan 'User' (Pemilik)
	ownerUser, err := client.User.Query().Where(user.AddressEQ(recipientAddress)).Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			log.Printf("User %s tidak ditemukan di DB. 'EventPass' akan diindeks tanpa pemilik.", recipientAddress)
			// (Alternatif: Buat user baru di sini, seperti di 'handleEventCreated')
		} else {
			log.Printf("Error query user %s: %v", recipientAddress, err)
		}
		return // Kita tidak bisa melanjutkan tanpa user
	}

	// Dapatkan 'Event' (Sumber)
	sourceEvent, err := client.Event.Query().Where(event.EventIDEQ(eventID)).Only(ctx)
	if err != nil {
		log.Printf("Gagal menemukan Event %d di DB. 'EventPass' akan diindeks tanpa tautan event.", eventID)
		return // Kita tidak bisa melanjutkan tanpa event
	}

	// --- 4. Buat (atau Cek) 'EventPass' ---

	// Cek dulu apakah 'EventPass' ini sudah ada
	_, err = client.EventPass.Query().
		Where(eventpass.PassIDEQ(passID)).
		Only(ctx)

	// Jika 'err' BUKAN nil (artinya 'Not Found' atau error lain)
	if err != nil {
		if ent.IsNotFound(err) {
			// Ini adalah alur yang baik (happy path), pass belum ada

			// Buat 'EventPass' baru
			newPass, createErr := client.EventPass.Create().
				SetPassID(passID).
				SetName(nameCadence.String()).
				SetDescription(descriptionCadence.String()).
				SetThumbnail(thumbnailCadence.String()).
				SetEventType(uint8(eventTypeCadence)).
				SetIsUsed(false).      // Set default
				SetOwner(ownerUser).   // <-- Tautkan ke User (Pemilik)
				SetEvent(sourceEvent). // <-- Tautkan ke Event (Sumber)
				Save(ctx)

			if createErr != nil {
				log.Printf("Gagal menyimpan 'EventPass' baru (ID: %d): %v", passID, createErr)
			} else {
				log.Printf("Berhasil mengindeks 'EventPass' baru (ID: %d) untuk user %s", newPass.PassID, ownerUser.Address)
			}

		} else {
			// Error database lain
			log.Printf("Error saat query EventPass %d: %v", passID, err)
		}
	} else {
		// Jika err == nil, berarti pass sudah ada
		log.Printf("EventPass (ID: %d) sudah ada di database, dilewati.", passID)
	}
}

func ProfileUpdated(ctx context.Context, ev flow.Event, client *ent.Client) {
	log.Println("Memproses event ProfileUpdated...")

	// --- 1. Parsing Event ---
	var Fields = ev.Value.FieldsMappedByName()

	// Ambil 'address' (Wajib)
	addressCadence, err := getCadenceField[cadence.Address](Fields, "address")
	if err != nil {
		log.Println("Gagal parsing 'address':", err)
		return
	}
	userAddress := addressCadence.String()

	// --- 2. Temukan User yang Akan Di-update ---
	// Event 'Updated' mengasumsikan 'User' sudah ada.
	user, err := client.User.Query().
		Where(user.AddressEQ(userAddress)).
		Only(ctx)

	if err != nil {
		// Jika user tidak ditemukan, ini adalah masalah (data tidak sinkron)
		log.Printf("Error: Menerima 'ProfileUpdated' untuk user %s yang tidak ada di DB: %v", userAddress, err)
		return
	}

	// --- 3. Buat 'Updater' ---
	// Kita akan membangun query 'update' secara bertahap
	updater := user.Update()

	// --- 4. Parsing & Set Field (Satu per Satu) ---

	// bio (String)
	if bioCadence, ok := Fields["bio"].(cadence.String); ok {
		updater.SetBio(bioCadence.String())
	}

	// nickname ((String)?)
	if nicknameCadence, ok := Fields["nickname"].(cadence.Optional); ok {
		if nicknameCadence.Value != nil {
			updater.SetNickname(nicknameCadence.Value.(cadence.String).String())
		}
	}

	// pfp ((String)?)
	if pfpCadence, ok := Fields["pfp"].(cadence.Optional); ok {
		if pfpCadence.Value != nil {
			updater.SetPfp(pfpCadence.Value.(cadence.String).String())
		}
	}

	// shortDescription ((String)?)
	if shortDescCadence, ok := Fields["shortDescription"].(cadence.Optional); ok {
		if shortDescCadence.Value != nil {
			updater.SetShortDescription(shortDescCadence.Value.(cadence.String).String())
		}
	}

	// bgImage ((String)?)
	if bgImageCadence, ok := Fields["bgImage"].(cadence.Optional); ok {
		if bgImageCadence.Value != nil {
			updater.SetBgImage(bgImageCadence.Value.(cadence.String).String())
		}
	}

	// socials ({String:String})
	// (Ini mengasumsikan Anda menambahkan 'field.JSON("socials", map[string]string{})' ke skema 'User' Anda)
	if socialsCadence, ok := Fields["socials"].(cadence.Dictionary); ok {
		goMap := convertCadenceDictToGoMap(socialsCadence)
		updater.SetSocials(goMap)
	}

	// highlightedEventPassIds ([(UInt64)?])
	// (Ini mengasumsikan Anda memiliki 'field.JSON("highlighted_event_pass_ids", []uint64{})' di skema 'User')
	if passIDsCadence, ok := Fields["highlightedEventPassIds"].(cadence.Array); ok {
		goSlice := convertCadenceArrayToUint64Slice(passIDsCadence)
		updater.SetHighlightedEventPassIds(goSlice)
	}

	// highlightedMomentID ((UInt64)?)
	// (Ini mengasumsikan Anda menambahkan 'field.Uint64("highlighted_moment_id").Optional()' ke skema 'User')
	if momentIDCadence, ok := Fields["highlightedMomentID"].(cadence.Optional); ok {
		if momentIDCadence.Value != nil {
			// 'ent' akan membuat 'SetHighlightedMomentID'
			updater.SetHighlightedMomentID(uint64(momentIDCadence.Value.(cadence.UInt64)))
		} else {
			// 'ent' akan membuat 'ClearHighlightedMomentID' atau 'SetHighlightedMomentIDNil'
			updater.ClearHighlightedMomentID()
		}
	}

	// --- 5. Jalankan Query Update ---
	_, err = updater.Save(ctx)
	if err != nil {
		log.Printf("Gagal mengupdate profil untuk user %s: %v", userAddress, err)
	} else {
		log.Printf("Berhasil mengupdate profil untuk user %s", userAddress)
	}
}

func ListingAvailable(ctx context.Context, ev flow.Event, client *ent.Client) {
	log.Println("Memproses event ListingAvailable...")

	// --- 1. Parsing Event ---
	var Fields = ev.Value.FieldsMappedByName()

	// Ambil semua field yang diperlukan (kita akan parse Tipe secara manual)
	listingIDCadence, err := getCadenceField[cadence.UInt64](Fields, "listingResourceID")
	if err != nil {
		log.Println("Gagal parsing 'listingResourceID':", err)
		return
	}
	nftIDCadence, err := getCadenceField[cadence.UInt64](Fields, "nftID")
	if err != nil {
		log.Println("Gagal parsing 'nftID':", err)
		return
	}
	sellerAddressCadence, err := getCadenceField[cadence.Address](Fields, "storefrontAddress")
	if err != nil {
		log.Println("Gagal parsing 'storefrontAddress':", err)
		return
	}
	priceCadence, _ := getCadenceField[cadence.UFix64](Fields, "salePrice")
	expiryCadence, _ := getCadenceField[cadence.UInt64](Fields, "expiry")

	// --- PERUBAHAN DI SINI: Ambil 'Type' sebagai String ---
	// Ambil field 'type' sebagai 'cadence.Value' mentah
	nftTypeField, ok := Fields["nftType"]
	if !ok {
		log.Println("Gagal parsing 'nftType': field tidak ada")
		return
	}
	// Panggil .String() di atasnya
	nftType := nftTypeField.String()

	vaultTypeField, ok := Fields["salePaymentVaultType"]
	if !ok {
		log.Println("Gagal parsing 'salePaymentVaultType': field tidak ada")
		return
	}
	// Panggil .String() di atasnya
	vaultType := vaultTypeField.String()
	// --- AKHIR PERUBAHAN ---

	// --- 2. Konversi Tipe Go ---
	listingID := uint64(listingIDCadence)
	nftID := uint64(nftIDCadence)
	sellerAddress := sellerAddressCadence.String()

	price, _ := strconv.ParseFloat(priceCadence.String(), 64)
	expiryTime := time.Unix(int64(expiryCadence), 0)

	// --- 3. Cek Duplikat ---
	// (Kode 'Cek Duplikat' Anda tetap sama)
	_, err = client.Listing.Query().
		Where(listing.ListingIDEQ(listingID)).
		Only(ctx)
	if err == nil {
		log.Printf("Listing ID %d sudah ada di database, dilewati.", listingID)
		return
	}
	if !ent.IsNotFound(err) {
		log.Printf("Error saat query Listing %d: %v", listingID, err)
		return
	}

	// --- 4. Dapatkan Relasi (Seller & NFT) ---

	// Dapatkan 'User' (Penjual)
	sellerUser, err := client.User.Query().Where(user.AddressEQ(sellerAddress)).Only(ctx)
	if err != nil {
		log.Printf("Gagal menemukan User (Penjual) %s: %v", sellerAddress, err)
		return
	}

	// --- PERUBAHAN DI SINI: Validasi Tipe NFT menggunakan 'strings.Contains' ---
	// (Ganti 'f8d...' dengan alamat Anda jika berbeda,
	//  tapi sebaiknya cek nama unik kontraknya saja)
	if !strings.Contains(nftType, ".NFTAccessory.") {
		log.Printf("Tipe NFT %s bukan 'NFTAccessory', dilewati.", nftType)
		return
	}
	// --- AKHIR PERUBAHAN ---

	nft, err := client.NFTAccessory.Query().Where(nftaccessory.NftIDEQ(nftID)).Only(ctx)
	if err != nil {
		log.Printf("Gagal menemukan NFTAccessory (ID: %d) di DB: %v", nftID, err)
		return
	}

	// --- 5. Buat 'Listing' Baru ---
	// (Kode 'Create' Anda tetap sama)
	newListing, createErr := client.Listing.Create().
		SetListingID(listingID).
		SetPrice(price).
		SetExpiry(expiryTime).
		SetPaymentVaultType(vaultType). // Simpan string tipe vault
		SetSeller(sellerUser).
		SetNftAccessory(nft).
		Save(ctx)

	if createErr != nil {
		log.Printf("Gagal menyimpan 'Listing' baru (ID: %d): %v", listingID, createErr)
	} else {
		log.Printf("Berhasil mengindeks 'Listing' baru (ID: %d) untuk NFT %d", newListing.ListingID, nft.NftID)
	}
}

func ListingCompleted(ctx context.Context, ev flow.Event, client *ent.Client) {
	log.Println("Memproses event ListingCompleted...")

	// --- 1. Parsing Event ---
	var Fields = ev.Value.FieldsMappedByName()

	// Ambil 'listingResourceID'
	listingIDCadence, err := getCadenceField[cadence.UInt64](Fields, "listingResourceID")
	if err != nil {
		log.Println("Gagal parsing 'listingResourceID':", err)
		return
	}

	// Ambil 'purchased'
	purchasedCadence, err := getCadenceField[cadence.Bool](Fields, "purchased")
	if err != nil {
		log.Println("Gagal parsing 'purchased':", err)
		return
	}

	// --- 2. Konversi Tipe Go ---
	listingID := uint64(listingIDCadence)
	wasPurchased := bool(purchasedCadence)

	// --- 3. Logika Bisnis ---

	// 'ListingCompleted' juga di-emit saat 'unlist' (membatalkan penjualan)
	// Kita hanya peduli jika 'purchased' adalah 'true'
	if !wasPurchased {
		log.Printf("Listing ID %d di-unlist (tidak dibeli), mengabaikan penghapusan.", listingID)
		// (Anda mungkin ingin 'handler' terpisah untuk 'unlist'
		//  jika Anda perlu memperbarui 'isListed' flag)
		return
	}

	// 4. Temukan 'Listing' di DB
	// (Ini menggunakan 'ListingIDEQ' dari skema 'Listing' Anda)
	listingRecord, err := client.Listing.Query().
		Where(listing.ListingIDEQ(listingID)).
		Only(ctx)

	if err != nil {
		if ent.IsNotFound(err) {
			log.Printf("Listing ID %d sudah dihapus, dilewati.", listingID)
		} else {
			log.Printf("Error query 'Listing' %d: %v", listingID, err)
		}
		return
	}

	// 5. HAPUS 'Listing' dari database
	err = client.Listing.DeleteOne(listingRecord).Exec(ctx)
	if err != nil {
		log.Printf("Gagal menghapus 'Listing' ID %d: %v", listingID, err)
	} else {
		log.Printf("Berhasil menghapus 'Listing' ID %d (terjual).", listingID)
	}
}

func NFTDeposited(ctx context.Context, ev flow.Event, client *ent.Client) {
	log.Println("Memproses event NonFungibleToken.Deposited...")

	// --- 1. Parsing Event ---
	var Fields = ev.Value.FieldsMappedByName()

	// Ambil ID NFT
	nftIDCadence, err := getCadenceField[cadence.UInt64](Fields, "id")
	if err != nil {
		log.Println("Gagal parsing 'id' (NFT ID):", err)
		return
	}

	// Ambil 'to' (Pemilik Baru)
	// Ini adalah opsional ((Address)?)
	recipientOptional, err := getCadenceField[cadence.Optional](Fields, "to")
	if err != nil || recipientOptional.Value == nil {
		log.Println("Gagal parsing 'to' (Recipient Address), atau 'to' adalah nil:", err)
		return // Kita tidak bisa update owner jika tidak tahu siapa 'to'
	}
	recipientAddressCadence := recipientOptional.Value.(cadence.Address)

	// Ambil Tipe NFT
	nftTypeField, ok := Fields["type"]
	if !ok {
		log.Println("Gagal parsing 'type' (NFT Type)")
		return
	}
	nftType := nftTypeField.String()

	// --- 2. Konversi Tipe Go ---
	nftID := uint64(nftIDCadence)
	newOwnerAddress := recipientAddressCadence.String()

	// --- 3. Dapatkan 'User' (Pemilik Baru) ---
	// Gunakan pola Get-or-Create
	newOwner, err := client.User.Query().Where(user.AddressEQ(newOwnerAddress)).Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			log.Printf("Pemilik baru %s tidak ditemukan")
		} else {
			log.Printf("Error query user %s: %v", newOwnerAddress, err)
			return
		}
	}

	// --- 4. Tentukan Tipe NFT & Update Owner ---

	// Cek apakah ini 'NFTAccessory'
	if strings.Contains(nftType, ".NFTAccessory.") {

		// Temukan Aksesori di DB
		accessory, err := client.NFTAccessory.Query().
			Where(nftaccessory.NftIDEQ(nftID)).
			Only(ctx)

		if err != nil {
			log.Printf("NFTAccessory ID %d tidak ditemukan di DB (mungkin ini mint?): %v", nftID, err)
			return // Minting harus ditangani oleh event 'Minted' Anda
		}

		// Update Owner
		_, err = accessory.Update().SetOwner(newOwner).Save(ctx)
		if err != nil {
			log.Printf("Gagal update owner untuk NFTAccessory %d: %v", nftID, err)
		} else {
			log.Printf("Berhasil transfer NFTAccessory %d ke %s", nftID, newOwnerAddress)
		}

		// Cek apakah ini 'NFTMoment'
	} else if strings.Contains(nftType, ".NFTMoment.") {

		// Temukan Momen di DB
		moment, err := client.NFTMoment.Query().
			Where(nftmoment.NftIDEQ(nftID)).
			Only(ctx)

		if err != nil {
			log.Printf("NFTMoment ID %d tidak ditemukan di DB (mungkin ini mint?): %v", nftID, err)
			return
		}

		// Update Owner
		_, err = moment.Update().SetOwner(newOwner).Save(ctx)
		if err != nil {
			log.Printf("Gagal update owner untuk NFTMoment %d: %v", nftID, err)
		} else {
			log.Printf("Berhasil transfer NFTMoment %d ke %s", nftID, newOwnerAddress)
		}
	}
	// (Abaikan jika bukan tipe NFT yang kita pedulikan)
}
