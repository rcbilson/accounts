import React from 'react';
import { DataGrid, GridColDef } from '@mui/x-data-grid';

export default function AggregateTable({items}) {
  const columns: GridColDef[] = [
    {
      field: 'category',
      headerName: 'Category',
      width: 150,
      editable: false,
    },
    {
      field: 'amount',
      headerName: 'Amount',
      type: 'number',
      width: 150,
      editable: false,
    },
  ];

  return (
  <div style={{ height: 400, width: 300 }}>
    <DataGrid
      rows={items}
      columns={columns}
      getRowId={(row) => row.category}
      autoPageSize
    />
  </div>
  );
}
