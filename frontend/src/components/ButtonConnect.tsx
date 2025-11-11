'use client';

import { useEffect } from 'react';
import {
  useFlowQuery,
  useFlowMutate,
  useFlowTransactionStatus,
  useFlowCurrentUser,
} from '@onflow/react-sdk';

export default function ButtonConnect() {
  const { user, authenticate, unauthenticate } = useFlowCurrentUser();

  const {
    mutate: increment,
    isPending: txPending,
    data: txId,
    error: txError,
  } = useFlowMutate();

  const { transactionStatus, error: txStatusError } = useFlowTransactionStatus({
    id: txId || '',
  });

  // useEffect(() => {
  //   if (txId && transactionStatus?.status === 3) {
  //     // Transaction is executed
  //     refetch(); // Refresh the counter
  //   }
  // }, [transactionStatus?.status, txId, refetch]);

  return (
    <div>
      {user?.loggedIn ? (
        <div>
          <p>Connected: {user.addr}</p>

          <button onClick={unauthenticate}>Disconnect</button>

          {transactionStatus?.statusString && transactionStatus?.status && (
            <p>
              Status: {transactionStatus.status >= 3 ? 'Successful' : 'Pending'}
            </p>
          )}

          {txError && <p>Error: {txError.message}</p>}

          {txStatusError && <p>Status Error: {txStatusError.message}</p>}
        </div>
      ) : (
        <button onClick={authenticate}>Connect Wallet</button>
      )}
    </div>
  );
}
