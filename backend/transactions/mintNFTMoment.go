package transactions

import (
	"backend/utils" // Asumsi dari file Anda sebelumnya (untuk WaitForSeal)
	"context"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/onflow/cadence"
	"github.com/onflow/flow-go-sdk"
	"github.com/onflow/flow-go-sdk/access"
	"github.com/onflow/flow-go-sdk/access/http" // Menggunakan klien HTTP
	"github.com/onflow/flow-go-sdk/crypto"
)

// Ini adalah skrip transaksi minting Anda
// Saya telah menambahkan '0x%s' agar kita bisa menyuntikkan alamat kontrak
const mintNFTMomentScriptTemplate = `
import NonFungibleToken from 0x%s
import NFTMoment from 0x%s
import MetadataViews from 0x%s

transaction(
    recipient: Address,
    name: String,
    description: String,
    thumbnail: String
) {

    let minter: &NFTMoment.NFTMinter
    let recipientCollectionRef: &{NonFungibleToken.Receiver}

    prepare(signer: auth(BorrowValue) &Account) {

        let collectionData = NFTMoment.resolveContractView(resourceType: nil, viewType: Type<MetadataViews.NFTCollectionData>()) as! MetadataViews.NFTCollectionData?
            ?? panic("Could not resolve NFTCollectionData view.")
        
        // Pinjam Minter milik 'signer' (admin)
        self.minter = signer.storage.borrow<&NFTMoment.NFTMinter>(from: NFTMoment.MinterStoragePath)
            ?? panic("Minter resource not found in signer's storage")

        // Pinjam Receiver milik 'recipient'
        self.recipientCollectionRef = getAccount(recipient).capabilities.borrow<&{NonFungibleToken.Receiver}>(collectionData.publicPath)
            ?? panic("Recipient's Receiver capability not found")
    }

    execute {
        // Minter melakukan minting dan deposit
        self.minter.mintNFT(
            recipient: self.recipientCollectionRef,
            name: name,
            description: description,
            thumbnail: thumbnail
        )
    }
}
`

func MintNFTMoment(
	recipientAddressString string,
	name string,
	description string,
	thumbnail string,
) error {

	// Muat .env
	err := godotenv.Load()
	if err != nil {
		log.Println("Peringatan: Error loading .env file:", err)
	}

	ctx := context.Background()
	var flowClient access.Client

	// Koneksi Flow ke Emulator HTTP port
	flowClient, err = http.NewClient(http.EmulatorHost)
	if err != nil {
		return fmt.Errorf("gagal membuat flow client: %w", err)
	}

	// 1. SIAPKAN SIGNER (ADMIN/MINTER)
	privateKeyHex := os.Getenv("PRIVATE_KEY") // Ambil dari .env
	if privateKeyHex == "" {
		return fmt.Errorf("PRIVATE_KEY tidak ditemukan di environment variables")
	}

	// Gunakan alamat minter dari konstanta
	minterFlowAddress := flow.HexToAddress(deployerAddress)
	platformKey, err := crypto.DecodePrivateKeyHex(crypto.ECDSA_P256, privateKeyHex)
	if err != nil {
		return fmt.Errorf("gagal decode private key: %w", err)
	}

	platformAccount, err := flowClient.GetAccount(ctx, minterFlowAddress)
	if err != nil {
		return fmt.Errorf("gagal mendapatkan akun minter %s: %w", minterFlowAddress.String(), err)
	}

	// Asumsi kita menggunakan key pertama (index 0)
	key := platformAccount.Keys[0]
	signer, err := crypto.NewInMemorySigner(platformKey, key.HashAlgo)
	if err != nil {
		return fmt.Errorf("gagal memuat signer: %w", err)
	}

	// 2. BUAT SKRIP TRANSAKSI
	// Kita suntikkan alamat minter (yang juga alamat deployer) 2x
	// 1x untuk 'NFTMoment' dan 1x untuk 'MetadataViews'
	script := []byte(fmt.Sprintf(mintNFTMomentScriptTemplate, deployerAddress, deployerAddress, deployerAddress))

	// 3. SIAPKAN ARGUMEN (4 Argumen)

	// Helper function untuk argumen String (sama seperti template Anda)
	makeStrArg := func(s string) (cadence.String, error) {
		val, err := cadence.NewString(s)
		if err != nil {
			return "", fmt.Errorf("gagal membuat argumen string '%s': %w", s, err)
		}
		return val, nil
	}

	// --- Buat Argumen ---
	recipientAddressArg := cadence.NewAddress(flow.HexToAddress(recipientAddressString))

	nameArg, err := makeStrArg(name)
	if err != nil {
		return err
	}

	descriptionArg, err := makeStrArg(description)
	if err != nil {
		return err
	}

	thumbnailArg, err := makeStrArg(thumbnail)
	if err != nil {
		return err
	}

	// 4. BUAT TRANSAKSI
	latestBlock, err := flowClient.GetLatestBlock(ctx, true)
	if err != nil {
		return fmt.Errorf("gagal mendapatkan block terbaru: %w", err)
	}

	tx := flow.NewTransaction().
		SetScript(script).
		SetReferenceBlockID(latestBlock.ID).
		SetPayer(minterFlowAddress). // Admin adalah 'Payer'
		SetProposalKey(minterFlowAddress, key.Index, key.SequenceNumber).
		AddAuthorizer(minterFlowAddress) // Admin adalah 'Authorizer'

	// 5. TAMBAHKAN ARGUMEN
	_ = tx.AddArgument(recipientAddressArg)
	_ = tx.AddArgument(nameArg)
	_ = tx.AddArgument(descriptionArg)
	_ = tx.AddArgument(thumbnailArg)

	// 6. TANDA TANGANI TRANSAKSI
	err = tx.SignEnvelope(minterFlowAddress, key.Index, signer)
	if err != nil {
		return fmt.Errorf("gagal menandatangani transaksi: %w", err)
	}

	// 7. KIRIM TRANSAKSI
	log.Println("Mengirim transaksi 'mint_nft_moment'...")
	err = flowClient.SendTransaction(ctx, *tx)
	if err != nil {
		return fmt.Errorf("gagal mengirim transaksi: %w", err)
	}

	// 8. TUNGGU HASILNYA (SEAL)
	// (Menggunakan 'utils.WaitForSeal' dari template Anda)
	result, err := utils.WaitForSeal(ctx, flowClient, tx.ID())
	if err != nil {
		log.Printf("Transaksi %s gagal: %v\n", tx.ID(), err)
		return fmt.Errorf("transaksi %s gagal: %w", tx.ID(), err)
	}

	log.Printf("Transaksi Mint NFTMoment Berhasil! 🔥 Status: %s. TX ID: %s", result.Status, tx.ID())
	return nil // Sukses
}
