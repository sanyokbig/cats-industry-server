INSERT INTO groups
(id, name) VALUES
  (1, 'default'),
  (2, 'industrialist'),
  (3, 'super')
ON CONFLICT DO NOTHING;

INSERT INTO roles
(id, name) VALUES
  (1, 'see_shared_jobs'),
  (2, 'share_jobs'),
  (3, 'add_industrial_character'),
  (4, 'edit_own_jobs'),
  (5, 'edit_all_jobs')
ON CONFLICT DO NOTHING;

INSERT INTO groups_roles
(group_id, role_id) VALUES
  (1, 1),
  (2, 1),
  (2, 2),
  (2, 3),
  (2, 4),
  (2, 5)
ON CONFLICT DO NOTHING;