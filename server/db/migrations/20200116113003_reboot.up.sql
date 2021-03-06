-- ################################################################################################@
-- Users
-- ################################################################################################@
CREATE TABLE users (
    id UUID NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL,

    avatar TEXT,
    username TEXT NOT NULL,
    display_name TEXT NOT NULL,
    email TEXT NOT NULL,
    first_name TEXT NOT NULL,
    last_name TEXT NOT NULL,
    is_admin BOOLEAN NOT NULL,
    disabled_at TIMESTAMP WITH TIME ZONE,
    auth_key_hash TEXT NOT NULL,
    bio TEXT NOT NULL,
    public_key BYTEA NOT NULL,
    encrypted_private_key BYTEA NOT NULL,
    private_key_nonce BYTEA NOT NULL,
    encrypted_master_key BYTEA NOT NULL,
    master_key_nonce BYTEA NOT NULL,
    state BIGINT NOT NULL,
    -- password_update_request_at TIMESTAMP WITH TIME ZONE,
    -- password_update_code_hash TEXT,

    PRIMARY KEY(id)
);
CREATE UNIQUE INDEX index_users_on_username ON users (username);
CREATE UNIQUE INDEX index_users_on_email ON users (email);
CREATE UNIQUE INDEX index_users_on_avatar ON users (avatar);


CREATE TABLE sessions (
    id UUID NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL,

    hash BYTEA NOT NULL,
    device_os TEXT NOT NULL,
    device_type TEXT NOT NULL,

    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,

    PRIMARY KEY(id)
);
CREATE INDEX index_sessions_on_user_id ON sessions (user_id);

CREATE TABLE pending_sessions (
  id UUID NOT NULL,
  created_at TIMESTAMP WITH TIME ZONE NOT NULL,
  updated_at TIMESTAMP WITH TIME ZONE NOT NULL,

  hash BYTEA NOT NULL,
  failed_attempts BIGINT NOT NULL,

  user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,

  PRIMARY KEY(id)
);

CREATE TABLE pending_users (
    id UUID NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL,

    email TEXT NOT NULL,
    display_name TEXT NOT NULL,
    code_hash TEXT NOT NULL,
    failed_attempts BIGINT NOT NULL,
    verified_at TIMESTAMP WITH TIME ZONE,

    PRIMARY KEY(id)
);


-- ################################################################################################@
-- Groups
-- ################################################################################################@

CREATE TABLE groups (
    id UUID NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL,

    avatar TEXT,
    name TEXT NOT NULL,
    description TEXT NOT NULL,
    state BIGINT NOT NULL,

    PRIMARY KEY(id)
);
CREATE UNIQUE INDEX index_groups_on_avatar ON groups (avatar);


