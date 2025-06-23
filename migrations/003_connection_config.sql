-- OpenVPN connection configuration
CREATE TABLE IF NOT EXISTS openvpn_configs (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    host VARCHAR(100) NOT NULL,
    username VARCHAR(100) NOT NULL,
    password TEXT NOT NULL,
    port INT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- LDAP connection configuration
CREATE TABLE IF NOT EXISTS ldap_configs (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    host VARCHAR(100) NOT NULL,
    port INT NOT NULL,
    bind_dn VARCHAR(100) NOT NULL,
    bind_password TEXT NOT NULL,
    base_dn VARCHAR(100) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);
