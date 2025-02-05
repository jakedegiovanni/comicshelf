import { Card, CardContent, CardMedia, Grid2, Typography } from '@mui/material';
import { useEffect, useState } from 'react';
import { Comic, Page } from '../comics/comics.js';

export default function WeeklyComics() {
  const [comics, setComics] = useState(new Page<Comic>(0, 0, 0, 0, []));

  useEffect(() => {
    let ignore = false;
    fetch(`${import.meta.env.VITE_FRONTEND_SERVER_HOST}/api/v1/comics`)
      .then((resp: globalThis.Response) => {
        if (!resp.ok) {
          console.warn(resp);
          throw new Error("couldn't fetch");
        }

        return resp.json();
      })
      .then(comics => {
        if (!ignore) {
          setComics(comics);
        }
      })
      .catch(console.warn);

    return () => {
      ignore = true;
    };
  }, []);

  return (
    <Grid2 container spacing={4} style={{ margin: '0 32px' }}>
      {comics.results.map((comic: Comic) => (
        <Grid2 key={comic.id} size={{ sm: 6, lg: 2 }}>
          <Card>
            <CardMedia
              component='img'
              image={comic.thumbnail}
              title={comic.title}
            />
            <CardContent>
              <Typography gutterBottom variant='h5' component='div'>
                {comic.title}
              </Typography>
            </CardContent>
          </Card>
        </Grid2>
      ))}
    </Grid2>
  );
}