CREATE TABLE groups_members (
    joined_at TIMESTAMP WITH TIME ZONE NOT NULL,
    inviter_id UUID NOT NULL REFERENCES users(id),

    role TEXT NOT NULL,
    encrypted_master_key BYTEA NOT NULL,
    master_key_nonce BYTEA NOT NULL,

    group_id UUID NOT NULL REFERENCES groups(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE
);

CREATE INDEX index_groups_members_on_group_id ON groups_members (group_id);
CREATE INDEX index_groups_members_on_user_id ON groups_members (user_id);
CREATE INDEX index_groups_members_on_inviter_id ON groups_members (inviter_id);


CREATE TABLE groups_invitations (
    id UUID NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL,

    ephemeral_public_key BYTEA NOT NULL,
    signature BYTEA NOT NULL,
    encrypted_master_key BYTEA NOT NULL,

    group_id UUID NOT NULL REFERENCES groups(id) ON DELETE CASCADE,
    invitee_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    inviter_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,

    PRIMARY KEY(id)
);
CREATE INDEX index_groups_invitations_on_group_id ON groups_invitations (group_id);
CREATE INDEX index_groups_invitations_on_invitee_id ON groups_invitations (invitee_id);
CREATE INDEX index_groups_invitations_on_inviter_id ON groups_invitations (inviter_id);


-- ################################################################################################@
-- Billing
-- ################################################################################################@
CREATE TABLE billing_plans (
    id UUID NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL,

    name TEXT NOT NULL,
    price BIGINT NOT NULL,
    description TEXT NOT NULL,
    stripe_id TEXT NOT NULL,
    product TEXT NOT NULL,
    storage BIGINT NOT NULL,

    PRIMARY KEY(id)
);
CREATE UNIQUE INDEX index_billing_plans_on_stripe_id ON billing_plans (stripe_id);
CREATE UNIQUE INDEX index_billing_plans_on_product ON billing_plans (product);


-- INSERT INTO billing_plans
-- 		(id, created_at, updated_at, name, description, stripe_id, price, product, storage)
-- 		VALUES ('42fb1c42-caca-418d-81f3-a6313c4a0a42', '2020-02-10T14:33:39+00:00', '2020-02-10T14:33:39+00:00',
--     'Free', '', 'plan_Gck0Zy9Qx5qaGk', 0, 'FREE', 100000000); -- 100MB


CREATE TABLE billing_customers (
    id UUID NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL,

    email TEXT NOT NULL,
    stripe_customer_id TEXT,
    stripe_subscription_id TEXT,
    used_storage BIGINT NOT NULL,
    subscription_updated_at TIMESTAMP WITH TIME ZONE NOT NULL,

    plan_id UUID NOT NULL REFERENCES billing_plans(id),
    user_id UUID REFERENCES users(id),
    group_id UUID REFERENCES groups(id),

    PRIMARY KEY(id)
);
CREATE INDEX index_billing_customers_on_stripe_customer_id ON billing_customers (stripe_customer_id);
CREATE INDEX index_billing_customers_on_user_id ON billing_customers (user_id);
CREATE INDEX index_billing_customers_on_group_id ON billing_customers (group_id);
CREATE INDEX index_billing_customers_on_plan_id ON billing_customers (plan_id);


CREATE TABLE billing_payment_methods (
    id UUID NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL,

    is_default BOOLEAN NOT NULL,
    stripe_id TEXT NOT NULL,
    card_last_4 TEXT NOT NULL,
    card_expiration_month BIGINT NOT NULL,
    card_expiration_year BIGINT NOT NULL,

    customer_id UUID NOT NULL REFERENCES billing_customers(id) ON DELETE CASCADE,

    PRIMARY KEY(id)
);
CREATE INDEX index_billing_payment_methods_on_customer_id ON billing_payment_methods (customer_id);


CREATE TABLE billing_invoices (
    id UUID NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL,

    amount BIGINT NOT NULL,
    stripe_id TEXT NOT NULL UNIQUE,
    stripe_hosted_url TEXT NOT NULL,
    stripe_pdf_url TEXT NOT NULL,
    paid_at TIMESTAMP WITH TIME ZONE NOT NULL,

    customer_id UUID NOT NULL REFERENCES billing_customers(id),

    PRIMARY KEY(id)
);
CREATE INDEX index_billing_invoices_on_customer_id ON billing_invoices (customer_id);


-- ################################################################################################@
-- Objects
-- ################################################################################################@
CREATE TABLE objects (
    id BYTEA NOT NULL,

    updated_at_state BIGINT NOT NULL,
    algorithm TEXT,
    encrypted_data BYTEA,
    encrypted_key BYTEA,
    nonce BYTEA,

    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    group_id UUID REFERENCES groups(id) ON DELETE CASCADE,

    PRIMARY KEY(id)
);
CREATE INDEX index_objects_on_user_id ON objects (user_id);
CREATE INDEX index_objects_on_group_id ON objects (group_id);
CREATE INDEX index_objects_on_updated_at_state ON objects (updated_at_state);
