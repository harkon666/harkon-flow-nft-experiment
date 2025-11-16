package transactions

import (
	// "backend/utils" // Asumsi dari file Anda sebelumnya (untuk WaitForSeal)
	// "context"
	"fmt"
	"strconv"

	// "log"
	// "os"

	// "github.com/joho/godotenv"
	"github.com/onflow/cadence"
	// "github.com/onflow/flow-go-sdk"
	// "github.com/onflow/flow-go-sdk/access"
	// "github.com/onflow/flow-go-sdk/access/http" // Menggunakan klien HTTP
	// "github.com/onflow/flow-go-sdk/crypto"
)

func MakeStrArg(s string) (cadence.String, error) {
	val, err := cadence.NewString(s)
	if err != nil {
		return "", fmt.Errorf("gagal membuat argumen boolean '%s': %w", s, err)
	}
	return val, nil
}

func MakeBoolArg(s string) (cadence.Bool, error) {
	boolVal, err := strconv.ParseBool(s)
	val := cadence.NewBool(boolVal)
	if err != nil {
		return false, fmt.Errorf("gagal membuat argumen boolean '%s': %w", s, err)
	}
	return val, nil
}

func MakeUInt64Arg(s string) (cadence.UInt64, error) {
	uint64Val, err := strconv.ParseUint(s, 10, 64)
	val := cadence.NewUInt64(uint64Val)
	if err != nil {
		return 0, fmt.Errorf("gagal membuat argumen UInt64 '%s': %w", s, err)
	}
	return val, nil
}

func MakeUInt8Arg(s string) (cadence.UInt8, error) {
	uintVal, err := strconv.ParseUint(s, 10, 64)
	val := cadence.NewUInt8(uint8(uintVal))
	if err != nil {
		return 0, fmt.Errorf("gagal membuat argumen UInt8 '%s': %w", s, err)
	}
	return val, nil
}
