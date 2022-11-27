import * as React from 'react';
import { TextField, Stack } from '@mui/material';

export default function QueryBuilder({querySpec, setQuerySpec}) {
  return (
    <Stack direction="row">
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
    </Stack>
  )
}
