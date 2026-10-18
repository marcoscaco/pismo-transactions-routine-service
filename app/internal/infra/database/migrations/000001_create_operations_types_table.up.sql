-- Provisionamento basico da tabela de operation_types
-- GENERATED ALWAYS AS IDENTITY serve para evitar que inserts manuais sejam feitos na tabela, garantindo que o ID seja sempre gerado pelo banco de dados.
CREATE TABLE IF NOT EXISTS operation_types (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    description VARCHAR(255) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Insere os tipos de operação padrão
INSERT INTO operation_types (description) VALUES
    ('Normal Purchase'),
    ('Purchase with installments'),
    ('Withdrawal'),
    ('Credit Voucher') ON CONFLICT DO NOTHING;