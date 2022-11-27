import React, { useEffect, useState } from 'react';
import Box from '@mui/material/Box';
import { DataGrid, GridColDef } from '@mui/x-data-grid';

const columns: GridColDef[] = [
  { field: 'id', headerName: 'ID', width: 90 },
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

export default function TransactionProvider({querySpec}) {
  const [error, setError] = useState(null);
  const [isLoaded, setIsLoaded] = useState(false);
  const [items, setItems] = useState([]);

  // Note: the empty deps array [] means
  // this useEffect will run once
  // similar to componentDidMount()
  useEffect(() => {
    console.log(querySpec)
    let specs = []
    if (querySpec.descrLike !== "") {
      specs.push("DescrLike=" + querySpec.descrLike);
    }
    if (querySpec.category !== "") {
      specs.push("Category=" + querySpec.category);
    }
    if (querySpec.subcategory !== "") {
      specs.push("Subcategory=" + querySpec.subcategory);
    }
    let path = "/api/transactions";
    if (specs.length > 0) {
      path += "?" + specs[0]
      specs.slice(1).forEach((e) => { path += "&" + e })
    }
    fetch(encodeURI(path))
      .then(res => res.json())
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

  if (error) {
    return <div>Error: {error.message}</div>;
  } else if (!isLoaded) {
    return <div>Loading...</div>;
  } else {
    return GridView(items);
  }
}

function GridView(rows) {
  return (
    <Box sx={{ height: '100%', width: '100%' }}>
      <DataGrid
        rows={rows}
        columns={columns}
        autoPageSize
        rowsPerPageOptions={[25]}
        checkboxSelection
        disableSelectionOnClick
        experimentalFeatures={{ newEditingApi: true }}
      />
    </Box>
  );
}
