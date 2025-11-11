'use client';

import React from 'react';
import { useGetAccessoryMetadata } from '../../hooks/scripts/useGetAccessoryMetadata'; // Sesuaikan path

interface AccessoryCardProps {
  id: number;
  ownerAddress: string;
}

const AccessoryCard: React.FC<AccessoryCardProps> = ({ id, ownerAddress }) => {
  // Panggil hook baru kita untuk mengambil data untuk ID ini
  const { display, isLoading, error } = useGetAccessoryMetadata({
    address: ownerAddress,
    id: id,
  });

  // Tampilan saat memuat data metadata
  if (isLoading) {
    return (
      <div className="pixel-card text-center animate-pulse">
        <div className="w-full aspect-square bg-gray-700 border-2 border-green-500 mb-6"></div>
        <p className="text-green-400 text-xl mb-6 pixel-text">Loading...</p>
      </div>
    );
  }

  // Tampilan jika gagal
  if (error) {
    return <div className="pixel-card text-center"><p>Error memuat ID: {id}</p></div>;
  }

  // Tampilan setelah data dimuat
  return (
    <div className="pixel-card text-center">
      <div className="w-full aspect-square bg-gray-900 border-2 border-green-500 mb-6 flex items-center justify-center overflow-hidden">
        {/* Tampilkan gambar jika ada, jika tidak, tampilkan emoji */}
        {display?.thumbnail?.url ? (
          <img 
            src={display.thumbnail.url} 
            alt={display.name}
            className="w-full h-full object-cover"
            style={{ imageRendering: 'pixelated' }} // Menjaga style retro
          />
        ) : (
          <div className="text-6xl">🎮</div>
        )}
      </div>
      <p className="text-green-400 text-md mb-1 pixel-text">
        {/* Tampilkan nama NFT, jika tidak ada, tampilkan ID */}
        {display?.name ?? `Accessory #${id}`}
      </p>
      <p className="text-green-400 text-sm mb-6 pixel-text">
        {/* Tampilkan nama NFT, jika tidak ada, tampilkan ID */}
        #{id ?? "id"}
      </p>
      <button className="pixel-button w-full">
        [ DETAIL ]
      </button>
    </div>
  );
};

export default AccessoryCard;