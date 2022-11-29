import React from 'react';
import Box from '@mui/material/Box';
import { DataGrid, GridColDef } from '@mui/x-data-grid';

import { transactionUpdate } from './Transaction.js';

const columns: GridColDef[] = [
  {
    field: 'date',
    headerName: 'Date',
    width: 150,
    editable: false,
  },
  {
    field: 'descr',
    headerName: 'Description',
    width: 450,
    editable: false,
  },
  {
    field: 'amount',
    headerName: 'Amount',
    type: 'number',
    width: 150,
    editable: false,
  },
  {
    field: 'category',
    headerName: 'Category',
    width: 250,
    editable: true,
  },
  {
    field: 'subcategory',
    headerName: 'Subcategory',
    width: 250,
    editable: true,
  },
];

export default function GridView({rows}) {

  const processRowUpdate = async (newRow) => {
    await transactionUpdate(newRow);
    return newRow
  }

  const handleProcessRowUpdateError = (error: Error) => {
    alert(error.message);
  }

  return (
    <Box sx={{ height: '100%', width: '100%' }}>
      <DataGrid
        rows={rows}
        columns={columns}
        autoPageSize
        disableSelectionOnClick
        experimentalFeatures={{ newEditingApi: true }}
        processRowUpdate={processRowUpdate}
        onProcessRowUpdateError={handleProcessRowUpdateError}
      />
    </Box>
  );
}
