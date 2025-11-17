/*
*
*  This is an example implementation of a Flow Non-Fungible Token
*  using the V2 standard.
*  It is not part of the official standard but it assumed to be
*  similar to how many NFTs would implement the core functionality.
*
*  This contract does not implement any sophisticated classification
*  system for its NFTs. It defines a simple NFT with minimal metadata.
*
*/

import "NonFungibleToken"
import "ViewResolver"
import "MetadataViews"

access(all) contract NFTAccessory: NonFungibleToken {

    /// Standard Paths
    access(all) let CollectionStoragePath: StoragePath
    access(all) let CollectionPublicPath: PublicPath

    /// Path where the minter should be stored
    /// The standard paths for the collection are stored in the collection resource type
    access(all) let MinterStoragePath: StoragePath

    //entitlement for sale
    access(all) entitlement Sale

    /// Event to show when an NFT is minted
    access(all) event Minted(
        type: String,
        id: UInt64,
        uuid: UInt64,
        minterAddress: Address?,
        minterUUID: UInt64,
        name: String,
        description: String,
        thumbnail: String,
    )

    /// We choose the name NFT here, but this type can have any name now
    /// because the interface does not require it to have a specific name any more
    access(all) resource NFT: NonFungibleToken.NFT {

        access(all) let id: UInt64

        /// From the Display metadata view
        access(all) let name: String
        access(all) let description: String
        access(all) let thumbnail: String
        access(all) let equipmentType: String
        access(all) var listingResouceId: UInt64?

        /// Generic dictionary of traits the NFT has
        access(self) let metadata: {String: AnyStruct}

        init(
            name: String,
            description: String,
            thumbnail: String,
            equipmentType: String,
            metadata: {String: AnyStruct},
        ) {
            self.id = self.uuid
            self.name = name
            self.description = description
            self.thumbnail = thumbnail
            self.metadata = metadata
            self.equipmentType = equipmentType
            self.listingResouceId = nil
        }

        /// createEmptyCollection creates an empty Collection
        /// and returns it to the caller so that they can own NFTs
        /// @{NonFungibleToken.Collection}
        access(all) fun createEmptyCollection(): @{NonFungibleToken.Collection} {
            return <-NFTAccessory.createEmptyCollection(nftType: Type<@NFTAccessory.NFT>())
        }

        access(contract) fun itemListed(_ listingResouceId: UInt64) {
          self.listingResouceId = listingResouceId
        }

        access(contract) fun itemUnlisted() {
          self.listingResouceId = nil
        }

        access(all) view fun getMetaDataByField(field: String): AnyStruct? {
          // Don't force a revert if the playID or field is invalid
          if let field = self.metadata[field] {
            return field
          } else {
            return nil
          }
        }

        access(all) view fun getViews(): [Type] {
            return [
                Type<MetadataViews.Display>(),
                Type<MetadataViews.ExternalURL>(),
                Type<MetadataViews.NFTCollectionData>(),
                Type<MetadataViews.NFTCollectionDisplay>(),
                Type<MetadataViews.Traits>()
            ]
        }

        access(all) fun resolveView(_ view: Type): AnyStruct? {
            switch view {
                case Type<MetadataViews.Display>():
                    return MetadataViews.Display(
                        name: self.name,
                        description: self.description,
                        thumbnail: MetadataViews.HTTPFile(
                            url: self.thumbnail
                        )
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
                    return NFTAccessory.resolveContractView(resourceType: Type<@NFTAccessory.NFT>(), viewType: Type<MetadataViews.NFTCollectionData>())
                case Type<MetadataViews.NFTCollectionDisplay>():
                    return NFTAccessory.resolveContractView(resourceType: Type<@NFTAccessory.NFT>(), viewType: Type<MetadataViews.NFTCollectionDisplay>())
                case Type<MetadataViews.Traits>():
                    // exclude mintedTime and foo to show other uses of Traits
                    let excludedTraits = ["mintedTime"]
                    let traitsView = MetadataViews.dictToTraits(dict: self.metadata, excludedNames: excludedTraits)

                    // mintedTime is a unix timestamp, we should mark it with a displayType so platforms know how to show it.
                    let mintedTimeTrait = MetadataViews.Trait(name: "mintedTime", value: self.metadata["mintedTime"]!, displayType: "Date", rarity: nil)
                    traitsView.addTrait(mintedTimeTrait)

                    //just for quick check from metadata will be removed soon
                    let listingResourceID = MetadataViews.Trait(name: "listingResouceId", value: self.listingResouceId, displayType: "Number", rarity: nil)
                    traitsView.addTrait(listingResourceID)
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
                    let contractLevel = NFTAccessory.resolveContractView(
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

    access(all) resource Collection: NonFungibleToken.Collection {
        /// dictionary of NFT conforming tokens
        /// NFT is a resource type with an `UInt64` ID field
        access(all) var ownedNFTs: @{UInt64: {NonFungibleToken.NFT}}

        init () {
            self.ownedNFTs <- {}
        }

        /// getSupportedNFTTypes returns a list of NFT types that this receiver accepts
        access(all) view fun getSupportedNFTTypes(): {Type: Bool} {
            let supportedTypes: {Type: Bool} = {}
            supportedTypes[Type<@NFTAccessory.NFT>()] = true
            return supportedTypes
        }

        /// Returns whether or not the given type is accepted by the collection
        /// A collection that can accept any type should just return true by default
        access(all) view fun isSupportedNFTType(type: Type): Bool {
            return type == Type<@NFTAccessory.NFT>()
        }

        /// withdraw removes an NFT from the collection and moves it to the caller
        access(NonFungibleToken.Withdraw) fun withdraw(withdrawID: UInt64): @{NonFungibleToken.NFT} {
            let token <- self.ownedNFTs.remove(key: withdrawID)
                ?? panic("NFTAccessory.Collection.withdraw: Could not withdraw an NFT with ID "
                        .concat(withdrawID.toString())
                        .concat(". Check the submitted ID to make sure it is one that this collection owns."))

            return <-token
        }

        access(Sale) fun itemListed(_ listingResouceId: UInt64, nftID: UInt64) {
          let nft = self.borrowNFT(nftID) as! &NFTAccessory.NFT
          nft.itemListed(listingResouceId)
        }

        access(Sale) fun itemUnlisted(nftID: UInt64) {
          let nft = self.borrowNFT(nftID) as! &NFTAccessory.NFT
          nft.itemUnlisted()
        }

        /// deposit takes a NFT and adds it to the collections dictionary
        /// and adds the ID to the id array
        access(all) fun deposit(token: @{NonFungibleToken.NFT}) {
            let token <- token as! @NFTAccessory.NFT
            let id = token.id

            // add the new token to the dictionary which removes the old one
            let oldToken <- self.ownedNFTs[token.id] <- token

            destroy oldToken

            // This code is for testing purposes only
            // Do not add to your contract unless you have a specific
            // reason to want to emit the NFTUpdated event somewhere
            // in your contract
            let authTokenRef = (&self.ownedNFTs[id] as auth(NonFungibleToken.Update) &{NonFungibleToken.NFT}?)!
            //authTokenRef.updateTransferDate(date: getCurrentBlock().timestamp)
            NFTAccessory.emitNFTUpdated(authTokenRef)
        }

        /// getIDs returns an array of the IDs that are in the collection
        access(all) view fun getIDs(): [UInt64] {
            return self.ownedNFTs.keys
        }

        /// Gets the amount of NFTs stored in the collection
        access(all) view fun getLength(): Int {
            return self.ownedNFTs.length
        }

        access(all) view fun borrowNFT(_ id: UInt64): &{NonFungibleToken.NFT}? {
            return &self.ownedNFTs[id]
        }

        /// Borrow the view resolver for the specified NFT ID
        access(all) view fun borrowViewResolver(id: UInt64): &{ViewResolver.Resolver}? {
            if let nft = &self.ownedNFTs[id] as &{NonFungibleToken.NFT}? {
                return nft as &{ViewResolver.Resolver}
            }
            return nil
        }

        /// createEmptyCollection creates an empty Collection of the same type
        /// and returns it to the caller
        /// @return A an empty collection of the same type
        access(all) fun createEmptyCollection(): @{NonFungibleToken.Collection} {
            return <-NFTAccessory.createEmptyCollection(nftType: Type<@NFTAccessory.NFT>())
        }
    }

    /// createEmptyCollection creates an empty Collection for the specified NFT type
    /// and returns it to the caller so that they can own NFTs
    access(all) fun createEmptyCollection(nftType: Type): @{NonFungibleToken.Collection} {
        return <- create Collection()
    }

    /// Function that returns all the Metadata Views implemented by a Non Fungible Token
    ///
    /// @return An array of Types defining the implemented views. This value will be used by
    ///         developers to know which parameter to pass to the resolveView() method.
    ///
    access(all) view fun getContractViews(resourceType: Type?): [Type] {
        return [
            Type<MetadataViews.NFTCollectionData>(),
            Type<MetadataViews.NFTCollectionDisplay>(),
            Type<MetadataViews.EVMBridgedMetadata>()
        ]
    }

    /// Function that resolves a metadata view for this contract.
    ///
    /// @param view: The Type of the desired view.
    /// @return A structure representing the requested view.
    ///
    access(all) fun resolveContractView(resourceType: Type?, viewType: Type): AnyStruct? {
        switch viewType {
            case Type<MetadataViews.NFTCollectionData>():
                let collectionData = MetadataViews.NFTCollectionData(
                    storagePath: self.CollectionStoragePath,
                    publicPath: self.CollectionPublicPath,
                    publicCollection: Type<&NFTAccessory.Collection>(),
                    publicLinkedType: Type<&NFTAccessory.Collection>(),
                    createEmptyCollectionFunction: (fun(): @{NonFungibleToken.Collection} {
                        return <-NFTAccessory.createEmptyCollection(nftType: Type<@NFTAccessory.NFT>())
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
                    name: "NFTAccessory",
                    symbol: "XMPL",
                    uri: MetadataViews.URI(
                        baseURI: nil, // setting baseURI as nil sets the given value as the uri field value
                        value: "https://example-nft.onflow.org/contract-metadata.json"
                    )
                )
        }
        return nil
    }

    /// Resource that an admin or something similar would own to be
    /// able to mint new NFTs
    ///
    access(all) resource NFTMinter {

        /// mintNFT mints a new NFT with a new ID
        /// and returns it to the calling context
        access(account) fun mintNFT(
            name: String,
            description: String,
            thumbnail: String,
            equipmentType: String,
            score: UFix64,
            max: UFix64,
            descriptionRarity: String
        ): @NFTAccessory.NFT {

            let metadata: {String: AnyStruct} = {}
            let currentBlock = getCurrentBlock()
            metadata["mintedBlock"] = currentBlock.height
            metadata["mintedTime"] = currentBlock.timestamp
            metadata["rarity"] = MetadataViews.Rarity(score: score, max: max, description: descriptionRarity)

            // create a new NFT
            var newNFT <- create NFT(
                name: name,
                description: description,
                thumbnail: thumbnail,
                equipmentType: equipmentType,
                metadata: metadata,
            )

            emit Minted(type: newNFT.getType().identifier,
                        id: newNFT.id,
                        uuid: newNFT.uuid,
                        minterAddress: self.owner?.address,
                        minterUUID: self.uuid,
                        name: name,
                        description: description,
                        thumbnail: thumbnail
                      )

            return <-newNFT
        }
    }

    init() {

        // Set the named paths
        self.CollectionStoragePath = /storage/NFTAccessoryCollection
        self.CollectionPublicPath = /public/NFTAccessoryCollection
        self.MinterStoragePath = /storage/NFTAccessoryMinter

        // Create a Minter resource and save it to storage
        let minter <- create NFTMinter()
        self.account.storage.save(<-minter, to: self.MinterStoragePath)
    }
}