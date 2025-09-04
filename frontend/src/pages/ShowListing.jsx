import React, { useState, useEffect } from "react";
import { API_BASE_URL } from '../config/api';
import { useParams, useNavigate } from "react-router-dom";
import {
  Typography,
  Spin,
  Card,
  List,
  Button,
  Modal,
  InputNumber,
  Space,
  Row,
  Col,
} from "antd";
import axios from "axios";
import dayjs from "dayjs";
import { useAuth } from "../context/AuthContext";
import SeatMapDisplay from "../components/SeatMapDisplay";
import { toast } from "react-toastify";

const { Title, Text } = Typography;

function ShowListing() {
  const { id: movieId } = useParams();
  const { token, user } = useAuth();
  const navigate = useNavigate();
  const [movie, setMovie] = useState(null);
  const [shows, setShows] = useState([]);
  const [theatres, setTheatres] = useState({});
  const [halls, setHalls] = useState({});
  const [loading, setLoading] = useState(true);

  // Booking modal state
  const [isBookingModalVisible, setIsBookingModalVisible] = useState(false);
  const [selectedShow, setSelectedShow] = useState(null);
  const [numSeats, setNumSeats] = useState(1);
  const [bookingLoading, setBookingLoading] = useState(false);
  const [currentHallDetails, setCurrentHallDetails] = useState(null);
  const [currentBookedSeats, setCurrentBookedSeats] = useState({});

  // Alternative shows modal state
  const [isAlternativesModalVisible, setIsAlternativesModalVisible] =
      useState(false);
  const [alternativeShows, setAlternativeShows] = useState([]);

  useEffect(() => {
    const fetchData = async () => {
      try {
        const [movieRes, showsRes, theatresRes, hallsRes] = await Promise.all([
          axios.get(`${API_BASE_URL}/movies/${movieId}`),
          axios.get(`${API_BASE_URL}/shows`),
          axios.get(`${API_BASE_URL}/theatres`),
          axios.get(`${API_BASE_URL}/halls`),
        ]);

        setMovie(movieRes.data?.data || movieRes.data);
        setShows(
            (showsRes.data?.data || showsRes.data || []).filter(
                (s) => s.movie_id === movieId
            )
        );

        setTheatres(
            (theatresRes.data?.data || theatresRes.data || []).reduce(
                (acc, t) => ({ ...acc, [t.id]: t }),
                {}
            )
        );
        setHalls(
            (hallsRes.data?.data || hallsRes.data || []).reduce(
                (acc, h) => ({ ...acc, [h.id]: h }),
                {}
            )
        );
      } catch (error) {
        toast.error(
            error.response?.data?.message || "Failed to load show details."
        );
      } finally {
        setLoading(false);
      }
    };

    fetchData();
  }, [movieId]);

  const handleBookClick = async (show) => {
    if (!user) {
      toast.error("Please log in to book tickets.");
      return;
    }
    setSelectedShow(show);
    setNumSeats(1);
    setBookingLoading(true);

    try {
      const hallResponse = await axios.get(
          `${API_BASE_URL}/halls/${show.hall_id}`
      );
      setCurrentHallDetails(hallResponse.data.data);

      const bookingsResponse = await axios.get(
          `${API_BASE_URL}/bookings?showId=${show.id}`,
          {
            headers: { Authorization: `Bearer ${token}` },
          }
      );
      const booked = {};
      (bookingsResponse.data.data || []).forEach((booking) => {
        booking.seat_ids.forEach((seatId) => {
          booked[seatId] = true;
        });
      });
      setCurrentBookedSeats(booked);

      setIsBookingModalVisible(true);
    } catch (error) {
      toast.error(
          error.response?.data?.message ||
          "Failed to load hall details or booked seats."
      );
    } finally {
      setBookingLoading(false);
    }
  };

  const handleBookingConfirm = async () => {
    if (!selectedShow) return;
    setBookingLoading(true);
    try {
      const bookingRequest = {
        movieId: selectedShow.movie_id,
        hallId: selectedShow.hall_id,
        time: selectedShow.time, // This should already be in UTC format from the server
        numSeats,
      };
      
      // Debug logging
      console.log('Booking request time:', selectedShow.time);

      await axios.post(`${API_BASE_URL}/bookings`, bookingRequest, {
        headers: { Authorization: `Bearer ${token}` },
      });

      setIsBookingModalVisible(false);
      setIsAlternativesModalVisible(false);
      toast.success(`Booking successful! Your seats are confirmed.`);
      navigate("/");
    } catch (error) {
      const errorData = error.response?.data;
      const alternatives = errorData?.data?.alternatives;

      if (error.response?.status === 409 && alternatives) {
        setAlternativeShows(alternatives);
        setIsAlternativesModalVisible(true);
        setIsBookingModalVisible(false);
      } else {
        toast.error(errorData?.message || "Booking failed. Please try again.");
      }
    } finally {
      setBookingLoading(false);
    }
  };

  if (loading) {
    return (
        <Spin size="large" style={{ display: "block", margin: "50px auto" }} />
    );
  }

  if (!movie) {
    return (
        <Title level={3} style={{ textAlign: "center" }}>
          Movie not found.
        </Title>
    );
  }

  return (
      <div style={{ padding: "20px" }}>
        <Row justify="center">
          <Col xs={24} lg={20}>
            <Card>
              <Title
                  level={2}
                  style={{ textAlign: "center", marginBottom: "24px" }}
              >
                {movie.title}
              </Title>
              <Title level={4}>Available Shows:</Title>
              {shows.length === 0 ? (
                  <Text>No shows available for this movie.</Text>
              ) : (
                  <List
                      grid={{
                        gutter: 16,
                        xs: 1,
                        sm: 2,
                        md: 2,
                        lg: 3,
                        xl: 4,
                      }}
                      dataSource={shows}
                      renderItem={(show) => {
                        const hall = halls[show.hall_id];
                        const theatre = hall ? theatres[hall.theatre_id] : null;
                        return (
                            <List.Item style={{ minWidth: 280 }}>
                              <Card
                                  hoverable
                                  style={{
                                    minHeight: 200,
                                    minWidth: 260,
                                    display: "flex",
                                    flexDirection: "column",
                                    justifyContent: "space-between",
                                  }}
                                  title={dayjs(show.time).format("YYYY-MM-DD HH:mm")}
                                  actions={[
                                    <Button
                                        type="primary"
                                        onClick={() => handleBookClick(show)}
                                    >
                                      Book Tickets
                                    </Button>,
                                  ]}
                              >
                                <p>
                                  <Text strong>Theatre:</Text>{" "}
                                  {theatre ? theatre.name : "N/A"}
                                </p>
                                <p>
                                  <Text strong>Hall:</Text>{" "}
                                  {hall ? hall.name : "N/A"}
                                </p>
                                <p>
                                  <Text strong>Price:</Text> ₹
                                  {show.price ? show.price.toFixed(2) : "N/A"}
                                </p>
                              </Card>
                            </List.Item>
                        );
                      }}
                  />
              )}
            </Card>
          </Col>
        </Row>

        {/* Booking Modal */}
        <Modal
            title={`Book tickets for ${movie.title}`}
            open={isBookingModalVisible}
            onOk={handleBookingConfirm}
            onCancel={() => setIsBookingModalVisible(false)}
            confirmLoading={bookingLoading}
            okText="Confirm Booking"
            width={currentHallDetails ? 800 : 520}
        >
          <Space direction="vertical" style={{ width: "100%" }}>
            <Text>Select the number of seats you want to book together.</Text>
            <InputNumber
                min={1}
                max={10}
                value={numSeats}
                onChange={setNumSeats}
                style={{ width: "100%" }}
            />
            {currentHallDetails && (
                <SeatMapDisplay
                    seatMap={currentHallDetails.seat_map}
                    bookedSeats={currentBookedSeats}
                />
            )}
          </Space>
        </Modal>

        {/* Alternative Shows Modal */}
        <Modal
            title="Alternative Shows Available"
            open={isAlternativesModalVisible}
            onCancel={() => setIsAlternativesModalVisible(false)}
            footer={[
              <Button key="back" onClick={() => setIsAlternativesModalVisible(false)}>
                Close
              </Button>,
            ]}
        >
          <Text>
            Sorry, we couldn't find {numSeats} seats together for the selected
            show. Here are some alternatives for the same day:
          </Text>
          <List
              style={{ marginTop: "20px" }}
              grid={{ gutter: 16, xs: 1, sm: 2, md: 2 }}
              dataSource={alternativeShows}
              renderItem={(altShow) => {
                const hall = halls[altShow.hall_id];
                const theatre = hall ? theatres[hall.theatre_id] : null;
                return (
                    <List.Item style={{ minWidth: 260 }}>
                      <Card
                          size="small"
                          title={dayjs(altShow.time).format("YYYY-MM-DD HH:mm")}
                      >
                        <p>
                          <Text strong>Theatre:</Text>{" "}
                          {theatre ? theatre.name : "N/A"}
                        </p>
                        <p>
                          <Text strong>Hall:</Text> {hall ? hall.name : "N/A"}
                        </p>
                        <Button
                            type="primary"
                            size="small"
                            onClick={() => handleBookClick(altShow)}
                            style={{ marginTop: "10px" }}
                        >
                          Book This Show
                        </Button>
                      </Card>
                    </List.Item>
                );
              }}
          />
        </Modal>
      </div>
  );
}

export default ShowListing;
