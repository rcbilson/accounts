import React, { useState } from 'react';
import { Stack } from '@mui/material';

import Table from './Table.js';
import QueryBuilder from './QueryBuilder.js';

export default function App() {
  const [querySpec, setQuerySpec] = useState({
    dateFrom: null,
    dateUntil: null,
    descrLike: "",
    category: "",
    subcategory: "",
  });

  const [totalValue, setTotalValue] = useState(0)

  return (
    <Stack sx={{ height: '100vh', width: '100%' }}>
      <QueryBuilder querySpec={querySpec} setQuerySpec={setQuerySpec} totalValue={totalValue} />
      <Table querySpec={querySpec} setTotalValue={setTotalValue} />
    </Stack>
  );
}
