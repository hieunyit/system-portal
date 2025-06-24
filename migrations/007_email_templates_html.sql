-- Use HTML for all email templates
UPDATE email_templates
SET body='<p>Dear {{.Username}},</p><p>Access control: {{.AccessControl}}<br>MAC address: {{.MACAddress}}<br>Expiration: {{.Expiration}}<br>Password: {{.Password}}</p><p>Đây là email tự động.</p>'
WHERE action='create_user_local';

UPDATE email_templates
SET body='<p>Dear {{.Username}},</p><p>Access control: {{.AccessControl}}<br>MAC address: {{.MACAddress}}<br>Expiration: {{.Expiration}}</p><p>Đây là email tự động.</p>'
WHERE action='create_user_ldap';

UPDATE email_templates
SET body='<p>Your account {{.Username}} has been enabled.</p>'
WHERE action='enable_user';

UPDATE email_templates
SET body='<p>Your account {{.Username}} has been disabled.</p>'
WHERE action='disable_user';

UPDATE email_templates
SET body='<p>Your account will expire in {{.Days}} day(s). Please renew.</p>'
WHERE action='expiration';

UPDATE email_templates
SET body='<p>Your one-time password for {{.Username}} has been reset.</p>'
WHERE action='reset_otp';

UPDATE email_templates
SET body='<p>Your password for {{.Username}} has been changed.</p>'
WHERE action='change_password';
