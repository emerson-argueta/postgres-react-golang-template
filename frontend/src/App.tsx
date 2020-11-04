import React from 'react';
import './App.css';
import { Provider } from 'react-redux';
import Store from './redux/Store';
import { Navbar } from './components/Navbar';

const App = () => {
  return (
    <Provider store={Store}>
      <Navbar />
      {/* TODO: add application component */}
    </Provider>
  );
}

export default App;
