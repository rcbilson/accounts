import React from 'react';
import Box from '@mui/material/Box';
import { DataGrid, GridColDef, GridActionsCellItem } from '@mui/x-data-grid';
import DeleteIcon from '@mui/icons-material/Delete';
import CheckIcon from '@mui/icons-material/Check';

export default function NewTable({items, onUpdate, onDelete, onAccept}) {
  function processRowUpdate(newRow, oldRow) {
    onUpdate(newRow);
    return newRow;
  }

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
    {
      field: 'actions',
      type: 'actions',
      width: 80,
      getActions: (params) => [
        <GridActionsCellItem
          icon={<CheckIcon />}
          label="Accept"
          onClick={() => onAccept(params.row)}
        />,
        <GridActionsCellItem
          icon={<DeleteIcon />}
          label="Delete"
          onClick={() => onDelete(params.id)}
        />,
      ]
    },
  ];

  return (
    <Box sx={{ height: '100%', width: '100%' }}>
      <DataGrid
        rows={items}
        columns={columns}
        autoPageSize
        disableSelectionOnClick
        experimentalFeatures={{ newEditingApi: true }}
        processRowUpdate={processRowUpdate}
      />
    </Box>
  );
}
