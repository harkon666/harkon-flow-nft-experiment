'use client';

import { useFlowQuery } from '@onflow/react-sdk';

// 1. Skrip Cadence Anda
const GET_LISTING_DETAILS_SCRIPT = `
  import "NFTStorefrontV2"

  access(all) fun main(account: Address, listingResourceID: UInt64): NFTStorefrontV2.ListingDetails {
    let storefrontRef = getAccount(account).capabilities.borrow<&{NFTStorefrontV2.StorefrontPublic}>(
            NFTStorefrontV2.StorefrontPublicPath
        ) ?? panic("Could not borrow public storefront from address")
    
    let listing = storefrontRef.borrowListing(listingResourceID: listingResourceID)
        ?? panic("No listing with that ID")
    
    return listing.getDetails()
  }
`;

// 2. Definisikan tipe data TypeScript yang sesuai dengan struct ListingDetails
//    Ini sangat penting untuk frontend Anda
export interface ListingDetails {
  listingResourceID: number;
  nftType: { typeID: string }; // Tipe NFT (misal: A.f8...NFTAccessories.NFT)
  nftID: number;
  salePrice: string; // UFix64 akan menjadi string (misal "1.00000000")
  commissionReceivers: any[]; // Bisa dibuat lebih spesifik jika perlu
  expiry: number;
  // ... (tambahkan field lain jika Anda perlukan)
}

interface UseGetListingDetailsProps {
  address: string;
  listingResourceID: number;
}

// 3. Hook kustom Anda
export function useGetListingDetails({ address, listingResourceID }: UseGetListingDetailsProps) {
  console.log(address, listingResourceID, 'hook')
  const { data, isLoading, error } = useFlowQuery({
    cadence: GET_LISTING_DETAILS_SCRIPT,
    args: (arg, t) => [
      arg(address, t.Address),
      arg(listingResourceID, t.UInt64) // Kirim ID sebagai string
    ],
    query: {
      enabled: !!address && listingResourceID != null,
      staleTime: 1000 * 60, // Cache detail selama 1 menit
    },
  });

  return {
    details: data, // Akan bertipe: ListingDetails | undefined
    isLoading,
    error
  };
}