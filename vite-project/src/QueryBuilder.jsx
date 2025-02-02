import * as React from 'react';
import { TextField, Stack } from '@mui/material';

export default function QueryBuilder({querySpec, setQuerySpec, totalValue}) {
  return (
    <Stack direction="row">
      <TextField
        id="dateFrom"
	label="On or after date"
        type="date"
        sx={{ width: 220 }}
        InputLabelProps={{
          shrink: true,
        }}
        onChange={(e) => setQuerySpec({
          ...querySpec,
          dateFrom: e.target.value,
        })}
      />
      <TextField
        id="dateUntil"
	label="Before date"
        type="date"
        sx={{ width: 220 }}
        InputLabelProps={{
          shrink: true,
        }}
        onChange={(e) => setQuerySpec({
          ...querySpec,
          dateUntil: e.target.value,
        })}
      />
      <TextField
        id="descrLike"
        label="Description"
        value={querySpec.descrLike}
        onChange={(e) => setQuerySpec({
          ...querySpec,
          descrLike: e.target.value,
        })}
     />
     <TextField
        id="category"
        label="Category"
        value={querySpec.category}
        onChange={(e) => setQuerySpec({
          ...querySpec,
          category: e.target.value,
        })}
      />
     <TextField
        id="subcategory"
        label="Subcategory"
        value={querySpec.subcategory}
        onChange={(e) => setQuerySpec({
          ...querySpec,
          subcategory: e.target.value,
        })}
      />
     <TextField
        id="totalValue"
        label="Total Amount"
        value={totalValue.toFixed(2)}
        disabled
      />
    </Stack>
  )
}
