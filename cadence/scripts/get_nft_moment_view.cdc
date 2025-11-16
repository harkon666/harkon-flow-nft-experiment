import "NFTMoment"
import "MetadataViews"
import "ViewResolver"

access(all) struct NFTView {
    access(all) let id: UInt64
    access(all) let uuid: UInt64
    access(all) let name: String
    access(all) let description: String
    access(all) let thumbnail: String
    access(all) let externalURL: String
    access(all) let collectionPublicPath: PublicPath
    access(all) let collectionStoragePath: StoragePath
    access(all) let collectionPublic: String
    access(all) let collectionPublicLinkedType: String
    access(all) let collectionName: String
    access(all) let collectionDescription: String
    access(all) let collectionExternalURL: String
    access(all) let collectionSquareImage: String
    access(all) let collectionBannerImage: String
    access(all) let traits: MetadataViews.Traits

    init(
        id: UInt64,
        uuid: UInt64,
        name: String,
        description: String,
        thumbnail: String,
        externalURL: String,
        collectionPublicPath: PublicPath,
        collectionStoragePath: StoragePath,
        collectionPublic: String,
        collectionPublicLinkedType: String,
        collectionName: String,
        collectionDescription: String,
        collectionExternalURL: String,
        collectionSquareImage: String,
        collectionBannerImage: String,
        traits: MetadataViews.Traits
    ) {
        self.id = id
        self.uuid = uuid
        self.name = name
        self.description = description
        self.thumbnail = thumbnail
        self.externalURL = externalURL
        self.collectionPublicPath = collectionPublicPath
        self.collectionStoragePath = collectionStoragePath
        self.collectionPublic = collectionPublic
        self.collectionPublicLinkedType = collectionPublicLinkedType
        self.collectionName = collectionName
        self.collectionDescription = collectionDescription
        self.collectionExternalURL = collectionExternalURL
        self.collectionSquareImage = collectionSquareImage
        self.collectionBannerImage = collectionBannerImage
        self.traits = traits
    }
}

access(all) fun main(address: Address, id: UInt64): NFTView {
    let account = getAccount(address)

    let collectionData = NFTMoment.resolveContractView(resourceType: nil, viewType: Type<MetadataViews.NFTCollectionData>()) as! MetadataViews.NFTCollectionData?
            ?? panic("Could not resolve NFTCollectionData view. The NFTMoment contract needs to implement the NFTCollectionData Metadata view in order to execute this transaction")

    let collection = account.capabilities.borrow<&NFTMoment.Collection>(
            collectionData.publicPath
    ) ?? panic("The account ".concat(address.toString()).concat(" does not have a NonFungibleToken Collection at ")
                .concat(collectionData.publicPath.toString())
                .concat(". The account must initialize their account with this collection first!"))

    let viewResolver = collection.borrowViewResolver(id: id) 
        ?? panic("Could not borrow resolver with given id ".concat(id.toString()))

    let nftView = MetadataViews.getNFTView(id: id, viewResolver : viewResolver)

    return NFTView(
        id: nftView.id,
        uuid: nftView.uuid,
        name: nftView.display!.name,
        description: nftView.display!.description,
        thumbnail: nftView.display!.thumbnail.uri(),
        externalURL: nftView.externalURL!.url,
        collectionPublicPath: nftView.collectionData!.publicPath,
        collectionStoragePath: nftView.collectionData!.storagePath,
        collectionPublic: nftView.collectionData!.publicCollection.identifier,
        collectionPublicLinkedType: nftView.collectionData!.publicLinkedType.identifier,
        collectionName: nftView.collectionDisplay!.name,
        collectionDescription: nftView.collectionDisplay!.description,
        collectionExternalURL: nftView.collectionDisplay!.externalURL.url,
        collectionSquareImage: nftView.collectionDisplay!.squareImage.file.uri(),
        collectionBannerImage: nftView.collectionDisplay!.bannerImage.file.uri(),
        traits: nftView.traits!,
    )
}