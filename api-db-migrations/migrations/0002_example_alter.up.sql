-- Ejemplo: un cambio incremental (ALTER) versionado.
-- En tu caso, cada versión nueva del servicio agrega su migración.

ALTER TABLE users
  ADD COLUMN phone VARCHAR(32) NULL AFTER full_name;
