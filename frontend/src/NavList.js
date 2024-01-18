import { Link } from "react-router-dom";

import List from '@mui/material/List';
import ListItem from '@mui/material/ListItem';
import ListItemButton from '@mui/material/ListItemButton';
import ListItemIcon from '@mui/material/ListItemIcon';
import ListItemText from '@mui/material/ListItemText';
import QueryStatsIcon from '@mui/icons-material/QueryStats';
import FiberNewIcon from '@mui/icons-material/FiberNew';
import HomeIcon from '@mui/icons-material/Home';

export default function NavList() {
  return (
      <List>
          <ListItem button key='home' component={Link} to="/" disablePadding>
            <ListItemButton>
              <ListItemIcon>
                <HomeIcon />
              </ListItemIcon>
              <ListItemText primary='Home' />
            </ListItemButton>
          </ListItem>
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
  );
}
