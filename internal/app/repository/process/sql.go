package process

const (
	QueryDropTableArguments = `
-- QuerystmtDropTableArguments
DROP TABLE IF EXISTS arguments;
`

	QueryDropTableBanks = `
-- QuerystmtDropTableBanks
DROP TABLE IF EXISTS banks;
`
	QueryDropTableSystemTrx = `
-- QueryDropTableSystemTrx
DROP TABLE IF EXISTS system_trx;
`
	QueryDropTableBankTrx = `
-- QueryDropTableBankTrx
DROP TABLE IF EXISTS bank_trx;
`
	QueryDropTableReconciliationMap = `
-- QueryDropTableReconciliationMap
DROP TABLE IF EXISTS reconciliation_map;
`
	QueryCreateTableArguments = `
-- QueryCreateTableArguments
CREATE TABLE IF NOT EXISTS arguments AS
SELECT
    DATETIME(DATE(?)) AS start
    , DATETIME(DATE(?)) AS end
;
`
	QueryCreateTableBanks = `
-- QueryCreateTableBanks
CREATE TABLE IF NOT EXISTS banks AS
SELECT
    key AS id
    , LOWER(value) AS bank_name
FROM json_each(
    ?
)
;
`

	QueryCreateTableSystemTrx = `
-- QueryCreateTableSystemTrx
CREATE TABLE IF NOT EXISTS system_trx (
	TrxID TEXT PRIMARY KEY,
	Amount FLOAT,
	Type TEXT,
	TransactionTime DATETIME,
	FilePath TEXT
);

CREATE INDEX IF NOT EXISTS system_trx_Amount_index ON system_trx (Amount);
`
	QueryCreateTableBankTrx = `
-- QueryCreateTableBankTrx
CREATE TABLE IF NOT EXISTS bank_trx (
	UniqueIdentifier TEXT PRIMARY KEY,
	Amount FLOAT,
	Type TEXT,
	Bank TEXT,
	Date DATE,
	FilePath TEXT
);

CREATE INDEX IF NOT EXISTS bank_trx_Date_Type_Amount_UniqueIdentifier_index ON bank_trx (Date, Type, Amount, UniqueIdentifier);
`
	QueryCreateTableReconciliationMap = `
-- QueryCreateTableReconciliationMap
CREATE TABLE IF NOT EXISTS reconciliation_map (
	TrxID TEXT PRIMARY KEY,
	UniqueIdentifier TEXT
);

CREATE UNIQUE INDEX IF NOT EXISTS reconciliation_map_UniqueIdentifier_index ON reconciliation_map (UniqueIdentifier);
`
	QueryInsertTableSystemTrx = `
-- QueryInsertTableSystemTrx
INSERT INTO system_trx (TrxID, Amount, Type, TransactionTime, FilePath)
	SELECT
	json_extract(j.value, '$.TrxID') AS TrxID
	 , json_extract(j.value, '$.Amount') AS Amount
	 , json_extract(j.value, '$.Type') AS Type
	 , json_extract(j.value, '$.TransactionTime') AS TransactionTime
	 , json_extract(j.value, '$.FilePath') AS FilePath
	FROM json_each(
	 ?
	) AS j
;
`

	QueryInsertTableBankTrx = `
-- QueryInsertTableBankTrx
INSERT INTO bank_trx (UniqueIdentifier, Date, Type, FilePath, Bank, Amount)
	SELECT
	json_extract(j.value, '$.UniqueIdentifier') AS UniqueIdentifier
	 , json_extract(j.value, '$.Date') AS Date
	 , json_extract(j.value, '$.Type') AS Type
	 , json_extract(j.value, '$.FilePath') AS FilePath
	 , json_extract(j.value, '$.Bank') AS Bank
	 , json_extract(j.value, '$.Amount') AS Amount
	FROM json_each(
	 ?
	) AS j
;
`

	QueryInsertTableReconciliationMap = `
-- QueryInsertTableReconciliationMap
WITH main_data AS (
    SELECT
        CAST(? AS FLOAT) AS MinAmount
        , CAST(? AS FLOAT) AS MaxAmount
)
INSERT INTO reconciliation_map(
    TrxID,
    UniqueIdentifier
)
SELECT
    TrxID
     , UniqueIdentifier
FROM (
         SELECT
             TRUE
              , ROW_NUMBER() OVER (PARTITION BY st.TrxID ORDER BY bt.UniqueIdentifier) AS r_system
              , ROW_NUMBER() OVER (PARTITION BY bt.UniqueIdentifier ORDER BY st.TrxID) AS r_bank
              , st.TrxID
              , bt.UniqueIdentifier
         FROM main_data md
        INNER JOIN system_trx st ON st.Amount >= md.MinAmount AND st.Amount < md.MaxAmount
        INNER JOIN bank_trx bt ON
            bt.Date = STRFTIME('%FT%TZ', DATE(st.TransactionTime))
            AND bt.Type = st.Type
            AND bt.Amount = st.Amount
     )
WHERE r_system = r_bank;
`
	QueryGetReconciliationSummary = `
-- QueryGetReconciliationSummary
SELECT
    COALESCE(main_data.total_system_trx, 0) AS total_system_trx
    , COALESCE(main_data.total_matched_trx, 0) AS total_matched_trx
    , COALESCE((main_data.total_system_trx - main_data.total_matched_trx), 0) AS total_not_matched_trx
    , COALESCE(main_data.sum_system_trx, 0) AS sum_system_trx
    , COALESCE(main_data.sum_matched_trx, 0) AS sum_matched_trx
    , COALESCE((main_data.sum_system_trx - main_data.sum_matched_trx), 0) AS sum_discrepancies_trx
FROM (
    SELECT
        COUNT(*) AS total_system_trx
        , SUM(
            CASE
                WHEN rm.UniqueIdentifier IS NOT NULL then 1
                ELSE 0
            END
        ) AS total_matched_trx
        , SUM(st.Amount) AS sum_system_trx
        , SUM(
        CASE
            WHEN rm.UniqueIdentifier IS NOT NULL then st.Amount
            ELSE 0
        END
        ) AS sum_matched_trx
    FROM system_trx st
    LEFT JOIN reconciliation_map rm ON rm.TrxID = st.TrxID
) main_data
;
`

	QueryGetMatchedTrx = `
-- QueryGetMatchedTrx
SELECT
    st.TrxID AS SystemTrxTrxID,
    bt.UniqueIdentifier AS BankTrxUniqueIdentifier,
    STRFTIME('%F %T', st.TransactionTime) AS SystemTrxTransactionTime,
    DATE(bt.Date) AS BankTrxDate,
    st.Type AS SystemTrxType,
    st.Amount AS SystemTrxAmount,
    CASE
        WHEN bt.Type == 'DEBIT' THEN bt.Amount * (-1)
        ELSE bt.Amount
    END AS BankTrxAmount,
    bt.Bank
FROM reconciliation_map rm
INNER JOIN system_trx st on rm.TrxID = st.TrxID
INNER JOIN bank_trx bt on rm.UniqueIdentifier = bt.UniqueIdentifier
;
`

	QueryGetNotMatchedSystemTrx = `
-- QueryGetNotMatchedSystemTrx
SELECT st.TrxID                              AS TrxID,
       STRFTIME('%F %T', st.TransactionTime) AS TransactionTime,
       st.Type                               AS Type,
       st.Amount                             AS Amount
FROM system_trx st
LEFT JOIN reconciliation_map rm on rm.TrxID = st.TrxID
WHERE rm.TrxID IS NULL
;
`
	QueryGetNotMatchedBankTrx = `
-- QueryGetNotMatchedBankTrx
SELECT
    bt.UniqueIdentifier AS UniqueIdentifier,
    bt.Bank AS Bank,
    STRFTIME('%F', bt.Date) AS Date,
    CASE
        WHEN bt.Type == 'DEBIT' THEN bt.Amount * (-1)
        ELSE bt.Amount
    END AS Amount
FROM bank_trx bt
LEFT JOIN reconciliation_map rm on rm.UniqueIdentifier = bt.UniqueIdentifier
WHERE rm.UniqueIdentifier IS NULL
;
`
)
