import "NFTMoment"
import "EventPass"

access(all) contract UserProfile {

    access(all) let ProfileStoragePath: StoragePath
    access(all) let ProfilePublicPath: PublicPath
    access(all) let VerifierStoragePath: StoragePath

    access(all) event ProfileUpdated(
      address: Address,
      nickname: String?,
      bio: String,
      socials: {String: String},
      pfp: String?,
      shortDescription: String?,
      bgImage: String?,
      highlightedEventPassIds: [UInt64?],
      highlightedMomentID: UInt64?
    )

    access(all) event UserVerified(
      address: Address
    )

    access(all) entitlement Edit 

    access(all) struct ProfileView {
        access(all) var nickname: String?
        access(all) var bio: String
        access(all) var socials: {String: String}
        access(all) var pfp: String?
        access(all) var isVerified: Bool
        access(all) var shortDescription: String?
        access(all) var bgImage: String?
        access(all) var highlightedEventPassIds: [UInt64?]
        access(all) var highlightedMomentID: UInt64?

        init(
          nickname: String?,
          bio: String?,
          socials: {String: String},
          pfp: String?,
          isVerified: Bool,
          shortDescription: String?,
          bgImage: String?,
          highlightedEventPassIds: [UInt64?],
          highlightedMomentID: UInt64?
        ) {
            self.nickname = nickname
            self.bio = bio ?? ""
            self.socials = socials
            self.pfp = pfp
            self.bgImage = bgImage
            self.isVerified = isVerified
            self.shortDescription = shortDescription
            self.highlightedEventPassIds = highlightedEventPassIds
            self.highlightedMomentID = highlightedMomentID
        }
    }

    access(all) resource Profile {
        access(all) let id: UInt64
        access(all) var nickname: String?
        access(all) var bio: String
        access(all) var socials: {String: String}
        access(all) var pfp: String?
        access(all) var isVerified: Bool
        access(all) var shortDescription: String?
        access(all) var bgImage: String?
        access(all) var highlightedEventPassIds: [UInt64?]
        access(all) var highlightedMomentID: UInt64?

        init() {
            self.id = self.uuid
            self.nickname = nil
            self.bio = ""
            self.socials = {}
            self.pfp = nil
            self.bgImage = nil
            self.isVerified = false
            self.shortDescription = nil
            self.highlightedEventPassIds = []
            self.highlightedMomentID = nil
        }

        access(Edit) fun updateProfile(
          nickname: String?,
          bio: String?,
          socials: {String: String},
          pfp: String?,
          shortDescription: String?,
          bgImage: String?,
          highlightedEventPassIds: [UInt64?],
          eventPassCollection: &EventPass.Collection?,
          momentRef: &NFTMoment.NFT?,
        ) {
          self.nickname = nickname
          self.bio = bio ?? ""
          self.socials = socials
          self.pfp = pfp
          self.shortDescription = shortDescription
          self.bgImage = bgImage

          if highlightedEventPassIds.length == 0 {
            self.highlightedEventPassIds = []
          } else {
            assert(highlightedEventPassIds.length <= 4, message: "max event pass is 4")
            assert(eventPassCollection!.owner!.address == self.owner!.address, message: "You are not own this Event Pass Collection")
            for eventPassIds in highlightedEventPassIds {
              assert(eventPassCollection!.borrowNFT(eventPassIds!) != nil, message: "Invalid EventPass ID")
            }

            self.highlightedEventPassIds = highlightedEventPassIds
          }

          if momentRef != nil {
            assert(momentRef!.owner!.address == self.owner!.address, message: "You are not own this Moment NFT")
          }
          self.highlightedMomentID = momentRef?.id
          
          emit ProfileUpdated(
            address: self.owner!.address,
            nickname: self.nickname,
            bio: self.bio,
            socials: self.socials,
            pfp: self.pfp,
            shortDescription: self.shortDescription,
            bgImage: self.bgImage,
            highlightedEventPassIds: self.highlightedEventPassIds,
            highlightedMomentID: self.highlightedMomentID
          )
        }

        access(contract) fun verified() {
          self.isVerified = true
        }

        access(all) fun getDetails(): ProfileView {
          return ProfileView(
            nickname: self.nickname,
            bio: self.bio,
            socials: self.socials,
            pfp: self.pfp,
            isVerified: self.isVerified,
            shortDescription: self.shortDescription,
            bgImage: self.bgImage,
            highlightedEventPassIds: self.highlightedEventPassIds,
            highlightedMomentID: self.highlightedMomentID,
          )
        }
    }

    access(all) fun createEmptyProfile(): @Profile {
        return <- create Profile()
    }

    access(all) resource Verifier {
      access(all) fun verifyUser(userProfileRef: &UserProfile.Profile) {
        userProfileRef.verified()
        emit UserVerified(address: userProfileRef.owner!.address)
      }
    }

    init() {
        self.ProfileStoragePath = /storage/MomentumUserProfile
        self.ProfilePublicPath = /public/MomentumUserProfile
        self.VerifierStoragePath = /storage/MomentUserVerifier

        let verifier <- create Verifier()
        self.account.storage.save(<-verifier, to: self.VerifierStoragePath)
    }
}