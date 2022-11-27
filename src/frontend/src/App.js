import React, { useState } from 'react';
import { LocalizationProvider } from '@mui/x-date-pickers';
import { AdapterLuxon } from '@mui/x-date-pickers/AdapterLuxon';
import { Stack } from '@mui/material';

import Table from './Table.js';
import QueryBuilder from './QueryBuilder.js';

export default function App() {
  const [querySpec, setQuerySpec] = useState({
    descrLike: "",
    category: "",
    subcategory: "",
  });

  return (
    <LocalizationProvider dateAdapter={AdapterLuxon}>
      <Stack sx={{ height: '100vh', width: '100%' }}>
        <QueryBuilder querySpec={querySpec} setQuerySpec={setQuerySpec} />
        <Table querySpec={querySpec} />
      </Stack>
    </LocalizationProvider>
  );
}
