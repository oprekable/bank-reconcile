package sample

const (
	QueryDropTableArguments = `
-- QueryDropTableArguments
DROP TABLE IF EXISTS arguments;
`

	QueryDropTableBanks = `
-- QueryDropTableBanks
DROP TABLE IF EXISTS banks;
`

	QueryDropTableBaseData = `
-- QueryDropTableBaseData
DROP TABLE IF EXISTS base_data;
`

	QueryCreateTableArguments = `
-- QueryCreateTableArguments
CREATE TABLE IF NOT EXISTS arguments AS
SELECT
    DATETIME(DATE(?)) AS start
    , DATETIME(DATE(?)) AS end
    , ABS(CAST(? AS INTEGER)) AS limit_trx_data
    , ABS(CAST(? AS INTEGER)) AS match_percentage
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
	QueryCreateTableBaseData = `
-- QueryCreateTableBaseData
CREATE TABLE IF NOT EXISTS base_data AS
WITH RECURSIVE generate_series(no, transactionTime, trxID, amount, type, bank_array, bank, count_bank, limit_data) AS (
    SELECT
        1,
        datetime(
                abs(random()) % (strftime('%s', DATETIME(DATETIME(a.end, '+1 day'), '-1 second')) - strftime('%s', a.start)) + strftime('%s', a.start),
                'unixepoch'
        ),
        lower(hex(randomblob(16))),
        CAST(
                FLOOR((ABS(random()) % (100000 - 1000) + 1000) / 100) * 100 AS FLOAT
        ),
        CASE
            WHEN ABS(random()) % 2 = 1 THEN 'DEBIT'
            ELSE 'CREDIT'
        END,
        (SELECT JSON_GROUP_ARRAY(b.bank_name) FROM banks b),
        JSON_EXTRACT(
                (SELECT JSON_GROUP_ARRAY(b.bank_name) FROM banks b),
                '$[' || cast(ABS(RANDOM()) % (SELECT COUNT(*) FROM banks) as text) || ']'
        ),
        (SELECT COUNT(*) FROM banks),
        (CASE
            WHEN a.match_percentage == 100 THEN a.limit_trx_data
            ELSE (a.limit_trx_data * 2)
        END)
    FROM arguments a
    UNION ALL
    SELECT
        no+1,
        datetime(
                abs(random()) % (strftime('%s', DATETIME(DATETIME(a.end, '+1 day'), '-1 second')) - strftime('%s', a.start)) + strftime('%s', a.start),
                'unixepoch'
        ),
        lower(hex(randomblob(16))),
        CAST(
                FLOOR((ABS(random()) % (100000 - 1000) + 1000) / 100) * 100 AS FLOAT
        ),
        CASE
            WHEN ABS(random()) % 2 = 1 THEN 'DEBIT'
            ELSE 'CREDIT'
        END,
        bank_array,
        JSON_EXTRACT(
                (bank_array),
                '$[' || cast(ABS(RANDOM()) % count_bank AS text) || ']'
        ),
        count_bank,
        limit_data
    FROM generate_series, arguments a
    WHERE no+1 <= limit_data
)
SELECT
    gs.trxID,
    gs.transactionTime,
    gs.amount,
    gs.type,
    gs.bank
FROM generate_series gs
;
`
	QueryCreateIndexTableBaseData = `
-- QueryCreateIndexTableBaseData
CREATE INDEX IF NOT EXISTS idx_base_data_trx_id ON base_data(trxID);
`

	QueryGetTrxData = `
-- QueryGetTrxData
WITH with_system_trx AS (
    SELECT
        *
        , ROW_NUMBER() OVER win    AS row_number
        , FLOOR((ROW_NUMBER() OVER win/ CAST(COUNT() OVER counter AS REAL)) * 100) AS percentile
    FROM (
         SELECT
             bd1.trxID
              , system_trx.trxID IS NOT NULL AS is_system_trx
         FROM base_data bd1
         LEFT JOIN (
             SELECT
                 bd2.trxID
             FROM base_data bd2
             ORDER BY random()
             LIMIT (SELECT limit_trx_data FROM arguments)
         ) system_trx ON bd1.trxID = system_trx.trxID
    ) tmp_data
    WINDOW
        win AS (ORDER BY tmp_data.is_system_trx DESC)
        , counter AS (PARTITION BY tmp_data.is_system_trx)
    ORDER BY tmp_data.is_system_trx DESC
)
, with_bank_trx_flagged AS (
    SELECT
        wst.trxID
         , wst.is_system_trx
         , wst.row_number
         , CASE
               WHEN a.match_percentage = 100 THEN wst.is_system_trx
               WHEN a.match_percentage < 100 AND a.match_percentage >= 0 THEN
                   CASE
                       WHEN wst.percentile <= a.match_percentage THEN TRUE
                       WHEN wst.percentile > a.match_percentage THEN FALSE
                       ELSE TRUE END
               ELSE FALSE
        END AS is_bank_trx
    FROM with_system_trx wst, arguments a
)
, final AS (
	SELECT
		wbtf.trxID
		, lower(bd.bank) || '-' || lower(hex(randomblob(16))) AS uniqueIdentifier
		, wbtf.is_system_trx
		, wbtf.is_bank_trx
		, bd.type
		, bd.bank
		, bd.amount
		, bd.transactionTime
		, DATE(bd.transactionTime)  AS date
		, wbtf.row_number
		, COUNT(CASE WHEN wbtf.is_bank_trx THEN 1 END) OVER () AS max_row_number_is_bank_trx
		, COUNT(CASE WHEN wbtf.is_system_trx THEN 1 END) OVER () AS max_row_number_is_system_trx
		, a.limit_trx_data
	FROM with_bank_trx_flagged wbtf
	INNER JOIN base_data bd ON bd.trxID = wbtf.trxID
	LEFT JOIN arguments a ON TRUE
)
SELECT
    final.trxID
    , final.uniqueIdentifier
    , final.is_system_trx
    , final.is_bank_trx
    , final.type
    , final.bank
    , final.amount
    , final.transactionTime
    , final.date
FROM final
WHERE final.row_number <= (
    CASE
       WHEN ((final.max_row_number_is_system_trx * 2) - final.max_row_number_is_bank_trx) < (final.limit_trx_data * 2)
           THEN ((final.max_row_number_is_system_trx * 2) - final.max_row_number_is_bank_trx)
       ELSE final.limit_trx_data
    END
)
ORDER BY final.row_number
;
`
)
