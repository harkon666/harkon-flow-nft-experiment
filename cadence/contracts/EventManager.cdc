import "MetadataViews"
import "EventPass"

access(all) contract EventManager {
    access(all) var events: {UInt64: Event}
    access(all) var nextEventID: UInt64
    access(all) let eventManagerStoragePath: StoragePath

    access(all) enum eventType:UInt8 {
      access(all) case online;
      access(all) case offline;
    }

    access(all) event EventCreated(
      eventID: UInt64,
      hostAddress: Address,
      eventName: String,
      description: String,
      thumbnailURL: String,
      eventType: UInt8,
      location: String,
      lat: Fix64,
      long: Fix64,
      startDate: UFix64,
      endDate: UFix64,
      quota: UInt64
    )

    access(all) event UserRegistered(
      eventID: UInt64,
      userAddress: Address
    )
    
    access(all) event UserCheckedIn(
      eventID: UInt64,
      userAddress: Address
    )

    access(all) struct EventDetails {
        access(all) let eventID: UInt64
        access(all) let hostAddress: Address
        access(all) let eventName: String
        access(all) let description: String
        access(all) let thumbnail: MetadataViews.HTTPFile // URL Gambar
        access(all) let eventPassImg: MetadataViews.HTTPFile?
        access(all) let eventType: EventManager.eventType // "online" atau "offline"
        access(all) let location: String // Bisa URL (online) atau Alamat (offline)
        access(all) let createdAt: UFix64 // Timestamp
        access(all) let startDate: UFix64 // Timestamp
        access(all) let endDate: UFix64   // Timestamp
        access(all) let quota: UInt64
        access(all) let attendeeCount: Int
        access(all) let attendees: { Address: Bool }

        init(
            eventID: UInt64,
            hostAddress: Address,
            eventName: String,
            description: String,
            thumbnail: MetadataViews.HTTPFile,
            eventPassImg: MetadataViews.HTTPFile?,
            eventType: UInt8,
            location: String,
            createdAt: UFix64,
            startDate: UFix64,
            endDate: UFix64,
            quota: UInt64,
            attendeeCount: Int,
            attendees: { Address: Bool }
        ) {
            self.eventID = eventID
            self.hostAddress = hostAddress
            self.eventName = eventName
            self.description = description
            self.thumbnail = thumbnail
            self.eventPassImg = eventPassImg
            self.eventType = EventManager.eventType(rawValue: eventType)!
            self.location = location
            self.createdAt= createdAt
            self.startDate = startDate
            self.endDate = endDate
            self.quota = quota
            self.attendeeCount = attendeeCount
            self.attendees = attendees
        }
    }

    access(all) struct Event {
        access(all) let eventID: UInt64
        access(all) let hostAddress: Address
        access(all) let eventName: String
        access(all) let description: String
        access(all) let thumbnail: MetadataViews.HTTPFile
        access(all) let eventPassImg: MetadataViews.HTTPFile?
        access(all) let eventType: EventManager.eventType
        access(all) let location: String //url if online, location if offline
        access(all) let lat: Fix64
        access(all) let long: Fix64
        access(all) let createdAt: UFix64
        access(all) let startDate: UFix64
        access(all) let endDate: UFix64
        access(all) let quota: UInt64
        access(self) var attendees: {Address: Bool} //if false, user is registered, if true user is checked in

        init(
            hostAddress: Address,
            eventName: String,
            description: String,
            thumbnail: MetadataViews.HTTPFile,
            eventPassImg: MetadataViews.HTTPFile?,
            eventType: UInt8,
            location: String,
            lat: Fix64,
            long: Fix64,
            startDate: UFix64,
            endDate: UFix64,
            quota: UInt64,
        ) {
            let currentBlock = getCurrentBlock()
            self.eventID = EventManager.nextEventID
            self.hostAddress = hostAddress
            self.eventName = eventName
            self.description = description
            self.thumbnail = thumbnail
            self.eventPassImg = eventPassImg
            self.eventType = EventManager.eventType(rawValue: eventType)!
            self.location = location
            self.lat = lat
            self.long = long
            self.createdAt = currentBlock.timestamp
            self.startDate = startDate
            self.endDate = endDate
            self.quota = quota
            self.attendees = {}

            EventManager.nextEventID = EventManager.nextEventID + 1
        }

        access(contract) fun checkIn(userAddress: Address) {
            pre {
                self.attendees[userAddress] == false : "user already checked in or not registered"
            }
            self.attendees[userAddress] = true
            emit UserCheckedIn(eventID: self.eventID, userAddress: userAddress)
        }

        access(all) fun registerEvent(userAddress: Address) {
          pre {
                self.attendees.keys.length < Int(self.quota) : "Event is full"
                self.attendees[userAddress] == nil : "User sudah register"
          }
          self.attendees[userAddress] = false;
          emit UserRegistered(eventID: self.eventID, userAddress: userAddress)
        }

        access(all) fun getDetails(): EventDetails {
            return EventDetails(
                eventID: self.eventID,
                hostAddress: self.hostAddress,
                eventName: self.eventName,
                description: self.description,
                thumbnail: self.thumbnail,
                eventPassImg: self.eventPassImg,
                eventType: self.eventType.rawValue,
                location: self.location,
                createdAt: self.createdAt,
                startDate: self.startDate,
                endDate: self.endDate,
                quota: self.quota,
                attendeeCount: self.attendees.keys.length,
                attendees: self.attendees
            )
        }
    }

    access(all) fun createEvent(
      hostAddress: Address,
      eventName: String,
      description: String,
      thumbnailURL: String,
      eventPassImg: String?,
      eventType: UInt8,
      location: String,
      lat: Fix64,
      long: Fix64,
      startDate: UFix64,
      endDate: UFix64,
      quota: UInt64
    ): UInt64 {
      let currentBlock = getCurrentBlock()
      let newEvent = Event(
          hostAddress: hostAddress,
          eventName: eventName,
          description: description,
          thumbnail: MetadataViews.HTTPFile(url: thumbnailURL),
          eventPassImg: eventPassImg != nil ? MetadataViews.HTTPFile(url: thumbnailURL) : nil,
          eventType: eventType,
          location: location,
          lat: lat,
          long: long,
          startDate: startDate,
          endDate: endDate,
          quota: quota,
      )
      let newID: UInt64 = newEvent.eventID
      
      EventManager.events[newID] = newEvent
      emit EventManager.EventCreated(
          eventID: newID,
          hostAddress: hostAddress,
          eventName: eventName,
          description: description,
          thumbnailURL: thumbnailURL,
          eventType: eventType,
          location: location,
          lat: lat,
          long: long,
          startDate: startDate,
          endDate: endDate,
          quota: quota
      )
      
      return newID
    }

    access(all) fun userRegisterForEvent(eventID: UInt64, userAccount: &Account) {
        let userAddress = userAccount.address

        let eventRef = &self.events[eventID] as &Event?
            ?? panic("Event tidak ditemukan")
        
        eventRef.registerEvent(userAddress: userAddress)
    }
    access(all) resource Admin {
        //backend gated, only admin to make sure user is attending the event
        access(all) fun checkInUserToEvent(eventID: UInt64, userAddress: Address) {
            let eventRef = &EventManager.events[eventID] as &Event?
                ?? panic("Event not found")
            eventRef.checkIn(userAddress: userAddress)
        }
    }

    access(all) fun getAllEventDetails(): [EventDetails] {
        var allDetails: [EventDetails] = []
        for id in self.events.keys {
          if let eventRef = &self.events[id] as &Event? {
            allDetails.append(eventRef.getDetails())
          }
        }
        return allDetails
    }
    
    access(all) fun getEventDetails(eventID: UInt64): EventDetails? {
        if let eventRef = &self.events[eventID] as &Event? {
            return eventRef.getDetails()
        }
        return nil
    }

    init() {
        self.events = {}
        self.nextEventID = 1
        self.eventManagerStoragePath = /storage/EventManagerAdmin
        
        self.account.storage.save(<- create Admin(), to: self.eventManagerStoragePath)
    }
}