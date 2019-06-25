--
-- PostgreSQL database dump
--

-- Dumped from database version 10.9
-- Dumped by pg_dump version 10.9

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

DROP DATABASE IF EXISTS wasim;
--
-- Name: wasim; Type: DATABASE; Schema: -; Owner: -
--

CREATE DATABASE wasim WITH TEMPLATE = template0 ENCODING = 'UTF8' LC_COLLATE = 'en_US.UTF-8' LC_CTYPE = 'en_US.UTF-8';


\connect wasim

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

--
-- Name: plpgsql; Type: EXTENSION; Schema: -; Owner: -
--

CREATE EXTENSION IF NOT EXISTS plpgsql WITH SCHEMA pg_catalog;


--
-- Name: EXTENSION plpgsql; Type: COMMENT; Schema: -; Owner: -
--

COMMENT ON EXTENSION plpgsql IS 'PL/pgSQL procedural language';


SET default_tablespace = '';

SET default_with_oids = false;

--
-- Name: execution_run; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.execution_run (
    op_num bigint NOT NULL,
    run_num integer,
    op_name text NOT NULL,
    arg_count integer,
    result_value bigint,
    op_code integer,
    memory_address bigint,
    local_id integer,
    from_global bigint,
    to_global bigint,
    base_value bigint,
    modifier_value bigint,
    arg_1 bigint,
    arg_2 bigint,
    arg_3 bigint,
    target bigint,
    condition bigint,
    function_name text,
    module_name text,
    program_counter bigint,
    stack_top bigint,
    value numeric,
    preserve_top boolean,
    discard integer,
    condition_met boolean,
    stack_length_start integer,
    stack_length_finish integer,
    function_id integer,
    stack_start jsonb,
    mem_image bytea,
    stack_finish jsonb,
    locals_start jsonb,
    locals_finish jsonb
);


--
-- Name: execution_run_functions; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.execution_run_functions (
    meta_id integer NOT NULL,
    run_num integer,
    function_name text,
    function_num integer,
    num_returns integer,
    num_params integer
);


--
-- Name: execution_run_metadata_meta_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.execution_run_metadata_meta_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: execution_run_metadata_meta_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.execution_run_metadata_meta_id_seq OWNED BY public.execution_run_functions.meta_id;


--
-- Name: execution_runs_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.execution_runs_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: execution_run_functions meta_id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.execution_run_functions ALTER COLUMN meta_id SET DEFAULT nextval('public.execution_run_metadata_meta_id_seq'::regclass);


--
-- PostgreSQL database dump complete
--

