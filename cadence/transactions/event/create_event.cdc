import "EventManager"
import "MetadataViews"

// create event in 'EventManager'

transaction(
    eventName: String,
    description: String,
    thumbnailURL: String,
    eventPassImg: String?,
    eventType: UInt8, // 0 = online, 1 = offline
    location: String,
    lat: Fix64,
    long: Fix64,
    startDate: UFix64,
    endDate: UFix64,
    quota: UInt64
) {

    prepare(signer: auth(BorrowValue) &Account) {
        let newID = EventManager.createEvent(
            hostAddress: signer.address,
            eventName: eventName,
            description: description,
            thumbnailURL: thumbnailURL,
            eventPassImg: eventPassImg,
            eventType: eventType,
            location: location,
            lat: lat,
            long: long,
            startDate: startDate,
            endDate: endDate,
            quota: quota
        )
        
        log("Event created successfully with ID: ".concat(newID.toString()))
    }
}