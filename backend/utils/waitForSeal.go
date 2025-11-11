package utils

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/onflow/flow-go-sdk"
	"github.com/onflow/flow-go-sdk/access"
)

func WaitForSeal(ctx context.Context, c access.Client, id flow.Identifier) (*flow.TransactionResult, error) {
	log.Println("Menunggu transaksi %s di-seal...\n", id)

	// Tentukan timeout agar tidak menunggu selamanya
	// 60 detik adalah waktu yang wajar untuk Testnet
	timeout := time.After(60 * time.Second)
	// Tentukan ticker untuk polling (cek setiap 1-2 detik)
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-timeout:
			// Jika timeout, kembalikan error
			return nil, fmt.Errorf("timeout (60s) menunggu transaksi %s", id.String())

		case <-ctx.Done():
			// Jika context API dibatalkan (misal: user menutup request)
			return nil, ctx.Err()

		case <-ticker.C:
			// Setiap 2 detik, kita cek status
			result, err := c.GetTransactionResult(ctx, id)

			// 1. Cek error Go (Network/Not Found)
			if err != nil {
				// Ini adalah 'Flow resource not found'
				// JANGAN KEMBALIKAN ERROR, kita anggap ini sementara
				log.Println("... (Menunggu tx %s diketahui jaringan: %v)", id.String(), err)
				// Lanjutkan ke iterasi loop berikutnya (coba lagi nanti)
				continue
			}

			// 2. Cek error Cadence (Fatal, transaksi gagal di-seal)
			if result.Error != nil {
				log.Println("Transaksi %s GAGAL di-seal (Error Cadence): %v\n", id, result.Error)
				return result, fmt.Errorf("transaksi gagal di chain: %w", result.Error)
			}

			// 3. Cek Status
			if result.Status == flow.TransactionStatusSealed {
				log.Println("\nTransaksi %s BERHASIL di-seal! Status: %s\n", id.String(), result.Status)
				return result, nil // SUKSES
			}

			log.Println("... (Status tx %s: %s)", id.String(), result.Status)
		}
	}
}
