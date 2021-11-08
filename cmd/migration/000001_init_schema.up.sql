CREATE DATABASE IF NOT EXISTS blocks;

CREATE TABLE IF NOT EXISTS blocks.block (
    `protocol`         String,
    `chain_id`         String,
    `hash`             String,
    `baker_fees`       UInt64,
    `level`            UInt64,
    `predecessor`      String,
    `priority`         Int64,
    `timestamp`        DateTime,
    `validation_pass`  UInt64,
    `validation_hash`  String,
    `fitness`          String,
    `signature`        String,
    `baker`            String,
    `cycle_num`        UInt64,
    `cycle_position`   UInt64,
    `consumed_gas`     String
) ENGINE = MergeTree()
  PARTITION BY toYYYYMMDD(timestamp)
  ORDER BY (level, hash);

CREATE TABLE IF NOT EXISTS blocks.block (
    Protocol        String,
    ChainID         String,
    Hash            String,
    Timestamp       DateTime,
    Header          String,
    Metadata        String,
    Operations      String
) engine = MergeTree()
      PARTITION BY toYYYYMMDD(Timestamp)
      ORDER BY (Hash)

