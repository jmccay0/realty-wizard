import { BrowserRouter as Router, Routes, Route } from 'react-router-dom';
import Home from './pages/Home';
import NewProject from './pages/NewProject';
import Wizard from './pages/Wizard';
import Dashboard from './pages/Dashboard';
import TrecExplainer from './pages/TrecExplainer';
import ServiceMarketplace from './pages/ServiceMarketplace';
import ServiceDetail from './pages/ServiceDetail';
import ServiceRequestForm from './pages/ServiceRequestForm';
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
          <Route path="/learn/trec-contract" element={<TrecExplainer />} />
          <Route path="/services" element={<ServiceMarketplace />} />
          <Route path="/services/:id" element={<ServiceDetail />} />
          <Route path="/services/:id/request" element={<ServiceRequestForm />} />
        </Routes>
      </div>
    </Router>
  );
}

export default App;
