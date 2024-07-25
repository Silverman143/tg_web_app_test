CREATE TABLE Users (
    id SERIAL PRIMARY KEY,
    telegram_id BIGINT NOT NULL UNIQUE,
    user_name VARCHAR(255) NOT NULL,
    lang VARCHAR(10),
    registration_date TIMESTAMP NOT NULL,
    avatar_url VARCHAR(255),
    stars_balance INT,                
    wallet_address VARCHAR(255),
    -- User key to invite other users
    referrer_key VARCHAR(255),      
    tasks_complete JSONB,
    global_rank INT DEFAULT 0
);

CREATE TABLE All_Time_Leaders (
    id SERIAL PRIMARY KEY,
    user_id BIGINT UNIQUE REFERENCES Users(id) ON DELETE CASCADE,
    stars INT DEFAULT 0
);

CREATE TABLE Referrals (
    id SERIAL PRIMARY KEY,
    referrer BIGINT REFERENCES Users(id) ON DELETE CASCADE,
    referral BIGINT REFERENCES Users(id) ON DELETE CASCADE,
    bonus NUMERIC,
    UNIQUE (referrer, referral)
);

CREATE TABLE Current_Week_Leaders (
    id SERIAL PRIMARY KEY,
    user_id BIGINT UNIQUE REFERENCES Users(id) ON DELETE CASCADE,
    stars INT DEFAULT 0,
    week_start DATE NOT NULL,
    week_end DATE NOT NULL
);

CREATE TABLE Historic_Week_Leaders (
    id SERIAL PRIMARY KEY,
    user_id BIGINT REFERENCES Users(id) ON DELETE CASCADE,
    stars INT DEFAULT 0,
    week_start DATE NOT NULL,
    week_end DATE NOT NULL
);

CREATE TABLE Current_Month_Leaders (
    id SERIAL PRIMARY KEY,
    user_id BIGINT UNIQUE REFERENCES Users(id) ON DELETE CASCADE,
    stars INT DEFAULT 0,
    month_start DATE NOT NULL,
    month_end DATE NOT NULL
);

CREATE TABLE Historic_Month_Leaders (
    id SERIAL PRIMARY KEY,
    user_id BIGINT REFERENCES Users(id) ON DELETE CASCADE,
    stars INT DEFAULT 0,
    month_start DATE NOT NULL,
    month_end DATE NOT NULL
);

CREATE TABLE Users_daily_bonuses (
    user_id BIGINT PRIMARY KEY REFERENCES Users(id) ON DELETE CASCADE,
    last_collected DATE,
    days_counter INT DEFAULT 0
);

-- Bonus days prices
CREATE TABLE Daily_bonuses_info (
    day INT UNIQUE,
    price INT
);

CREATE TABLE Tasks (
    id SERIAL PRIMARY KEY,
    daily BOOLEAN,
    type VARCHAR(255),
    description TEXT,
    price INT,
    telegram_sub BOOLEAN,
    start_date TIMESTAMP,
    end_date TIMESTAMP,
    is_active BOOLEAN
);

CREATE TABLE Tasks_updates (
    id SERIAL PRIMARY KEY,
    task INT REFERENCES Tasks(id) ON DELETE CASCADE,
    user_id BIGINT REFERENCES Users(id) ON DELETE CASCADE,
    date TIMESTAMP NOT NULL,
    status VARCHAR(255)
);

CREATE TABLE Transactions (
    id SERIAL PRIMARY KEY,
    type VARCHAR(255),
    user_id BIGINT REFERENCES Users(id) ON DELETE CASCADE,
    date TIMESTAMP NOT NULL,
    amount NUMERIC,
    currency VARCHAR(10),
    status VARCHAR(255)
);
