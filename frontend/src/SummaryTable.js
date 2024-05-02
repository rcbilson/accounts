import React from 'react';
import { DataGrid, GridColDef } from '@mui/x-data-grid';

export default function SummaryTable({summary, showAmounts}) {
  if (!summary || !summary.amounts) {
    return <div>Loading...</div>
  }
  let columns: GridColDef[] = [];
  columns.push({
      field: 'category',
      headerName: 'Category',
      width: 150,
      editable: false,
    });
  if (showAmounts) {
    columns.push({
        field: 'amount',
        headerName: 'Amount',
        type: 'number',
        width: 150,
        editable: false,
      });
  }
  columns.push({
      field: 'percent',
      headerName: '%',
      type: 'number',
      width: 100,
      editable: false,
    });

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
  <div style={{ height: 400 }}>
    <DataGrid
      rows={showAmounts ? rows : withPercent}
      columns={columns}
      getRowId={(row) => row.category}
      autoPageSize
    />
  </div>
  );
}
