import React, { useEffect, useState } from 'react';
import { Stack } from '@mui/material';

import PermanentDrawer from './PermanentDrawer.js';
import CategoryWidget from './CategoryWidget.js';
import DateBuilder from './DateBuilder.js';

export default function HomePage() {
  const [querySpec, setQuerySpec] = useState({limit: 10});

  return (
    <Stack direction='row'>
      <PermanentDrawer />
      <Stack>
        <DateBuilder querySpec={querySpec} setQuerySpec={setQuerySpec} />
        <CategoryWidget querySpec={querySpec}/>
      </Stack>
    </Stack>
  )
}
