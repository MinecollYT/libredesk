-- name: get-teams
SELECT id, created_at, updated_at, name, uuid, disabled from teams;

-- name: get-team
SELECT id, name, uuid from teams where id = $1;

-- name: get-team-members
SELECT u.id, u.uuid, t.id as team_id
FROM users u
JOIN team_members tm ON tm.user_id = u.id
JOIN teams t ON t.id = tm.team_id
WHERE t.name = $1;

-- name: insert-team
INSERT INTO teams (name) values($1);

-- name: update-team
UPDATE teams set name = $2, updated_at = now() where id = $1;