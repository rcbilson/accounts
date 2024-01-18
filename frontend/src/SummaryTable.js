import React from 'react';
import { DataGrid, GridColDef } from '@mui/x-data-grid';

export default function SummaryTable({summary}) {
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
    {
      field: 'percent',
      headerName: '%',
      type: 'number',
      width: 100,
      editable: false,
    },
  ];

  const withPercent = summary.amounts.map((x) => {
    return {...x, percent: x.amount / summary.income * 100}
  })
  const rows=[
    {
      category: 'Income',
      amount: summary.income
    },
    ...withPercent
  ];
  return (
  <div style={{ height: 400, width: 400 }}>
    <DataGrid
      rows={rows}
      columns={columns}
      getRowId={(row) => row.category}
      autoPageSize
    />
  </div>
  );
}
