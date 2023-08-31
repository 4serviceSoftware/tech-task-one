--
-- PostgreSQL database dump
--

-- Dumped from database version 14.1
-- Dumped by pg_dump version 14.1

-- Started on 2023-08-24 00:17:48 EEST

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
-- TOC entry 210 (class 1259 OID 40973)
-- Name: nodes; Type: TABLE; Schema: public; Owner: kbnq
--

CREATE TABLE public.nodes (
    id integer NOT NULL,
    name character varying(255) NOT NULL,
    parent_id integer NOT NULL
);


--ALTER TABLE public.nodes OWNER TO kbnq;

--
-- TOC entry 209 (class 1259 OID 40972)
-- Name: nodes_id_seq; Type: SEQUENCE; Schema: public; Owner: kbnq
--

ALTER TABLE public.nodes ALTER COLUMN id ADD GENERATED BY DEFAULT AS IDENTITY (
    SEQUENCE NAME public.nodes_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1
);


--
-- TOC entry 3431 (class 2606 OID 40977)
-- Name: nodes nodes_pkey; Type: CONSTRAINT; Schema: public; Owner: kbnq
--

ALTER TABLE ONLY public.nodes
    ADD CONSTRAINT nodes_pkey PRIMARY KEY (id);


-- Completed on 2023-08-24 00:17:50 EEST

--
-- PostgreSQL database dump complete
--

