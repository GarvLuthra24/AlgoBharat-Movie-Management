import React from 'react';
import { BrowserRouter as Router, Routes, Route } from 'react-router-dom';
import { AuthProvider, useAuth } from './context/AuthContext';
import { ToastContainer } from 'react-toastify';
import 'react-toastify/dist/ReactToastify.css';
import './App.css';
import './styles/custom.css';
import Navbar from './components/Navbar';
import { Layout, Typography } from 'antd';

// Pages
import Home from './pages/Home';
import Movies from './pages/Movies';
import Theatres from './pages/Theatres';
import Login from './pages/Login';
import Register from './pages/Register';
import MovieManagement from './pages/MovieManagement';
import TheatreManagement from './pages/TheatreManagement';
import ShowManagement from './pages/ShowManagement';
import ShowListing from './pages/ShowListing';
import AnalyticsDashboard from './pages/AnalyticsDashboard';
import UserManagement from './pages/UserManagement';

const { Content, Footer } = Layout;
const { Title } = Typography;

const PrivateRoute = ({ children, adminOnly }) => {
  const { user, isAdmin } = useAuth();
  if (!user) {
    return <Login />;
  }
  if (adminOnly && !isAdmin) {
    return <div style={{ textAlign: 'center', padding: '50px', color: '#ff4d4f', fontSize: '18px', fontWeight: 'bold' }}>Access Denied: Admins only</div>;
  }
  return children;
};

function AppContent() {
  return (
    <Layout style={{ minHeight: '100vh' }}>
      <Navbar />
      <Content style={{ padding: '24px' }}>
        <Routes>
          <Route path="/" element={<Home />} />
          <Route path="/movies" element={<Movies />} />
          <Route path="/movies/:id" element={<ShowListing />} />
          <Route path="/theatres" element={<Theatres />} />
          <Route path="/login" element={<Login />} />
          <Route path="/register" element={<Register />} />
          <Route path="/admin/movies" element={
            <PrivateRoute adminOnly={true}>
              <MovieManagement />
            </PrivateRoute>
          } />
          <Route path="/admin/theatres" element={
            <PrivateRoute adminOnly={true}>
              <TheatreManagement />
            </PrivateRoute>
          } />
          <Route path="/admin/shows" element={
            <PrivateRoute adminOnly={true}>
              <ShowManagement />
            </PrivateRoute>
          } />
          <Route path="/admin/analytics" element={
            <PrivateRoute adminOnly={true}>
              <AnalyticsDashboard />
            </PrivateRoute>
          } />
          <Route path="/admin/users" element={
            <PrivateRoute adminOnly={true}>
              <UserManagement />
            </PrivateRoute>
          } />
          <Route path="/admin-dashboard" element={
            <PrivateRoute adminOnly={true}>
              <Title level={2}>Admin Dashboard</Title>
            </PrivateRoute>
          } />
        </Routes>
      </Content>
      <Footer style={{ textAlign: 'center' }}>
        AlgoBharat ©2024 Created by You
      </Footer>
      <ToastContainer position="bottom-right" autoClose={3000} hideProgressBar={false} newestOnTop={false} closeOnClick rtl={false} pauseOnFocusLoss draggable pauseOnHover />
    </Layout>
  );
}

function App() {
  return (
    <Router>
      <AuthProvider>
        <AppContent />
      </AuthProvider>
    </Router>
  );
}

export default App;
