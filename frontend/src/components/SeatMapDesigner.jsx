import React from 'react';
import { InputNumber, Button, Space, Typography, Row, Col, message, Divider } from 'antd';
import { MinusCircleOutlined, PlusOutlined } from '@ant-design/icons';

const { Text } = Typography;

const SeatMapDesigner = ({ value, onChange }) => {
  const seatMap = value || {};

  const updateParent = (newSeatMap) => {
    if (typeof onChange === 'function') {
      onChange(newSeatMap);
    }
  };

  const handleAddRow = () => {
    const newRowNum = Object.keys(seatMap).length > 0 ? Math.max(...Object.keys(seatMap).map(k => parseInt(k))) + 1 : 1;
    const newSeatMap = {
      ...seatMap,
      [newRowNum.toString()]: [2, 2, 2], // Default to 3 columns with 2 seats each
    };
    updateParent(newSeatMap);
  };

  const handleRemoveRow = (rowKey) => {
    const newSeatMap = { ...seatMap };
    delete newSeatMap[rowKey];
    updateParent(newSeatMap);
  };

  const handleColumnChange = (rowKey, colIndex, numSeats) => {
    if (numSeats < 2) {
      message.error('Minimum 2 seats required per column.');
      return;
    }
    const newSeatMap = { ...seatMap };
    newSeatMap[rowKey][colIndex] = numSeats;
    updateParent(newSeatMap);
  };

  const rows = Object.keys(seatMap).sort((a, b) => parseInt(a) - parseInt(b));

  return (
    <div style={{ border: '1px solid #d9d9d9', padding: '16px', borderRadius: '8px' }}>
      <Text strong>Define Hall Layout</Text>
      <Divider style={{ margin: '12px 0' }} />
      {rows.map((rowKey) => (
        <Row key={rowKey} gutter={[8, 8]} align="middle" style={{ marginBottom: '10px' }}>
          <Col flex="60px">
            <Text strong>Row {rowKey}:</Text>
          </Col>
          {seatMap[rowKey].map((numSeats, colIndex) => (
            <Col key={colIndex} flex="auto">
              <Space direction="vertical" size="small">
                <Text type="secondary">Column {colIndex + 1}</Text>
                <InputNumber
                  min={2}
                  value={numSeats}
                  onChange={(val) => handleColumnChange(rowKey, colIndex, val)}
                  style={{ width: '100%' }}
                />
              </Space>
            </Col>
          ))}
          <Col flex="32px">
            <Button
              type="text"
              danger
              icon={<MinusCircleOutlined />}
              onClick={() => handleRemoveRow(rowKey)}
            />
          </Col>
        </Row>
      ))}
      <Button type="dashed" onClick={handleAddRow} block icon={<PlusOutlined />}>
        Add Row
      </Button>

      <Divider style={{ margin: '20px 0' }}>Visual Preview</Divider>
      <div style={{ backgroundColor: '#fafafa', padding: '10px', borderRadius: '4px' }}>
        <div style={{ width: '100%', backgroundColor: '#333', color: '#fff', textAlign: 'center', padding: '5px', marginBottom: '10px', borderRadius: '2px' }}>Screen</div>
        {rows.map((rowKey) => (
          <div key={rowKey} style={{ display: 'flex', marginBottom: '5px', alignItems: 'center', justifyContent: 'center' }}>
            <Text style={{ marginRight: '10px', fontWeight: 'bold', minWidth: '30px' }}>R{rowKey}</Text>
            {seatMap[rowKey].map((numSeats, colIndex) => (
              <div key={colIndex} style={{ display: 'flex', border: '1px solid #ccc', margin: '0 5px', padding: '2px', borderRadius: '4px' }}>
                {Array.from({ length: numSeats }).map((_, seatIndex) => (
                  <div
                    key={seatIndex}
                    style={{
                      width: '20px',
                      height: '20px',
                      backgroundColor: '#1890ff',
                      margin: '2px',
                      borderRadius: '3px',
                    }}
                    title={`Row ${rowKey}, Col ${colIndex + 1}, Seat ${seatIndex + 1}`}
                  ></div>
                ))}
              </div>
            ))}
          </div>
        ))}
      </div>
    </div>
  );
};

export default SeatMapDesigner;
