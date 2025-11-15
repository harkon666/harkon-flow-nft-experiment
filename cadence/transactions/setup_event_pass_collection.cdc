/// This transaction is what an account would run
/// to set itself up to receive NFTs

import "NonFungibleToken"
import "EventPass"
import "MetadataViews"

transaction {

    prepare(signer: auth(BorrowValue, IssueStorageCapabilityController, PublishCapability, SaveValue, UnpublishCapability) &Account) {
        
        let collectionData = EventPass.resolveContractView(resourceType: nil, viewType: Type<MetadataViews.NFTCollectionData>()) as! MetadataViews.NFTCollectionData?
            ?? panic("Could not resolve NFTCollectionData view. The EventPass contract needs to implement the NFTCollectionData Metadata view in order to execute this transaction")

        // Return early if the account already has a collection
        if signer.storage.borrow<&EventPass.Collection>(from: collectionData.storagePath) != nil {
            return
        }

        // Create a new empty collection
        let collection <- EventPass.createEmptyCollection(nftType: Type<@EventPass.NFT>())

        // save it to the account
        signer.storage.save(<-collection, to: collectionData.storagePath)

        // create a public capability for the collection
        signer.capabilities.unpublish(collectionData.publicPath)
        let collectionCap = signer.capabilities.storage.issue<&EventPass.Collection>(collectionData.storagePath)
        signer.capabilities.publish(collectionCap, at: collectionData.publicPath)
    }
}