-- Update create user email templates with detailed information
UPDATE email_templates
SET body='Dear {{.Username}},\nAccess control: {{.AccessControl}}\nMAC address: {{.MACAddress}}\nExpiration: {{.Expiration}}\nPassword: {{.Password}}\n\n\u0110\u00E2y l\u00E0 email t\u1EF1 \u0111\u1ED9ng.'
WHERE action='create_user_local';

UPDATE email_templates
SET body='Dear {{.Username}},\nAccess control: {{.AccessControl}}\nMAC address: {{.MACAddress}}\nExpiration: {{.Expiration}}\n\n\u0110\u00E2y l\u00E0 email t\u1EF1 \u0111\u1ED9ng.'
WHERE action='create_user_ldap';
