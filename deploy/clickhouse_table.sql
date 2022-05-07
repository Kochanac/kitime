
CREATE TABLE user_video_times (
    `user_id` UInt32,
    `event_time` DateTime,
    `event_type` Enum8('watch' = 0, 'scroll' = 1),
    `video_id` UInt32,
    `video_timestamp` UInt32
)
ENGINE = MergeTree
ORDER BY event_time;


CREATE TABLE user_video_times_queue (
    `user_id` UInt32,
    `event_time` DateTime,
    `event_type` Enum8('watch' = 0, 'scroll' = 1),
    `video_id` UInt32,
    `video_timestamp` UInt32
)
ENGINE = Kafka('kafka-cluster-kafka-brokers.infra:9092', 'vobla-topic', 'clickhouse_consumer',
    'JSONEachRow') settings kafka_thread_per_consumer = 0, kafka_num_consumers = 1;

CREATE MATERIALIZED VIEW user_video_times_view TO user_video_times AS
    SELECT * from user_video_times_queue;
