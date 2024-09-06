-- user_resources source

CREATE VIEW user_resources AS
SELECT s.uuid as session_id, s.user_id, c.id  as campaign_id, a.id as adventure_id 
FROM sessions s 
LEFT JOIN campaigns c ON s.user_id = c.user_id 
LEFT JOIN adventures a ON c.id = a.campaign_id;