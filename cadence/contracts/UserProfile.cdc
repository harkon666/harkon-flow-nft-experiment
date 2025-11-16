access(all) contract UserProfile {

    access(all) let ProfileStoragePath: StoragePath
    access(all) let ProfilePublicPath: PublicPath

    access(all) event ProfileCreated(address: Address)
    access(all) event ProfileUpdated(address: Address)
    access(all) event PFPSet(address: Address, pfpLink: String)

    access(all) entitlement Edit 

    access(all) resource Profile {
        access(all) var bio: String
        access(all) var socials: {String: String}
        access(all) var pfp: String?

        access(all) fun setBio(newBio: String) {
            self.bio = newBio
            emit ProfileUpdated(address: self.owner!.address)
        }

        access(Edit) fun setSocial(key: String, value: String) {
            self.socials[key] = value
            emit ProfileUpdated(address: self.owner!.address)
        }

        access(Edit) fun setPFP(
          pfpLink: String
        ) {
            self.pfp = pfpLink
            emit PFPSet(address: self.owner!.address, pfpLink: pfpLink)
        }

        access(Edit) fun removePFP() {
            self.pfp = nil
            emit ProfileUpdated(address: self.owner!.address)
        }

        init(userAddress: Address) {
            self.bio = ""
            self.socials = {}
            self.pfp = nil
            emit ProfileCreated(address: userAddress)
        }
    }

    access(all) fun createEmptyProfile(_ userAddress: Address): @Profile {
        return <- create Profile(userAddress: userAddress)
    }

    init() {
        self.ProfileStoragePath = /storage/MomentumUserProfile
        self.ProfilePublicPath = /public/MomentumUserProfile
    }
}