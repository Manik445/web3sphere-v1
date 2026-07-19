-- ==============================================================================
-- Web3Sphere Initial Database Schema (Rollback)
-- ==============================================================================

-- 7. Notifications & Jobs
DROP TABLE IF EXISTS api_keys CASCADE;
DROP TABLE IF EXISTS system_jobs CASCADE;
DROP TABLE IF EXISTS email_queue CASCADE;
DROP TABLE IF EXISTS notification_templates CASCADE;
DROP TABLE IF EXISTS notifications CASCADE;

-- 6. Escrow & Finance
DROP TABLE IF EXISTS crypto_transactions CASCADE;
DROP TABLE IF EXISTS payment_transactions CASCADE;
DROP TABLE IF EXISTS ledger_entries CASCADE;
DROP TABLE IF EXISTS ledger_accounts CASCADE;
DROP TABLE IF EXISTS wallet_transactions CASCADE;
DROP TABLE IF EXISTS wallets CASCADE;
DROP TABLE IF EXISTS supported_tokens CASCADE;
DROP TABLE IF EXISTS blockchain_networks CASCADE;
DROP TABLE IF EXISTS escrow_transactions CASCADE;
DROP TABLE IF EXISTS escrow_accounts CASCADE;

-- 5. Projects & Hiring
DROP TABLE IF EXISTS contracts CASCADE;
DROP TABLE IF EXISTS applications CASCADE;
DROP TABLE IF EXISTS tasks CASCADE;
DROP TABLE IF EXISTS project_members CASCADE;
DROP TABLE IF EXISTS projects CASCADE;

-- 4. Companies & Freelancers
DROP TABLE IF EXISTS user_skills CASCADE;
DROP TABLE IF EXISTS skills CASCADE;
DROP TABLE IF EXISTS freelancer_profiles CASCADE;
DROP TABLE IF EXISTS company_members CASCADE;
DROP TABLE IF EXISTS companies CASCADE;

-- 3. Auditing & System Logs
DROP TABLE IF EXISTS activity_logs CASCADE;
DROP TABLE IF EXISTS audit_logs CASCADE;

-- 2. Core Users & Authentication
DROP TABLE IF EXISTS user_devices CASCADE;
DROP TABLE IF EXISTS user_sessions CASCADE;
DROP TABLE IF EXISTS temp_data CASCADE;
DROP TABLE IF EXISTS user_info CASCADE;
DROP TABLE IF EXISTS users CASCADE;

-- 1. Reference Data & Configuration
DROP TABLE IF EXISTS feature_flags CASCADE;
DROP TABLE IF EXISTS config CASCADE;
DROP TABLE IF EXISTS country_info CASCADE;

-- Extensions
DROP EXTENSION IF EXISTS "uuid-ossp";
