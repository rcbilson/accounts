import React, { useState } from 'react';
import { Stack } from '@mui/material';

import PermanentDrawer from './PermanentDrawer.js';
import CategoryWidget from './CategoryWidget.js';
import SummaryWidget from './SummaryWidget.js';
import DateBuilder from './DateBuilder.js';

export default function HomePage() {
  const [querySpec, setQuerySpec] = useState({limit: 10});

  return (
    <Stack direction='row'>
      <PermanentDrawer />
      <Stack>
        <DateBuilder querySpec={querySpec} setQuerySpec={setQuerySpec} />
        <Stack direction='row'>
          <SummaryWidget querySpec={querySpec}/>
          <CategoryWidget querySpec={querySpec}/>
        </Stack>
      </Stack>
    </Stack>
  )
}
