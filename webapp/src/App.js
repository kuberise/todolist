import useSWR, {mutate} from 'swr';
import { useState } from 'react';
import '@fontsource/inter';

import * as React from 'react';
import { CssVarsProvider } from '@mui/joy/styles';
import Sheet from '@mui/joy/Sheet';
import FormControl from '@mui/joy/FormControl';
import Input from '@mui/joy/Input';
import Button from '@mui/joy/Button';
import IconButton from '@mui/joy/IconButton';



import Add from '@mui/icons-material/Add';
import AddTaskIcon from '@mui/icons-material/AddTask';
import LabelIcon from '@mui/icons-material/Label';
import RemoveCircleOutlineIcon from '@mui/icons-material/RemoveCircleOutline';

import List from '@mui/joy/List';
import ListItem from '@mui/joy/ListItem';
import ListItemContent from '@mui/joy/ListItemContent';
import ListItemDecorator from '@mui/joy/ListItemDecorator';


import Typography from '@mui/joy/Typography';

// created function to handle API request
const fetcher = (...args) => fetch(...args).then((res) => res.json());



const poster = (title) => {
  if(title !== "") {
    fetch(`${process.env.REACT_APP_BACKEND_URL}`, { method: 'POST', body: JSON.stringify({item: title}) })
    .then(() => {
      mutate(`${process.env.REACT_APP_BACKEND_URL}`)
    }).catch(e => console.log(e))
  }
}

const remover = (title) => {
  if(title !== "") {
    fetch(`${process.env.REACT_APP_BACKEND_URL}/delete?item=${title}`, { method: 'DELETE' })
    .then(() => {
      mutate(`${process.env.REACT_APP_BACKEND_URL}`)
    })
    .catch(e => console.log(e))
  }
}

function TODOList() {

  const {
    data: todos,
    error,
    isValidating,
  } = useSWR(`${process.env.REACT_APP_BACKEND_URL}`, fetcher);

  if (error) return <div className='failed'>failed to load</div>;
  if (isValidating) return <div className="Loading">Loading...</div>;

  return (
    <List
        aria-labelledby="list"
      >
      {
        todos.map((todo, index) => (
          <ListItem 
          key={index} 
          
          variant='outlined'
          sx={{
            mx: 'auto', // margin left & right
            my: 0.4, // margin top & bottom
            py: 0, // padding top & bottom
            px: 0.4, // padding left & right
            borderRadius: 'sm',
          }}
          >
          <ListItemDecorator>
          <LabelIcon fontSize="small" />
          </ListItemDecorator>
          <ListItemContent>
            <Typography level="title-sm">{todo}</Typography>
          </ListItemContent>
          <ListItemDecorator>
          <IconButton 
          size="sm" 
          variant='outlined'
          onClick={() => remover(todo)}
           ><RemoveCircleOutlineIcon /></IconButton>
          </ListItemDecorator>
        </ListItem>
        ))
      }
      </List>
  )
}

export default function App() {

  const [title, setTitle] = useState("")  

  return (
    <CssVarsProvider>
      <Sheet
  sx={{
    width: 600,
    mx: 'auto', // margin left & right
    my: 4, // margin top & bottom
    py: 3, // padding top & bottom
    px: 2, // padding left & right
    display: 'flex',
    flexDirection: 'column',
    alignItems: 'center',
    gap: 2,
    borderRadius: 'sm',
    boxShadow: 'md',
  }}
>
  TODO List
  <TODOList />
  <FormControl>
  <Input
  name="New Item"
  value={title}
  onChange={(event) =>
    setTitle(event.target.value)
  }
  placeholder="do some cool stuff..."
  startDecorator={<AddTaskIcon />}
  endDecorator={<Button onClick={() => {
    poster(title)
    setTitle("")
  }} variant='outlined' startDecorator={<Add />} >
  Add
</Button>}>
</Input>
</FormControl>
</Sheet>
    </CssVarsProvider>
  );
}

