import React, { useState, useEffect } from "react";
import { API_BASE_URL } from '../config/api';
import axios from "axios";
import { toast } from "react-toastify";
import { Card, Typography, Row, Col, Spin, List } from "antd";
import { HomeOutlined } from "@ant-design/icons";

const { Title, Text } = Typography;

function Theatres() {
  const [theatres, setTheatres] = useState([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const fetchTheatres = async () => {
      try {
        const response = await axios.get(`${API_BASE_URL}/theatres`);
        setTheatres(response.data.data || []);
      } catch (error) {
        console.error("Failed to fetch theatres:", error);
        toast.error(error.response?.data?.message || "Failed to load theatres.");
      } finally {
        setLoading(false);
      }
    };

    fetchTheatres();
  }, []);

  if (loading) {
    return (
        <Spin size="large" style={{ display: "block", margin: "50px auto" }} />
    );
  }

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
                <HomeOutlined /> Available Theatres
              </Title>
              <Text
                  type="secondary"
                  style={{
                    display: "block",
                    textAlign: "center",
                    marginBottom: "24px",
                  }}
              >
                Browse all theatres and their details.
              </Text>

              {theatres.length === 0 ? (
                  <Text>No theatres available.</Text>
              ) : (
                  <List
                      grid={{
                        gutter: 16,
                        xs: 1,
                        sm: 2,
                        md: 3,
                        lg: 4,
                        xl: 5,
                        xxl: 6,
                      }}
                      dataSource={theatres}
                      renderItem={(theatre) => (
                          <List.Item style={{ minWidth: 220 }}>
                            <Card
                                hoverable
                                style={{
                                  minHeight: 160,
                                  minWidth: 220,
                                  display: "flex",
                                  flexDirection: "column",
                                  justifyContent: "center",
                                  borderRadius: "8px",
                                }}
                                title={theatre.name}
                            >
                              <p>
                                <Text strong>ID:</Text> {theatre.id}
                              </p>
                            </Card>
                          </List.Item>
                      )}
                  />
              )}
            </Card>
          </Col>
        </Row>
      </div>
  );
}

export default Theatres;
