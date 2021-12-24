CREATE TABLE IF NOT EXISTS blocks.transactions_mi (
    `block_hash`        String,
    `hash`              String,
    `source`            String,
    `destination`       String,
    `fee`               String,
    `amount`            String
) ENGINE = MergeTree()
      ORDER BY (block_hash, hash);