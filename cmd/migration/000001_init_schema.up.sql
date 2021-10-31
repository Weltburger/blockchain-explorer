CREATE DATABASE IF NOT EXISTS blocks;

SET flatten_nested = 0;

CREATE TABLE IF NOT EXISTS blocks.block (
    Protocol        String,
    ChainID         String,
    Hash            String,
    Timestamp       DateTime,
    Header          Nested
    (
        Level                       UInt64,
        Proto                       Int64,
        Predecessor                 String,
        Timestamp                   String,
        ValidationPass              Int64,
        OperationsHash              String,
        Fitness                     Array(String),
        Context                     String,
        Priority                    Int64,
        ProofOfWorkNonce            String,
        LiquidityBakingEscapeVote   UInt8,
        Signature                   String
    ),
    Metadata        Nested
    (
        Protocol                    String,
        NextProtocol                String,
        TestChainStatus             Tuple(Status String),
        MaxOperationsTTL            UInt64,
        MaxOperationDataLength      UInt64,
        MaxBlockHeaderLength        UInt64,
        MaxOperationListLength      Array(Nested(MaxSize UInt64, MaxOp UInt64)),
        Baker                       String,
        LevelInfo                   Nested
        (
            Level                   UInt64,
            LevelPosition           UInt64,
            Cycle                   UInt64,
            CyclePosition           UInt64,
            ExpectedCommitment      UInt8
        ),
        VotingPeriodInfo            Nested
        (
            VotingPeriod            Nested
            (
                VIndex              UInt64,
                Kind                String,
                StartPosition       UInt64
            ),
            Position                UInt64,
            Remaining               UInt64
        ),
        NonceHash                   Nullable(String),
        ConsumedGas                 String,
        Deactivated                 Array(String),
        BalanceUpdates              Nested
        (
            Kind                    String,
            Contract                Nullable(String),
            Change                  String,
            Origin                  String,
            Category                Nullable(String),
            Delegate                Nullable(String),
            Cycle                   Nullable(UInt64)
        ),
        LiquidityBakingEscapeEma    UInt64,
        ImplicitOperationsResults   Nested
        (
            Kind                    String,
            Storage                 Nested
            (
                SInt                String,
                SBytes              String
            ),
            BalanceUpdates          Nested
            (
                Kind                String,
                Contract            String,
                Change              String,
                Origin              String
            ),
            ConsumedGas             String,
            ConsumedMilligas        String,
            StorageSize             String
        )
    ),
    Operations        String

) engine = MergeTree()
      PARTITION BY toYYYYMMDD(Timestamp)
      ORDER BY (Hash)

