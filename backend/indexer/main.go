package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/onflow/flow-go-sdk"
	"github.com/onflow/flow-go-sdk/access/grpc"
	"google.golang.org/grpc/credentials/insecure"

	grpcOpts "google.golang.org/grpc"

	"backend/utils"

	_ "github.com/jackc/pgx/v5/stdlib"
)

const (
	ContractAddress = "1b7f070ebf7d0431" // Alamat tempat kontrak di-deploy
)

var (
	FlowCapabilityControllerIssued = "flow.StorageCapabilityControllerIssued"
	NFTMomentMinted                = "A.f8d6e0586b0a20c7.NFTMoment.Minted"
	NFTAccessoryMinted             = "A.f8d6e0586b0a20c7.AccessoryPack.AccessoryDistributed"
	NFTMomentEquipAccessory        = "A.f8d6e0586b0a20c7.NFTMoment.AccessoryEquipped"
	NFTMomentUnequipAccessory      = "A.f8d6e0586b0a20c7.NFTMoment.AccessoryUnequipped"
)

func main() {
	ctx := context.Background()
	// Load .env file if it exists (optional, environment variables can be set by Docker/system)
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: .env file not found, using environment variables from system:", err)
	}

	client := utils.Open(os.Getenv("DATABASE_URL"))
	if err := client.Schema.Create(ctx); err != nil {
		log.Fatal(err)
	}

	grpcClient, err := grpc.NewBaseClient(
		grpc.EmulatorHost,
		grpcOpts.WithTransportCredentials(insecure.NewCredentials()),
	)

	if err != nil {
		log.Println("Gagal terhubung ke emulator gRPC:", err)
	}

	grpcBlock, err := grpcClient.GetLatestBlockHeader(ctx, true)

	if err != nil {
		log.Println("Gagal gRPC get latest block:", err)
	}
	fmt.Println("Block ID:", grpcBlock.ID.String(), grpcBlock.Height)

	dataCh, errCh, initErr := grpcClient.SubscribeEventsByBlockHeight(
		ctx,
		0,
		flow.EventFilter{
			EventTypes: []string{
				NFTMomentMinted, NFTAccessoryMinted, NFTMomentEquipAccessory, NFTMomentUnequipAccessory,
				FlowCapabilityControllerIssued,
			},
		},
	)
	if initErr != nil {
		// handle init error
		log.Println("Gagal subscribe ke event init err:", initErr.Error())
	}

	for {
		select {
		case <-ctx.Done():
			// graceful shutdown
			return
		case data, ok := <-dataCh:
			if !ok {
				panic("data subscription closed")
			}
			for _, ev := range data.Events {
				fmt.Println("Type:", ev.Type)

				switch ev.Type {
				case FlowCapabilityControllerIssued:
					utils.HandleCapabilityIssued(ctx, ev, client)
				case NFTMomentMinted:
					utils.NFTMomentMinted(ctx, ev, client)
				case NFTAccessoryMinted:
					utils.NFTAccessoryMinted(ctx, ev, client)
				case NFTMomentEquipAccessory:
					utils.NFTMomentEquipAccessory(ctx, ev, client)
				case NFTMomentUnequipAccessory:
					utils.NFTMomentUnequipAccessory(ctx, ev, client)
				}
			}
		case err := <-errCh:
			if err != nil {
				// handle streaming error (log, reconnect / exponential back-off)
				log.Println("errorCh")
			}
		}
	}
}
