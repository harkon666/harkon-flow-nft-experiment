/// This script uses the NFTMinter resource to mint a new NFT
/// It must be run with the account that has the minter resource
/// stored in /storage/NFTMinter
///
/// The royalty arguments indicies must be aligned

import "NonFungibleToken"
import "EventPass"
import "MetadataViews"

transaction(
    recipient: Address,
    name: String,
    description: String,
    thumbnail: String
) {

    /// local variable for storing the minter reference
    let minter: &EventPass.NFTMinter

    /// Reference to the receiver's collection
    let recipientCollectionRef: &EventPass.Collection

    prepare(signer: auth(BorrowValue) &Account) {

        let collectionData = EventPass.resolveContractView(resourceType: nil, viewType: Type<MetadataViews.NFTCollectionData>()) as! MetadataViews.NFTCollectionData?
            ?? panic("Could not resolve NFTCollectionData view. The EventPass contract needs to implement the NFTCollectionData Metadata view in order to execute this transaction")
        
        // borrow a reference to the NFTMinter resource in storage
        self.minter = signer.storage.borrow<&EventPass.NFTMinter>(from: EventPass.MinterStoragePath)
            ?? panic("The signer does not store an EventPass.Minter object at the path "
                     .concat(EventPass.MinterStoragePath.toString())
                     .concat("The signer must initialize their account with this minter resource first!"))

        // Borrow the recipient's public NFT collection reference
        self.recipientCollectionRef = getAccount(recipient).capabilities.borrow<&EventPass.Collection>(collectionData.publicPath)
            ?? panic("The recipient does not have a NonFungibleToken Receiver at "
                    .concat(collectionData.publicPath.toString())
                    .concat(" that is capable of receiving an NFT.")
                    .concat("The recipient must initialize their account with this collection and receiver first!"))
    }

    execute {
        // Mint the NFT and deposit it to the recipient's collection
        self.minter.mintNFT(
            recipient: self.recipientCollectionRef,
            name: name,
            description: description,
            thumbnail: thumbnail,
            eventType: 1,
            eventID: 1
        )        
    }

}