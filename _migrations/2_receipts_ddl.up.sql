CREATE TABLE receipts (
    id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    description varchar(255) NOT NULL,
    value numeric(10,2) NOT NULL,
    date date NOT NULL,
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now()
);