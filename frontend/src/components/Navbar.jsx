import React from 'react';
import { Layout, Menu, Typography, Button } from 'antd';
import { Link } from 'react-router-dom';
import { useAuth } from '../context/AuthContext';

const { Header } = Layout;
const { Title } = Typography;

function Navbar() {
  const { user, logout, isAdmin } = useAuth();

  return (
    <Header style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
      <div className="logo" style={{ color: 'white', fontSize: '20px', fontWeight: 'bold' }}>
        AlgoBharat
      </div>
      <Menu theme="dark" mode="horizontal" defaultSelectedKeys={['/']} style={{ flex: 1, justifyContent: 'flex-end' }}>
        <Menu.Item key="/">
          <Link to="/">Home</Link>
        </Menu.Item>
        <Menu.Item key="/movies">
          <Link to="/movies">Movies</Link>
        </Menu.Item>
        <Menu.Item key="/theatres">
          <Link to="/theatres">Theatres</Link>
        </Menu.Item>
        {isAdmin && (
          <Menu.Item key="/admin/movies">
            <Link to="/admin/movies">Movie Management</Link>
          </Menu.Item>
        )}
        {isAdmin && (
          <Menu.Item key="/admin/theatres">
            <Link to="/admin/theatres">Theatre Management</Link>
          </Menu.Item>
        )}
        {isAdmin && (
          <Menu.Item key="/admin/shows">
            <Link to="/admin/shows">Show Management</Link>
          </Menu.Item>
        )}
        {isAdmin && (
          <Menu.Item key="/admin/analytics">
            <Link to="/admin/analytics">Analytics Dashboard</Link>
          </Menu.Item>
        )}
        {isAdmin && (
          <Menu.Item key="/admin/users">
            <Link to="/admin/users">User Management</Link>
          </Menu.Item>
        )}
        {user ? (
          <Menu.Item key="logout" onClick={logout}>
            Logout ({user.username})
          </Menu.Item>
        ) : (
          <>
            <Menu.Item key="/login">
              <Link to="/login">Login</Link>
            </Menu.Item>
            <Menu.Item key="/register">
              <Link to="/register">Register</Link>
            </Menu.Item>
          </>
        )}
      </Menu>
    </Header>
  );
}

export default Navbar;
