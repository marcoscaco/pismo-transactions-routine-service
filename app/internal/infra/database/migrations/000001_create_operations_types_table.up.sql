-- Provisionamento basico da tabela de operation_types
-- GENERATED ALWAYS AS IDENTITY serve para evitar que inserts manuais sejam feitos na tabela, garantindo que o ID seja sempre gerado pelo banco de dados.
CREATE TABLE IF NOT EXISTS operation_types (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    description VARCHAR(255) NOT NULL,
    amount_multiplier SMALLINT NOT NULL CHECK (amount_multiplier IN (-1, 1)),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Insere os tipos de operação padrão
INSERT INTO operation_types (description, amount_multiplier) VALUES
    ('Normal Purchase', -1),
    ('Purchase with installments', -1),
    ('Withdrawal', -1),
    ('Credit Voucher', 1) ON CONFLICT DO NOTHING;