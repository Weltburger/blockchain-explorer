CREATE TABLE IF NOT EXISTS blocks.transactions (
    `block_hash` String,
    `hash` String,
    `branch` String,
    `source` String,
    `destination` String,
    `fee` String,
    `counter` String,
    `gas_limit` String,
    `amount` String,
    `consumed_milligas` String,
    `storage_size` String,
    `signature` String
) ENGINE = MergeTree()
      ORDER BY (branch, hash);