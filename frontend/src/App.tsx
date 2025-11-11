import { useState } from 'react';
import { Wallet, Package, ShoppingCart, Check, X } from 'lucide-react';
import {
  useFlowCurrentUser,
} from '@onflow/react-sdk'
import ButtonConnect from './components/ButtonConnect';
import InventorySection from './components/sections/InventorySection';
import YourSalesSection from './components/sections/YourSalesSection';

interface AccessoryListing {
  id: number;
  name: string;
  price: string;
  image: string;
}

interface OwnedAccessory {
  id: string;
  name: string;
  image: string;
  quantity: number;
}

interface MomentNFT {
  id: string;
  name: string;
  image: string;
  equippedAccessories: string[];
}

function App() {
  const [isConnected, setIsConnected] = useState(false);
  const [showGachaResult, setShowGachaResult] = useState(false);
  const [gachaResult, setGachaResult] = useState('');
  const [isShaking, setIsShaking] = useState(false);
  const [ownedAccessories, setOwnedAccessories] = useState<OwnedAccessory[]>([
    { id: '1', name: 'Bingkai Emas', image: '🖼️', quantity: 2 },
    { id: '2', name: 'Stiker Langka', image: '⭐', quantity: 1 },
    { id: '3', name: 'Filter Retro', image: '📷', quantity: 3 },
  ]);
  const [userMoments, setUserMoments] = useState<MomentNFT[]>([
    { id: '1', name: 'Moment #001', image: '🎮', equippedAccessories: [] },
    { id: '2', name: 'Moment #002', image: '🎲', equippedAccessories: [] },
  ]);
  const [selectedMoment, setSelectedMoment] = useState<string>('1');
  const [selectedAccessoryToEquip, setSelectedAccessoryToEquip] = useState<string>('');

  const mockListings: AccessoryListing[] = [
    { id: 1, name: 'Bingkai Emas', price: '1.5', image: '🖼️' },
    { id: 2, name: 'Stiker Langka', price: '2.0', image: '⭐' },
    { id: 3, name: 'Filter Retro', price: '0.8', image: '📷' },
    { id: 4, name: 'Efek Pixel', price: '1.2', image: '🎨' },
    { id: 5, name: 'Badge Limited', price: '3.0', image: '🏆' },
    { id: 6, name: 'Stempel VIP', price: '2.5', image: '💎' },
  ];

  const { user, authenticate, unauthenticate } = useFlowCurrentUser();

  const handleBuyGacha = () => {
    setIsShaking(true);
    setShowGachaResult(false);

    setTimeout(() => {
      setIsShaking(false);
      const accessories = ['Bingkai Emas', 'Stiker Langka', 'Filter Retro', 'Efek Pixel', 'Badge Limited', 'Stempel VIP'];
      const randomAccessory = accessories[Math.floor(Math.random() * accessories.length)];
      setGachaResult(randomAccessory);
      setShowGachaResult(true);

      const newAccessory = ownedAccessories.find(a => a.name === randomAccessory);
      if (newAccessory) {
        setOwnedAccessories(ownedAccessories.map(a =>
          a.id === newAccessory.id ? { ...a, quantity: a.quantity + 1 } : a
        ));
      } else {
        const accessoryData = mockListings.find(a => a.name === randomAccessory);
        if (accessoryData) {
          setOwnedAccessories([...ownedAccessories, {
            id: String(accessoryData.id),
            name: accessoryData.name,
            image: accessoryData.image,
            quantity: 1
          }]);
        }
      }
    }, 1000);
  };

  const handleEquipAccessory = () => {
    if (!selectedAccessoryToEquip || !selectedMoment) return;

    const currentMoment = userMoments.find(m => m.id === selectedMoment);
    if (currentMoment && !currentMoment.equippedAccessories.includes(selectedAccessoryToEquip)) {
      setUserMoments(userMoments.map(m =>
        m.id === selectedMoment
          ? { ...m, equippedAccessories: [...m.equippedAccessories, selectedAccessoryToEquip] }
          : m
      ));
      setSelectedAccessoryToEquip('');
    }
  };

  const handleRemoveAccessory = (accessoryId: string) => {
    setUserMoments(userMoments.map(m =>
      m.id === selectedMoment
        ? { ...m, equippedAccessories: m.equippedAccessories.filter(a => a !== accessoryId) }
        : m
    ));
  };

  return (
    <div className="min-h-screen bg-black text-white">
      <style>{`
        @import url('https://fonts.googleapis.com/css2?family=Press+Start+2P&family=Courier+Prime&display=swap');

        * {
          image-rendering: pixelated;
          image-rendering: -moz-crisp-edges;
          image-rendering: crisp-edges;
        }

        h1, h2, h3, button, .pixel-text {
          font-family: 'Press Start 2P', cursive;
          text-transform: uppercase;
        }

        body, p, span {
          font-family: 'Courier Prime', monospace;
        }

        .pixel-border {
          border: 3px solid #00ef8b;
          border-radius: 0;
          box-shadow: 0 0 0 3px black, 0 0 0 6px #00ef8b;
        }

        .pixel-button {
          font-family: 'Press Start 2P', cursive;
          background: #00ef8b;
          color: black;
          border: 3px solid #00ef8b;
          border-radius: 0;
          padding: 16px 24px;
          cursor: pointer;
          transition: all 0.1s;
          text-transform: uppercase;
          font-size: 12px;
        }

        .pixel-button:hover {
          background: black;
          color: #00ef8b;
          box-shadow: 0 0 20px #00ef8b;
        }

        .pixel-button:active {
          transform: translate(2px, 2px);
        }

        .pixel-card {
          background: #0a0a0a;
          border: 2px solid #00ef8b;
          border-radius: 0;
          padding: 16px;
          transition: all 0.2s;
        }

        .pixel-card:hover {
          box-shadow: 0 0 30px rgba(0, 239, 139, 0.3);
          transform: translateY(-4px);
        }

        .shake {
          animation: shake 0.5s infinite;
        }

        @keyframes shake {
          0%, 100% { transform: rotate(0deg); }
          25% { transform: rotate(-5deg); }
          75% { transform: rotate(5deg); }
        }

        .glow {
          text-shadow: 0 0 10px #00ef8b, 0 0 20px #00ef8b, 0 0 30px #00ef8b;
        }

        .scan-line {
          position: relative;
          overflow: hidden;
        }

        .scan-line::before {
          content: '';
          position: absolute;
          top: 0;
          left: 0;
          width: 100%;
          height: 2px;
          background: rgba(0, 239, 139, 0.1);
          animation: scan 4s linear infinite;
        }

        @keyframes scan {
          0% { transform: translateY(0); }
          100% { transform: translateY(100vh); }
        }

        .pixel-img {
          width: 100%;
          height: 100%;
          object-fit: cover;
          image-rendering: pixelated;
        }

        .modal-open .scan-line::before {
          animation: none;
        }

        .modal-open .pixel-card:hover {
          transform: none;
          box-shadow: none;
        }
      `}</style>

      <div className="scan-line">
        <header className="border-b-4 border-green-500 bg-black sticky top-0 z-50">
          <div className="container mx-auto px-4 py-4 flex justify-between items-center">
            <div className="flex gap-8">
              <a className="text-green-500 text-xl glow">Harkon-NFT-Eksperimen</a>
              <a className="text-green-500 text-xl glow">My Inventory</a>
            </div>
            <ButtonConnect />
            {/* <button
              onClick={authenticate}
              className="pixel-button flex items-center gap-2"
            >
              <Wallet size={16} />
              {isConnected ? '[ Connected ]' : '[ Connect Wallet ]'}
            </button> */}
          </div>
        </header>
        
        <section className="container mx-auto px-4 py-16 text-center">
          <div className="max-w-3xl mx-auto">
            <h2 className="text-3xl text-green-500 mb-8 glow">Welcome to Harkon-NFT</h2>
            <p className="text-lg text-green-300 leading-relaxed">
              Pixelated moments from your favorite events. Mint your Moment,
              customize it with accessories, and trade on the market.
              Built on Flow Blockchain with retro vibes.
            </p>
          </div>
        </section>
        
        <InventorySection />

        <section className="container mx-auto px-4 py-16">
          <h2 className="text-2xl text-green-500 text-center mb-12 glow">Mint Your Moment</h2>
          <div className="max-w-md mx-auto">
            <div className="pixel-card text-center">
              <div className="w-full aspect-square bg-gray-900 border-2 border-green-500 mb-6 flex items-center justify-center">
                <div className="text-6xl">🎮</div>
              </div>
              <p className="text-green-400 text-xl mb-6 pixel-text">Price: 5.0 FLOW</p>
              <button className="pixel-button w-full">
                [ MINT MOMENT ]
              </button>
            </div>
          </div>
        </section>

        <section className="container mx-auto px-4 py-16 bg-gradient-to-b from-black to-gray-900">
          <h2 className="text-2xl text-green-500 text-center mb-12 glow">Gacha Accessory Pack</h2>
          <div className="max-w-md mx-auto">
            <div className="pixel-card text-center">
              <div className={`w-full aspect-square bg-gray-900 border-2 border-green-500 mb-6 flex items-center justify-center ${isShaking ? 'shake' : ''}`}>
                <Package size={120} className="text-green-500" />
              </div>
              <p className="text-green-400 text-xl mb-6 pixel-text">Price: 1.0 FLOW</p>
              <button
                onClick={handleBuyGacha}
                className="pixel-button w-full mb-4"
              >
                [ BUY GACHA PACK ]
              </button>

              {showGachaResult && (
                <div className="mt-6 p-4 border-2 border-green-500 bg-black">
                  <p className="text-green-500 mb-2 pixel-text text-sm">You got:</p>
                  <p className="text-green-400 text-lg glow">{gachaResult}!</p>
                </div>
              )}
            </div>
          </div>
        </section>

        <YourSalesSection />

        <section className="container mx-auto px-4 py-16">
          <h2 className="text-2xl text-green-500 text-center mb-12 glow">Accessory Marketplace</h2>
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
            {mockListings.map((listing) => (
              <div key={listing.id} className="pixel-card">
                <div className="w-full aspect-square bg-gray-900 border-2 border-green-500 mb-4 flex items-center justify-center">
                  <div className="text-6xl">{listing.image}</div>
                </div>
                <h3 className="text-green-400 text-sm mb-3 pixel-text">{listing.name}</h3>
                <p className="text-green-300 mb-4">{listing.price} FLOW</p>
                <button className="pixel-button w-full text-xs flex items-center justify-center gap-2">
                  <ShoppingCart size={12} />
                  [ BUY ]
                </button>
              </div>
            ))}
          </div>
        </section>

        <section className="container mx-auto px-4 py-16 bg-gradient-to-b from-black to-gray-900">
          <h2 className="text-2xl text-green-500 text-center mb-12 glow">Equip Accessory to Moment</h2>
          <div className="max-w-4xl mx-auto">
            <div className="grid grid-cols-1 lg:grid-cols-2 gap-8">
              <div className="pixel-card">
                <h3 className="text-green-400 text-sm mb-6 pixel-text">Select Your Moment</h3>
                <div className="grid grid-cols-2 gap-4 mb-6">
                  {userMoments.map((moment) => (
                    <button
                      key={moment.id}
                      onClick={() => setSelectedMoment(moment.id)}
                      className={`pixel-card p-4 text-center transition-all ${selectedMoment === moment.id ? 'border-4 border-green-500 box-shadow: 0 0 20px rgba(0, 239, 139, 0.5)' : 'border-2'}`}
                    >
                      <div className="text-5xl mb-2">{moment.image}</div>
                      <p className="text-green-300 text-xs">{moment.name}</p>
                    </button>
                  ))}
                </div>
              </div>

              <div className="pixel-card">
                <h3 className="text-green-400 text-sm mb-6 pixel-text">Current Equipment</h3>
                {selectedMoment ? (
                  <div>
                    {userMoments.find(m => m.id === selectedMoment)?.equippedAccessories.length === 0 ? (
                      <p className="text-green-300 mb-6">No accessories equipped yet</p>
                    ) : (
                      <div className="grid grid-cols-3 gap-3 mb-6">
                        {userMoments.find(m => m.id === selectedMoment)?.equippedAccessories.map((accId) => {
                          const accessory = ownedAccessories.find(a => a.id === accId);
                          return accessory ? (
                            <div key={accId} className="relative">
                              <div className="pixel-card p-2 text-center border-2 border-green-500">
                                <div className="text-3xl">{accessory.image}</div>
                              </div>
                              <button
                                onClick={() => handleRemoveAccessory(accId)}
                                className="absolute top-0 right-0 bg-black border-2 border-green-500 text-green-500 p-0 w-6 h-6 flex items-center justify-center text-xs hover:bg-green-500 hover:text-black"
                              >
                                <X size={12} />
                              </button>
                            </div>
                          ) : null;
                        })}
                      </div>
                    )}
                  </div>
                ) : null}
              </div>
            </div>

            <div className="mt-8 pixel-card p-6">
              <div className="flex flex-col gap-4">
                <p className="text-green-300">
                  {selectedAccessoryToEquip ? (
                    <span>Selected: <span className="glow">{ownedAccessories.find(a => a.id === selectedAccessoryToEquip)?.name}</span></span>
                  ) : (
                    <span className="text-green-400">Select an accessory from your inventory above</span>
                  )}
                </p>
                <button
                  onClick={handleEquipAccessory}
                  disabled={!selectedAccessoryToEquip || !selectedMoment}
                  className="pixel-button w-full text-sm flex items-center justify-center gap-2 disabled:opacity-50 disabled:cursor-not-allowed"
                >
                  <Check size={16} />
                  [ EQUIP ACCESSORY ]
                </button>
              </div>
            </div>
          </div>
        </section>

        <footer className="border-t-4 border-green-500 bg-black py-8 mt-16">
          <div className="container mx-auto px-4 text-center">
            <p className="text-green-500 text-sm">
              © 2025 Harkon-NFT | Built on Flow Blockchain | Retro Web3
            </p>
          </div>
        </footer>
      </div>
    </div>
  );
}

export default App;
