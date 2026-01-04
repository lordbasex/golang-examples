-- Insertar datos de ejemplo en las tablas existentes.
-- Esta migración añade usuarios, campañas y sus relaciones.

-- Insertar usuarios de ejemplo
INSERT INTO users (email, full_name, phone) VALUES
  ('juan.perez@example.com', 'Juan Pérez', '+54-11-4555-1234'),
  ('maria.gonzalez@example.com', 'María González', '+54-11-4555-5678'),
  ('carlos.rodriguez@example.com', 'Carlos Rodríguez', '+54-11-4555-9012'),
  ('ana.martinez@example.com', 'Ana Martínez', NULL),
  ('luis.fernandez@example.com', 'Luis Fernández', '+54-11-4555-3456');

-- Insertar campañas de ejemplo
INSERT INTO campaigns (name, status) VALUES
  ('Campaña de Lanzamiento 2026', 'active'),
  ('Promoción de Verano', 'active'),
  ('Black Friday 2026', 'draft'),
  ('Newsletter Mensual', 'active'),
  ('Webinar Técnico', 'completed');

-- Insertar relaciones campaign_members
-- Campaña 1: Juan (owner), María (admin), Carlos (member)
INSERT INTO campaign_members (campaign_id, user_id, role) VALUES
  (1, 1, 'owner'),
  (1, 2, 'admin'),
  (1, 3, 'member');

-- Campaña 2: María (owner), Ana (admin)
INSERT INTO campaign_members (campaign_id, user_id, role) VALUES
  (2, 2, 'owner'),
  (2, 4, 'admin');

-- Campaña 3: Carlos (owner), Luis (member)
INSERT INTO campaign_members (campaign_id, user_id, role) VALUES
  (3, 3, 'owner'),
  (3, 5, 'member');

-- Campaña 4: Ana (owner), Juan (member), María (member)
INSERT INTO campaign_members (campaign_id, user_id, role) VALUES
  (4, 4, 'owner'),
  (4, 1, 'member'),
  (4, 2, 'member');

-- Campaña 5: Luis (owner), Carlos (admin), Ana (member)
INSERT INTO campaign_members (campaign_id, user_id, role) VALUES
  (5, 5, 'owner'),
  (5, 3, 'admin'),
  (5, 4, 'member');

