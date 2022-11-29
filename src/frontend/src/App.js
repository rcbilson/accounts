import React, { useEffect, useState } from 'react';
import { Stack } from '@mui/material';

import Table from './Table.js';
import QueryBuilder from './QueryBuilder.js';
import { transactionQuery } from './Transaction.js'

export default function App() {
  const [querySpec, setQuerySpec] = useState({
    dateFrom: null,
    dateUntil: null,
    descrLike: "",
    category: "",
    subcategory: "",
  });

  const [totalValue, setTotalValue] = useState(0)

  const [error, setError] = useState(null);
  const [isLoaded, setIsLoaded] = useState(false);
  const [items, setItems] = useState([]);

  useEffect(() => {
    transactionQuery(querySpec)
      .then(
        (result) => {
          setIsLoaded(true);
          setItems(result);
          setTotalValue(result.reduce((a, c) => a + parseFloat(c.amount), 0));
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
      <Stack sx={{ height: '100vh', width: '100%' }}>
        <QueryBuilder querySpec={querySpec} setQuerySpec={setQuerySpec} totalValue={totalValue} />
        <Table rows={items} />
      </Stack>
    );
  }
}
