'use client';

import React, { useState, useEffect } from 'react';
import { X, Check } from 'lucide-react'; // Asumsi Anda menggunakan lucide-react
import { useGetAccessoryIDs } from '@/hooks/scripts/useGetNFTAccessory'; // Asumsi
import { useGetAccessoryMetadata } from '@/hooks/scripts/useGetAccessoryMetadata';
import { useFlowCurrentUser } from '@onflow/react-sdk';
import { useEquipNFT } from '@/hooks/transactions/useEquipNFT';
import { useCheckEquipment } from '@/hooks/scripts/useCheckEquipment';
import StackImage from '../StackImage';
import { resolveIpfsUrl } from '@/lib/utils';


const AccessoryChoice: React.FC<{ 
  accessoryId: number; 
  ownerAddress: string; 
  isSelected: boolean; 
  onSelect: () => void;
}> = ({ accessoryId, ownerAddress, isSelected, onSelect }) => {
  
  // Ambil metadata untuk 1 aksesori ini
  const { display, isLoading } = useGetAccessoryMetadata({
    address: ownerAddress,
    id: accessoryId
  });

  if (isLoading) {
    return <div className="pixel-card p-2 animate-pulse bg-gray-700 aspect-square"></div>;
  }
  return (
    <button
      onClick={onSelect}
      className={`p-2 text-center transition-all ${isSelected ? 'border-4 border-green-500 glow-light' : 'border-2'}`}
    >
      <div className="text-3xl">{display?.thumbnail ? <img src={resolveIpfsUrl(display.thumbnail)} /> : '✨'}</div>
      <p className="text-green-300 text-xs truncate">{display?.name ?? `ID: ${accessoryId}`}</p>
      <p className="text-green-300 text-xs truncate">{display?.id ?? `ID: ${accessoryId}`}</p>
    </button>
  );
};

interface Moment {
  id: number;
  name: string;
  image: string; // (Berdasarkan kode Anda, ini emoji)
  equippedAccessories: number[];
  thumbnail: string;
}

// Ini adalah semua props yang dibutuhkan modal
interface EquipModalProps {
  isOpen: boolean;
  onClose: () => void;
  moment: Moment; // Momen yang akan di-equip
  onRemove: (momentId: number, accessoryId: number) => void;
  onTransactionSuccess: () => void;
  ownerAddress: string;
}
// --- Akhir Tipe Data ---

const EquipModal: React.FC<EquipModalProps> = ({
  isOpen,
  onClose,
  moment,
  onTransactionSuccess,
  ownerAddress
}) => {
  const { user } = useFlowCurrentUser();
  
  // State HANYA untuk aksesori yang DIPILIH
  const [selectedAccessoryToEquip, setSelectedAccessoryToEquip] = useState<number | null>(null);

  // --- LOGIKA FETCHING BARU DI SINI ---
  // 1. Ambil daftar ID aksesori milik pengguna
  const { 
    data: accessoryIDs, 
    isLoading: isLoadingIDs, 
    error: errorIDs 
  } = useGetAccessoryIDs({
    address: user?.addr ?? ""
  });

  const { equippedFrameData, isLoading: isLoadingFrame, refetch } = useCheckEquipment({
    address: ownerAddress,
    momentId: moment.id
  });
  // ------------------------------------

  // Reset state saat modal ditutup
  useEffect(() => {
    if (!isOpen) {
      setSelectedAccessoryToEquip(null);
      document.body.classList.remove('modal-open');
    }

    if (isOpen) {
      // Tambahkan class ke body
      document.body.classList.add('modal-open');
    }

    // Fungsi 'cleanup' (pembersihan)
    // Ini akan berjalan saat modal ditutup (isOpen = false)
    return () => {
      document.body.classList.remove('modal-open');
    };
  }, [isOpen]);

  const { equip, isPending: isEquipPending, isSealed, error } = useEquipNFT();
  // const { remove, isPending: isRemovePending } = useRemoveNFT();
  useEffect(() => {
    // Jika 'equip' atau 'remove' selesai, panggil callback
    if (isSealed) {
      onTransactionSuccess();
      refetch()
    }
  }, [isSealed, onTransactionSuccess]);

  

  const handleEquipAccessory = () => {
    if (moment && selectedAccessoryToEquip) {
      equip({nftMomentId: moment.id, nftAccessoryId: selectedAccessoryToEquip});
    }
  };

  // const handleRemoveAccessory = (accessoryId: number) => {
  //   if (moment) {
  //     onRemove(moment.id, accessoryId);
  //   }
  // };

  if (!isOpen) {
    return null;
  }
  console.log(equippedFrameData, 'woi equip')
  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center bg-black bg-opacity-75" onClick={onClose}>
      <div
        className="relative z-60 w-full max-w-4xl bg-gray-900 border-2 border-green-500 p-6 glow-light"
        onClick={(e) => e.stopPropagation()}
      >
        <button onClick={onClose} className="absolute top-2 right-2 text-green-500 hover:text-green-300">
          <X size={24} />
        </button>
        <h2 className="text-2xl text-green-500 text-center mb-8 glow">
          Equip: {moment.name}
        </h2>

        <div className="grid grid-cols-1 lg:grid-cols-2 gap-8">
          
          {/* Kolom 1: Momen yang sedang di-equip (sudah dipilih) */}
          <div className="pixel-card">
            <StackImage frame={equippedFrameData?.thumbnail?.url || ""} moment={moment?.thumbnail}/>
            <h3 className="text-green-400 text-sm mb-6 pixel-text">Currently Equipped</h3>
            <div className="flex justify-center mb-4">
              <div className="pixel-card p-4 text-center border-4 border-green-500 glow-light w-1/2">
                {equippedFrameData ? 
                <>
                    <p className="text-green-300 text-xs truncate">{equippedFrameData?.name}</p>
                </> : "No Equipped Frame"}
              </div>
            </div>
          </div>

          {/* Kolom 2: Pilih Aksesori (Sekarang Fetching di Sini) */}
          <div className="pixel-card">
            <h3 className="text-green-400 text-sm mb-6 pixel-text">Select Accessory to Equip</h3>
            
            {isLoadingIDs && <p>Loading accessories...</p>}
            {errorIDs && <p className="text-red-500">Error loading accessories.</p>}
            
            <div className="grid grid-cols-3 gap-4 mb-6 max-h-48 overflow-y-auto">
              {!isLoadingIDs && accessoryIDs && user?.addr && (
                accessoryIDs.map((id) => (
                  <AccessoryChoice
                    key={id}
                    accessoryId={id}
                    ownerAddress={user.addr ?? ""}
                    isSelected={selectedAccessoryToEquip === id}
                    onSelect={() => setSelectedAccessoryToEquip(id)}
                  />
                ))
              )}
            </div>
          </div>
        </div>

        {/* Tombol Aksi */}
        <div className="mt-8 pixel-card p-6">
          <button
            onClick={handleEquipAccessory}
            disabled={!selectedAccessoryToEquip || !moment}
            className="pixel-button w-full text-sm ...">
            <Check size={16} />
            {isEquipPending ? "...Equipping" : error ? "[ ERROR ]" : `[ EQUIP SELECTED ACCESSORY ]`}
          </button>
        </div>
        
      </div>
    </div>
  );
};

export default EquipModal;