import React, { useEffect, useState } from 'react';
import { Stack } from '@mui/material';

import NewTable from './NewTable.js';
import DragAndDrop from './DragAndDrop.js';
import * as Transaction from './Transaction.js'

export default function NewPage() {
  const [error, setError] = useState(null);
  const [isLoaded, setIsLoaded] = useState(false);
  const [items, setItems] = useState([]);

  useEffect(() => {
    Transaction.Query({ state: "new" })
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
  }, [])

  const handleUpdate = async (newRow) => {
    await Transaction.Update(newRow);
    setItems((prevItems) => prevItems.filter((item) => item.id !== newRow.id));
    return newRow
  }

  const handleDelete = async (id) => {
    await Transaction.Delete(id)
    setItems((prevItems) => prevItems.filter((item) => item.id !== id));
  }

  if (error) {
    return <div>Error: {error.message}</div>;
  } else if (!isLoaded) {
    return <div>Loading...</div>;
  } else {
    return (
      <DragAndDrop>
        <Stack sx={{ height: '100vh', width: '100%' }}>
          <NewTable items={items} onUpdate={handleUpdate} onDelete={handleDelete} />
        </Stack>
      </DragAndDrop>
    );
  }
}
