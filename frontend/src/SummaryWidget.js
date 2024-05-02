import React, { useEffect, useState } from 'react';

import * as Transaction from './Transaction.js';
import SummaryTable from './SummaryTable.js';

export default function SummaryWidget({querySpec, showAmounts=true}) {
  const [error, setError] = useState(null);
  const [isLoaded, setIsLoaded] = useState(false);
  const [summary, setSummary] = useState({});

  useEffect(() => {
    Transaction.Summary(querySpec)
      .then(
        (result) => {
          setIsLoaded(true);
          setSummary(result);
        },
        // Note: it's important to handle errors here
        // instead of a catch() block so that we don't swallow
        // exceptions from actual bugs in components.
        (error) => {
          setIsLoaded(true);
          setError(error);
        }
      )
  }, [querySpec])

  if (error) {
    return <div>Error: {error.message}</div>;
  } else if (!isLoaded) {
    return <div>Loading...</div>;
  } else {
    return (
      <SummaryTable summary={summary} showAmounts={showAmounts} />
    )
  }
}
