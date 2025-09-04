import React, { useState, useEffect } from "react";
import axios from "axios";
import { Link } from "react-router-dom";
import { toast } from "react-toastify";
import { Card, Typography, Row, Col, Spin } from "antd";
import { VideoCameraOutlined } from "@ant-design/icons";
import { API_BASE_URL } from "../config/api";

const { Title, Text } = Typography;

function Movies() {
  const [movies, setMovies] = useState([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const fetchMovies = async () => {
      try {
        const response = await axios.get(`${API_BASE_URL}/movies`);
        setMovies(response.data.data || []);
      } catch (error) {
        console.error("Failed to fetch movies:", error);
        toast.error(error.response?.data?.message || "Failed to load movies.");
      } finally {
        setLoading(false);
      }
    };

    fetchMovies();
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
                <VideoCameraOutlined /> Available Movies
              </Title>
              <Text
                  type="secondary"
                  style={{
                    display: "block",
                    textAlign: "center",
                    marginBottom: "24px",
                  }}
              >
                Browse movies and check their showtimes.
              </Text>

              <Row gutter={[16, 16]}>
                {movies.map((movie) => (
                    <Col
                        xs={24}
                        sm={12}
                        md={8}
                        lg={6}
                        xl={4}
                        key={movie.id}
                        style={{ display: "flex" }}
                    >
                      <Link
                          to={`/movies/${movie.id}`}
                          style={{ width: "100%", textDecoration: "none" }}
                      >
                        <Card
                            hoverable
                            style={{
                              width: "100%",
                              borderRadius: "10px",
                              overflow: "hidden",
                              boxShadow: "0 4px 8px rgba(0,0,0,0.1)",
                              transition: "all 0.3s ease",
                              minHeight: "160px",
                              display: "flex",
                              flexDirection: "column",
                              justifyContent: "space-between",
                            }}
                            bodyStyle={{ padding: "16px" }}
                        >
                          <Title
                              level={4}
                              style={{
                                marginBottom: "5px",
                                whiteSpace: "nowrap",
                                overflow: "hidden",
                                textOverflow: "ellipsis",
                              }}
                          >
                            {movie.title}
                          </Title>
                          <Text type="secondary">
                            Duration: {movie.duration_minutes} mins
                          </Text>
                          <p
                              style={{
                                marginTop: "10px",
                                color: "#1890ff",
                                fontWeight: "bold",
                              }}
                          >
                            View Showtimes →
                          </p>
                        </Card>
                      </Link>
                    </Col>
                ))}
              </Row>
            </Card>
          </Col>
        </Row>
      </div>
  );
}

export default Movies;
