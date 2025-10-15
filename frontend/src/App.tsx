import { BrowserRouter as Router, Routes, Route } from 'react-router-dom';
import Home from './pages/Home';
import NewProject from './pages/NewProject';
import Wizard from './pages/Wizard';
import Dashboard from './pages/Dashboard';
import './index.css';

function App() {
  return (
    <Router>
      <div style={{ minHeight: '100vh' }}>
        <Routes>
          <Route path="/" element={<Home />} />
          <Route path="/new" element={<NewProject />} />
          <Route path="/project/:id" element={<Dashboard />} />
          <Route path="/project/:id/wizard" element={<Wizard />} />
        </Routes>
      </div>
    </Router>
  );
}

export default App;
