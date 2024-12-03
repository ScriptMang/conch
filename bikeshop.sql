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
-- Name: addresses; Type: TABLE; Schema: public; Owner: <username>
--

CREATE TABLE public.addresses (
    id integer NOT NULL,
    user_id integer NOT NULL,
    address character varying(255) NOT NULL
);


ALTER TABLE public.addresses OWNER TO <username>;

--
-- Name: addresses_id_seq; Type: SEQUENCE; Schema: public; Owner: <username>
--

CREATE SEQUENCE public.addresses_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.addresses_id_seq OWNER TO <username>;

--
-- Name: addresses_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: <username>
--

ALTER SEQUENCE public.addresses_id_seq OWNED BY public.addresses.id;


--
-- Name: invoices; Type: TABLE; Schema: public; Owner: <username>
--

CREATE TABLE public.invoices (
    id integer NOT NULL,
    user_id integer NOT NULL,
    product character varying(255) NOT NULL,
    category character varying(255) NOT NULL,
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
    password bytea
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
-- Name: passwords_user_id_seq; Type: SEQUENCE; Schema: public; Owner: <username>
--

CREATE SEQUENCE public.passwords_user_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.passwords_user_id_seq OWNER TO <username>;

--
-- Name: passwords_user_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: <username>
--

ALTER SEQUENCE public.passwords_user_id_seq OWNED BY public.passwords.user_id;


--
-- Name: users; Type: TABLE; Schema: public; Owner: <username>
--

CREATE TABLE public.users (
    id integer NOT NULL,
    username character varying(80) NOT NULL,
    fname character varying(80) NOT NULL,
    lname character varying(80) NOT NULL,
    addr_id integer NOT NULL
);


ALTER TABLE public.users OWNER TO <username>;

--
-- Name: users_id_seq; Type: SEQUENCE; Schema: public; Owner: <username>
--

CREATE SEQUENCE public.users_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.users_id_seq OWNER TO <username>;

--
-- Name: users_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: <username>
--

ALTER SEQUENCE public.users_id_seq OWNED BY public.users.id;


--
-- Name: addresses id; Type: DEFAULT; Schema: public; Owner: <username>
--

ALTER TABLE ONLY public.addresses ALTER COLUMN id SET DEFAULT nextval('public.addresses_id_seq'::regclass);


--
-- Name: invoices id; Type: DEFAULT; Schema: public; Owner: <username>
--

ALTER TABLE ONLY public.invoices ALTER COLUMN id SET DEFAULT nextval('public.invoices_id_seq'::regclass);


--
-- Name: passwords id; Type: DEFAULT; Schema: public; Owner: <username>
--

ALTER TABLE ONLY public.passwords ALTER COLUMN id SET DEFAULT nextval('public.passwords_id_seq'::regclass);


--
-- Name: passwords user_id; Type: DEFAULT; Schema: public; Owner: <username>
--

ALTER TABLE ONLY public.passwords ALTER COLUMN user_id SET DEFAULT nextval('public.passwords_user_id_seq'::regclass);


--
-- Name: users id; Type: DEFAULT; Schema: public; Owner: <username>
--

ALTER TABLE ONLY public.users ALTER COLUMN id SET DEFAULT nextval('public.users_id_seq'::regclass);


--
-- Data for Name: addresses; Type: TABLE DATA; Schema: public; Owner: <username>
--

COPY public.addresses (id, user_id, address) FROM stdin;
\.


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
-- Data for Name: users; Type: TABLE DATA; Schema: public; Owner: <username>
--

COPY public.users (id, username, fname, lname, addr_id) FROM stdin;
\.


--
-- Name: addresses_id_seq; Type: SEQUENCE SET; Schema: public; Owner: <username>
--

SELECT pg_catalog.setval('public.addresses_id_seq', 1, false);


--
-- Name: invoices_id_seq; Type: SEQUENCE SET; Schema: public; Owner: <username>
--

SELECT pg_catalog.setval('public.invoices_id_seq', 3, true);


--
-- Name: passwords_id_seq; Type: SEQUENCE SET; Schema: public; Owner: <username>
--

SELECT pg_catalog.setval('public.passwords_id_seq', 2, true);


--
-- Name: passwords_user_id_seq; Type: SEQUENCE SET; Schema: public; Owner: <username>
--

SELECT pg_catalog.setval('public.passwords_user_id_seq', 2, true);


--
-- Name: users_id_seq; Type: SEQUENCE SET; Schema: public; Owner: <username>
--

SELECT pg_catalog.setval('public.users_id_seq', 2, true);


--
-- Name: addresses addresses_pkey; Type: CONSTRAINT; Schema: public; Owner: <username>
--

ALTER TABLE ONLY public.addresses
    ADD CONSTRAINT addresses_pkey PRIMARY KEY (id);


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
-- Name: users users_pkey; Type: CONSTRAINT; Schema: public; Owner: <username>
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_pkey PRIMARY KEY (id);


--
-- Name: users users_username_key; Type: CONSTRAINT; Schema: public; Owner: <username>
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_username_key UNIQUE (username);


--
-- Name: addresses addresses_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: <username>
--

ALTER TABLE ONLY public.addresses
    ADD CONSTRAINT addresses_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: passwords fk_user; Type: FK CONSTRAINT; Schema: public; Owner: <username>
--

ALTER TABLE ONLY public.passwords
    ADD CONSTRAINT fk_user FOREIGN KEY (user_id) REFERENCES public.users(id) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: invoices invoices_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: <username>
--

ALTER TABLE ONLY public.invoices
    ADD CONSTRAINT invoices_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: users users_addr_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: <username>
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_addr_id_fkey FOREIGN KEY (addr_id) REFERENCES public.addresses(id);


--
-- PostgreSQL database dump complete
--

