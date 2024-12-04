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
    username character varying(80) NOT NULL
);


ALTER TABLE public.usernames OWNER TO <username>;

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

ALTER SEQUENCE public.users_id_seq OWNED BY public.usernames.id;


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
-- Name: usercontacts id; Type: DEFAULT; Schema: public; Owner: <username>
--

ALTER TABLE ONLY public.usercontacts ALTER COLUMN id SET DEFAULT nextval('public.usercontacts_id_seq'::regclass);


--
-- Name: usernames id; Type: DEFAULT; Schema: public; Owner: <username>
--

ALTER TABLE ONLY public.usernames ALTER COLUMN id SET DEFAULT nextval('public.users_id_seq'::regclass);


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
-- Name: usercontacts_id_seq; Type: SEQUENCE SET; Schema: public; Owner: <username>
--

SELECT pg_catalog.setval('public.usercontacts_id_seq', 1, false);


--
-- Name: users_id_seq; Type: SEQUENCE SET; Schema: public; Owner: <username>
--

SELECT pg_catalog.setval('public.users_id_seq', 2, true);


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
-- Name: usercontacts usercontacts_pkey; Type: CONSTRAINT; Schema: public; Owner: <username>
--

ALTER TABLE ONLY public.usercontacts
    ADD CONSTRAINT usercontacts_pkey PRIMARY KEY (id);


--
-- Name: usernames users_pkey; Type: CONSTRAINT; Schema: public; Owner: <username>
--

ALTER TABLE ONLY public.usernames
    ADD CONSTRAINT users_pkey PRIMARY KEY (id);


--
-- Name: usernames users_username_key; Type: CONSTRAINT; Schema: public; Owner: <username>
--

ALTER TABLE ONLY public.usernames
    ADD CONSTRAINT users_username_key UNIQUE (username);


--
-- Name: passwords fk_user; Type: FK CONSTRAINT; Schema: public; Owner: <username>
--

ALTER TABLE ONLY public.passwords
    ADD CONSTRAINT fk_user FOREIGN KEY (user_id) REFERENCES public.usernames(id) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: invoices invoices_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: <username>
--

ALTER TABLE ONLY public.invoices
    ADD CONSTRAINT invoices_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.usernames(id) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: usercontacts usercontacts_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: <username>
--

ALTER TABLE ONLY public.usercontacts
    ADD CONSTRAINT usercontacts_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.usernames(id) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- PostgreSQL database dump complete
--

