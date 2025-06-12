SELECT 
    id, 
    title, 
    scheduled_at,
    is_active,
    datetime('now') as now_utc,
    datetime('now', 'localtime') as now_local,
    datetime(scheduled_at, 'localtime') as scheduled_local
FROM questions;
