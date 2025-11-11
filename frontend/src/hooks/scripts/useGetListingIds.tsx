'use client';

import { useFlowQuery } from '@onflow/react-sdk';

// 1. Skrip Cadence Anda
const GET_LISTING_IDS_SCRIPT = `
  import "NFTStorefrontV2"

  access(all) fun main(account: Address): [UInt64] {
    return getAccount(account).capabilities.borrow<&{NFTStorefrontV2.StorefrontPublic}>(
            NFTStorefrontV2.StorefrontPublicPath
        )?.getListingIDs()
        ?? [] // Kembalikan array kosong jika tidak ada storefront
  }
`;

interface UseGetListingIDsProps {
  address: string;
}

// 2. Hook kustom Anda
export function useGetListingIDs({ address }: UseGetListingIDsProps) {

  const { data, isLoading, error, refetch } = useFlowQuery({
    cadence: GET_LISTING_IDS_SCRIPT,
    args: (arg, t) => [
      arg(address, t.Address)
    ],
    query: {
      // Hanya jalankan jika 'address' ada
      enabled: !!address, 
    },
  });

  return {
    listingIDs: data, // Akan bertipe: number[] | undefined
    isLoading,
    error,
    refetch
  };
}