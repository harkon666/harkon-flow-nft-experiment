'use client';

import { useFlowQuery } from '@onflow/react-sdk';
import * as t from "@onflow/types";
import { arg } from "@onflow/fcl";

// 1. Definisikan string Cadence Anda
const GET_ACCESSORY_IDS_SCRIPT = `
  import "NonFungibleToken"
  import "NFTMoment"

  access(all) fun main(address: Address): [UInt64] {
    let account = getAccount(address)

    let collectionRef = account.capabilities.borrow<&{NonFungibleToken.Collection}>(
            NFTMoment.CollectionPublicPath
        ) ?? panic("Tidak bisa meminjam koleksi publik di path: ".concat(NFTMoment.CollectionPublicPath.toString()))

    return collectionRef.getIDs()
  }
`;

// 2. Tentukan Tipe untuk argumen hook
interface UseGetAccessoryIDsProps {
  address: string // Terima alamat (atau null jika belum login)
}

// 3. Tentukan Tipe untuk data yang dikembalikan
//    Cadence [UInt64] akan menjadi number[] di JavaScript
type AccessoryIDQueryResult = number[];

// 4. Buat hook kustom Anda
export function useGetMomentIDs({ address }: UseGetAccessoryIDsProps) {

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
      // Ini PENTING:
      // Hanya jalankan query (enabled: true)
      // JIKA 'address' tidak null (!!address).
      enabled: !!address, 
    },
  });

  // 6. Kembalikan (return) hasilnya
  return {
    data: data, // Akan bertipe: number[] | undefined
    isLoading,
    error,
    refetch
  };
}