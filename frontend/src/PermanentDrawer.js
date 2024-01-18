import Drawer from '@mui/material/Drawer';

import NavList from './NavList.js';

export default function PermanentDrawer() {
  return (
    <Drawer
        sx={{
          width: 250,
          flexShrink: 0,
          '& .MuiDrawer-paper': {
            width: 250,
            boxSizing: 'border-box',
          },
        }}
        variant="permanent"
        anchor="left"
    >
      <NavList />
    </Drawer>
  )
}
