--
-- PostgreSQL database dump
--

-- Dumped from database version 16.2 (Homebrew)
-- Dumped by pg_dump version 16.2 (Homebrew)

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
-- Name: invoices; Type: TABLE; Schema: public; 
--

CREATE TABLE public.invoices (
    fname character varying(255),
    lname character varying(255),
    product character varying(255),
    price numeric(5,2),
    quantity integer,
    category character varying(255),
    shipping character varying(255)
);


ALTER TABLE public.invoices OWNER TO user;

--
-- Data for Name: invoices; Type: TABLE DATA; Schema: public; 
--

COPY public.invoices (fname, lname, product, price, quantity, category, shipping) FROM stdin;
Dante	Ferges	Safety Googles	15.99	3	Safety Equipment	423 Elm St, Chicago IL 60629
Michael	Wither	Lubricant	11.99	1	Maintenance	230 Furginson Rd, Oklahoma OK 73102
Georgei	Ventalin	Door Hinges	12.50	5	Home improvement	495 Durvington Ave, Topeka KS 66603
Edart	Muskrat	Wrench	24.99	1	Plumbing	654 Ulysses Ave, Oklahoma OK 73103
Abra	Katern	DiscoBall	19.99	6	Party	829 Sherbert St, Portland ME 04102
Charles	Tarly	Zombie Book	14.99	2	Fiction	134 Pluton St, Boston MA 02108
\.


--
-- PostgreSQL database dump complete
--

