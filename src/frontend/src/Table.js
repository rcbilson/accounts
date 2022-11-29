import Box from '@mui/material/Box';
import { DataGrid, GridColDef } from '@mui/x-data-grid';

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
  return (
    <Box sx={{ height: '100%', width: '100%' }}>
      <DataGrid
        rows={rows}
        columns={columns}
        autoPageSize
        disableSelectionOnClick
        experimentalFeatures={{ newEditingApi: true }}
      />
    </Box>
  );
}
