import { createRoot } from 'react-dom/client';
import { StrictMode } from 'react';
import App from './app.jsx';
import '@fontsource/roboto/300.css';
import '@fontsource/roboto/400.css';
import '@fontsource/roboto/500.css';
import '@fontsource/roboto/700.css';
import { CssBaseline, ThemeProvider } from '@mui/material';
import theme from './theme.jsx';

const rootElement = document.getElementById('root');
if (rootElement === null) throw new Error('');

createRoot(rootElement).render(
  <StrictMode>
    <ThemeProvider theme={theme}>
      <CssBaseline />
      <App />
    </ThemeProvider>
  </StrictMode>,
);
