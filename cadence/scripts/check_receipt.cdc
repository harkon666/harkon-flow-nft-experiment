import "AccessoryPack"

access(all) fun main(address: Address): Bool {
  let account = getAccount(address)
  if account.storage.type(at: AccessoryPack.ReceiptStoragePath) != nil {
    return false
  }
  return true
}