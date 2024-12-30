--
-- PostgreSQL database dump
--

-- Dumped from database version 16.6 (Homebrew)
-- Dumped by pg_dump version 16.6 (Homebrew)

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: invoices; Type: TABLE; Schema: public; Owner: <username>
--

CREATE TABLE public.invoices (
    id integer NOT NULL,
    user_id integer NOT NULL,
    product character varying(80) NOT NULL,
    category character varying(80) NOT NULL,
    price numeric(5,2) NOT NULL,
    quantity integer NOT NULL
);


ALTER TABLE public.invoices OWNER TO <username>;

--
-- Name: invoices_id_seq; Type: SEQUENCE; Schema: public; Owner: <username>
--

CREATE SEQUENCE public.invoices_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.invoices_id_seq OWNER TO <username>;

--
-- Name: invoices_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: <username>
--

ALTER SEQUENCE public.invoices_id_seq OWNED BY public.invoices.id;


--
-- Name: passwords; Type: TABLE; Schema: public; Owner: <username>
--

CREATE TABLE public.passwords (
    id integer NOT NULL,
    user_id integer NOT NULL,
    password bytea NOT NULL
);


ALTER TABLE public.passwords OWNER TO <username>;

--
-- Name: passwords_id_seq; Type: SEQUENCE; Schema: public; Owner: <username>
--

CREATE SEQUENCE public.passwords_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.passwords_id_seq OWNER TO <username>;

--
-- Name: passwords_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: <username>
--

ALTER SEQUENCE public.passwords_id_seq OWNED BY public.passwords.id;


--
-- Name: tokens; Type: TABLE; Schema: public; Owner: <username>
--

CREATE TABLE public.tokens (
    id integer NOT NULL,
    user_id integer NOT NULL,
    token character varying(80)
);


ALTER TABLE public.tokens OWNER TO <username>;

--
-- Name: tokens_id_seq; Type: SEQUENCE; Schema: public; Owner: <username>
--

CREATE SEQUENCE public.tokens_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.tokens_id_seq OWNER TO <username>;

--
-- Name: tokens_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: <username>
--

ALTER SEQUENCE public.tokens_id_seq OWNED BY public.tokens.id;


--
-- Name: usercontacts; Type: TABLE; Schema: public; Owner: <username>
--

CREATE TABLE public.usercontacts (
    id integer NOT NULL,
    user_id integer NOT NULL,
    fname character varying(80) NOT NULL,
    lname character varying(80) NOT NULL,
    address character varying(80) NOT NULL
);


ALTER TABLE public.usercontacts OWNER TO <username>;

--
-- Name: usercontacts_id_seq; Type: SEQUENCE; Schema: public; Owner: <username>
--

CREATE SEQUENCE public.usercontacts_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.usercontacts_id_seq OWNER TO <username>;

--
-- Name: usercontacts_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: <username>
--

ALTER SEQUENCE public.usercontacts_id_seq OWNED BY public.usercontacts.id;


--
-- Name: usernames; Type: TABLE; Schema: public; Owner: <username>
--

CREATE TABLE public.usernames (
    id integer NOT NULL,
    username character varying(255) NOT NULL
);


ALTER TABLE public.usernames OWNER TO <username>;

--
-- Name: usernames_id_seq; Type: SEQUENCE; Schema: public; Owner: <username>
--

CREATE SEQUENCE public.usernames_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.usernames_id_seq OWNER TO <username>;

--
-- Name: usernames_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: <username>
--

ALTER SEQUENCE public.usernames_id_seq OWNED BY public.usernames.id;


--
-- Name: invoices id; Type: DEFAULT; Schema: public; Owner: <username>
--

ALTER TABLE ONLY public.invoices ALTER COLUMN id SET DEFAULT nextval('public.invoices_id_seq'::regclass);


--
-- Name: passwords id; Type: DEFAULT; Schema: public; Owner: <username>
--

ALTER TABLE ONLY public.passwords ALTER COLUMN id SET DEFAULT nextval('public.passwords_id_seq'::regclass);


--
-- Name: tokens id; Type: DEFAULT; Schema: public; Owner: <username>
--

ALTER TABLE ONLY public.tokens ALTER COLUMN id SET DEFAULT nextval('public.tokens_id_seq'::regclass);


--
-- Name: usercontacts id; Type: DEFAULT; Schema: public; Owner: <username>
--

ALTER TABLE ONLY public.usercontacts ALTER COLUMN id SET DEFAULT nextval('public.usercontacts_id_seq'::regclass);


--
-- Name: usernames id; Type: DEFAULT; Schema: public; Owner: <username>
--

ALTER TABLE ONLY public.usernames ALTER COLUMN id SET DEFAULT nextval('public.usernames_id_seq'::regclass);


--
-- Data for Name: invoices; Type: TABLE DATA; Schema: public; Owner: <username>
--

COPY public.invoices (id, user_id, product, category, price, quantity) FROM stdin;
\.


