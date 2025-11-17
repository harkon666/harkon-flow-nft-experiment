import "NonFungibleToken"
import "NFTAccessory"
import "ViewResolver"
import "MetadataViews"
import "EventPass"

access(all) contract NFTMoment: NonFungibleToken {
    // Path standar untuk menyimpan data
    access(all) let CollectionStoragePath: StoragePath
    access(all) let CollectionPublicPath: PublicPath
    access(all) let MinterStoragePath: StoragePath

    // Event standar
    access(all) event Withdraw(id: UInt64, from: Address?)
    access(all) event Deposit(id: UInt64, to: Address?)
    // Event kustom
    access(all) event Minted(recipient: Address, id: UInt64, name: String, description: String, thumbnail: String)
    access(all) event AccessoryEquipped(NftMomentId: UInt64, NftAccessoryId: UInt64?, prevNFTAccessoryId: UInt64?)
    access(all) event AccessoryUnequipped(NftMomentId: UInt64, NftAccessoryId: UInt64?)
    
    access(all) entitlement Equip
    //custom metadataview
    access(all) struct NFTMomentEquipmentMetadataView {
      access(all) let id: UInt64
      access(all) let equippedFrame: &NFTAccessory.NFT?
      init(
        equippedFrame: &NFTAccessory.NFT?,
        id: UInt64,
      ) {
        self.id = id
        self.equippedFrame = equippedFrame
      }
    }

    access(all) enum Tier: UInt8 {
        access(all) case community
        access(all)case pro
    }
    // 4. RESOURCE NFT
    // Ini adalah "benda" NFT Anda
    access(all) resource NFT: NonFungibleToken.NFT {
        access(all) let id: UInt64
        
        // --- Metadata Anda ---
        // Ini adalah data yang Anda simpan ON-CHAIN
        access(all) let name: String
        access(all) let description: String
        access(all) let thumbnail: String
        access(self) let metadata: {String: AnyStruct}
        access(all) var equippedFrame: @NFTAccessory.NFT?
        access(all) var tier: String

        init(
            name: String,
            description: String,
            thumbnail: String,
            metadata: {String: AnyStruct},
            tier: String
        ) {
            self.id = self.uuid // ID unik dibuat otomatis
            self.name = name
            self.description = description
            self.thumbnail = thumbnail
            self.metadata = metadata
            self.equippedFrame<-nil
            self.tier = tier
        }

        access(all) fun createEmptyCollection(): @{NonFungibleToken.Collection} {
            return <-NFTMoment.createEmptyCollection(nftType: Type<@NFTMoment.NFT>())
        }

        access(contract) fun equipFrame(frameNFT: @NFTAccessory.NFT): @NFTAccessory.NFT? {
          let accessoryNFT <- frameNFT
          if self.equippedFrame != nil {
            
            let prevEquippedAccessory <- self.equippedFrame <- accessoryNFT
            emit NFTMoment.AccessoryEquipped(NftMomentId: self.id, NftAccessoryId: self.equippedFrame?.id, prevNFTAccessoryId: prevEquippedAccessory?.id)
            return <-prevEquippedAccessory
          } else {

            let oldEqippedAccessories <- self.equippedFrame <-accessoryNFT
            destroy oldEqippedAccessories
            emit NFTMoment.AccessoryEquipped(NftMomentId: self.id, NftAccessoryId: self.equippedFrame?.id, prevNFTAccessoryId: nil)
            
            return nil
          }
        }

        access(contract) fun unequipFrame(): @NFTAccessory.NFT {
          pre {
            self.equippedFrame != nil: "no accessory equipped"
          }
          let unequippedAccessory <- self.equippedFrame <- nil
          emit NFTMoment.AccessoryUnequipped(NftMomentId: self.id, NftAccessoryId: unequippedAccessory?.id)
          return <- unequippedAccessory as! @NFTAccessory.NFT
        }

        access(all) view fun getViews(): [Type] {
            return [
                Type<MetadataViews.Display>(),
                Type<NFTMomentEquipmentMetadataView>(),
                Type<MetadataViews.ExternalURL>(),
                Type<MetadataViews.NFTCollectionData>(),
                Type<MetadataViews.NFTCollectionDisplay>(),
                Type<MetadataViews.Traits>()
            ]
        }

        // resolveView() adalah tempat Anda "menata" data Anda
        // agar sesuai dengan standar
        access(all) fun resolveView(_ view: Type): AnyStruct? {
            switch view {
                case Type<MetadataViews.Display>():
                    return MetadataViews.Display(
                        name: self.name,
                        description: self.description,
                        thumbnail: MetadataViews.HTTPFile(
                            url: self.thumbnail
                        ),
                    )
                case Type<NFTMomentEquipmentMetadataView>():
                    // Pinjam referensi ke 'equippedFrame' yang 
                    // tersimpan di resource NFTMoment ini.
                    // (Asumsi Anda punya: var equippedFrame: @NFTAccessories.NFT?)
                    let frameRef = &self.equippedFrame as &NFTAccessory.NFT?
                    
                    // Buat dan kembalikan struct kustom Anda
                    return NFTMomentEquipmentMetadataView(
                        equippedFrame: frameRef,
                        id: self.id,
                    )
                case Type<MetadataViews.Editions>():
                    // There is no max number of NFTs that can be minted from this contract
                    // so the max edition field value is set to nil
                    let editionInfo = MetadataViews.Edition(name: "Example NFT Edition", number: self.id, max: nil)
                    let editionList: [MetadataViews.Edition] = [editionInfo]
                    return MetadataViews.Editions(
                        editionList
                    )
                case Type<MetadataViews.Serial>():
                    return MetadataViews.Serial(
                        self.id
                    )
                case Type<MetadataViews.ExternalURL>():
                    return MetadataViews.ExternalURL("https://example-nft.onflow.org/".concat(self.id.toString()))
                case Type<MetadataViews.NFTCollectionData>():
                    return NFTMoment.resolveContractView(resourceType: Type<@NFTMoment.NFT>(), viewType: Type<MetadataViews.NFTCollectionData>())
                case Type<MetadataViews.NFTCollectionDisplay>():
                    return NFTMoment.resolveContractView(resourceType: Type<@NFTMoment.NFT>(), viewType: Type<MetadataViews.NFTCollectionDisplay>())
                case Type<MetadataViews.Traits>():
                    // exclude mintedTime and foo to show other uses of Traits
                    let excludedTraits = ["mintedTime", "frame"]
                    let traitsView = MetadataViews.dictToTraits(dict: self.metadata, excludedNames: excludedTraits)

                    // mintedTime is a unix timestamp, we should mark it with a displayType so platforms know how to show it.
                    let mintedTimeTrait = MetadataViews.Trait(name: "mintedTime", value: self.metadata["mintedTime"]!, displayType: "Date", rarity: nil)
                    traitsView.addTrait(mintedTimeTrait)

                    let frameRef: &NFTAccessory.NFT? = &self.equippedFrame
                    var frameTraitRarity: MetadataViews.Rarity? = nil
                    let frameRawRarityValue: AnyStruct?? = frameRef?.getMetaDataByField(field: "rarity")
                    let frameRarity = NFTMoment.unwrapRarity(rawRarityValue: frameRawRarityValue)
                    
                    let frameTrait = MetadataViews.Trait(name: "frame", value: frameRef?.name, displayType: nil, rarity: frameRarity)
                    traitsView.addTrait(frameTrait)
                    // if frameRef != nil {
                    //   let frameTraitRarity = frameRef.
                    // } 
                    // foo is a trait with its own rarity
                    // let fooTraitRarity = MetadataViews.Rarity(score: 10.0, max: 100.0, description: "Common")
                    // let fooTrait = MetadataViews.Trait(name: "foo", value: self.metadata["foo"], displayType: nil, rarity: fooTraitRarity)
                    // traitsView.addTrait(fooTrait)

                    return traitsView
                case Type<MetadataViews.EVMBridgedMetadata>():
                    // Implementing this view gives the project control over how the bridged NFT is represented as an
                    // ERC721 when bridged to EVM on Flow via the public infrastructure bridge.
                    // NOTE: If your NFT is a cross-VM NFT, meaning you control both your Cadence & EVM contracts and
                    //      registered your custom association with the VM bridge, it's recommended you use the 
                    //      CrossVMMetadata.EVMBytesMetadata view to define and pass metadata as EVMBytes into your
                    //      EVM contract at the time of bridging into EVM. For more information about cross-VM NFTs,
                    //      see FLIP-318: https://github.com/onflow/flips/issues/318

                    // Get the contract-level name and symbol values
                    let contractLevel = NFTMoment.resolveContractView(
                            resourceType: nil,
                            viewType: Type<MetadataViews.EVMBridgedMetadata>()
                        ) as! MetadataViews.EVMBridgedMetadata?

                    if let contractMetadata = contractLevel {
                        // Compose the token-level URI based on a base URI and the token ID, pointing to a JSON file. This
                        // would be a file you've uploaded and are hosting somewhere - in this case HTTP, but this could be
                        // IPFS, S3, a data URL containing the JSON directly, etc.
                        let baseURI = "https://example-nft.onflow.org/token-metadata/"
                        let uriValue = self.id.toString().concat(".json")

                        return MetadataViews.EVMBridgedMetadata(
                            name: contractMetadata.name,
                            symbol: contractMetadata.symbol,
                            uri: MetadataViews.URI(
                                baseURI: baseURI, // defining baseURI results in a concatenation of baseURI and value
                                value: self.id.toString().concat(".json")
                            )
                        )
                    } else {
                        return nil
                    }
            }
            return nil
        }
    }

    access(all) fun resolveContractView(resourceType: Type?, viewType: Type): AnyStruct? {
        switch viewType {
            case Type<MetadataViews.NFTCollectionData>():
                let collectionData = MetadataViews.NFTCollectionData(
                    storagePath: self.CollectionStoragePath,
                    publicPath: self.CollectionPublicPath,
                    publicCollection: Type<&NFTMoment.Collection>(),
                    publicLinkedType: Type<&NFTMoment.Collection>(),
                    createEmptyCollectionFunction: (fun(): @{NonFungibleToken.Collection} {
                        return <-NFTMoment.createEmptyCollection(nftType: Type<@NFTMoment.NFT>())
                    })
                )
                return collectionData
            case Type<MetadataViews.NFTCollectionDisplay>():
                let media = MetadataViews.Media(
                    file: MetadataViews.HTTPFile(
                        url: "https://assets.website-files.com/5f6294c0c7a8cdd643b1c820/5f6294c0c7a8cda55cb1c936_Flow_Wordmark.svg"
                    ),
                    mediaType: "image/svg+xml"
                )
                return MetadataViews.NFTCollectionDisplay(
                    name: "The Example Collection",
                    description: "This collection is used as an example to help you develop your next Flow NFT.",
                    externalURL: MetadataViews.ExternalURL("https://example-nft.onflow.org"),
                    squareImage: media,
                    bannerImage: media,
                    socials: {
                        "twitter": MetadataViews.ExternalURL("https://twitter.com/flow_blockchain")
                    }
                )
            case Type<MetadataViews.EVMBridgedMetadata>():
                // Implementing this view gives the project control over how the bridged NFT is represented as an ERC721
                // when bridged to EVM on Flow via the public infrastructure bridge.

                // Compose the contract-level URI. In this case, the contract metadata is located on some HTTP host,
                // but it could be IPFS, S3, a data URL containing the JSON directly, etc.
                return MetadataViews.EVMBridgedMetadata(
                    name: "NFTMoment",
                    symbol: "XMPL",
                    uri: MetadataViews.URI(
                        baseURI: nil, // setting baseURI as nil sets the given value as the uri field value
                        value: "https://example-nft.onflow.org/contract-metadata.json"
                    )
                )
        }
        return nil
    }

    access(all) resource Collection: NonFungibleToken.Collection {
        
        access(all) var ownedNFTs: @{UInt64: {NonFungibleToken.NFT}}
        access(all) var isUsedFreeMint: Bool 
        init() {
            self.ownedNFTs <- {}
            self.isUsedFreeMint = false
        }

        access(all) view fun getSupportedNFTTypes(): {Type: Bool} {
            let supportedTypes: {Type: Bool} = {}
            supportedTypes[Type<@NFTMoment.NFT>()] = true
            return supportedTypes
        }

        access(all) view fun isSupportedNFTType(type: Type): Bool {
            return type == Type<@NFTMoment.NFT>()
        }


        //in this function only can withdraw if there is no equipped accessories
        access(NonFungibleToken.Withdraw) fun withdraw(withdrawID: UInt64): @{NonFungibleToken.NFT} {
            let nft = self.borrowNFT(withdrawID) as! &NFTMoment.NFT
            assert(nft.equippedFrame == nil, message: "Please unequip your accessory")
            let token <- self.ownedNFTs.remove(key: withdrawID)
                ?? panic("NFTMoment.Collection.withdraw: Could not withdraw an NFT with ID "
                        .concat(withdrawID.toString())
                        .concat(". Check the submitted ID to make sure it is one that this collection owns."))
            return <-token
        }

        access(all) fun deposit(token: @{NonFungibleToken.NFT}) {
            let token <- token as! @NFTMoment.NFT
            let id = token.id

            // add the new token to the dictionary which removes the old one
            let oldToken <- self.ownedNFTs[token.id] <- token

            destroy oldToken
        }

        access(contract) fun useFreeMint() {
          pre {
            self.isUsedFreeMint == false: "Free Mint already used"
          }
          self.isUsedFreeMint = true
        }

        access(Equip) fun equipFrame(momentNFTID: UInt64, frameNFT: @NFTAccessory.NFT): @NFTAccessory.NFT? {
          pre {
            frameNFT.listingResouceId == nil: "frameNFT is listed for sale, please unlist frameNFT"
          }
          let nft: &NFTMoment.NFT = self.borrowNFT(momentNFTID) as! &NFTMoment.NFT
          return <-nft.equipFrame(frameNFT: <-frameNFT)
        }

        access(Equip) fun unequipFrame(momentNFTID: UInt64): @NFTAccessory.NFT {
          let nft: &NFTMoment.NFT = self.borrowNFT(momentNFTID) as! &NFTMoment.NFT
          return <- nft.unequipFrame()
        }

        access(all) view fun getIDs(): [UInt64] {
          return self.ownedNFTs.keys
        }

        access(all) view fun getLength(): Int {
          return self.ownedNFTs.length
        }

        access(all) view fun borrowNFT(_ id: UInt64): &{NonFungibleToken.NFT}? {
            return &self.ownedNFTs[id]
        }

        access(all) view fun borrowViewResolver(id: UInt64): &{ViewResolver.Resolver}? {
          if let nft = &self.ownedNFTs[id] as &{NonFungibleToken.NFT}? {
            return nft as &{ViewResolver.Resolver}
          }
          return nil
        }

        access(all) fun createEmptyCollection(): @{NonFungibleToken.Collection} {
            return <-NFTMoment.createEmptyCollection(nftType: Type<@NFTMoment.NFT>())
        }
    }

    access(all) fun unwrapRarity(rawRarityValue: AnyStruct??): MetadataViews.Rarity? {
      if let unwrappedOnce: AnyStruct? = rawRarityValue {

          if let rarity = unwrappedOnce as? MetadataViews.Rarity {
              
              return rarity
          }
      }
      return nil
    }

    access(all) fun createEmptyCollection(nftType: Type): @{NonFungibleToken.Collection} {
        return <- create Collection()
    }

    access(all) view fun getContractViews(resourceType: Type?): [Type] {
        return [
            Type<MetadataViews.NFTCollectionData>(),
            Type<MetadataViews.NFTCollectionDisplay>(),
            Type<MetadataViews.EVMBridgedMetadata>()
        ]
    }

    access(all) fun applyTier(_ tier: Tier): String {
      switch tier {
        case Tier.community:
          return "community"
        case Tier.pro:
          return "pro"
        default:
          return "not a tier"
      }
    }

    access(all) resource NFTMinter {

        access(all) fun freeMint(
            recipient: &NFTMoment.Collection,
            name: String,
            description: String,
            thumbnail: String,
        ) {
            let metadata: {String: AnyStruct} = {}
            let currentBlock = getCurrentBlock()
            let appliedTier = NFTMoment.applyTier(Tier(rawValue: 0)!)
            metadata["tier"] = appliedTier
            metadata["mintedBlock"] = currentBlock.height
            metadata["mintedTime"] = currentBlock.timestamp

            let newNFT <- create NFT(
                name: name,
                description: description,
                thumbnail: thumbnail,
                metadata: metadata,
                tier: appliedTier
            )

            let id = newNFT.id
            emit Minted(recipient: recipient.owner!.address, id: id, name: name, description: description, thumbnail: thumbnail)

            recipient.deposit(token: <-newNFT)
        }

        access(all) fun mintNFTWithEventPass(
            recipient: &NFTMoment.Collection,
            recipientPass: &EventPass.NFT,
            name: String,
            description: String,
            thumbnail: String,
            tier: UInt8
        ) {
            pre {
              recipient.isUsedFreeMint == false: "Free Mint already used"
            }
            recipientPass.useEventPass()

            let metadata: {String: AnyStruct} = {}
            let currentBlock = getCurrentBlock()
            let appliedTier = NFTMoment.applyTier(Tier(rawValue: tier)!)
            metadata["tier"] = appliedTier
            metadata["mintedBlock"] = currentBlock.height
            metadata["mintedTime"] = currentBlock.timestamp

            let newNFT <- create NFT(
                name: name,
                description: description,
                thumbnail: thumbnail,
                metadata: metadata,
                tier: appliedTier
            )

            let id = newNFT.id

            emit Minted(recipient: recipient.owner!.address, id: id, name: name, description: description, thumbnail: thumbnail)

            recipient.deposit(token: <-newNFT)
        }
    }

    init() {        
        self.CollectionStoragePath = /storage/NFTMomentCollection
        self.CollectionPublicPath = /public/NFTMomentReceiver
        self.MinterStoragePath = /storage/NFTMomentMinter

        let minter <- create NFTMinter()
        self.account.storage.save(<-minter, to: self.MinterStoragePath)
    }
}