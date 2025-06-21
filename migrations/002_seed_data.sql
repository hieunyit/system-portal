-- Insert default groups
INSERT INTO groups (name, display_name, description)
VALUES
    ('admin', 'Administrator', 'Full access to all features'),
    ('support', 'Support Staff', 'Limited OpenVPN user management')
ON CONFLICT (name) DO NOTHING;

-- Insert default permissions
INSERT INTO permissions (resource, action, description) VALUES
    ('portal', 'manage_users', 'Manage portal users'),
    ('portal', 'view_users', 'View portal users'),
    ('openvpn', 'manage_users', 'Manage OpenVPN users'),
    ('openvpn', 'view_status', 'View OpenVPN status')
ON CONFLICT (resource, action) DO NOTHING;

-- Assign permissions to groups
INSERT INTO group_permissions (group_id, permission_id)
SELECT g.id, p.id FROM groups g, permissions p WHERE g.name = 'admin'
ON CONFLICT DO NOTHING;

-- Support group permissions (limited)
INSERT INTO group_permissions (group_id, permission_id)
SELECT g.id, p.id FROM groups g
JOIN permissions p ON (p.resource, p.action) IN (( 'openvpn', 'manage_users' ), ( 'openvpn', 'view_status' ))
WHERE g.name = 'support'
ON CONFLICT DO NOTHING;

-- Create initial admin user
INSERT INTO users (username, email, password_hash, full_name, group_id)
VALUES (
    'admin',
    'admin@company.com',
    '$2a$14$8K1p/a0dL2LkzCKXNP7rVufDhZLCYLWJwONWtdVBXvhX7nVHsP.5K',
    'System Administrator',
    (SELECT id FROM groups WHERE name = 'admin')
)
ON CONFLICT (username) DO NOTHING;

-- Create initial support user
INSERT INTO users (username, email, password_hash, full_name, group_id)
VALUES (
    'support',
    'support@company.com',
    '$2a$14$8K1p/a0dL2LkzCKXNP7rVufDhZLCYLWJwONWtdVBXvhX7nVHsP.5K',
    'Support Staff',
    (SELECT id FROM groups WHERE name = 'support')
)
ON CONFLICT (username) DO NOTHING;
