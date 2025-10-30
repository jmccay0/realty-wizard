import { BrowserRouter as Router, Routes, Route } from 'react-router-dom';
import { AuthProvider } from './contexts/AuthContext';
import ProtectedRoute from './components/ProtectedRoute';
import Home from './pages/Home';
import NewProject from './pages/NewProject';
import Wizard from './pages/Wizard';
import Dashboard from './pages/Dashboard';
import EnterContract from './pages/EnterContract';
import ViewDocuments from './pages/ViewDocuments';
import Login from './pages/Login';
import Register from './pages/Register';
import './index.css';

function App() {
  return (
    <Router>
      <AuthProvider>
        <div style={{ minHeight: '100vh' }}>
          <Routes>
            {/* Public routes */}
            <Route path="/login" element={<Login />} />
            <Route path="/register" element={<Register />} />

            {/* Protected routes */}
            <Route
              path="/"
              element={
                <ProtectedRoute>
                  <Home />
                </ProtectedRoute>
              }
            />
            <Route
              path="/new"
              element={
                <ProtectedRoute>
                  <NewProject />
                </ProtectedRoute>
              }
            />
            <Route
              path="/project/:id"
              element={
                <ProtectedRoute>
                  <Dashboard />
                </ProtectedRoute>
              }
            />
            <Route
              path="/project/:id/wizard"
              element={
                <ProtectedRoute>
                  <Wizard />
                </ProtectedRoute>
              }
            />
            <Route
              path="/project/:id/contract/new"
              element={
                <ProtectedRoute>
                  <EnterContract />
                </ProtectedRoute>
              }
            />
            <Route
              path="/project/:id/documents"
              element={
                <ProtectedRoute>
                  <ViewDocuments />
                </ProtectedRoute>
              }
            />
          </Routes>
        </div>
      </AuthProvider>
    </Router>
  );
}

export default App;
