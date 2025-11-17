/// This transaction is what an account would run but will be fixed when release
/// because the code still hardcode and there is wild return

import "NonFungibleToken"
import "NFTMoment"
import "NFTAccessory"
import "UserProfile"
import "EventPass"
import "MetadataViews"

transaction {

    prepare(signer: auth(BorrowValue, IssueStorageCapabilityController, PublishCapability, SaveValue, UnpublishCapability) &Account) {
        //NFT Moment collection setup
        let collectionNFTMomentData = NFTMoment.resolveContractView(resourceType: nil, viewType: Type<MetadataViews.NFTCollectionData>()) as! MetadataViews.NFTCollectionData?
            ?? panic("Could not resolve NFTCollectionData view. The NFTMoment contract needs to implement the NFTCollectionData Metadata view in order to execute this transaction")

        if signer.storage.borrow<&NFTMoment.Collection>(from: collectionNFTMomentData.storagePath) != nil {
            return
        }

        let collectionNFTMoment <- NFTMoment.createEmptyCollection(nftType: Type<@NFTMoment.NFT>())

        signer.storage.save(<-collectionNFTMoment, to: collectionNFTMomentData.storagePath)

        signer.capabilities.unpublish(collectionNFTMomentData.publicPath)
        let collectionNFTMomentCap = signer.capabilities.storage.issue<&NFTMoment.Collection>(collectionNFTMomentData.storagePath)
        signer.capabilities.publish(collectionNFTMomentCap, at: collectionNFTMomentData.publicPath)

        //NFT Accessory collection setup 
        let collectionNFTAccessoryData = NFTAccessory.resolveContractView(resourceType: nil, viewType: Type<MetadataViews.NFTCollectionData>()) as! MetadataViews.NFTCollectionData?
            ?? panic("Could not resolve NFTCollectionData view. The NFTAccessory contract needs to implement the NFTCollectionData Metadata view in order to execute this transaction")

        if signer.storage.borrow<&NFTAccessory.Collection>(from: collectionNFTAccessoryData.storagePath) != nil {
            return
        }

        let collectionNFTAccessory <- NFTAccessory.createEmptyCollection(nftType: Type<@NFTAccessory.NFT>())

        signer.storage.save(<-collectionNFTAccessory, to: collectionNFTAccessoryData.storagePath)

        signer.capabilities.unpublish(collectionNFTAccessoryData.publicPath)

        let collectionNFTAccessoryCap = signer.capabilities.storage.issue<&NFTAccessory.Collection>(collectionNFTAccessoryData.storagePath)
        signer.capabilities.publish(collectionNFTAccessoryCap, at: collectionNFTAccessoryData.publicPath)

        //event pass collection setup
        let collectionEventPassData = EventPass.resolveContractView(resourceType: nil, viewType: Type<MetadataViews.NFTCollectionData>()) as! MetadataViews.NFTCollectionData?
            ?? panic("Could not resolve NFTCollectionData view. The EventPass contract needs to implement the NFTCollectionData Metadata view in order to execute this transaction")

        if signer.storage.borrow<&EventPass.Collection>(from: collectionEventPassData.storagePath) != nil {
            return
        }

        let collectionEventPass <- EventPass.createEmptyCollection(nftType: Type<@EventPass.NFT>())

        signer.storage.save(<-collectionEventPass, to: collectionEventPassData.storagePath)

        signer.capabilities.unpublish(collectionEventPassData.publicPath)
        let collectionEventPassCap = signer.capabilities.storage.issue<&EventPass.Collection>(collectionEventPassData.storagePath)
        signer.capabilities.publish(collectionEventPassCap, at: collectionEventPassData.publicPath)

        // user profile setup

        let profile: @UserProfile.Profile <- UserProfile.createEmptyProfile()
        signer.storage.save(<-profile, to: UserProfile.ProfileStoragePath)


        signer.capabilities.unpublish(UserProfile.ProfilePublicPath)
        let userProfileCap = signer.capabilities.storage.issue<&UserProfile.Profile>(UserProfile.ProfileStoragePath)
        signer.capabilities.publish(userProfileCap, at: UserProfile.ProfilePublicPath)
    }
}