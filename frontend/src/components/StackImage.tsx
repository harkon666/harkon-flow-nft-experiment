'use client'; // (Jika Anda menggunakan Next.js App Router, tambahkan ini)

import React from 'react';
import BingkaiKayu from '@/assets/Bingkai Kayu.png' 
import Gajah from "@/assets/gajah.png"
import RainEffect from "@/assets/rain-effect.gif"
import { resolveIpfsUrl } from '@/lib/utils';

interface StackImageProps {
  moment: string
  frame: string
}

const StackImage: React.FC<StackImageProps> = ({ frame, moment }) => {
  console.log(frame, 'frame')
  return (
    <section className="container mx-auto px-4 py-8">
      
      <div className="max-w-md mx-auto">
        
        {/* --- 4. AREA PREVIEW (KANVAS) --- */}
        {/* Di sinilah keajaiban CSS terjadi */}
        <div className="relative w-full aspect-square bg-gray-100 rounded-xl overflow-hidden border-2 border-gray-200 shadow-lg">
          {/* <img
            src={RainEffect}
            alt="effect accessory"
            // 'pointer-events-none' agar bingkai tidak bisa diklik
            className="absolute inset-0 w-full h-full object-cover z-20 pointer-events-none"
          /> */}
          <img
            src={resolveIpfsUrl(moment)}
            alt="Preview Momen"
            className="absolute inset-0 w-full h-full object-cover z-10"
          />

          {frame ?
            <img
              src={resolveIpfsUrl(frame)}
              alt="Bingkai Kayu"
              // 'pointer-events-none' agar bingkai tidak bisa diklik
              className="absolute inset-0 w-full h-full object-cover z-20 pointer-events-none"
            />
          : null}
        </div>
      </div>
    </section>
  );
};

export default StackImage;