import React from 'react';
import { Typography, Space } from 'antd';

const { Text } = Typography;

const SeatMapDisplay = ({ seatMap, bookedSeats = {} }) => {
  if (!seatMap || Object.keys(seatMap).length === 0) {
    return <Text>No seat map available for this hall.</Text>;
  }

  const rows = Object.keys(seatMap).sort((a, b) => parseInt(a) - parseInt(b));

  return (
    <div className="seat-map-display-container">
      <div className="seat-map-display-screen">Screen</div>
      {rows.map((rowKey) => (
        <div key={rowKey} className="seat-display-row">
          <Text className="seat-display-row-label">R{rowKey}</Text>
          {seatMap[rowKey].map((numSeats, colIndex) => (
            <div key={colIndex} className="seat-display-col-group">
              {Array.from({ length: numSeats }).map((_, seatIndex) => {
                const seatId = `${rowKey}-${colIndex + 1}-${seatIndex + 1}`;
                const isBooked = bookedSeats[seatId];
                return (
                  <div
                    key={seatIndex}
                    className={`seat-display-seat ${isBooked ? 'booked' : 'available'}`}
                    title={`Row ${rowKey}, Col ${colIndex + 1}, Seat ${seatIndex + 1}`}
                  ></div>
                );
              })}
            </div>
          ))}
        </div>
      ))}
      <div className="seat-display-legend">
        <Space>
          <div className="seat-display-legend-color available"></div>
          <Text>Available</Text>
          <div className="seat-display-legend-color booked"></div>
          <Text>Booked</Text>
        </Space>
      </div>
    </div>
  );
};

export default SeatMapDisplay;
