//this transaction will be executed in backend
//there is useFreeMint and can be executed once in lifetime


import "NonFungibleToken"
import "NFTMoment"
import "MetadataViews"
import "EventPass"

transaction(
    recipient: Address,
    name: String,
    description: String,
    thumbnail: String,
) {

    /// local variable for storing the minter reference
    let minter: &NFTMoment.NFTMinter
    let adminEventPass: &EventPass.NFTMinter

    /// Reference to the receiver's collection
    let recipientCollectionRef: &NFTMoment.Collection

    prepare(signer: auth(BorrowValue) &Account) {

        let collectionData = NFTMoment.resolveContractView(resourceType: nil, viewType: Type<MetadataViews.NFTCollectionData>()) as! MetadataViews.NFTCollectionData?
            ?? panic("Could not resolve NFTCollectionData view. The NFTMoment contract needs to implement the NFTCollectionData Metadata view in order to execute this transaction")
        let collectionEventData = EventPass.resolveContractView(resourceType: nil, viewType: Type<MetadataViews.NFTCollectionData>()) as! MetadataViews.NFTCollectionData?
            ?? panic("Could not resolve NFTCollectionData view. The EventPass contract needs to implement the NFTCollectionData Metadata view in order to execute this transaction")
        
        self.minter = signer.storage.borrow<&NFTMoment.NFTMinter>(from: NFTMoment.MinterStoragePath)
            ?? panic("The signer does not store an NFTMoment.Minter object at the path "
                     .concat(NFTMoment.MinterStoragePath.toString())
                     .concat("The signer must initialize their account with this minter resource first!"))
        self.adminEventPass = signer.storage.borrow<&EventPass.NFTMinter>(from: EventPass.MinterStoragePath)
            ?? panic("The signer does not store an EventPass.Minter object at the path "
                     .concat(NFTMoment.MinterStoragePath.toString())
                     .concat("The signer must initialize their account with this minter resource first!"))

        // Borrow the recipient's public NFT collection reference
        self.recipientCollectionRef = getAccount(recipient).capabilities.borrow<&NFTMoment.Collection>(collectionData.publicPath)
            ?? panic("The recipient does not have a NonFungibleToken Receiver at "
                    .concat(collectionData.publicPath.toString())
                    .concat(" that is capable of receiving an NFT.")
                    .concat("The recipient must initialize their account with this collection and receiver first!"))
    }

    execute {
        // Mint the NFT and deposit it to the recipient's collection
        self.minter.freeMint(
            recipient: self.recipientCollectionRef,
            name: name,
            description: description,
            thumbnail: thumbnail,
        )        
    }
}