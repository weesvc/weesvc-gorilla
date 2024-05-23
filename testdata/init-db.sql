CREATE TABLE IF NOT EXISTS places(
    id INTEGER PRIMARY KEY,
    name TEXT UNIQUE NOT NULL,
    description TEXT,
    latitude REAL,
    longitude REAL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
);

INSERT INTO places (id, name, description, latitude, longitude, created_at, updated_at)
VALUES
    (1, 'Mount Rushmore', 'Mount Rushmore National Memorial, SD, USA', 43.88031, -103.45387, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
    (2, 'Bellagio Fountains', 'Fountains of Bellagio, NV, USA', 36.11274, -115.17430, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
    (3, 'MCO', 'Orlando International Airport, FL, USA', 28.42461, -81.31075, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
    (4, 'Hoover Dam', 'Hoover Dam, Nevada, USA', 36.01604, -114.73783, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
    (5, 'Red Rocks', 'Red Rocks Park and Amphitheatre, CO, USA', 39.66551, -105.20531, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
    (6, 'MIA', 'Miami International Airport, FL, USA',	25.79516, -80.27959, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
    (7, 'Unknown', '',	0.00000, 0.00000, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
    (8, 'Grand Canyon', 'Grand Canyon National Park, AZ, USA', 36.26603, -112.36380, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
    (9, 'Hollywood Studios', 'Disney''s Hollywood Studios, FL, USA', 28.35801, -81.55918, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
    (10, 'ORD', 'O''Hare International Airport, Chicago, IL, USA', 41.97861, -87.90472, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
    (11, 'Chrysler Building', 'Chrysler Building, New York, NY, USA', 40.751652, -73.975311, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
    (12, 'Moscone Center', 'Moscone Center, San Francisco, CA, USA', 37.784172, -122.401558, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
    (13, 'US Capitol', 'United States Capitol, Washington DC, USA', 38.889805, -77.009056, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
    (14, 'Cayman Islands', 'George Town, Cayman Islands', 19.292997, -81.366806, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
    (15, 'Tokyo Tower', 'Tokyo Tower, Tokyo, Japan', 35.658581, 139.745438, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
    (16, 'HND', 'Tokyo Haneda Airport, Tokyo, Japan', 35.553333, 139.781113, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
    (17, 'Tulum', 'Tulum, Quintana Roo, Mexico', 20.214788, -87.430588, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
    (18, 'Carlsten Fortress', 'Carlsten Fortress, Marstrand, Sweden', 57.886124, 11.578453, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
    (19, 'Kungsbron Bridge', 'Kungsbron, Stockholm, Sweden', 59.332848, 18.053135, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
    (20, 'Sherlock Holmes Museum', 'The Sherlock Holmes Museum, London, UK', 51.523788, -0.158611, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

