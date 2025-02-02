import * as React from 'react';
import { TextField, Stack } from '@mui/material';

export default function QueryBuilder({querySpec, setQuerySpec}) {
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
    </Stack>
  )
}
