import React, { useEffect, useState } from 'react';
import FormControl from '@mui/material/FormControl';
import InputLabel from '@mui/material/InputLabel';
import MenuItem from '@mui/material/MenuItem';
import Select from '@mui/material/Select';
import { Stack } from '@mui/material';

import * as Transaction from './Transaction.js';
import SummaryChart from './SummaryChart.jsx';

export default function SummaryChartWidget({querySpec}) {
  const [error, setError] = useState(null);
  const [isLoaded, setIsLoaded] = useState(false);
  const [chart, setChart] = useState({});
  const [chartType, setChartType] = React.useState('cumulativePercent');

  const handleChange = (event) => {
    setChartType(event.target.value);
  };

  useEffect(() => {
    const summaryChartSpec = {
      "dateFrom": querySpec.dateFrom,
      "dateUntil": querySpec.dateUntil,
      "chartType": chartType,
    }

    Transaction.SummaryChart(summaryChartSpec)
      .then(
        (result) => {
          setIsLoaded(true);
          setChart(result);
        },
        // Note: it's important to handle errors here
        // instead of a catch() block so that we don't swallow
        // exceptions from actual bugs in components.
        (error) => {
          setIsLoaded(true);
          setError(error);
        }
      )
  }, [querySpec, chartType])

  if (error) {
    return <div>Error: {error.message}</div>;
  } else if (!isLoaded) {
    return <div>Loading...</div>;
  } else {
    return (
      <Stack direction='row' alignItems='flex-start'>
        <SummaryChart dataset={chart.amounts} />
        <FormControl>
          <InputLabel id="summary-chart-select-label">Chart type</InputLabel>
          <Select
            labelId="summary-chart-select-label"
            id="summary-chart-select"
            value={chartType}
            label="Chart Type"
            onChange={handleChange}
          >
            <MenuItem value={'monthDollar'}>By Month Dollar</MenuItem>
            <MenuItem value={'monthPercent'}>By Month Percent</MenuItem>
            <MenuItem value={'cumulativeDollar'}>Cumulative Dollar</MenuItem>
            <MenuItem value={'cumulativePercent'}>Cumulative Percent</MenuItem>
          </Select>
        </FormControl>
      </Stack>
    )
  }
}