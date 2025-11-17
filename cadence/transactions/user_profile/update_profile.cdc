

import "UserProfile"
import "EventPass"
import "NFTMoment"

transaction(
  highlightedEventPassIds: [UInt64?],
  momentID: UInt64?
) {

    prepare(signer: auth(BorrowValue) &Account) {
        //NFT Moment collection setup
        let userProfile = signer.storage.borrow<auth(UserProfile.Edit) &UserProfile.Profile>(
          from: UserProfile.ProfileStoragePath
        ) ?? panic("User Profile Ressource not found")

        let eventPassCollection = signer.storage.borrow<&EventPass.Collection>(from: EventPass.CollectionStoragePath)
        let momentCollection = signer.storage.borrow<&NFTMoment.Collection>(from: NFTMoment.CollectionStoragePath)
        let momentRef = momentID != nil ? momentCollection?.borrowNFT(momentID!) as! &NFTMoment.NFT : nil
        userProfile.updateProfile(
          nickname: "test",
          bio: "test",
          socials: {"test": "test"},
          pfp: "test",
          shortDescription: "test",
          bgImage: "test",
          highlightedEventPassIds: highlightedEventPassIds,
          eventPassCollection: eventPassCollection,
          momentRef: momentRef,
        )
    }
}