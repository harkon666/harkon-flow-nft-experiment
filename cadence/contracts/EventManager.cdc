/**
EventManager.cdc

Kontrak ini bertindak sebagai "Papan Buletin" atau "Registri" on-chain
untuk semua event yang dibuat oleh platform Harkon-NFT.
Ia menyimpan detail event, tetapi TIDAK menangani tiket (EventPass).
*/

import "MetadataViews" // Kita mungkin perlu ini untuk data gambar

access(all) contract EventManager {

    // 1. STATE KONTRAK
    // di-indeks berdasarkan ID unik.
    access(all) var events: {UInt64: Event}
    
    // Ini adalah 'counter' untuk memastikan ID event selalu unik
    access(all) var nextEventID: UInt64

    access(all) enum eventType:UInt8 {
      access(all) case online;
      access(all) case offline;
    }

    access(all) event EventCreated(
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

    // 2. STRUCT DATA PUBLIK
    // Ini adalah 'struct' data bersih yang akan dibaca oleh
    // skrip frontend (useFlowQuery) Anda.
    // Ini lebih aman daripada mengekspos resource secara langsung.
    access(all) struct EventDetails {
        access(all) let eventID: UInt64
        access(all) let hostAddress: Address // Siapa penyelenggara event (misal: RPN)
        access(all) let eventName: String
        access(all) let description: String
        access(all) let thumbnail: MetadataViews.HTTPFile // URL Gambar
        access(all) let eventType: EventManager.eventType // "online" atau "offline"
        access(all) let location: String // Bisa URL (online) atau Alamat (offline)
        access(all) let createdAt: UFix64 // Timestamp
        access(all) let startDate: UFix64 // Timestamp
        access(all) let endDate: UFix64   // Timestamp
        access(all) let quota: UInt64 // Kapasitas maks.
        access(all) let attendeeCount: Int // Jumlah yang sudah check-in
        access(all) let attendees: { Address: Bool }

        init(
            eventID: UInt64,
            hostAddress: Address,
            eventName: String,
            description: String,
            thumbnail: MetadataViews.HTTPFile,
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

    // 3. RESOURCE EVENT (Data Inti)
    // Ini adalah 'resource' yang disimpan di storage KONTRAK
    access(all) struct Event {
        access(all) let eventID: UInt64
        access(all) let hostAddress: Address
        access(all) let eventName: String
        access(all) let description: String
        access(all) let thumbnail: MetadataViews.HTTPFile
        access(all) let eventType: EventManager.eventType
        access(all) let location: String
        access(all) let lat: Fix64
        access(all) let long: Fix64
        access(all) let createdAt: UFix64
        access(all) let startDate: UFix64
        access(all) let endDate: UFix64
        access(all) let quota: UInt64

        // DAFTAR PESERTA (ATTENDEE)
        // Ini adalah "daftar tamu" on-chain!
        // Alamat -> Status Check-in (true)
        access(self) var attendees: {Address: Bool} //if false, user is registered, if true user is checked in

        init(
            hostAddress: Address,
            eventName: String,
            description: String,
            thumbnail: MetadataViews.HTTPFile,
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
            self.eventType = EventManager.eventType(rawValue: eventType)!
            self.location = location
            self.lat = lat
            self.long = long
            self.createdAt = currentBlock.timestamp
            self.startDate = startDate
            self.endDate = endDate
            self.quota = quota
            self.attendees = {}

            // Naikkan 'counter' untuk event berikutnya
            EventManager.nextEventID = EventManager.nextEventID + 1
        }

        // --- FUNGSI INTERNAL (Hanya bisa dipanggil oleh Admin) ---

        // Fungsi untuk mendaftarkan check-in
        access(contract) fun checkIn(userAddress: Address) {
            pre {
                self.attendees[userAddress] == false : "User sudah check-in"
            }
            self.attendees[userAddress] = true
        }

        access(all) fun registerEvent(userAddress: Address) {
          pre {
                self.attendees.keys.length < Int(self.quota) : "Event sudah penuh"
                self.attendees[userAddress] == nil : "User sudah register"
          }
          self.attendees[userAddress] = false;
        }

        // Fungsi untuk membaca detail (dipanggil oleh Admin)
        access(all) fun getDetails(): EventDetails {
            return EventDetails(
                eventID: self.eventID,
                hostAddress: self.hostAddress,
                eventName: self.eventName,
                description: self.description,
                thumbnail: self.thumbnail,
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

    access(all) fun userRegisterForEvent(eventID: UInt64, userAccount: &Account) {
        // Dapatkan alamat signer
        let userAddress = userAccount.address

        // Pinjam event dari storage kontrak
        let eventRef = &self.events[eventID] as &Event?
            ?? panic("Event tidak ditemukan")
        
        // Panggil fungsi register internal di resource Event
        eventRef.registerEvent(userAddress: userAddress)
        emit UserRegistered(
          eventID: eventID,
          userAddress: userAddress
        )
    }

    // 4. RESOURCE ADMIN (Untuk Backend Go Anda)
    // Ini adalah "kunci" untuk mengelola event
    access(all) resource Admin {
        
        // Fungsi untuk membuat event baru
        access(all) fun createEvent(
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
        ): UInt64 {
            let currentBlock = getCurrentBlock()
            let newEvent = Event(
                hostAddress: hostAddress,
                eventName: eventName,
                description: description,
                thumbnail: MetadataViews.HTTPFile(url: thumbnailURL),
                eventType: eventType,
                location: location,
                lat: lat,
                long: long,
                startDate: startDate,
                endDate: endDate,
                quota: quota,
            )
            emit EventManager.EventCreated(
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
            let newID: UInt64 = newEvent.eventID
            
            // Simpan event baru ke 'database' kontrak
            EventManager.events[newID] = newEvent
            
            return newID
        }

        // Fungsi untuk melakukan check-in (Offline atau Online)
        access(all) fun checkInUserToEvent(eventID: UInt64, userAddress: Address) {
            // Pinjam event dari storage kontrak
            let eventRef = &EventManager.events[eventID] as &Event?
                ?? panic("Event tidak ditemukan")
            
            // Panggil fungsi check-in internal
            eventRef.checkIn(userAddress: userAddress)
        }
    }

    // --- FUNGSI PUBLIK (Untuk Skrip Frontend) ---
    
    // Mengembalikan SEMUA event
    access(all) fun getAllEventDetails(): [EventDetails] {
        var allDetails: [EventDetails] = []
        for id in self.events.keys {
          if let eventRef = &self.events[id] as &Event? {
            allDetails.append(eventRef.getDetails())
          }
        }
        return allDetails
    }
    
    // Mengembalikan satu event
    access(all) fun getEventDetails(eventID: UInt64): EventDetails? {
        if let eventRef = &self.events[eventID] as &Event? {
            return eventRef.getDetails()
        }
        return nil
    }

    // ---
    init() {
        self.events = {}
        self.nextEventID = 1
        
        // Simpan Admin Resource di storage deployer
        self.account.storage.save(<- create Admin(), to: /storage/EventManagerAdmin)
    }
}