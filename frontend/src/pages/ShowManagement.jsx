import React, { useState, useEffect } from 'react';
import { API_BASE_URL } from '../config/api';
import axios from 'axios';
import { useAuth } from '../context/AuthContext';
import dayjs from 'dayjs';
import { toast } from 'react-toastify';
import {
  Card,
  Button,
  Modal,
  Table,
  Space,
  Select,
  Form,
  InputNumber,
  Typography,
  Row,
  Col,
  DatePicker,
} from 'antd';
import { PlusOutlined, EditOutlined, DeleteOutlined, VideoCameraOutlined } from '@ant-design/icons';

const { Title, Text } = Typography;
const { Option } = Select;

function ShowManagement() {
  const { token } = useAuth();
  const [shows, setShows] = useState([]);
  const [movies, setMovies] = useState([]);
  const [halls, setHalls] = useState([]);
  const [theatres, setTheatres] = useState([]);
  const [loading, setLoading] = useState(false);
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [editingShow, setEditingShow] = useState(null);

  const [formValues, setFormValues] = useState({
    movie_id: '',
    theatre_id: '',
    hall_id: '',
    time: '',
    price: 0,
  });

  const fetchData = async () => {
    setLoading(true);
    try {
      const [showsRes, moviesRes, theatresRes, hallsRes] = await Promise.all([
        axios.get(`${API_BASE_URL}/shows`, { headers: { Authorization: `Bearer ${token}` } }),
        axios.get(`${API_BASE_URL}/movies`, { headers: { Authorization: `Bearer ${token}` } }),
        axios.get(`${API_BASE_URL}/theatres`, { headers: { Authorization: `Bearer ${token}` } }),
        axios.get(`${API_BASE_URL}/halls`, { headers: { Authorization: `Bearer ${token}` } }),
      ]);

      setShows(showsRes.data.data || []);
      setMovies(moviesRes.data.data || []);
      setTheatres(theatresRes.data.data || []);
      setHalls(hallsRes.data.data || []);
    } catch (error) {
      toast.error(error.response?.data?.message || 'Failed to fetch necessary data.');
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    if (token) fetchData();
  }, [token]);

  const handleAdd = () => {
    setEditingShow(null);
    setFormValues({ movie_id: '', theatre_id: '', hall_id: '', time: '', price: 0 });
    setIsModalOpen(true);
  };

  const handleEdit = (record) => {
    const hallOfShow = halls.find((h) => h.id === record.hall_id);
    const theatreIdOfShow = hallOfShow ? hallOfShow.theatre_id : '';

    setEditingShow(record);
    setFormValues({
      movie_id: record.movie_id,
      theatre_id: theatreIdOfShow,
      hall_id: record.hall_id,
      time: dayjs(record.time).toISOString(),
      price: record.price,
    });
    setIsModalOpen(true);
  };

  const handleDelete = async (id) => {
    try {
      await axios.delete(`${API_BASE_URL}/shows/${id}`, {
        headers: { Authorization: `Bearer ${token}` },
      });
      toast.success('Show deleted successfully!');
      fetchData();
    } catch (error) {
      toast.error(error.response?.data?.message || 'Failed to delete show.');
    }
  };

  const saveShow = async () => {
    try {
      const payload = {
        ...formValues,
        time: new Date(formValues.time).toISOString(),
        price: Number(formValues.price),
      };

      if (editingShow) {
        await axios.put(`${API_BASE_URL}/shows/${editingShow.id}`, payload, {
          headers: { Authorization: `Bearer ${token}` },
        });
      } else {
        await axios.post(`${API_BASE_URL}/shows`, payload, {
          headers: { Authorization: `Bearer ${token}` },
        });
      }
      toast.success(`Show ${editingShow ? 'updated' : 'added'} successfully!`);
      setIsModalOpen(false);
      fetchData();
    } catch (error) {
      toast.error(error.response?.data?.message || 'Operation failed. Please check your input.');
    }
  };

  const columns = [
    { title: 'ID', dataIndex: 'id', key: 'id', width: 80 },
    {
      title: 'Movie',
      dataIndex: 'movie_id',
      key: 'movie_id',
      render: (movieId) => movies.find((m) => m.id === movieId)?.title || movieId,
    },
    {
      title: 'Hall',
      dataIndex: 'hall_id',
      key: 'hall_id',
      render: (hallId) => halls.find((h) => h.id === hallId)?.name || hallId,
    },
    {
      title: 'Time',
      dataIndex: 'time',
      key: 'time',
      render: (time) => dayjs(time).format('YYYY-MM-DD HH:mm'),
    },
    {
      title: 'Price',
      dataIndex: 'price',
      key: 'price',
      render: (price) => `₹ ${price.toFixed(2)}`,
    },
    {
      title: 'Action',
      key: 'action',
      render: (_, record) => (
          <Space>
            <Button icon={<EditOutlined />} onClick={() => handleEdit(record)}>
              Edit
            </Button>
            <Button
                icon={<DeleteOutlined />}
                danger
                onClick={() => handleDelete(record.id)}
            >
              Delete
            </Button>
          </Space>
      ),
    },
  ];

  return (
      <div style={{ padding: '20px' }}>
        <Row justify="center">
          <Col xs={24} style={{ maxWidth: '90%' }}>
            <Card
                bordered={false}
                style={{
                  boxShadow: '0 2px 8px rgba(0,0,0,0.1)',
                  borderRadius: '12px',
                  padding: '24px',
                }}
            >
              <Title level={3} style={{ textAlign: 'center', marginBottom: '8px' }}>
                <VideoCameraOutlined /> Show Management
              </Title>
              <Text
                  type="secondary"
                  style={{
                    display: 'block',
                    textAlign: 'center',
                    marginBottom: '24px',
                  }}
              >
                Add, edit or remove shows for movies in theatres and halls.
              </Text>

              <Button
                  type="primary"
                  icon={<PlusOutlined />}
                  onClick={handleAdd}
                  style={{ marginBottom: '16px' }}
              >
                Add Show
              </Button>

              <Table
                  columns={columns}
                  dataSource={shows}
                  loading={loading}
                  rowKey="id"
                  bordered
                  pagination={{ pageSize: 8 }}
              />
            </Card>
          </Col>
        </Row>

        {/* Add/Edit Show Modal */}
        <Modal
            title={editingShow ? 'Edit Show' : 'Add Show'}
            open={isModalOpen}
            onOk={saveShow}
            onCancel={() => setIsModalOpen(false)}
            okText="Save"
            centered
        >
          <Form layout="vertical" name="show_form">
            <Form.Item name="movie_id" label="Movie" required>
              <Select
                  placeholder="Select a movie"
                  value={formValues.movie_id}
                  onChange={(e) => setFormValues({ ...formValues, movie_id: e })}
              >
                {movies.map((movie) => (
                    <Option key={movie.id} value={movie.id}>
                      {movie.title}
                    </Option>
                ))}
              </Select>
            </Form.Item>

            <Form.Item name="theatre_id" label="Theatre" required>
              <Select
                  placeholder="Select a theatre"
                  value={formValues.theatre_id}
                  onChange={(value) => {
                    setFormValues({ ...formValues, theatre_id: value, hall_id: '' });
                  }}
              >
                {theatres.map((theatre) => (
                    <Option key={theatre.id} value={theatre.id}>
                      {theatre.name}
                    </Option>
                ))}
              </Select>
            </Form.Item>

            <Form.Item name="hall_id" label="Hall" required>
              <Select
                  placeholder="Select a hall"
                  disabled={!formValues.theatre_id}
                  value={formValues.hall_id}
                  onChange={(e) => setFormValues({ ...formValues, hall_id: e })}
              >
                {halls
                    .filter((hall) => hall.theatre_id === formValues.theatre_id)
                    .map((hall) => (
                        <Option key={hall.id} value={hall.id}>
                          {hall.name}
                        </Option>
                    ))}
              </Select>
            </Form.Item>

            <Form.Item name="time" label="Show Time" required>
              <DatePicker
                  showTime
                  format="YYYY-MM-DD HH:mm"
                  style={{ width: '100%' }}
                  value={formValues.time ? dayjs(formValues.time) : null}
                  onChange={(e) =>
                      setFormValues({ ...formValues, time: e ? e.toISOString() : '' })
                  }
              />
            </Form.Item>

            <Form.Item name="price" label="Ticket Price" required>
              <InputNumber
                  min={0.01}
                  step={0.01}
                  style={{ width: '100%' }}
                  value={formValues.price}
                  onChange={(e) => setFormValues({ ...formValues, price: e })}
              />
            </Form.Item>
          </Form>
        </Modal>
      </div>
  );
}

export default ShowManagement;
