--
-- PostgreSQL database dump
--

-- Dumped from database version 16.3 (Homebrew)
-- Dumped by pg_dump version 16.3 (Homebrew)

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
-- Name: invoices; Type: TABLE; Schema: public; Owner: user
--

CREATE TABLE public.invoices (
    id integer NOT NULL,
    fname character varying(50) NOT NULL,
    lname character varying(50) NOT NULL,
    product character varying(255) NOT NULL,
    price numeric(4,2) NOT NULL,
    quantity integer NOT NULL,
    category character varying(255) NOT NULL,
    shipping character varying(255) NOT NULL
);


ALTER TABLE public.invoices OWNER TO user;

--
-- Name: invoices_id_seq; Type: SEQUENCE; Schema: public; Owner: user
--

CREATE SEQUENCE public.invoices_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.invoices_id_seq OWNER TO user;

--
-- Name: invoices_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: user
--

ALTER SEQUENCE public.invoices_id_seq OWNED BY public.invoices.id;


--
-- Name: invoices id; Type: DEFAULT; Schema: public; Owner: user
--

ALTER TABLE ONLY public.invoices ALTER COLUMN id SET DEFAULT nextval('public.invoices_id_seq'::regclass);


--
-- Data for Name: invoices; Type: TABLE DATA; Schema: public; Owner: user
--

COPY public.invoices (id, fname, lname, product, price, quantity, category, shipping) FROM stdin;
1	Dante	Ferges	Safety Goggles	15.99	3	Safety Equipment	423 Elm St, Chicago IL 60629
2	Michael	Wither	Lubricant	11.99	1	Maintenance	230 Furginson Rd, Oklahoma OK 73102
3	Georgei	Ventalin	Door Hinges	12.50	5	Home Improvement	495 Durvington Ave, Topeka KS 66603
4	Edart	Muskrat	Wrench	24.99	1	Plumbing	654 Ulysses Ave, Oklahoma OK 73103
5	Abra	Katern	DiscoBall	19.99	6	Party	829 Sherbert St, Portland ME 04102
6	Charles	Tarly	Zombie Book	14.99	2	Fiction	134 Pluton St, Boston MA 02108
\.


--
-- Name: invoices_id_seq; Type: SEQUENCE SET; Schema: public; Owner: user
--

SELECT pg_catalog.setval('public.invoices_id_seq', 6, true);


--
-- Name: invoices invoices_pkey; Type: CONSTRAINT; Schema: public; Owner: user
--

ALTER TABLE ONLY public.invoices
    ADD CONSTRAINT invoices_pkey PRIMARY KEY (id);


--
-- PostgreSQL database dump complete
--

