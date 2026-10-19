-- Cria a tavela de operações de transações
CREATE TABLE IF NOT EXISTS transactions_operations (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    account_id BIGINT NOT NULL,
    operation_type_id BIGINT NOT NULL,
    amount NUMERIC(15, 2) NOT NULL,
    event_date TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_account FOREIGN KEY (account_id) REFERENCES accounts(id),
    CONSTRAINT fk_operation_type FOREIGN KEY (operation_type_id) REFERENCES operation_types(id)
);

-- Cria um índice para melhorar a performance das consultas por conta
CREATE INDEX IF NOT EXISTS idx_transactions_account_id ON transactions_operations(account_id);