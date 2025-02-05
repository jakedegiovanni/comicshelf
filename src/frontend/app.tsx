import { AppBar, Box, IconButton, Toolbar, Typography } from '@mui/material';
import { Menu } from '@mui/icons-material';
import WeeklyComics from './weeklyComics.jsx';

export default function App() {
  return (
    <Box>
      <AppBar position='sticky' style={{ marginBottom: '64px' }}>
        <Toolbar>
          <IconButton
            edge='start'
            color='inherit'
            aria-label='menu'
            sx={{ mr: 2 }}
          >
            <Menu />
          </IconButton>
          <Typography
            variant='h6'
            color='inherit'
            component='div'
            sx={{ flexGrow: 1 }}
          >
            Weekly Comics
          </Typography>
        </Toolbar>
      </AppBar>
      <WeeklyComics />
    </Box>
  );
}
