import React, { useEffect, useState } from 'react';

import * as Transaction from './Transaction.js';
import SummaryChart from './SummaryChart.jsx';

export default function SummaryChartWidget({querySpec}) {
  const [error, setError] = useState(null);
  const [isLoaded, setIsLoaded] = useState(false);
  const [chart, setChart] = useState({});

  useEffect(() => {
    Transaction.SummaryChart(querySpec)
      .then(
        (result) => {
          setIsLoaded(true);
          setChart(result);
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
      <SummaryChart dataset={chart.amounts} />
    )
  }
}