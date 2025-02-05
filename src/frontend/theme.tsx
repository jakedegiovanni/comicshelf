import { createTheme } from '@mui/material';
import { red, blueGrey, amber } from '@mui/material/colors';

const theme = createTheme({
  cssVariables: true,
  palette: {
    primary: {
      main: blueGrey.A400,
    },
    secondary: {
      main: amber.A400,
    },
    error: {
      main: red.A400,
    },
  },
});

export default theme;
