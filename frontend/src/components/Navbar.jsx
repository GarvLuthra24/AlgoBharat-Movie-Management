import React, { useState } from "react";
import { Layout, Menu, Drawer, Button } from "antd";
import { Link, useLocation } from "react-router-dom";
import { MenuOutlined } from "@ant-design/icons";
import { useAuth } from "../context/AuthContext";

const { Header } = Layout;

function Navbar() {
    const { user, logout, isAdmin } = useAuth();
    const location = useLocation();
    const [open, setOpen] = useState(false);

    // Menu items
    const menuItems = [
        { key: "/", label: <Link to="/">Home</Link> },
        { key: "/movies", label: <Link to="/movies">Movies</Link> },
        { key: "/theatres", label: <Link to="/theatres">Theatres</Link> },
        ...(isAdmin
            ? [
                { key: "/admin/movies", label: <Link to="/admin/movies">Movie Management</Link> },
                { key: "/admin/theatres", label: <Link to="/admin/theatres">Theatre Management</Link> },
                { key: "/admin/shows", label: <Link to="/admin/shows">Show Management</Link> },
                { key: "/admin/analytics", label: <Link to="/admin/analytics">Analytics Dashboard</Link> },
                { key: "/admin/users", label: <Link to="/admin/users">User Management</Link> },
            ]
            : []),
        ...(user
            ? [{ key: "logout", label: <span onClick={logout}>Logout ({user.username})</span> }]
            : [
                { key: "/login", label: <Link to="/login">Login</Link> },
                { key: "/register", label: <Link to="/register">Register</Link> },
            ]),
    ];

    return (
        <Header
            style={{
                background: "#0d1117", // GitHub-dark vibe
                display: "flex",
                justifyContent: "space-between",
                alignItems: "center",
                padding: "0 20px",
            }}
        >
            {/* Logo */}
            <div
                style={{
                    color: "#fff",
                    fontSize: "22px",
                    fontWeight: "bold",
                    letterSpacing: "1px",
                }}
            >
                AlgoBharat
            </div>

            {/* Desktop Menu */}
            <div className="desktop-menu" style={{ flex: 1, display: "flex", justifyContent: "flex-end" }}>
                <Menu
                    theme="dark"
                    mode="horizontal"
                    selectedKeys={[location.pathname]}
                    items={menuItems}
                    style={{
                        flex: 1,
                        justifyContent: "flex-end",
                        background: "transparent",
                        borderBottom: "none",
                    }}
                />
            </div>

            {/* Mobile Hamburger */}
            <div className="mobile-menu" style={{ display: "none" }}>
                <Button
                    type="text"
                    icon={<MenuOutlined style={{ color: "white", fontSize: "22px" }} />}
                    onClick={() => setOpen(true)}
                />
                <Drawer
                    title="Navigation"
                    placement="right"
                    onClose={() => setOpen(false)}
                    open={open}
                    bodyStyle={{ padding: 0 }}
                >
                    <Menu
                        mode="inline"
                        selectedKeys={[location.pathname]}
                        items={menuItems}
                        onClick={() => setOpen(false)}
                    />
                </Drawer>
            </div>

            {/* Cool styles */}
            <style>{`
        /* Responsive toggle */
        @media (max-width: 768px) {
          .desktop-menu { display: none !important; }
          .mobile-menu { display: block !important; }
        }

        /* Menu item styles */
        .ant-menu-dark.ant-menu-horizontal > .ant-menu-item {
          transition: all 0.3s ease;
          margin: 0 5px;
          border-radius: 6px;
        }

        /* Hover effect */
        .ant-menu-dark.ant-menu-horizontal > .ant-menu-item:hover {
          background: #1f2937 !important; /* subtle gray */
        }

        /* Active item (highlighted) */
        .ant-menu-dark.ant-menu-horizontal > .ant-menu-item-selected {
          background: #2563eb !important; /* bright blue */
          font-weight: 600;
        }
      `}</style>
        </Header>
    );
}

export default Navbar;
