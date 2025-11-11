'use client';

import { useFlowQuery } from '@onflow/react-sdk';

// 1. Definisikan string Cadence Anda
const GET_ACCESSORY_IDS_SCRIPT = `
  import "NonFungibleToken"
  import "NFTAccessory"

  access(all) fun main(address: Address): [UInt64] {
    let account = getAccount(address)

    // Pastikan path (NFTAccessory.CollectionPublicPath)
    // didefinisikan di file konfigurasi FCL Anda
    let collectionRef = account.capabilities.borrow<&{NonFungibleToken.Collection}>(
            NFTAccessory.CollectionPublicPath
        ) ?? panic("Tidak bisa meminjam koleksi publik di path: ".concat(NFTAccessory.CollectionPublicPath.toString()))

    return collectionRef.getIDs()
  }
`;

// 2. Tentukan Tipe untuk argumen hook
interface UseGetAccessoryIDsProps {
  address: string // Terima alamat (atau null jika belum login)
}

// 3. Tentukan Tipe untuk data yang dikembalikan
//    Cadence [UInt64] akan menjadi number[] di JavaScript

// 4. Buat hook kustom Anda
export function useGetAccessoryIDs({ address }: UseGetAccessoryIDsProps) {

  // 5. Gunakan useFlowQuery di dalamnya
  const { data, isLoading, error, refetch } = useFlowQuery({
    // Cadence script
    cadence: GET_ACCESSORY_IDS_SCRIPT,
    
    // Argumen untuk script
    args: (arg, t) => [
      arg(address, t.Address)
    ],
    
    // Konfigurasi query
    query: {
      enabled: !!address, 
    },
  });

  // 6. Kembalikan (return) hasilnya
  return {
    data: data,
    isLoading,
    error,
    refetch
  };
}