import React, { useState, useEffect } from "react";
import { API_BASE_URL } from '../config/api';
import {
  Table,
  Button,
  Modal,
  Form,
  Input,
  Typography,
  Space,
  Row,
  Col,
  Card,
  Divider,
} from "antd";
import {
  PlusOutlined,
  EditOutlined,
  DeleteOutlined,
  ClusterOutlined,
  HomeOutlined,
} from "@ant-design/icons";
import axios from "axios";
import { useAuth } from "../context/AuthContext";
import SeatMapDesigner from "../components/SeatMapDesigner";
import { toast } from "react-toastify";

const { Title, Text } = Typography;

function TheatreManagement() {
  const { token } = useAuth();
  const [theatres, setTheatres] = useState([]);
  const [loading, setLoading] = useState(false);
  const [isTheatreModalOpen, setIsTheatreModalOpen] = useState(false);
  const [editingTheatre, setEditingTheatre] = useState(null);

  const [isHallModalOpen, setIsHallModalOpen] = useState(false);
  const [selectedTheatre, setSelectedTheatre] = useState(null);
  const [halls, setHalls] = useState([]);
  const [isAddEditHallModalOpen, setIsAddEditHallModalOpen] = useState(false);
  const [editingHall, setEditingHall] = useState(null);

  const [currentSeatMap, setCurrentSeatMap] = useState({});

  const [form] = Form.useForm();
  const [hallForm] = Form.useForm();

  const fetchTheatres = async () => {
    setLoading(true);
    try {
      const response = await axios.get(`${API_BASE_URL}/theatres`, {
        headers: { Authorization: `Bearer ${token}` },
      });
      setTheatres(response.data.data || []);
    } catch (error) {
      toast.error(error.response?.data?.message || "Failed to fetch theatres.");
    } finally {
      setLoading(false);
    }
  };

  const fetchHalls = async (theatreId) => {
    setLoading(true);
    try {
      const response = await axios.get(
          `${API_BASE_URL}/halls?theatreId=${theatreId}`,
          { headers: { Authorization: `Bearer ${token}` } }
      );
      setHalls(response.data.data || []);
    } catch (error) {
      toast.error(error.response?.data?.message || "Failed to fetch halls.");
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    if (token) fetchTheatres();
  }, [token]);

  // --- Theatre Modal ---
  const handleAddTheatre = () => {
    setEditingTheatre(null);
    form.resetFields();
    setIsTheatreModalOpen(true);
  };

  const handleEditTheatre = (record) => {
    setEditingTheatre(record);
    form.setFieldsValue(record);
    setIsTheatreModalOpen(true);
  };

  const handleDeleteTheatre = async (id) => {
    try {
      await axios.delete(`${API_BASE_URL}/theatres/${id}`, {
        headers: { Authorization: `Bearer ${token}` },
      });
      toast.success("Theatre deleted successfully!");
      fetchTheatres();
    } catch (error) {
      toast.error(error.response?.data?.message || "Failed to delete theatre.");
    }
  };

  const handleOkTheatre = async () => {
    try {
      const values = await form.validateFields();
      if (editingTheatre) {
        await axios.put(
            `${API_BASE_URL}/theatres/${editingTheatre.id}`,
            values,
            { headers: { Authorization: `Bearer ${token}` } }
        );
      } else {
        await axios.post(`${API_BASE_URL}/theatres`, values, {
          headers: { Authorization: `Bearer ${token}` },
        });
      }
      toast.success(
          `Theatre ${editingTheatre ? "updated" : "added"} successfully!`
      );
      setIsTheatreModalOpen(false);
      fetchTheatres();
    } catch (error) {
      toast.error(
          error.response?.data?.message ||
          "Operation failed. Please check your input."
      );
    }
  };

  // --- Hall List Modal ---
  const handleManageHalls = (theatre) => {
    setSelectedTheatre(theatre);
    fetchHalls(theatre.id);
    setIsHallModalOpen(true);
  };

  const handleCloseHallModal = () => {
    setIsHallModalOpen(false);
    setSelectedTheatre(null);
    setHalls([]);
  };

  // --- Add/Edit Hall Modal ---
  const handleAddHall = () => {
    setEditingHall(null);
    hallForm.resetFields();
    setCurrentSeatMap({});
    setIsAddEditHallModalOpen(true);
  };

  const handleEditHall = (record) => {
    setEditingHall(record);
    hallForm.setFieldsValue({ name: record.name });
    setCurrentSeatMap(record.seat_map || {});
    setIsAddEditHallModalOpen(true);
  };

  const handleDeleteHall = async (id) => {
    try {
      await axios.delete(`${API_BASE_URL}/halls/${id}`, {
        headers: { Authorization: `Bearer ${token}` },
      });
      toast.success("Hall deleted successfully!");
      fetchHalls(selectedTheatre.id);
    } catch (error) {
      toast.error(error.response?.data?.message || "Failed to delete hall.");
    }
  };

  const handleOkAddEditHall = async () => {
    try {
      const formValues = await hallForm.validateFields();
      if (!currentSeatMap || Object.keys(currentSeatMap).length === 0) {
        toast.error("Please design the seat map by adding at least one row.");
        return;
      }

      const payload = {
        ...formValues,
        seat_map: currentSeatMap,
        theatre_id: selectedTheatre.id,
      };

      if (editingHall) {
        await axios.put(
            `${API_BASE_URL}/halls/${editingHall.id}`,
            payload,
            { headers: { Authorization: `Bearer ${token}` } }
        );
      } else {
        await axios.post(`${API_BASE_URL}/halls`, payload, {
          headers: { Authorization: `Bearer ${token}` },
        });
      }

      toast.success(`Hall ${editingHall ? "updated" : "added"} successfully!`);
      setIsAddEditHallModalOpen(false);
      handleCloseHallModal();
    } catch (error) {
      toast.error(error.response?.data?.message || "Hall operation failed.");
    }
  };

  // --- Columns ---
  const theatreColumns = [
    { title: "ID", dataIndex: "id", key: "id", width: 80 },
    { title: "Name", dataIndex: "name", key: "name" },
    {
      title: "Action",
      key: "action",
      render: (_, record) => (
          <Space>
            <Button icon={<EditOutlined />} onClick={() => handleEditTheatre(record)}>
              Edit
            </Button>
            <Button
                icon={<DeleteOutlined />}
                danger
                onClick={() => handleDeleteTheatre(record.id)}
            >
              Delete
            </Button>
            <Button
                icon={<ClusterOutlined />}
                type="dashed"
                onClick={() => handleManageHalls(record)}
            >
              Manage Halls
            </Button>
          </Space>
      ),
    },
  ];

  const hallColumns = [
    { title: "ID", dataIndex: "id", key: "id", width: 80 },
    { title: "Name", dataIndex: "name", key: "name" },
    {
      title: "Seat Map",
      dataIndex: "seat_map",
      key: "seat_map",
      render: (seat_map) => (
          <Text type="secondary">{Object.keys(seat_map || {}).length} rows</Text>
      ),
    },
    {
      title: "Action",
      key: "action",
      render: (_, record) => (
          <Space>
            <Button icon={<EditOutlined />} onClick={() => handleEditHall(record)}>
              Edit
            </Button>
            <Button
                icon={<DeleteOutlined />}
                danger
                onClick={() => handleDeleteHall(record.id)}
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
              <Title level={3} style={{ textAlign: "center", marginBottom: "8px" }}>
                <HomeOutlined /> Theatre Management
              </Title>
              <Text
                  type="secondary"
                  style={{
                    display: "block",
                    textAlign: "center",
                    marginBottom: "24px",
                  }}
              >
                Add, edit or remove theatres and manage their halls with seating
                layouts.
              </Text>

              <Button
                  type="primary"
                  icon={<PlusOutlined />}
                  onClick={handleAddTheatre}
                  style={{ marginBottom: "16px" }}
              >
                Add Theatre
              </Button>

              <Table
                  columns={theatreColumns}
                  dataSource={theatres}
                  loading={loading}
                  rowKey="id"
                  bordered
                  pagination={{ pageSize: 8 }}
              />
            </Card>
          </Col>
        </Row>

        {/* Add/Edit Theatre Modal */}
        <Modal
            title={editingTheatre ? "Edit Theatre" : "Add Theatre"}
            open={isTheatreModalOpen}
            onOk={handleOkTheatre}
            onCancel={() => setIsTheatreModalOpen(false)}
            confirmLoading={loading}
            centered
        >
          <Form form={form} layout="vertical" name="theatre_form">
            <Form.Item
                name="name"
                label="Theatre Name"
                rules={[{ required: true, message: "Please input the theatre name!" }]}
            >
              <Input placeholder="Enter theatre name" />
            </Form.Item>
          </Form>
        </Modal>

        {/* Halls Modal */}
        <Modal
            title={`Halls for ${selectedTheatre?.name || ""}`}
            open={isHallModalOpen}
            onCancel={handleCloseHallModal}
            footer={null}
            width={850}
            centered
        >
          <Button
              type="primary"
              icon={<PlusOutlined />}
              onClick={handleAddHall}
              style={{ marginBottom: "16px" }}
          >
            Add Hall
          </Button>
          <Table
              columns={hallColumns}
              dataSource={halls}
              loading={loading}
              rowKey="id"
              bordered
              pagination={{ pageSize: 6 }}
          />
        </Modal>

        {/* Add/Edit Hall Modal */}
        <Modal
            key={isAddEditHallModalOpen ? "hall-modal-open" : "hall-modal-closed"}
            title={editingHall ? "Edit Hall" : "Add Hall"}
            open={isAddEditHallModalOpen}
            onOk={handleOkAddEditHall}
            onCancel={() => setIsAddEditHallModalOpen(false)}
            confirmLoading={loading}
            width={700}
            centered
        >
          <Form form={hallForm} layout="vertical">
            <Form.Item
                name="name"
                label="Hall Name"
                rules={[{ required: true, message: "Please input the hall name!" }]}
            >
              <Input placeholder="e.g., Audi 1, IMAX Screen" />
            </Form.Item>
          </Form>
          <Divider orientation="left">Seat Map Designer</Divider>
          <SeatMapDesigner value={currentSeatMap} onChange={setCurrentSeatMap} />
        </Modal>
      </div>
  );
}

export default TheatreManagement;
