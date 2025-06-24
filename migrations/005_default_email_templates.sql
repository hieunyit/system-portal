-- Insert default email templates
INSERT INTO email_templates (action, subject, body)
VALUES
    ('create_user_local', 'OpenVPN account {{.Username}} created', 'Your account {{.Username}} has been created. Password: {{.Password}}'),
    ('create_user_ldap', 'OpenVPN account {{.Username}} created', 'Your LDAP account {{.Username}} is now active.'),
    ('enable_user', 'OpenVPN account {{.Username}} enabled', 'Your account {{.Username}} has been enabled.'),
    ('disable_user', 'OpenVPN account {{.Username}} disabled', 'Your account {{.Username}} has been disabled.'),
    ('expiration', 'OpenVPN account expiring in {{.Days}} day(s)', 'Your account will expire in {{.Days}} day(s). Please renew.'),
    ('reset_otp', 'OTP reset for {{.Username}}', 'Your one-time password for {{.Username}} has been reset.'),
    ('change_password', 'Password changed for {{.Username}}', 'Your password for {{.Username}} has been changed.')
ON CONFLICT (action) DO NOTHING;
