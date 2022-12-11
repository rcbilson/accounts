import * as React from 'react';
import { Link } from "react-router-dom";
import Box from '@mui/material/Box';
import Drawer from '@mui/material/Drawer';
import IconButton from '@mui/material/IconButton';
import List from '@mui/material/List';
import ListItem from '@mui/material/ListItem';
import ListItemButton from '@mui/material/ListItemButton';
import ListItemIcon from '@mui/material/ListItemIcon';
import ListItemText from '@mui/material/ListItemText';
import QueryStatsIcon from '@mui/icons-material/QueryStats';
import FiberNewIcon from '@mui/icons-material/FiberNew';
import MenuIcon from '@mui/icons-material/Menu';

export default function TemporaryDrawer() {
  const [drawerOpen, setDrawerOpen] = React.useState(false);

  const toggleDrawer =
    (open) =>
      (event) => {
        if (
          event.type === 'keydown' &&
          ((event).key === 'Tab' ||
            (event).key === 'Shift')
        ) {
          return;
        }

        setDrawerOpen(open);
      };

  const list = (
    <Box
      sx={{ width: 250 }}
      role="presentation"
      onClick={toggleDrawer(false)}
      onKeyDown={toggleDrawer(false)}
    >
      <List>
          <ListItem button key='query' component={Link} to="/query" disablePadding>
            <ListItemButton>
              <ListItemIcon>
                <QueryStatsIcon />
              </ListItemIcon>
              <ListItemText primary='Query' />
            </ListItemButton>
          </ListItem>
          <ListItem button key='new' component={Link} to="/new" disablePadding>
            <ListItemButton>
              <ListItemIcon>
                <FiberNewIcon />
              </ListItemIcon>
              <ListItemText primary='New' />
            </ListItemButton>
          </ListItem>
      </List>
    </Box>
  );

  return (
    <div>
          <IconButton onClick={toggleDrawer(true)} color="primary" aria-label="menu">
            <MenuIcon />
          </IconButton>
          <Drawer
            anchor='left'
            open={drawerOpen}
            onClose={toggleDrawer(false)}
          >
            {list}
          </Drawer>
    </div>
  );
}
