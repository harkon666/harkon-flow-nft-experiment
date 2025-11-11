'use client';

import React from 'react';
// Hook dari File 2
import { useGetListingDetails } from '@/hooks/scripts/useGetListingDetails';
// Hook yang sudah kita buat sebelumnya untuk mengambil metadata
import { useGetAccessoryMetadata } from '@/hooks/scripts/useGetAccessoryMetadata'; // Sesuaikan path

interface ListingCardProps {
  listingResourceID: number;
  ownerAddress: string;
}

const ListingCard: React.FC<ListingCardProps> = ({ listingResourceID, ownerAddress }) => {
  // 1. Ambil detail Listing (Harga, NFT ID)
  const { details, isLoading: isLoadingDetails, error: errorDetails } = useGetListingDetails({
    address: ownerAddress,
    listingResourceID: Number(listingResourceID),
  });

  // 2. Ambil metadata NFT (Nama, Gambar)
  //    Kita gunakan 'nftID' dari hasil hook pertama
  const { display, isLoading: isLoadingDisplay, error: errorDisplay } = useGetAccessoryMetadata({
    address: ownerAddress,
    id: details?.nftID ?? 0, // '!' aman digunakan di sini berkat 'enabled' di bawah
    queryOptions: {
      // Hanya jalankan hook ini JIKA hook pertama sudah selesai & punya nftID
      enabled: true
    }
  });

  // Tampilan saat memuat
  if (isLoadingDetails || isLoadingDisplay) {
    return (
      <div className="pixel-card text-center animate-pulse">
        <div className="w-full aspect-square bg-gray-700 border-2 border-green-500 mb-6"></div>
        <div className="h-4 bg-gray-700 rounded w-1/2 mx-auto mb-4"></div>
        <div className="h-4 bg-gray-700 rounded w-1/4 mx-auto mb-6"></div>
        <div className="h-10 bg-gray-700 rounded w-full"></div>
      </div>
    );
  }

  // Tampilan jika error
  if (errorDetails || errorDisplay) {
    return <div className="pixel-card text-center"><p>Error memuat listing...</p></div>;
  }

  // Tampilan sukses
  return (
    <div className="pixel-card text-center">
      <div className="w-full aspect-square bg-gray-900 border-2 border-green-500 mb-6 flex items-center justify-center overflow-hidden">
        {display?.thumbnail?.url ? (
          <img 
            src={display.thumbnail.url} 
            alt={display.name}
            className="w-full h-full object-cover"
            style={{ imageRendering: 'pixelated' }}
          />
        ) : (
          <div className="text-6xl">🎮</div>
        )}
      </div>
      
      {/* Tampilkan NAMA dari metadata */}
      <p className="text-green-400 text-xl mb-2 pixel-text">
        {display?.name ?? `Accessory #${details?.nftID}`}
      </p>
      <p className="text-green-400 text-sm mb-2 pixel-text">
        #{details?.nftID ?? "#"}
      </p>

      {/* Tampilkan HARGA dari listing */}
      <p className="text-white text-lg mb-6 pixel-text">
        {details?.salePrice} FLOW
      </p>

      {/* Tombol ini akan memerlukan transaksi 'unlist' */}
      <button className="pixel-button w-full bg-red-500 hover:bg-red-400">
        [ UNLIST ]
      </button>
    </div>
  );
};

export default ListingCard;