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
    </Stack>
  )
}
