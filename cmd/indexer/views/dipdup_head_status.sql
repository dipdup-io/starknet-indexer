CREATE OR REPLACE VIEW dipdup_head_status AS
SELECT
    name,
    CASE
        WHEN last_time < NOW() - interval '15 minutes' THEN 'OUTDATED'
        ELSE 'OK'
    END AS status,
    last_time,
    last_height
FROM
    state;