import React, { useEffect, useState } from 'react';
import { Stack } from '@mui/material';

import QueryTable from './QueryTable.jsx';
import QueryBuilder from './QueryBuilder.jsx';
import * as Transaction from './Transaction.js';
import TemporaryDrawer from './TemporaryDrawer.jsx';

export default function QueryPage() {
  const [querySpec, setQuerySpec] = useState({
    dateFrom: null,
    dateUntil: null,
    descrLike: "",
    category: "",
    subcategory: "",
  });

  const [error, setError] = useState(null);
  const [isLoaded, setIsLoaded] = useState(false);
  const [items, setItems] = useState([]);

  useEffect(() => {
    Transaction.Query(querySpec)
      .then(
        (result) => {
          setIsLoaded(true);
          setItems(result);
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

  const handleUpdate = async (newRow) => {
    await Transaction.Update(newRow);
    setItems((prevItems) => prevItems.map((item) => {
      if (item.id === newRow.id) {
        return newRow;
      }
      return item;
    }))
  }

  const handleDelete = (id) => {
    Transaction.Delete(id)
    setItems((prevItems) => prevItems.filter((item) => item.id !== id));
  }

  if (error) {
    return <div>Error: {error.message}</div>;
  } else if (!isLoaded) {
    return <div>Loading...</div>;
  } else {
    const totalValue = items.reduce((a, c) => a + parseFloat(c.amount), 0);
    return (
      <Stack sx={{ height: '100vh', width: '100%' }}>
        <Stack direction='row'>
          <TemporaryDrawer />
          <QueryBuilder querySpec={querySpec} setQuerySpec={setQuerySpec} totalValue={totalValue} />
        </Stack>
        <QueryTable items={items} onUpdate={handleUpdate} onDelete={handleDelete} />
      </Stack>
    );
  }
}
