import React, { useState, useEffect, useCallback } from 'react';
import { API_BASE_URL } from '../config/api';
import { useParams, Link, useNavigate } from 'react-router-dom';
import { Typography, Spin, Button, message, Row, Col, InputNumber, Card, Space, List, Alert, Divider, Result } from 'antd';
import axios from 'axios';
import { useAuth } from '../context/AuthContext';
import dayjs from 'dayjs';

const { Title, Text } = Typography;

function Booking() {
  const { showId } = useParams();
  const { token } = useAuth();
  const navigate = useNavigate();

  const [show, setShow] = useState(null);
  const [movie, setMovie] = useState(null);
  const [hall, setHall] = useState(null);
  const [loading, setLoading] = useState(true);
  const [bookingLoading, setBookingLoading] = useState(false);
  
  const [numSeatsToBook, setNumSeatsToBook] = useState(1);
  
  // State for results
  const [bookingResult, setBookingResult] = useState(null);
  const [alternativeShows, setAlternativeShows] = useState([]);

  const fetchData = useCallback(async () => {
    if (!token) {
      message.error('You must be logged in to book tickets.');
      navigate('/login');
      return;
    }
    setLoading(true);
    setBookingResult(null); // Reset results when fetching new show data
    setAlternativeShows([]);
    try {
      const showResponse = await axios.get(`${API_BASE_URL}/shows/${showId}`, { headers: { Authorization: `Bearer ${token}` } });
      setShow(showResponse.data);

      const movieResponse = await axios.get(`${API_BASE_URL}/movies/${showResponse.data.movie_id}`, { headers: { Authorization: `Bearer ${token}` } });
      setMovie(movieResponse.data);

      const hallResponse = await axios.get(`${API_BASE_URL}/halls/${showResponse.data.hall_id}`, { headers: { Authorization: `Bearer ${token}` } });
      setHall(hallResponse.data);

    } catch (error) {
      message.error('Failed to load show details. It may no longer exist.');
    } finally {
      setLoading(false);
    }
  }, [showId, token, navigate]);

  useEffect(() => {
    fetchData();
  }, [fetchData]);

  const handleBookTickets = async () => {
    if (!numSeatsToBook || numSeatsToBook <= 0) {
      message.error('Please enter a valid number of seats.');
      return;
    }

    setBookingLoading(true);
    setBookingResult(null);
    setAlternativeShows([]);

    try {
      const bookingRequest = {
        movieId: movie.id,
        hallId: hall.id,
        time: show.time,
        numSeats: numSeatsToBook,
      };

      const response = await axios.post(`${API_BASE_URL}/bookings`, bookingRequest, {
        headers: { Authorization: `Bearer ${token}` },
      });

      setBookingResult(response.data);

    } catch (error) {
      const errorMessage = error.response?.data?.message || error.response?.data || 'An unexpected error occurred.';
      message.error(errorMessage);

      if (error.response?.status === 409 && error.response?.data?.alternatives) {
        setAlternativeShows(error.response.data.alternatives);
      } else {
        setAlternativeShows([]);
      }
    } finally {
      setBookingLoading(false);
    }
  };

  if (loading) {
    return <Spin size="large" style={{ display: 'block', margin: '50px auto' }} />;
  }

  if (!show || !movie || !hall) {
    return <Result status="warning" title="Show details not found." subTitle="The show you are looking for may no longer be available." />;
  }

  return (
    <Row justify="center" style={{ padding: '20px' }}>
      <Col xs={24} md={18} lg={14} xl={12}>
        <Card>
          <Title level={2} style={{ textAlign: 'center' }}>Book Tickets</Title>
          <Title level={4} style={{ textAlign: 'center', marginBottom: '24px' }}>{movie.title}</Title>
          
          <p><Text strong>Theatre:</Text> {hall.name}</p>
          <p><Text strong>Show Time:</Text> {dayjs(show.time).format('dddd, MMMM D, YYYY h:mm A')}</p>
          
          <Divider />

          <Space direction="vertical" size="large" style={{ width: '100%', marginTop: '20px' }}>
            <Space align="center">
              <Text strong style={{ fontSize: '16px' }}>Number of Seats:</Text>
              <InputNumber
                min={1}
                max={20}
                value={numSeatsToBook}
                onChange={setNumSeatsToBook}
                size="large"
              />
            </Space>
            <Button
              type="primary"
              onClick={handleBookTickets}
              loading={bookingLoading}
              block
              size="large"
            >
              Find Seats & Book Now
            </Button>
          </Space>

          {bookingResult && (
            <Result
              style={{ marginTop: '20px' }}
              status="success"
              title="Booking Successful!"
              subTitle={`Booking ID: ${bookingResult.id}`}
              extra={[
                <Text key="seats" strong>Your Seats: {bookingResult.seat_ids.join(', ')}</Text>,
              ]}
            />
          )}

          {alternativeShows.length > 0 && (
            <Alert
              style={{ marginTop: '20px' }}
              type="warning"
              showIcon
              message="Seats Not Available Together"
              description={
                <div>
                  <p>Sorry, we couldn't find {numSeatsToBook} contiguous seats for this show. Here are some alternatives for the same day:</p>
                  <List
                    dataSource={alternativeShows}
                    renderItem={altShow => (
                      <List.Item>
                        <Space>
                          <Text>Time: {dayjs(altShow.time).format('h:mm A')}</Text>
                          <Text>Hall: {altShow.hall_id}</Text>
                          <Link to={`/book/${altShow.id}`}>
                            <Button>Try this show</Button>
                          </Link>
                        </Space>
                      </List.Item>
                    )}
                  />
                </div>
              }
            />
          )}
        </Card>
      </Col>
    </Row>
  );
}

export default Booking;
