'use client';

import { useFlowQuery } from '@onflow/react-sdk';
import * as t from "@onflow/types";
import { arg } from "@onflow/fcl";

// 1. Skrip Cadence (Sudah saya modifikasi agar mengembalikan data, bukan referensi)
//    Ini akan mengembalikan struct Display dari bingkai, atau nil
const CHECK_EQUIPMENT_SCRIPT = `
  import "NFTMoment"
  import "NFTAccessory"
  import "MetadataViews"
  import "ViewResolver"

  access(all) fun main(address: Address, id: UInt64): MetadataViews.Display? {
    
    let account = getAccount(address)
    // Asumsi Anda punya ViewResolver.ResolverCollection di path ini
    let collection = account.capabilities.borrow<&{ViewResolver.ResolverCollection}>(
          NFTMoment.CollectionPublicPath
      ) ?? panic("Tidak bisa meminjam koleksi")

    let resolver = collection.borrowViewResolver(id: id) 
        ?? panic("Tidak bisa meminjam resolver")
    
    // Minta view kustom Anda
    let view = resolver.resolveView(Type<NFTMoment.NFTMomentEquipmentMetadataView>())

    // Buka 'view' kustom itu
    if let equipmentView = view as! NFTMoment.NFTMomentEquipmentMetadataView? {
        
        // Jika ada bingkai terpasang (frameRef != nil)
        if let frameRef = equipmentView.equippedFrame {
            
            // 'frameRef' adalah &NFTAccessories.NFT
            // Kita 'cast' (ubah) dia menjadi Resolver agar kita bisa membaca metadatanya
            if let frameResolver = frameRef as? &{ViewResolver.Resolver} {
                 
                 // Ambil data Display dari bingkai itu
                 let displayView = frameResolver.resolveView(Type<MetadataViews.Display>())
                 
                 // Kembalikan datanya!
                 return displayView as! MetadataViews.Display?
            }
        }
    }
    
    // Kembalikan nil jika tidak ada bingkai terpasang
    return nil
  }
`;

// 2. Tipe data yang kita harapkan (dari jawaban sebelumnya)
interface DisplayView {
  name: string;
  description: string;
  thumbnail: { url: string }; 
}

// 3. Tipe untuk props hook
interface UseCheckEquipmentProps {
  address: string;
  momentId: number | null; // ID dari MOMEN yang ingin dicek
}

// 4. Hook kustom Anda
export function useCheckEquipment({ address, momentId }: UseCheckEquipmentProps) {

  const { data, isLoading, error, refetch } = useFlowQuery({
    cadence: CHECK_EQUIPMENT_SCRIPT,
    
    args: (arg, t) => [
      arg(address, t.Address),
      arg(momentId ?? 0, t.UInt64) // Konversi JS 'number' ke Cadence 'UInt64'
    ],
    
    query: {
      // Hanya jalankan jika kita punya alamat DAN ID Momen
      enabled: !!address && momentId != null, 
      staleTime: 1000 * 60, // Cache data selama 1 menit
    },
  });

  return {
    // Kembalikan data dengan nama yang jelas
    equippedFrameData: data, // Akan bertipe: DisplayView | null | undefined
    isLoading,
    error,
    refetch
  };
}