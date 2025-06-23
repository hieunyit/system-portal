-- Insert default groups
INSERT INTO groups (name, display_name, description)
VALUES
    ('admin', 'Administrator', 'Full access to all features'),
    ('support', 'Support Staff', 'Limited OpenVPN user management')
ON CONFLICT (name) DO NOTHING;

-- Insert default permissions
INSERT INTO permissions (resource, action, description) VALUES
    -- Portal permissions
    ('portal', 'manage_users', 'Manage portal users'),
    ('portal', 'view_users', 'View portal users'),
    -- OpenVPN user permissions
    ('openvpn', 'view_users', 'View OpenVPN users'),
    ('openvpn', 'create_users', 'Create OpenVPN users'),
    ('openvpn', 'edit_users', 'Edit OpenVPN users'),
    ('openvpn', 'delete_users', 'Delete OpenVPN users'),
    -- OpenVPN group permissions
    ('openvpn', 'view_groups', 'View OpenVPN groups'),
    ('openvpn', 'manage_groups', 'Manage OpenVPN groups'),
    -- OpenVPN status permissions
    ('openvpn', 'view_status', 'View OpenVPN status')
ON CONFLICT (resource, action) DO NOTHING;

-- Assign permissions to groups
-- Admin gets all permissions
INSERT INTO group_permissions (group_id, permission_id)
SELECT g.id, p.id FROM groups g, permissions p
WHERE g.name = 'admin'
ON CONFLICT DO NOTHING;

-- Support group permissions (limited)
-- Support group has limited permissions
INSERT INTO group_permissions (group_id, permission_id)
SELECT g.id, p.id FROM groups g
JOIN permissions p ON (p.resource, p.action) IN (
    ( 'openvpn', 'view_users' ),
    ( 'openvpn', 'create_users' ),
    ( 'openvpn', 'edit_users' ),
    ( 'openvpn', 'view_groups' ),
    ( 'openvpn', 'view_status' )
)
WHERE g.name = 'support'
ON CONFLICT DO NOTHING;

-- Create initial admin user with low-cost bcrypt hash
INSERT INTO users (username, email, password_hash, full_name, group_id)
VALUES (
    'admin',
    'admin@company.com',
    crypt('admin123', gen_salt('bf', 10)),
    'System Administrator',
    (SELECT id FROM groups WHERE name = 'admin')
)
ON CONFLICT (username) DO NOTHING;

-- Create initial support user with low-cost bcrypt hash
INSERT INTO users (username, email, password_hash, full_name, group_id)
VALUES (
    'support',
    'support@company.com',
    crypt('support123', gen_salt('bf', 10)),
    'Support Staff',
    (SELECT id FROM groups WHERE name = 'support')
)
ON CONFLICT (username) DO NOTHING;
