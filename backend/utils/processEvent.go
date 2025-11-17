package utils

import ( // Dibutuhkan jika Anda akan melakukan operasi DB
	"backend/ent"
	"backend/ent/nftaccessory"
	"backend/ent/nftmoment"
	"backend/ent/user"
	"context"
	"fmt"
	"log"
	"strings"

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

func HandleCapabilityIssued(ctx context.Context, ev flow.Event, client *ent.Client) {
	log.Println("Memproses event StorageCapabilityControllerIssued...")

	// Dapatkan semua field dari event
	Fields := ev.Value.FieldsMappedByName()
	cadenceTypeString := Fields["type"].String()

	if !strings.Contains(cadenceTypeString, "A.f8d6e0586b0a20c7.UserProfile.Profile") {
		return
	}

	// 4. CEK ANDA: Apakah ini event untuk UserProfile?
	//    Kita cek apakah string-nya mengandung ".UserProfile."

	log.Println("Event UserProfile terdeteksi. Memproses...")

	// 5. Ambil & Parse 'address' field
	ownerAddressCadence, err := getCadenceField[cadence.Address](Fields, "address")

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
		if ent.IsNotFound(err) {
			log.Println("User baru terdeteksi, membuat:", ownerAddress)
			user, err := client.User.Create().SetAddress(ownerAddress).Save(ctx)
			if err != nil {
				log.Println("failed creating user", err)
			} else {
				log.Println("user created", user)
				nftMinted, err := client.NFTMoment.Create().
					SetName(name).
					SetDescription(description).
					SetThumbnail(thumbnail).
					SetNftID(uint64(idNftCadence)).
					SetOwnerID(user.ID).
					Save(ctx)
				if err != nil {
					log.Println("error when create insert NFT")
				} else {
					log.Println("nft minted", nftMinted)
				}
			}
		} else {
			log.Println("Error query")
		}

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
			log.Println("User baru terdeteksi, membuat:", ownerAddress)
			user, err := client.User.Create().SetAddress(ownerAddress).Save(ctx)
			if err != nil {
				log.Println("failed creating user", err)
			} else {
				log.Println("user created", user)
				nftMinted, err := client.NFTAccessory.Create().
					SetName(name).
					SetDescription(description).
					SetThumbnail(thumbnail).
					SetNftID(uint64(idNftCadence)).
					SetOwnerID(user.ID).
					SetEquipmentType(equipmentTypeCadence.String()).
					Save(ctx)
				if err != nil {
					log.Println("error when create insert NFT")
				} else {
					log.Println("nft minted", nftMinted)
				}
			}
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
