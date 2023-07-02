CREATE TABLE public.transactions (
  id SERIAL PRIMARY KEY,
  booking_date DATE NOT NULL,
  valuta_date DATE NOT NULL,
  recipient TEXT,
  booking_text TEXT NOT NULL,
  purpose TEXT,
  balance FLOAT NOT NULL,
  amount FLOAT NOT NULL,
  category_id INTEGER,
  hash TEXT NOT NULL
);
CREATE INDEX ON public.transactions(category_id);
CREATE INDEX ON public.transactions(booking_date);

CREATE TABLE public.categories (
  id SERIAL PRIMARY KEY,
  name TEXT NOT NULL,
  description TEXT,
  color TEXT
);

CREATE TABLE public.category_rules (
  id SERIAL PRIMARY KEY,
  category_id INTEGER NOT NULL,
  description TEXT,
  regex TEXT NOT NULL,
  mapping_field TEXT NOT NULL
);
CREATE INDEX ON public.category_rules(category_id);
