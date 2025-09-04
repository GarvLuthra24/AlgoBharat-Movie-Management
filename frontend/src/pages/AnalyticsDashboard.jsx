import React, { useState, useEffect } from 'react';
import { API_BASE_URL } from '../config/api';
import { Typography, Select, Card, Spin, Row, Col } from 'antd';
import axios from 'axios';
import { useAuth } from '../context/AuthContext';
import { toast } from 'react-toastify';
import {AccountBookOutlined} from "@ant-design/icons";

const { Title, Text } = Typography;
const { Option } = Select;

function AnalyticsDashboard() {
  const { token } = useAuth();
  const [movies, setMovies] = useState([]);
  const [selectedMovie, setSelectedMovie] = useState(null);
  const [revenue, setRevenue] = useState(null);
  const [loading, setLoading] = useState(false);

  useEffect(() => {
    const fetchMovies = async () => {
      setLoading(true);
      try {
        const response = await axios.get(`${API_BASE_URL}/movies`, {
          headers: { Authorization: `Bearer ${token}` },
        });
        setMovies(response.data.data || []);
      } catch (error) {
        console.error('Failed to fetch movies:', error);
        toast.error(error.response?.data?.message || 'Failed to load movies for analytics.');
      } finally {
        setLoading(false);
      }
    };

    fetchMovies();
  }, [token]);

  const fetchMovieRevenue = async (movieId) => {
    setLoading(true);
    try {
      const response = await axios.get(`${API_BASE_URL}/analytics/movies/${movieId}/revenue`, {
        headers: { Authorization: `Bearer ${token}` },
      });
      setRevenue(response.data.data?.total_revenue || null);
    } catch (error) {
      console.error('Failed to fetch movie revenue:', error);
      toast.error(error.response?.data?.message || 'Failed to fetch revenue for the selected movie.');
      setRevenue(null);
    } finally {
      setLoading(false);
    }
  };

  const handleMovieChange = (value) => {
    setSelectedMovie(value);
    if (value) fetchMovieRevenue(value);
    else setRevenue(null);
  };

  return (
      <div style={{ padding: '2em 0', background: '#f0f2f5', minHeight: '100vh' }}>
        <Row justify="center">
          <Col xs={24} lg={21}>
            <Card
                bordered={false}
                style={{
                  borderRadius: '1em',
                  padding: '2em',
                  width: '90%',
                  margin: '0 auto',
                }}
            >

              <Title level={3} style={{ textAlign: 'center', marginBottom: '8px' }}>
                <AccountBookOutlined /> Movie Revenue Analytics
              </Title>

              <Card
                  style={{
                    marginBottom: '2em',
                    borderRadius: '0.8em',
                    padding: '1.5em',
                  }}
                  bodyStyle={{ padding: '1em' }}
              >
                <Text strong>🍿 Select a Movie:</Text>
                <Select
                    showSearch
                    placeholder="Select a movie"
                    optionFilterProp="children"
                    onChange={handleMovieChange}
                    value={selectedMovie}
                    style={{ width: '100%', marginTop: '0.8em' }}
                    loading={loading}
                    filterOption={(input, option) =>
                        option.children.toLowerCase().includes(input.toLowerCase())
                    }
                >
                  {movies.map((movie) => (
                      <Option key={movie.id} value={movie.id}>
                        {movie.title}
                      </Option>
                  ))}
                </Select>
              </Card>

              {loading && selectedMovie && (
                  <div style={{ textAlign: 'center', margin: '3em 0' }}>
                    <Spin size="large" />
                  </div>
              )}

              {revenue !== null && !loading && (
                  <Card
                      style={{
                        marginBottom: '1.5em',
                        borderRadius: '0.8em',
                        padding: '1.5em',
                        textAlign: 'center',
                      }}
                  >
                    <Title level={4}>
                      💰 Total Revenue for {movies.find((m) => m.id === selectedMovie)?.title}:
                    </Title>
                    <Text strong style={{ fontSize: '2em', color: '#1890ff' }}>
                      ₹ {revenue.toFixed(2)}
                    </Text>
                  </Card>
              )}

              {selectedMovie && revenue === null && !loading && (
                  <Text style={{ display: 'block', textAlign: 'center', marginTop: '1.5em' }}>
                    No revenue data available for this movie yet.
                  </Text>
              )}
            </Card>
          </Col>
        </Row>
      </div>
  );
}

export default AnalyticsDashboard;
