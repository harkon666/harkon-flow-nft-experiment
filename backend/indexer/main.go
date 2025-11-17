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
	NFTMomentMinted                = "A.1bb6b1e0a5170088.NFTMoment.Minted"
	NFTAccessoryMinted             = "A.1bb6b1e0a5170088.AccessoryPack.AccessoryDistributed"
	NFTMomentEquipAccessory        = "A.1bb6b1e0a5170088.NFTMoment.AccessoryEquipped"
	NFTMomentUnequipAccessory      = "A.1bb6b1e0a5170088.NFTMoment.AccessoryUnequipped"
	EventCreated                   = "A.1bb6b1e0a5170088.EventManager.EventCreated"
	ProfileUpdated                 = "A.1bb6b1e0a5170088.UserProfile.ProfileUpdated"
	UserRegisteredEvent            = "A.1bb6b1e0a5170088.EventManager.UserRegistered"
	UserCheckedInEvent             = "A.1bb6b1e0a5170088.EventManager.UserCheckedIn"
	EventPassMinted                = "A.1bb6b1e0a5170088.EventPass.Minted"
	ListingAvailable               = "A.2d55b98eb200daef.NFTStorefrontV2.ListingAvailable"
	ListingCompleted               = "A.2d55b98eb200daef.NFTStorefrontV2.ListingCompleted"
	NFTDeposited                   = "A.631e88ae7f1d7c20.NonFungibleToken.Deposited"
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
		grpc.TestnetHost,
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
		290706055,
		flow.EventFilter{
			EventTypes: []string{
				NFTMomentMinted, NFTAccessoryMinted, NFTMomentEquipAccessory, NFTMomentUnequipAccessory,
				FlowCapabilityControllerIssued, EventCreated, UserRegisteredEvent, UserCheckedInEvent,
				EventPassMinted, ProfileUpdated, ListingAvailable, NFTDeposited, ListingCompleted,
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
				case EventCreated:
					utils.EventCreated(ctx, ev, client)
				case UserRegisteredEvent:
					utils.UserRegistered(ctx, ev, client)
				case UserCheckedInEvent:
					utils.UserCheckedIn(ctx, ev, client)
				case EventPassMinted:
					utils.EventPassMinted(ctx, ev, client)
				case ProfileUpdated:
					utils.ProfileUpdated(ctx, ev, client)
				case ListingAvailable:
					utils.ListingAvailable(ctx, ev, client)
				case ListingCompleted:
					utils.ListingCompleted(ctx, ev, client)
				case NFTDeposited:
					utils.NFTDeposited(ctx, ev, client)
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
