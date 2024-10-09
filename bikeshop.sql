--
-- PostgreSQL database dump
--

-- Dumped from database version 16.4 (Homebrew)
-- Dumped by pg_dump version 16.4 (Homebrew)

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
    usr_id integer,
    fname character varying(50) NOT NULL,
    lname character varying(50) NOT NULL,
    product character varying(255) NOT NULL,
    price numeric(5,2) NOT NULL,
    quantity integer NOT NULL,
    category character varying(255) NOT NULL,
    shipping character varying(255) NOT NULL
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
-- Name: users; Type: TABLE; Schema: public; Owner: <username>
--

CREATE TABLE public.users (
    id integer NOT NULL,
    username character varying(80),
    password character varying(255)
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
-- Name: invoices id; Type: DEFAULT; Schema: public; Owner: <username>
--

ALTER TABLE ONLY public.invoices ALTER COLUMN id SET DEFAULT nextval('public.invoices_id_seq'::regclass);


--
-- Name: users id; Type: DEFAULT; Schema: public; Owner: <username>
--

ALTER TABLE ONLY public.users ALTER COLUMN id SET DEFAULT nextval('public.users_id_seq'::regclass);


--
-- Data for Name: invoices; Type: TABLE DATA; Schema: public; Owner: <username>
--

COPY public.invoices (id, usr_id, fname, lname, product, price, quantity, category, shipping) FROM stdin;
\.


--
-- Data for Name: users; Type: TABLE DATA; Schema: public; Owner: <username>
--

COPY public.users (id, username, password) FROM stdin;
\.


--
-- Name: invoices_id_seq; Type: SEQUENCE SET; Schema: public; Owner: <username>
--

SELECT pg_catalog.setval('public.invoices_id_seq', 1, false);


--
-- Name: users_id_seq; Type: SEQUENCE SET; Schema: public; Owner: <username>
--

SELECT pg_catalog.setval('public.users_id_seq', 1, false);


--
-- Name: invoices invoices_pkey; Type: CONSTRAINT; Schema: public; Owner: <username>
--

ALTER TABLE ONLY public.invoices
    ADD CONSTRAINT invoices_pkey PRIMARY KEY (id);


--
-- Name: users users_pkey; Type: CONSTRAINT; Schema: public; Owner: <username>
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_pkey PRIMARY KEY (id);


--
-- Name: invoices invoices_usr_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: <username>
--

ALTER TABLE ONLY public.invoices
    ADD CONSTRAINT invoices_usr_id_fkey FOREIGN KEY (usr_id) REFERENCES public.users(id) ON DELETE CASCADE;


--
-- PostgreSQL database dump complete
--

