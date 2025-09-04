import React from 'react';
import { Typography, Row, Col, Card, Button, Space } from 'antd';
import { Link } from 'react-router-dom';
import { PlayCircleOutlined, CalendarOutlined, EnvironmentOutlined } from '@ant-design/icons';

const { Title, Paragraph } = Typography;

function Home() {
  return (
    <div style={{ textAlign: 'center', padding: '50px 20px' }}>
      <Row justify="center" align="middle" style={{ minHeight: 'calc(100vh - 200px)' }}>
        <Col xs={24} md={18} lg={14}>
          <Card style={{ borderRadius: '15px', overflow: 'hidden', boxShadow: '0 10px 30px rgba(0,0,0,0.1)' }}>
            <div style={{ padding: '40px 20px', background: 'linear-gradient(135deg, #1890ff 0%, #722ed1 100%)', color: '#fff' }}>
              <Title level={1} style={{ color: '#fff', marginBottom: '10px', fontSize: '3.5em', fontWeight: 800 }}>
                Welcome to AlgoBharat Cinemas
              </Title>
              <Paragraph style={{ color: 'rgba(255,255,255,0.8)', fontSize: '1.2em', maxWidth: '600px', margin: '0 auto 30px' }}>
                Your ultimate destination for seamless movie ticket booking. Discover the latest blockbusters and secure your seats with ease.
              </Paragraph>
              <Space size="large">
                <Button type="primary" size="large" icon={<PlayCircleOutlined />} style={{ borderRadius: '30px', height: '50px', padding: '0 30px', background: '#fff', borderColor: '#fff', color: '#1890ff', fontWeight: 600 }}>
                  <Link to="/movies" style={{ color: '#1890ff' }}>Browse Movies</Link>
                </Button>
                <Button type="default" size="large" icon={<EnvironmentOutlined />} style={{ borderRadius: '30px', height: '50px', padding: '0 30px', background: 'rgba(255,255,255,0.2)', borderColor: 'rgba(255,255,255,0.5)', color: '#fff', fontWeight: 600 }}>
                  <Link to="/theatres" style={{ color: '#fff' }}>Find Theatres</Link>
                </Button>
              </Space>
            </div>
            <div style={{ padding: '30px 20px', backgroundColor: '#f0f2f5' }}>
              <Title level={3} style={{ color: '#2c3e50', marginBottom: '20px' }}>Why Choose Us?</Title>
              <Row gutter={[32, 32]} justify="center">
                <Col xs={24} md={8}>
                  <Card hoverable style={{ borderRadius: '10px', boxShadow: '0 5px 15px rgba(0,0,0,0.05)' }}>
                    <CalendarOutlined style={{ fontSize: '3em', color: '#1890ff', marginBottom: '10px' }} />
                    <Title level={4}>Easy Booking</Title>
                    <Paragraph>Book your favorite movie tickets in just a few clicks, anytime, anywhere.</Paragraph>
                  </Card>
                </Col>
                <Col xs={24} md={8}>
                  <Card hoverable style={{ borderRadius: '10px', boxShadow: '0 5px 15px rgba(0,0,0,0.05)' }}>
                    <PlayCircleOutlined style={{ fontSize: '3em', color: '#722ed1', marginBottom: '10px' }} />
                    <Title level={4}>Latest Movies</Title>
                    <Paragraph>Stay updated with the newest releases and trending films.</Paragraph>
                  </Card>
                </Col>
                <Col xs={24} md={8}>
                  <Card hoverable style={{ borderRadius: '10px', boxShadow: '0 5px 15px rgba(0,0,0,0.05)' }}>
                    <EnvironmentOutlined style={{ fontSize: '3em', color: '#52c41a', marginBottom: '10px' }} />
                    <Title level={4}>Nearby Theatres</Title>
                    <Paragraph>Locate theatres near you and check showtimes instantly.</Paragraph>
                  </Card>
                </Col>
              </Row>
            </div>
          </Card>
        </Col>
      </Row>
    </div>
  );
}

export default Home;
