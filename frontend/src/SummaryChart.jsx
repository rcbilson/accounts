import * as React from 'react';
import { BarChart } from '@mui/x-charts/BarChart';
//import { axisClasses } from '@mui/x-charts/ChartsAxis';
const chartSetting = {
/*
  yAxis: [
    {
      label: 'rainfall (mm)',
    },
  ],
*/
  width: 900,
  height: 300,
/*
  sx: {
    [`.${axisClasses.left} .${axisClasses.label}`]: {
      transform: 'translate(-20px, 0)',
    },
  },
*/
};

// the incoming data is in the form 
//   [{"category":"Charity","amount":"108.00","month":"2025-01"},{"category":"Expenses","amount":"2248.73","month":"2025-01"}]
// convert it to
//   [{"month":"2025-01","Charity":108.00,"Expenses":2248.73}]
function convertData(dataset) {
  const result = [];
  dataset.forEach((e) => {
    const month = e.month;
    const category = e.category;
    const amount = parseFloat(e.amount);
    let row = result.find((r) => r.month === month);
    if (!row) {
      row = {month: month};
      result.push(row);
    }
    row[category] = amount;
  });
  return result;
}

export default function SummaryChart({dataset}) {
  return (
    <BarChart
      dataset={convertData(dataset)}
      xAxis={[{ scaleType: 'band', dataKey: 'month' }]}
      series={[
        { dataKey: 'Expenses', label: 'Expenses', stack: 'stack' },
        { dataKey: 'Travel', label: 'Travel', stack: 'stack' },
        { dataKey: 'Charity', label: 'Charity', stack: 'stack' },
        { dataKey: 'HouseProject', label: 'House Project', stack: 'stack' },
        { dataKey: 'Savings', label: 'Savings', stack: 'stack' },
      ]}
      {...chartSetting}
    />
  );
}
