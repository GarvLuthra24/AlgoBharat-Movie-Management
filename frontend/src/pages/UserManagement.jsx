import React, { useState, useEffect } from 'react';
import { API_BASE_URL } from '../config/api';
import { Table, Button, Modal, Form, Select, Typography, Input, Row, Col, Card, Spin } from 'antd';
import axios from 'axios';
import { useAuth } from '../context/AuthContext';
import { toast } from 'react-toastify';

const { Title, Text } = Typography;
const { Option } = Select;

function UserManagement() {
  const { token, user, isAdmin } = useAuth();
  const [users, setUsers] = useState([]);
  const [loading, setLoading] = useState(false);
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [editingUser, setEditingUser] = useState(null);
  const [form] = Form.useForm();

  const fetchUsers = async () => {
    setLoading(true);
    try {
      if (!token) {
        toast.info('Authentication token not found. Please log in.');
        setLoading(false);
        return;
      }
      const response = await axios.get(`${API_BASE_URL}/users`, {
        headers: { Authorization: `Bearer ${token}` },
      });
      if (response.data?.status?.success) {
        setUsers(response.data.data || []);
      } else {
        toast.error(response.data?.message || 'Failed to fetch users.');
        setUsers([]);
      }
    } catch (error) {
      toast.error(error.response?.data?.message || 'Failed to fetch users.');
      setUsers([]);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    if (token && isAdmin) fetchUsers();
    else setLoading(false);
  }, [token, isAdmin]);

  const handleEdit = (record) => {
    setEditingUser(record);
    form.setFieldsValue(record);
    setIsModalOpen(true);
  };

  const handleOk = async () => {
    setLoading(true);
    try {
      const values = await form.validateFields();
      if (editingUser) {
        await axios.put(`${API_BASE_URL}/users/${editingUser.id}/role`, { role: values.role }, {
          headers: { Authorization: `Bearer ${token}` },
        });
        toast.success('User role updated successfully!');
      }
      setIsModalOpen(false);
      fetchUsers();
    } catch (error) {
      toast.error(error.response?.data?.message || 'Operation failed. Please check your input.');
    } finally {
      setLoading(false);
    }
  };

  const handleCancel = () => setIsModalOpen(false);

  const columns = [
    { title: 'ID', dataIndex: 'id', key: 'id', width: 100 },
    { title: 'Username', dataIndex: 'username', key: 'username' },
    { title: 'Role', dataIndex: 'role', key: 'role', width: 120 },
    {
      title: 'Action',
      key: 'action',
      width: 150,
      render: (_, record) => (
          <Button type="link" onClick={() => handleEdit(record)}>Edit Role</Button>
      ),
    },
  ];

  // Conditional rendering
  if (!user) {
    return (
        <div style={{ textAlign: 'center', padding: '50px', color: '#ff4d4f', fontSize: '18px', fontWeight: 'bold' }}>
          Please log in to access User Management.
        </div>
    );
  }

  if (!isAdmin) {
    return (
        <div style={{ textAlign: 'center', padding: '50px', color: '#ff4d4f', fontSize: '18px', fontWeight: 'bold' }}>
          Access Denied: Only administrators can view this page.
        </div>
    );
  }

  return (
      <div style={{ padding: '20px' }}>
        <Row justify="center">
          <Col xs={24} lg={20} xl={16}>
            <Card>
              <Title level={2} style={{ textAlign: 'center' }}>User Management</Title>
              {loading && users.length === 0 ? (
                  <Spin size="large" style={{ display: 'block', margin: '50px auto' }} />
              ) : (
                  <Table
                      columns={columns}
                      dataSource={users}
                      loading={loading}
                      rowKey="id"
                      pagination={{ pageSize: 10 }}
                      scroll={{ x: true }}
                      locale={{ emptyText: 'No users found.' }}
                  />
              )}
            </Card>
          </Col>
        </Row>

        <Modal
            key={isModalOpen ? 'user-modal-open' : 'user-modal-closed'}
            title="Edit User Role"
            open={isModalOpen}
            onOk={handleOk}
            onCancel={handleCancel}
            confirmLoading={loading}
            centered
        >
          <Form form={form} layout="vertical" name="user_role_form">
            <Form.Item
                name="username"
                label="Username"
            >
              <Input disabled />
            </Form.Item>
            <Form.Item
                name="role"
                label="Role"
                rules={[{ required: true, message: 'Please select a role!' }]}
            >
              <Select placeholder="Select a role">
                <Option value="user">User</Option>
                <Option value="admin">Admin</Option>
              </Select>
            </Form.Item>
          </Form>
        </Modal>
      </div>
  );
}

export default UserManagement;
