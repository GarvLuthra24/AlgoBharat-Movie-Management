import React, { useState, useEffect } from "react";
import { API_BASE_URL } from '../config/api';
import axios from "axios";
import { useAuth } from "../context/AuthContext";
import { toast } from "react-toastify";
import {
  Card,
  Button,
  Modal,
  Table,
  Space,
  Input,
  Form,
  InputNumber,
  Typography,
  Row,
  Col,
} from "antd";
import {
  PlusOutlined,
  EditOutlined,
  DeleteOutlined,
  VideoCameraOutlined,
} from "@ant-design/icons";

const { Title, Text } = Typography;

function MovieManagement() {
  const { token } = useAuth();
  const [movies, setMovies] = useState([]);
  const [loading, setLoading] = useState(false);
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [editingMovie, setEditingMovie] = useState(null);
  const [formValues, setFormValues] = useState({
    title: "",
    duration_minutes: 120,
  });

  const fetchMovies = async () => {
    setLoading(true);
    try {
      const response = await axios.get(`${API_BASE_URL}/movies`, {
        headers: { Authorization: `Bearer ${token}` },
      });
      setMovies(response.data.data || []);
    } catch (error) {
      toast.error(error.response?.data?.message || "Failed to fetch movies.");
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    if (token) fetchMovies();
  }, [token]);

  const handleAdd = () => {
    setEditingMovie(null);
    setFormValues({ title: "", duration_minutes: 120 });
    setIsModalOpen(true);
  };

  const handleEdit = (movie) => {
    setEditingMovie(movie);
    setFormValues({
      title: movie.title,
      duration_minutes: movie.duration_minutes,
    });
    setIsModalOpen(true);
  };

  const handleDelete = async (id) => {
    try {
      await axios.delete(`${API_BASE_URL}/movies/${id}`, {
        headers: { Authorization: `Bearer ${token}` },
      });
      toast.success("Movie deleted successfully!");
      fetchMovies();
    } catch (error) {
      toast.error(error.response?.data?.message || "Failed to delete movie.");
    }
  };

  const submitForm = async () => {
    if (!formValues.title || !formValues.duration_minutes) {
      toast.error("Please fill all fields");
      return;
    }
    try {
      if (editingMovie) {
        await axios.put(
            `${API_BASE_URL}/movies/${editingMovie.id}`,
            formValues,
            { headers: { Authorization: `Bearer ${token}` } }
        );
      } else {
        await axios.post(`${API_BASE_URL}/movies`, formValues, {
          headers: { Authorization: `Bearer ${token}` },
        });
      }
      toast.success(`Movie ${editingMovie ? "updated" : "added"} successfully!`);
      setIsModalOpen(false);
      fetchMovies();
    } catch (error) {
      toast.error(
          error.response?.data?.message ||
          "Operation failed. Please check your input."
      );
    }
  };

  const columns = [
    { title: "ID", dataIndex: "id", key: "id", width: 80 },
    {
      title: "Title",
      dataIndex: "title",
      key: "title",
      render: (text) => <Text strong>{text}</Text>,
    },
    {
      title: "Duration (min)",
      dataIndex: "duration_minutes",
      key: "duration_minutes",
      width: 140,
    },
    {
      title: "Action",
      key: "action",
      width: 150,
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
      <div style={{ padding: "20px" }}>
        <Row justify="center">
          <Col xs={24} style={{ maxWidth: "90%" }}>
            <Card
                bordered={false}
                style={{
                  boxShadow: "0 2px 8px rgba(0,0,0,0.1)",
                  borderRadius: "12px",
                  padding: "24px",
                }}
            >
              <Title
                  level={3}
                  style={{ textAlign: "center", marginBottom: "8px" }}
              >
                <VideoCameraOutlined /> Movie Management
              </Title>
              <Text
                  type="secondary"
                  style={{
                    display: "block",
                    textAlign: "center",
                    marginBottom: "24px",
                  }}
              >
                Manage your movie catalog. Add, edit, or remove movies easily.
              </Text>

              <Button
                  type="primary"
                  icon={<PlusOutlined />}
                  onClick={handleAdd}
                  style={{ marginBottom: "16px" }}
              >
                Add Movie
              </Button>

              <Table
                  columns={columns}
                  dataSource={movies}
                  loading={loading}
                  rowKey="id"
                  pagination={{ pageSize: 8 }}
                  bordered
              />
            </Card>
          </Col>
        </Row>

        <Modal
            title={editingMovie ? "Edit Movie" : "Add Movie"}
            open={isModalOpen}
            onOk={submitForm}
            onCancel={() => setIsModalOpen(false)}
            okText="Save"
            cancelText="Cancel"
            destroyOnClose
            centered
        >
          <Form layout="vertical">
            <Form.Item
                label="Movie Title"
                required
                tooltip="Enter the full movie name"
            >
              <Input
                  value={formValues.title}
                  onChange={(e) =>
                      setFormValues({ ...formValues, title: e.target.value })
                  }
                  placeholder="Enter movie title"
              />
            </Form.Item>
            <Form.Item
                label="Duration (minutes)"
                required
                tooltip="Enter movie duration in minutes"
            >
              <InputNumber
                  min={1}
                  value={formValues.duration_minutes}
                  onChange={(value) =>
                      setFormValues({ ...formValues, duration_minutes: value })
                  }
                  style={{ width: "100%" }}
                  placeholder="Enter duration in minutes"
              />
            </Form.Item>
          </Form>
        </Modal>
      </div>
  );
}

export default MovieManagement;
