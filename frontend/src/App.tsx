import React from 'react';
import './App.css';
import { Provider } from 'react-redux';
import Store from './redux/Store';
import { AppNavBar } from './components/AppNavBar/AppNavBar';
import { ChurchFundManaging } from './components/ChurchFundManaging/ChurchFundManaging';

const App = () => {
  return (
    <Provider store={Store}>
      <AppNavBar />
      <ChurchFundManaging key="churchfundmanaging" />
    </Provider>
  );
}

export default App;
