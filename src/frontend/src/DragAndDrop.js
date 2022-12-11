import React, { useState } from 'react';
import * as Transaction from './Transaction.js';
import { Typography } from '@mui/material';

const DragAndDrop = ({refresh}) => {
  const [currentFile, setCurrentFile] = useState(null);
  const [error, setError] = useState(null);

  const importFiles = async (files) => {
    setError(null);
    try {
      await files.reduce(async (promise, file) => {
        await promise;
        setCurrentFile(file.name);
        const contents = await file.text();
        await Transaction.Import(contents)
      }, Promise.resolve())
      setCurrentFile(null);
      refresh();
    } catch (e) {
      setError(e.message)
    }
  }

  const handleDragEnter = e => {
    e.preventDefault();
    e.stopPropagation();
  };
  const handleDragLeave = e => {
    e.preventDefault();
    e.stopPropagation();
  };
  const handleDragOver = e => {
    e.preventDefault();
    e.stopPropagation();
  };
  const handleDrop = async e => {
    e.preventDefault();
    e.stopPropagation();
    const files = [...e.dataTransfer.files];
    //console.log(files);
    importFiles(files);
  };
  let statusBar;
  if (error) {
    if (currentFile) {
      statusBar =  `An error occurred while importing ${currentFile}: ${error}`
    } else {
      statusBar =  `An error occurred: ${error}`
    }
  } else if (currentFile) {
    statusBar = `Importing ${currentFile}`
  } else {
    statusBar = 'Drop a csv file here to import.'
  }
  return (
    <Typography
      sx={{width: '100%'}}
      align="center"
      onDrop={e => handleDrop(e)}
      onDragOver={e => handleDragOver(e)}
      onDragEnter={e => handleDragEnter(e)}
      onDragLeave={e => handleDragLeave(e)}
    >
      {statusBar}
    </Typography>
  );
};
export default DragAndDrop;