--
-- Data for Name: passwords; Type: TABLE DATA; Schema: public; Owner: <username>
--

COPY public.passwords (id, user_id, password) FROM stdin;
\.


--
-- Data for Name: tokens; Type: TABLE DATA; Schema: public; Owner: <username>
--

COPY public.tokens (id, user_id, token) FROM stdin;
\.


--
-- Data for Name: usercontacts; Type: TABLE DATA; Schema: public; Owner: <username>
--

COPY public.usercontacts (id, user_id, fname, lname, address) FROM stdin;
\.


--
-- Data for Name: usernames; Type: TABLE DATA; Schema: public; Owner: <username>
--

COPY public.usernames (id, username) FROM stdin;
\.


--
-- Name: invoices_id_seq; Type: SEQUENCE SET; Schema: public; Owner: <username>
--

SELECT pg_catalog.setval('public.invoices_id_seq', 1, false);


--
-- Name: passwords_id_seq; Type: SEQUENCE SET; Schema: public; Owner: <username>
--

SELECT pg_catalog.setval('public.passwords_id_seq', 1, false);


--
-- Name: tokens_id_seq; Type: SEQUENCE SET; Schema: public; Owner: <username>
--

SELECT pg_catalog.setval('public.tokens_id_seq', 1, false);


--
-- Name: usercontacts_id_seq; Type: SEQUENCE SET; Schema: public; Owner: <username>
--

SELECT pg_catalog.setval('public.usercontacts_id_seq', 1, false);


--
-- Name: usernames_id_seq; Type: SEQUENCE SET; Schema: public; Owner: <username>
--

SELECT pg_catalog.setval('public.usernames_id_seq', 1, false);


--
-- Name: invoices invoices_id_user_id_product_category_price_quantity_key; Type: CONSTRAINT; Schema: public; Owner: <username>
--

ALTER TABLE ONLY public.invoices
    ADD CONSTRAINT invoices_id_user_id_product_category_price_quantity_key UNIQUE (id, user_id, product, category, price, quantity);


--
-- Name: invoices invoices_pkey; Type: CONSTRAINT; Schema: public; Owner: <username>
--

ALTER TABLE ONLY public.invoices
    ADD CONSTRAINT invoices_pkey PRIMARY KEY (id);


--
-- Name: passwords passwords_password_key; Type: CONSTRAINT; Schema: public; Owner: <username>
--

ALTER TABLE ONLY public.passwords
    ADD CONSTRAINT passwords_password_key UNIQUE (password);


--
-- Name: passwords passwords_pkey; Type: CONSTRAINT; Schema: public; Owner: <username>
--

ALTER TABLE ONLY public.passwords
    ADD CONSTRAINT passwords_pkey PRIMARY KEY (id);


--
-- Name: tokens tokens_pkey; Type: CONSTRAINT; Schema: public; Owner: <username>
--

ALTER TABLE ONLY public.tokens
    ADD CONSTRAINT tokens_pkey PRIMARY KEY (id);


--
-- Name: tokens user_id_unique; Type: CONSTRAINT; Schema: public; Owner: <username>
--

ALTER TABLE ONLY public.tokens
    ADD CONSTRAINT user_id_unique UNIQUE (user_id);


--
-- Name: usercontacts usercontacts_fname_lname_address_key; Type: CONSTRAINT; Schema: public; Owner: <username>
--

ALTER TABLE ONLY public.usercontacts
    ADD CONSTRAINT usercontacts_fname_lname_address_key UNIQUE (fname, lname, address);


--
-- Name: usercontacts usercontacts_pkey; Type: CONSTRAINT; Schema: public; Owner: <username>
--

ALTER TABLE ONLY public.usercontacts
    ADD CONSTRAINT usercontacts_pkey PRIMARY KEY (id);


--
-- Name: usernames usernames_pkey; Type: CONSTRAINT; Schema: public; Owner: <username>
--

ALTER TABLE ONLY public.usernames
    ADD CONSTRAINT usernames_pkey PRIMARY KEY (id);


--
-- Name: usernames usernames_username_key; Type: CONSTRAINT; Schema: public; Owner: <username>
--

ALTER TABLE ONLY public.usernames
    ADD CONSTRAINT usernames_username_key UNIQUE (username);


--
-- Name: invoices invoices_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: <username>
--

ALTER TABLE ONLY public.invoices
    ADD CONSTRAINT invoices_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.usernames(id) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: passwords passwords_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: <username>
--

ALTER TABLE ONLY public.passwords
    ADD CONSTRAINT passwords_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.usernames(id) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: tokens tokens_fk; Type: FK CONSTRAINT; Schema: public; Owner: <username>
--

ALTER TABLE ONLY public.tokens
    ADD CONSTRAINT tokens_fk FOREIGN KEY (user_id) REFERENCES public.usernames(id) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: usercontacts usercontacts_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: <username>
--

ALTER TABLE ONLY public.usercontacts
    ADD CONSTRAINT usercontacts_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.usernames(id) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- PostgreSQL database dump complete
--

