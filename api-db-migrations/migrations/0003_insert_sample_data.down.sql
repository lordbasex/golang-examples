-- Revertir la inserción de datos de ejemplo.
-- Se eliminan en orden inverso debido a las foreign keys.

-- Eliminar relaciones campaign_members
DELETE FROM campaign_members WHERE campaign_id IN (1, 2, 3, 4, 5);

-- Eliminar campañas
DELETE FROM campaigns WHERE id IN (1, 2, 3, 4, 5);

-- Eliminar usuarios
DELETE FROM users WHERE email IN (
  'juan.perez@example.com',
  'maria.gonzalez@example.com',
  'carlos.rodriguez@example.com',
  'ana.martinez@example.com',
  'luis.fernandez@example.com'
);

