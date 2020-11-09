import React from 'react';
import './App.css';
import { Provider } from 'react-redux';
import Store from './redux/Store';
import { Navbar } from './components/Navbar';
import { CommunityGaolTracker } from './components/CommunityGoalTracker';

const App = () => {

  return (
    <Provider store={Store}>
      <Navbar />
      <CommunityGaolTracker />
    </Provider>
  );
}

export default App;
