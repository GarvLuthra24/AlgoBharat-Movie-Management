import React, { useState, useEffect } from "react";
import { API_BASE_URL } from "../config/api";
import axios from "axios";
import { toast } from "react-toastify";
import { Card, Typography, Row, Col, Spin } from "antd";
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
        <div style={{ padding: "30px" }}>
            <Row justify="center">
                <Col xs={24} style={{ maxWidth: "1200px", width: "100%" }}>
                    <div
                        style={{
                            textAlign: "center",
                            marginBottom: "30px",
                        }}
                    >
                        <Title level={2}>
                            <HomeOutlined /> Available Theatres
                        </Title>
                        <Text type="secondary" style={{ fontSize: "16px" }}>
                            Browse all theatres and their details.
                        </Text>
                    </div>

                    {theatres.length === 0 ? (
                        <Text>No theatres available.</Text>
                    ) : (
                        <Row gutter={[24, 24]} justify="center">
                            {theatres.map((theatre) => (
                                <Col
                                    key={theatre.id}
                                    xs={24}
                                    sm={12}
                                    md={8}
                                    lg={6}
                                    style={{ display: "flex", justifyContent: "center" }}
                                >
                                    <Card
                                        hoverable
                                        style={{
                                            width: "100%",
                                            maxWidth: 320,
                                            minHeight: 180,
                                            borderRadius: "16px",
                                            boxShadow: "0 4px 12px rgba(0,0,0,0.1)",
                                            transition: "all 0.3s ease",
                                            display: "flex",
                                            flexDirection: "column",
                                            justifyContent: "center",
                                            textAlign: "center",
                                        }}
                                        bodyStyle={{ padding: "20px" }}
                                    >
                                        <Title level={4} style={{ marginBottom: "12px" }}>
                                            {theatre.name}
                                        </Title>
                                        <Text type="secondary" style={{ fontSize: "14px" }}>
                                            ID: {theatre.id}
                                        </Text>
                                    </Card>
                                </Col>
                            ))}
                        </Row>
                    )}
                </Col>
            </Row>
        </div>
    );
}

export default Theatres;
