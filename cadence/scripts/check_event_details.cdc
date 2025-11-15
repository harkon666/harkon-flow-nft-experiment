import "EventManager"

// Ini adalah SKRIP (read-only), bukan transaksi
// Ia mengembalikan struct EventDetails atau nil

access(all) fun main(eventID: UInt64): EventManager.EventDetails? {
    
    // Panggil fungsi publik 'getEventDetails' di kontrak
    // Fungsi ini aman untuk dipanggil siapa saja
    return EventManager.getEventDetails(eventID: eventID)
}