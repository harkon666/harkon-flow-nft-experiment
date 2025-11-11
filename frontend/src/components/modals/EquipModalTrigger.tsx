'use client';

import React, { useState } from 'react';
import EquipModal from './EquipModal';

// Props-nya sekarang jauh lebih sedikit
interface EquipModalTriggerProps {
  moment: any; // Momen SPESIFIK yang di-equip
  onTransactionSuccess: () => void;
  ownerAddress: string
}

const EquipModalTrigger: React.FC<EquipModalTriggerProps> = ({
  moment,
  onTransactionSuccess,
  ownerAddress
}) => {
  const [isModalOpen, setIsModalOpen] = useState(false);
  
  return (
    <>
      <button
        onClick={() => setIsModalOpen(true)}
        className="pixel-button w-full text-xs mt-2"
      >
        {/* {isEquipPending ? 'Loading...' : '[ EQUIP ]'} */}
        [ EQUIP ]
      </button>

      {/* Modal sekarang dirender DI SINI,
        dan ia akan mengambil datanya sendiri saat 'isModalOpen' menjadi true 
      */}
      {isModalOpen && (
        <div className="fixed">
          <EquipModal
            isOpen={isModalOpen}
            onClose={() => setIsModalOpen(false)}
            onTransactionSuccess={onTransactionSuccess}
            moment={moment}
            ownerAddress={ownerAddress}
          />
        </div>
      )}
    </>
  );
};

export default EquipModalTrigger;