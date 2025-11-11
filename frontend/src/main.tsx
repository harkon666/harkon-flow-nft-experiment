import { StrictMode } from 'react';
import { createRoot } from 'react-dom/client';
import App from './App.tsx';
import './index.css';
import { FlowProvider } from '@onflow/react-sdk';
import { config } from '@onflow/fcl';
import flowJson from '../../flow.json';

config()
  .put('accessNode.api', 'http://localhost:8888')
  .put('discovery.wallet', 'http://localhost:8701/fcl/authn');

createRoot(document.getElementById('root')!).render(
  <StrictMode>
    <FlowProvider
      config={{
        accessNodeUrl: 'http://localhost:8888',
        flowNetwork: 'emulator',
        discoveryWallet: 'https://fcl-discovery.onflow.org/emulator/authn',
      }}
      flowJson={flowJson}
    >
      <App />
    </FlowProvider>
  </StrictMode>
);
