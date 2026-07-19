-- ==============================================================================
-- Web3Sphere Initial Seed Data
-- ==============================================================================

-- Core roles
INSERT INTO config (key, value, type, description) VALUES
('app.maintenance_mode', 'false', 'boolean', 'Whether the application is in maintenance mode'),
('app.signup_enabled', 'true', 'boolean', 'Whether new user signups are allowed'),
('app.default_currency', 'USDC', 'string', 'Default currency for platform transactions'),
('app.platform_fee_percent', '2.5', 'float', 'Default platform fee percentage for contracts')
ON CONFLICT (key) DO NOTHING;

-- Initial feature flags
INSERT INTO feature_flags (key, enabled, description) VALUES
('escrow_v1', true, 'Enable V1 of the smart contract escrow system'),
('ai_matching', false, 'Enable AI-powered talent matching'),
('crypto_payments', true, 'Enable crypto payments for contracts')
ON CONFLICT (key) DO NOTHING;

-- Initial Blockchain Networks
INSERT INTO blockchain_networks (name, chain_id, rpc_url, explorer_url, native_currency, is_testnet, enabled) VALUES
('Ethereum Mainnet', 1, 'https://mainnet.infura.io/v3/', 'https://etherscan.io', 'ETH', false, true),
('Polygon Mainnet', 137, 'https://polygon-rpc.com', 'https://polygonscan.com', 'MATIC', false, true),
('Ethereum Sepolia', 11155111, 'https://sepolia.infura.io/v3/', 'https://sepolia.etherscan.io', 'ETH', true, true),
('Polygon Amoy', 80002, 'https://rpc-amoy.polygon.technology', 'https://www.oklink.com/amoy', 'MATIC', true, true)
ON CONFLICT (chain_id) DO NOTHING;

-- Sample Countries (Subset to avoid huge file, more can be added)
INSERT INTO country_info (country_name, iso2, iso3, phone_code, currency, currency_symbol, region, sub_region, enabled) VALUES
('United States', 'US', 'USA', '+1', 'USD', '$', 'Americas', 'Northern America', true),
('United Kingdom', 'GB', 'GBR', '+44', 'GBP', '£', 'Europe', 'Northern Europe', true),
('Canada', 'CA', 'CAN', '+1', 'CAD', '$', 'Americas', 'Northern America', true),
('Australia', 'AU', 'AUS', '+61', 'AUD', '$', 'Oceania', 'Australia and New Zealand', true),
('Germany', 'DE', 'DEU', '+49', 'EUR', '€', 'Europe', 'Western Europe', true),
('India', 'IN', 'IND', '+91', 'INR', '₹', 'Asia', 'Southern Asia', true),
('Singapore', 'SG', 'SGP', '+65', 'SGD', '$', 'Asia', 'South-Eastern Asia', true),
('United Arab Emirates', 'AE', 'ARE', '+971', 'AED', 'د.إ', 'Asia', 'Western Asia', true)
ON CONFLICT (iso2) DO NOTHING;

-- Sample Skills
INSERT INTO skills (name, category) VALUES
('Solidity', 'Smart Contracts'),
('Rust', 'Smart Contracts'),
('Go', 'Backend'),
('Node.js', 'Backend'),
('React', 'Frontend'),
('Next.js', 'Frontend'),
('Web3.js', 'Blockchain Integration'),
('Ethers.js', 'Blockchain Integration'),
('Hardhat', 'Testing & Deployment'),
('Foundry', 'Testing & Deployment'),
('Security Auditing', 'Security'),
('DeFi Architecture', 'Architecture')
ON CONFLICT (name) DO NOTHING;
