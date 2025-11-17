

import "UserProfile"
import "EventPass"
import "NFTMoment"

transaction(
  nickname: String?,
  bio: String?,
  socials: {String: String},
  pfp: String?,
  shortDescription: String?,
  bgImage: String?,
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
          nickname: nickname,
          bio: bio,
          socials: socials,
          pfp: pfp,
          shortDescription: shortDescription,
          bgImage: bgImage,
          highlightedEventPassIds: highlightedEventPassIds,
          eventPassCollection: eventPassCollection,
          momentRef: momentRef,
        )
    }
}