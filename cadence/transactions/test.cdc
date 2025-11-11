import "NonFungibleToken"
import "MetadataViews"
import "NFTAccessory"
/// Can pass in any contract address and name and NFT type name
/// This lets you choose the token you want to send because
/// the transaction gets the metadata from the provided contract.
///
/// @param to: The address to transfer the token to
/// @param id: The id of token to transfer
/// @param nftTypeIdentifier: The type identifier name of the NFT type you want to transfer
            /// Ex: "A.0b2a3299cc857e29.TopShot.NFT"
///
transaction(address:Address, id: UInt64) {

    // NFTCollectionData struct to get paths from
    let collectionData: &NFTAccessory.Collection

    prepare(signer: auth(BorrowValue) &Account) {
        
        self.collectionData = getAccount(address).capabilities.borrow<&NFTAccessory.Collection>(NFTAccessory.CollectionPublicPath)
            ?? panic("No ressource collection in storage")
        
        let nft = self.collectionData.borrowNFT(id) as! auth(NFTAccessory.Sale) &NFTAccessory.NFT
        // nft.itemListed(10)
        // self.collectionData = signer.storage.borrow<&NFTAccessory.Collection>(from: NFTAccessory.CollectionPublicPath)
        //     ?? panic("No ressource collection in storage")

        // self.collectionData = signer.storage.borrow<auth(NFTAccessory.Sale) &NFTAccessory.Collection>(from: NFTAccessory.CollectionStoragePath)
        //     ?? panic("No ressource collection in storage")
        
        // self.collectionData.itemUnlisted(nftID: id)
    }
}