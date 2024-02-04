CREATE VIEW user_role_view AS
SELECT u.*, r.role_id, r.role_name
FROM "user" AS u
LEFT JOIN "user_role" AS ur ON u.user_id = ur.user_id
LEFT JOIN "role" AS r ON ur.role_id = r.role_id;