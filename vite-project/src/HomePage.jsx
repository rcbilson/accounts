import React, { useState } from 'react';
import { Stack } from '@mui/material';
import { Typography } from '@mui/material';

import PermanentDrawer from './PermanentDrawer.jsx';
import CategoryWidget from './CategoryWidget.jsx';
import SummaryWidget from './SummaryWidget.jsx';
import DateBuilder from './DateBuilder.jsx';

function queryDate(y, m, d) {
  while (m <= 0) {
    y -= 1;
    m += 12;
  }
  return `${y}-${m.toString().padStart(2, '0')}-${d.toString().padStart(2, '0')}`
}

function summarySince(date) {
  return (
    <Stack>
      <Typography>since {date}</Typography>
      <SummaryWidget showAmounts={false} querySpec={{dateFrom: date}}/>
    </Stack>
  )
}

export default function HomePage() {
  const [querySpec, setQuerySpec] = useState({limit: 10});
  const now = new Date();
  const thisYear = now.getFullYear();
  const thisMonth = now.getMonth();
  const qLastMonth = queryDate(thisYear, thisMonth-1, 1);
  const qLast3Month = queryDate(thisYear, thisMonth-3, 1);
  const qLast6Month = queryDate(thisYear, thisMonth-6, 1);
  const qLast12Month = queryDate(thisYear-1, thisMonth, 1);

  return (
    <Stack direction='row'>
      <PermanentDrawer />
      <Stack>
        <DateBuilder querySpec={querySpec} setQuerySpec={setQuerySpec} />
        <Stack direction='row'>
          <SummaryWidget querySpec={querySpec}/>
          <CategoryWidget querySpec={querySpec}/>
        </Stack>
        <hr/>
        <Stack direction='row'>
          {summarySince(qLastMonth)}
          {summarySince(qLast3Month)}
          {summarySince(qLast6Month)}
          {summarySince(qLast12Month)}
        </Stack>
      </Stack>
    </Stack>
  )
}
